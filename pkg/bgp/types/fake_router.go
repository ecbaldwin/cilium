// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package types

import (
	"context"
	"net/netip"
)

type FakeRouter struct {
	paths    map[string]*Path
	policies map[string]*RoutePolicy
	// neighbors holds the configured peers, keyed by peer name.
	neighbors map[string]*Neighbor
	// discovered simulates the addresses a real router resolves for unnumbered
	// peers via IPv6 ND, keyed by interface name.
	discovered map[string]netip.Addr
}

func NewFakeRouter() Router {
	return &FakeRouter{
		paths:      make(map[string]*Path),
		policies:   make(map[string]*RoutePolicy),
		neighbors:  make(map[string]*Neighbor),
		discovered: make(map[string]netip.Addr),
	}
}

func (f *FakeRouter) Stop(ctx context.Context, r StopRequest) {}

func (f *FakeRouter) AddNeighbor(ctx context.Context, n *Neighbor) error {
	f.neighbors[n.Name] = n
	return nil
}

func (f *FakeRouter) UpdateNeighbor(ctx context.Context, n *Neighbor) error {
	f.neighbors[n.Name] = n
	return nil
}

func (f *FakeRouter) RemoveNeighbor(ctx context.Context, n *Neighbor) error {
	delete(f.neighbors, n.Name)
	return nil
}

// DiscoverPeerAddress simulates the router resolving the address of an unnumbered
// peer configured on the provided interface.
func (f *FakeRouter) DiscoverPeerAddress(iface string, addr netip.Addr) {
	f.discovered[iface] = addr
}

func (f *FakeRouter) ResetNeighbor(ctx context.Context, r ResetNeighborRequest) error {
	return nil
}

func (f *FakeRouter) ResetAllNeighbors(ctx context.Context, r ResetAllNeighborsRequest) error {
	return nil
}

func (f *FakeRouter) AdvertisePath(ctx context.Context, p PathRequest) (PathResponse, error) {
	path := p.Path
	f.paths[path.NLRI.String()] = path
	return PathResponse{path}, nil
}

func (f *FakeRouter) WithdrawPath(ctx context.Context, p PathRequest) error {
	path := p.Path
	delete(f.paths, path.NLRI.String())
	return nil
}

func (f *FakeRouter) AddRoutePolicy(ctx context.Context, p RoutePolicyRequest) error {
	f.policies[p.Policy.Name] = p.Policy
	return nil
}

func (f *FakeRouter) RemoveRoutePolicy(ctx context.Context, p RoutePolicyRequest) error {
	delete(f.policies, p.Policy.Name)
	return nil
}

func (f *FakeRouter) GetPeerState(ctx context.Context, r *GetPeerStateRequest) (*GetPeerStateResponse, error) {
	var res GetPeerStateResponse
	for _, n := range f.neighbors {
		state := PeerState{
			Name:      n.Name,
			Address:   n.Address,
			Interface: n.Interface,
		}
		// An unnumbered peer only has the address the router discovered for it.
		if !state.Address.IsValid() && n.Interface != "" {
			state.Address = f.discovered[n.Interface]
		}
		res.Peers = append(res.Peers, state)
	}
	return &res, nil
}

func (f *FakeRouter) GetPeerStateLegacy(ctx context.Context) (GetPeerStateLegacyResponse, error) {
	return GetPeerStateLegacyResponse{}, nil
}

func (f *FakeRouter) GetRoutes(ctx context.Context, r *GetRoutesRequest) (*GetRoutesResponse, error) {
	var routes []*Route
	for _, path := range f.paths {
		routes = append(routes, &Route{
			Prefix: path.NLRI.String(),
			Paths:  []*Path{path},
		})
	}
	return &GetRoutesResponse{Routes: routes}, nil
}

func (f *FakeRouter) GetRoutePolicies(ctx context.Context) (*GetRoutePoliciesResponse, error) {
	var policies []*RoutePolicy
	for _, policy := range f.policies {
		policies = append(policies, policy)
	}
	return &GetRoutePoliciesResponse{Policies: policies}, nil
}

func (f *FakeRouter) GetBGP(ctx context.Context) (GetBGPResponse, error) {
	return GetBGPResponse{}, nil
}
