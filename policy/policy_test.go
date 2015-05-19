// Copyright (C) 2014,2015 Nippon Telegraph and Telephone Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package policy

import (
	log "github.com/Sirupsen/logrus"
	"github.com/osrg/gobgp/config"
	"github.com/osrg/gobgp/packet"
	"github.com/osrg/gobgp/table"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestPrefixCalcurateNoRange(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	// creatae path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("10.0.0.1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam := []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65001})}
	aspath := bgp.NewPathAttributeAsPath(aspathParam)
	nexthop := bgp.NewPathAttributeNextHop("10.0.0.1")
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{origin, aspath, nexthop, med}
	nlri := []bgp.NLRInfo{*bgp.NewNLRInfo(24, "10.10.0.101")}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	path := table.ProcessMessage(updateMsg, peer)[0]
	// test
	pl1, _ := NewPrefix(net.ParseIP("10.10.0.0"), 24, "")
	match1 := ipPrefixCalculate(path, pl1)
	assert.Equal(t, false, match1)
	pl2, _ := NewPrefix(net.ParseIP("10.10.0.101"), 24, "")
	match2 := ipPrefixCalculate(path, pl2)
	assert.Equal(t, true, match2)
	pl3, _ := NewPrefix(net.ParseIP("10.10.0.0"), 16, "21..24")
	match3 := ipPrefixCalculate(path, pl3)
	assert.Equal(t, true, match3)
}

func TestPrefixCalcurateAddress(t *testing.T) {
	// creatae path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("10.0.0.1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam := []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65001})}
	aspath := bgp.NewPathAttributeAsPath(aspathParam)
	nexthop := bgp.NewPathAttributeNextHop("10.0.0.1")
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{origin, aspath, nexthop, med}
	nlri := []bgp.NLRInfo{*bgp.NewNLRInfo(24, "10.10.0.101")}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	path := table.ProcessMessage(updateMsg, peer)[0]
	// test
	pl1, _ := NewPrefix(net.ParseIP("10.11.0.0"), 16, "21..24")
	match1 := ipPrefixCalculate(path, pl1)
	assert.Equal(t, false, match1)
	pl2, _ := NewPrefix(net.ParseIP("10.10.0.0"), 16, "21..24")
	match2 := ipPrefixCalculate(path, pl2)
	assert.Equal(t, true, match2)
}

func TestPrefixCalcurateLength(t *testing.T) {
	// creatae path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("10.0.0.1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam := []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65001})}
	aspath := bgp.NewPathAttributeAsPath(aspathParam)
	nexthop := bgp.NewPathAttributeNextHop("10.0.0.1")
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{origin, aspath, nexthop, med}
	nlri := []bgp.NLRInfo{*bgp.NewNLRInfo(24, "10.10.0.101")}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	path := table.ProcessMessage(updateMsg, peer)[0]
	// test
	pl1, _ := NewPrefix(net.ParseIP("10.10.64.0"), 24, "21..24")
	match1 := ipPrefixCalculate(path, pl1)
	assert.Equal(t, false, match1)
	pl2, _ := NewPrefix(net.ParseIP("10.10.64.0"), 16, "21..24")
	match2 := ipPrefixCalculate(path, pl2)
	assert.Equal(t, true, match2)
}

func TestPrefixCalcurateLengthRange(t *testing.T) {
	// creatae path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("10.0.0.1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam := []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65001})}
	aspath := bgp.NewPathAttributeAsPath(aspathParam)
	nexthop := bgp.NewPathAttributeNextHop("10.0.0.1")
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{origin, aspath, nexthop, med}
	nlri := []bgp.NLRInfo{*bgp.NewNLRInfo(24, "10.10.0.101")}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	path := table.ProcessMessage(updateMsg, peer)[0]
	// test
	pl1, _ := NewPrefix(net.ParseIP("10.10.0.0"), 16, "21..23")
	match1 := ipPrefixCalculate(path, pl1)
	assert.Equal(t, false, match1)
	pl2, _ := NewPrefix(net.ParseIP("10.10.0.0"), 16, "25..26")
	match2 := ipPrefixCalculate(path, pl2)
	assert.Equal(t, false, match2)
	pl3, _ := NewPrefix(net.ParseIP("10.10.0.0"), 16, "21..24")
	match3 := ipPrefixCalculate(path, pl3)
	assert.Equal(t, true, match3)
}

func TestPrefixCalcurateNoRangeIPv6(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	// creatae path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("2001::192:168:50:1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam := []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65001})}
	aspath := bgp.NewPathAttributeAsPath(aspathParam)
	mpnlri := []bgp.AddrPrefixInterface{bgp.NewIPv6AddrPrefix(64, "2001:123:123:1::")}
	mpreach := bgp.NewPathAttributeMpReachNLRI("2001::192:168:50:1", mpnlri)
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{mpreach, origin, aspath, med}
	nlri := []bgp.NLRInfo{}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	path := table.ProcessMessage(updateMsg, peer)[0]
	// test
	pl1, _ := NewPrefix(net.ParseIP("2001:123:123::"), 48, "")
	match1 := ipPrefixCalculate(path, pl1)
	assert.Equal(t, false, match1)
	pl2, _ := NewPrefix(net.ParseIP("2001:123:123:1::"), 64, "")
	match2 := ipPrefixCalculate(path, pl2)
	assert.Equal(t, true, match2)
	pl3, _ := NewPrefix(net.ParseIP("2001:123:123::"), 48, "64..80")
	match3 := ipPrefixCalculate(path, pl3)
	assert.Equal(t, true, match3)
}

func TestPrefixCalcurateAddressIPv6(t *testing.T) {
	// creatae path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("2001::192:168:50:1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam := []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65001})}
	aspath := bgp.NewPathAttributeAsPath(aspathParam)
	mpnlri := []bgp.AddrPrefixInterface{bgp.NewIPv6AddrPrefix(64, "2001:123:123:1::")}
	mpreach := bgp.NewPathAttributeMpReachNLRI("2001::192:168:50:1", mpnlri)
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{mpreach, origin, aspath, med}
	nlri := []bgp.NLRInfo{}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	path := table.ProcessMessage(updateMsg, peer)[0]
	// test
	pl1, _ := NewPrefix(net.ParseIP("2001:123:128::"), 48, "64..80")
	match1 := ipPrefixCalculate(path, pl1)
	assert.Equal(t, false, match1)
	pl2, _ := NewPrefix(net.ParseIP("2001:123:123::"), 48, "64..80")
	match2 := ipPrefixCalculate(path, pl2)
	assert.Equal(t, true, match2)
}

func TestPrefixCalcurateLengthIPv6(t *testing.T) {
	// creatae path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("2001::192:168:50:1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam := []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65001})}
	aspath := bgp.NewPathAttributeAsPath(aspathParam)
	mpnlri := []bgp.AddrPrefixInterface{bgp.NewIPv6AddrPrefix(64, "2001:123:123:1::")}
	mpreach := bgp.NewPathAttributeMpReachNLRI("2001::192:168:50:1", mpnlri)
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{mpreach, origin, aspath, med}
	nlri := []bgp.NLRInfo{}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	path := table.ProcessMessage(updateMsg, peer)[0]
	// test
	pl1, _ := NewPrefix(net.ParseIP("2001:123:123:64::"), 64, "64..80")
	match1 := ipPrefixCalculate(path, pl1)
	assert.Equal(t, false, match1)
	pl2, _ := NewPrefix(net.ParseIP("2001:123:123:64::"), 48, "64..80")
	match2 := ipPrefixCalculate(path, pl2)
	assert.Equal(t, true, match2)
}

func TestPrefixCalcurateLengthRangeIPv6(t *testing.T) {
	// creatae path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("2001::192:168:50:1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam := []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65001})}
	aspath := bgp.NewPathAttributeAsPath(aspathParam)
	mpnlri := []bgp.AddrPrefixInterface{bgp.NewIPv6AddrPrefix(64, "2001:123:123:1::")}
	mpreach := bgp.NewPathAttributeMpReachNLRI("2001::192:168:50:1", mpnlri)
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{mpreach, origin, aspath, med}
	nlri := []bgp.NLRInfo{}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	path := table.ProcessMessage(updateMsg, peer)[0]
	// test
	pl1, _ := NewPrefix(net.ParseIP("2001:123:123::"), 48, "62..63")
	match1 := ipPrefixCalculate(path, pl1)
	assert.Equal(t, false, match1)
	pl2, _ := NewPrefix(net.ParseIP("2001:123:123::"), 48, "65..66")
	match2 := ipPrefixCalculate(path, pl2)
	assert.Equal(t, false, match2)
	pl3, _ := NewPrefix(net.ParseIP("2001:123:123::"), 48, "63..65")
	match3 := ipPrefixCalculate(path, pl3)
	assert.Equal(t, true, match3)
}

func TestPolicyNotMatch(t *testing.T) {
	// creatae path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("10.0.0.1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam := []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65001})}
	aspath := bgp.NewPathAttributeAsPath(aspathParam)
	nexthop := bgp.NewPathAttributeNextHop("10.0.0.1")
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{origin, aspath, nexthop, med}
	nlri := []bgp.NLRInfo{*bgp.NewNLRInfo(24, "10.10.0.101")}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	path := table.ProcessMessage(updateMsg, peer)[0]
	// create policy
	ps := config.PrefixSet{
		PrefixSetName: "ps1",
		PrefixList: []config.Prefix{
			config.Prefix{
				Address:         net.ParseIP("10.3.0.0"),
				Masklength:      16,
				MasklengthRange: "21..24",
			}},
	}
	ns := config.NeighborSet{
		NeighborSetName: "ns1",
		NeighborInfoList: []config.NeighborInfo{
			config.NeighborInfo{
				Address: net.ParseIP("10.0.0.1"),
			}},
	}
	ds := config.DefinedSets{
		PrefixSetList:   []config.PrefixSet{ps},
		NeighborSetList: []config.NeighborSet{ns},
	}
	s := config.Statement{
		Name: "statement1",
		Conditions: config.Conditions{
			MatchPrefixSet:   "ps1",
			MatchNeighborSet: "ns1",
			MatchSetOptions:  config.MATCH_SET_OPTIONS_TYPE_ALL,
		},
		Actions: config.Actions{
			AcceptRoute: false,
			RejectRoute: true,
		},
	}
	pd := config.PolicyDefinition{"pd1", []config.Statement{s}}
	pl := config.RoutingPolicy{ds, []config.PolicyDefinition{pd}}
	//test
	df := pl.DefinedSets
	p := NewPolicy(pl.PolicyDefinitionList[0], df)
	match, pType, newPath := p.Apply(path)
	assert.Equal(t, false, match)
	assert.Equal(t, ROUTE_TYPE_NONE, pType)
	assert.Equal(t, nil, newPath)
}

func TestPolicyMatchAndReject(t *testing.T) {
	// creatae path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("10.0.0.1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam := []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65001})}
	aspath := bgp.NewPathAttributeAsPath(aspathParam)
	nexthop := bgp.NewPathAttributeNextHop("10.0.0.1")
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{origin, aspath, nexthop, med}
	nlri := []bgp.NLRInfo{*bgp.NewNLRInfo(24, "10.10.0.101")}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	path := table.ProcessMessage(updateMsg, peer)[0]
	// create policy
	ps := config.PrefixSet{
		PrefixSetName: "ps1",
		PrefixList: []config.Prefix{
			config.Prefix{
				Address:         net.ParseIP("10.10.0.0"),
				Masklength:      16,
				MasklengthRange: "21..24",
			}},
	}
	ns := config.NeighborSet{
		NeighborSetName: "ns1",
		NeighborInfoList: []config.NeighborInfo{
			config.NeighborInfo{
				Address: net.ParseIP("10.0.0.1"),
			}},
	}
	ds := config.DefinedSets{
		PrefixSetList:   []config.PrefixSet{ps},
		NeighborSetList: []config.NeighborSet{ns},
	}
	s := config.Statement{
		Name: "statement1",
		Conditions: config.Conditions{
			MatchPrefixSet:   "ps1",
			MatchNeighborSet: "ns1",
			MatchSetOptions:  config.MATCH_SET_OPTIONS_TYPE_ALL,
		},
		Actions: config.Actions{
			AcceptRoute: false,
			RejectRoute: true,
		},
	}
	pd := config.PolicyDefinition{"pd1", []config.Statement{s}}
	pl := config.RoutingPolicy{ds, []config.PolicyDefinition{pd}}
	//test
	df := pl.DefinedSets
	p := NewPolicy(pl.PolicyDefinitionList[0], df)
	match, pType, newPath := p.Apply(path)
	assert.Equal(t, true, match)
	assert.Equal(t, ROUTE_TYPE_REJECT, pType)
	assert.Equal(t, nil, newPath)
}

func TestPolicyMatchAndAccept(t *testing.T) {
	// creatae path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("10.0.0.1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam := []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65001})}
	aspath := bgp.NewPathAttributeAsPath(aspathParam)
	nexthop := bgp.NewPathAttributeNextHop("10.0.0.1")
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{origin, aspath, nexthop, med}
	nlri := []bgp.NLRInfo{*bgp.NewNLRInfo(24, "10.10.0.101")}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	path := table.ProcessMessage(updateMsg, peer)[0]
	// create policy
	ps := config.PrefixSet{
		PrefixSetName: "ps1",
		PrefixList: []config.Prefix{
			config.Prefix{
				Address:         net.ParseIP("10.10.0.0"),
				Masklength:      16,
				MasklengthRange: "21..24",
			}},
	}
	ns := config.NeighborSet{
		NeighborSetName: "ns1",
		NeighborInfoList: []config.NeighborInfo{
			config.NeighborInfo{
				Address: net.ParseIP("10.0.0.1"),
			}},
	}
	ds := config.DefinedSets{
		PrefixSetList:   []config.PrefixSet{ps},
		NeighborSetList: []config.NeighborSet{ns},
	}
	s := config.Statement{
		Name: "statement1",
		Conditions: config.Conditions{
			MatchPrefixSet:   "ps1",
			MatchNeighborSet: "ns1",
			MatchSetOptions:  config.MATCH_SET_OPTIONS_TYPE_ALL,
		},
		Actions: config.Actions{
			AcceptRoute: true,
			RejectRoute: false,
		},
	}
	pd := config.PolicyDefinition{"pd1", []config.Statement{s}}
	pl := config.RoutingPolicy{ds, []config.PolicyDefinition{pd}}
	//test
	df := pl.DefinedSets
	p := NewPolicy(pl.PolicyDefinitionList[0], df)
	match, pType, newPath := p.Apply(path)
	assert.Equal(t, true, match)
	assert.Equal(t, ROUTE_TYPE_ACCEPT, pType)
	assert.Equal(t, path, newPath)
}

func TestPolicyRejectOnlyPrefixSet(t *testing.T) {
	// creatae path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("10.0.1.1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam := []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65001})}
	aspath := bgp.NewPathAttributeAsPath(aspathParam)
	nexthop := bgp.NewPathAttributeNextHop("10.0.1.1")
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{origin, aspath, nexthop, med}
	nlri := []bgp.NLRInfo{*bgp.NewNLRInfo(24, "10.10.1.101")}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	path1 := table.ProcessMessage(updateMsg, peer)[0]

	peer = &table.PeerInfo{AS: 65002, Address: net.ParseIP("10.0.2.2")}
	origin = bgp.NewPathAttributeOrigin(0)
	aspathParam = []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65002})}
	aspath = bgp.NewPathAttributeAsPath(aspathParam)
	nexthop = bgp.NewPathAttributeNextHop("10.0.2.2")
	med = bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes = []bgp.PathAttributeInterface{origin, aspath, nexthop, med}
	nlri = []bgp.NLRInfo{*bgp.NewNLRInfo(24, "10.9.2.102")}
	withdrawnRoutes = []bgp.WithdrawnRoute{}
	updateMsg = bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	path2 := table.ProcessMessage(updateMsg, peer)[0]

	// create policy
	ps := config.PrefixSet{
		PrefixSetName: "ps1",
		PrefixList: []config.Prefix{
			config.Prefix{
				Address:         net.ParseIP("10.10.1.0"),
				Masklength:      16,
				MasklengthRange: "21..24",
			}},
	}
	ds := config.DefinedSets{
		PrefixSetList:   []config.PrefixSet{ps},
		NeighborSetList: []config.NeighborSet{},
	}
	s := config.Statement{
		Name: "statement1",
		Conditions: config.Conditions{
			MatchPrefixSet:   "ps1",
			MatchNeighborSet: "ns1",
			MatchSetOptions:  config.MATCH_SET_OPTIONS_TYPE_ALL,
		},
		Actions: config.Actions{
			AcceptRoute: false,
			RejectRoute: true,
		},
	}
	pd := config.PolicyDefinition{"pd1", []config.Statement{s}}
	pl := config.RoutingPolicy{ds, []config.PolicyDefinition{pd}}
	//test
	df := pl.DefinedSets
	p := NewPolicy(pl.PolicyDefinitionList[0], df)
	match, pType, newPath := p.Apply(path1)
	assert.Equal(t, true, match)
	assert.Equal(t, ROUTE_TYPE_REJECT, pType)
	assert.Equal(t, nil, newPath)

	match2, pType2, newPath2 := p.Apply(path2)
	assert.Equal(t, false, match2)
	assert.Equal(t, ROUTE_TYPE_NONE, pType2)
	assert.Equal(t, nil, newPath2)
}

func TestPolicyRejectOnlyNeighborSet(t *testing.T) {
	// creatae path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("10.0.1.1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam := []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65001})}
	aspath := bgp.NewPathAttributeAsPath(aspathParam)
	nexthop := bgp.NewPathAttributeNextHop("10.0.1.1")
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{origin, aspath, nexthop, med}
	nlri := []bgp.NLRInfo{*bgp.NewNLRInfo(24, "10.10.1.101")}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	path1 := table.ProcessMessage(updateMsg, peer)[0]

	peer = &table.PeerInfo{AS: 65002, Address: net.ParseIP("10.0.2.2")}
	origin = bgp.NewPathAttributeOrigin(0)
	aspathParam = []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65002})}
	aspath = bgp.NewPathAttributeAsPath(aspathParam)
	nexthop = bgp.NewPathAttributeNextHop("10.0.2.2")
	med = bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes = []bgp.PathAttributeInterface{origin, aspath, nexthop, med}
	nlri = []bgp.NLRInfo{*bgp.NewNLRInfo(24, "10.10.2.102")}
	withdrawnRoutes = []bgp.WithdrawnRoute{}
	updateMsg = bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	path2 := table.ProcessMessage(updateMsg, peer)[0]

	// create policy
	ns := config.NeighborSet{
		NeighborSetName: "ns1",
		NeighborInfoList: []config.NeighborInfo{
			config.NeighborInfo{
				Address: net.ParseIP("10.0.1.1"),
			}},
	}
	ds := config.DefinedSets{
		PrefixSetList:   []config.PrefixSet{},
		NeighborSetList: []config.NeighborSet{ns},
	}
	s := config.Statement{
		Name: "statement1",
		Conditions: config.Conditions{
			MatchPrefixSet:   "ps1",
			MatchNeighborSet: "ns1",
			MatchSetOptions:  config.MATCH_SET_OPTIONS_TYPE_ALL,
		},
		Actions: config.Actions{
			AcceptRoute: false,
			RejectRoute: true,
		},
	}
	pd := config.PolicyDefinition{"pd1", []config.Statement{s}}
	pl := config.RoutingPolicy{ds, []config.PolicyDefinition{pd}}
	//test
	df := pl.DefinedSets
	p := NewPolicy(pl.PolicyDefinitionList[0], df)
	match, pType, newPath := p.Apply(path1)
	assert.Equal(t, true, match)
	assert.Equal(t, ROUTE_TYPE_REJECT, pType)
	assert.Equal(t, nil, newPath)

	match2, pType2, newPath2 := p.Apply(path2)
	assert.Equal(t, false, match2)
	assert.Equal(t, ROUTE_TYPE_NONE, pType2)
	assert.Equal(t, nil, newPath2)
}

func TestPolicyDifferentRoutefamilyOfPathAndPolicy(t *testing.T) {
	// creatae path ipv4
	peerIPv4 := &table.PeerInfo{AS: 65001, Address: net.ParseIP("10.0.0.1")}
	originIPv4 := bgp.NewPathAttributeOrigin(0)
	aspathParamIPv4 := []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65001})}
	aspathIPv4 := bgp.NewPathAttributeAsPath(aspathParamIPv4)
	nexthopIPv4 := bgp.NewPathAttributeNextHop("10.0.0.1")
	medIPv4 := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributesIPv4 := []bgp.PathAttributeInterface{originIPv4, aspathIPv4, nexthopIPv4, medIPv4}
	nlriIPv4 := []bgp.NLRInfo{*bgp.NewNLRInfo(24, "10.10.0.101")}
	withdrawnRoutesIPv4 := []bgp.WithdrawnRoute{}
	updateMsgIPv4 := bgp.NewBGPUpdateMessage(withdrawnRoutesIPv4, pathAttributesIPv4, nlriIPv4)
	pathIPv4 := table.ProcessMessage(updateMsgIPv4, peerIPv4)[0]
	// creatae path ipv6
	peerIPv6 := &table.PeerInfo{AS: 65001, Address: net.ParseIP("2001::192:168:50:1")}
	originIPv6 := bgp.NewPathAttributeOrigin(0)
	aspathParamIPv6 := []bgp.AsPathParamInterface{bgp.NewAsPathParam(2, []uint16{65001})}
	aspathIPv6 := bgp.NewPathAttributeAsPath(aspathParamIPv6)
	mpnlriIPv6 := []bgp.AddrPrefixInterface{bgp.NewIPv6AddrPrefix(64, "2001:123:123:1::")}
	mpreachIPv6 := bgp.NewPathAttributeMpReachNLRI("2001::192:168:50:1", mpnlriIPv6)
	medIPv6 := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributesIPv6 := []bgp.PathAttributeInterface{mpreachIPv6, originIPv6, aspathIPv6, medIPv6}
	nlriIPv6 := []bgp.NLRInfo{}
	withdrawnRoutesIPv6 := []bgp.WithdrawnRoute{}
	updateMsgIPv6 := bgp.NewBGPUpdateMessage(withdrawnRoutesIPv6, pathAttributesIPv6, nlriIPv6)
	pathIPv6 := table.ProcessMessage(updateMsgIPv6, peerIPv6)[0]
	// create policy
	psIPv4 := config.PrefixSet{
		PrefixSetName: "psIPv4",
		PrefixList: []config.Prefix{
			config.Prefix{
				Address:         net.ParseIP("10.10.0.0"),
				Masklength:      16,
				MasklengthRange: "21..24",
			}},
	}
	nsIPv4 := config.NeighborSet{
		NeighborSetName: "nsIPv4",
		NeighborInfoList: []config.NeighborInfo{
			config.NeighborInfo{
				Address: net.ParseIP("10.0.0.1"),
			}},
	}
	psIPv6 := config.PrefixSet{
		PrefixSetName: "psIPv6",
		PrefixList: []config.Prefix{
			config.Prefix{
				Address:         net.ParseIP("2001:123:123::"),
				Masklength:      48,
				MasklengthRange: "64..80",
			}},
	}
	nsIPv6 := config.NeighborSet{
		NeighborSetName: "nsIPv6",
		NeighborInfoList: []config.NeighborInfo{
			config.NeighborInfo{
				Address: net.ParseIP("2001::192:168:50:1"),
			}},
	}
	ds := config.DefinedSets{
		PrefixSetList:   []config.PrefixSet{psIPv4, psIPv6},
		NeighborSetList: []config.NeighborSet{nsIPv4, nsIPv6},
	}
	stIPv4 := config.Statement{
		Name: "statement1",
		Conditions: config.Conditions{
			MatchPrefixSet:   "psIPv4",
			MatchNeighborSet: "nsIPv4",
			MatchSetOptions:  config.MATCH_SET_OPTIONS_TYPE_ALL,
		},
		Actions: config.Actions{
			AcceptRoute: false,
			RejectRoute: true,
		},
	}
	stIPv6 := config.Statement{
		Name: "statement2",
		Conditions: config.Conditions{
			MatchPrefixSet:   "psIPv6",
			MatchNeighborSet: "nsIPv6",
			MatchSetOptions:  config.MATCH_SET_OPTIONS_TYPE_ALL,
		},
		Actions: config.Actions{
			AcceptRoute: false,
			RejectRoute: true,
		},
	}
	pd := config.PolicyDefinition{"pd1", []config.Statement{stIPv4, stIPv6}}
	pl := config.RoutingPolicy{ds, []config.PolicyDefinition{pd}}
	//test
	df := pl.DefinedSets
	p := NewPolicy(pl.PolicyDefinitionList[0], df)
	match1, pType1, newPath1 := p.Apply(pathIPv4)
	assert.Equal(t, true, match1)
	assert.Equal(t, ROUTE_TYPE_REJECT, pType1)
	assert.Equal(t, nil, newPath1)

	match2, pType2, newPath2 := p.Apply(pathIPv6)
	assert.Equal(t, true, match2)
	assert.Equal(t, ROUTE_TYPE_REJECT, pType2)
	assert.Equal(t, nil, newPath2)
}

func TestAsPathLengthConditionEvaluate(t *testing.T) {
	// setup
	// create path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("10.0.0.1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam := []bgp.AsPathParamInterface{
		bgp.NewAsPathParam(2, []uint16{65001, 65000, 65004, 65005}),
		bgp.NewAsPathParam(1, []uint16{65001, 65000, 65004, 65005}),
	}
	aspath := bgp.NewPathAttributeAsPath(aspathParam)
	nexthop := bgp.NewPathAttributeNextHop("10.0.0.1")
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{origin, aspath, nexthop, med}
	nlri := []bgp.NLRInfo{*bgp.NewNLRInfo(24, "10.10.0.101")}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	table.UpdatePathAttrs4ByteAs(updateMsg.Body.(*bgp.BGPUpdate))
	path := table.ProcessMessage(updateMsg, peer)[0]

	// create match condition
	asPathLength := config.AsPathLength{
		Operator: "eq",
		Value:    8,
	}
	c := NewAsPathLengthCondition(asPathLength)

	// test
	assert.Equal(t, true, c.evaluate(path))

	// create match condition
	asPathLength = config.AsPathLength{
		Operator: "ge",
		Value:    3,
	}
	c = NewAsPathLengthCondition(asPathLength)

	// test
	assert.Equal(t, true, c.evaluate(path))

	// create match condition
	asPathLength = config.AsPathLength{
		Operator: "le",
		Value:    3,
	}
	c = NewAsPathLengthCondition(asPathLength)

	// test
	assert.Equal(t, false, c.evaluate(path))
}

func TestAsPathLengthConditionWithOtherCondition(t *testing.T) {
	// setup
	// create path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("10.0.0.1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam := []bgp.AsPathParamInterface{
		bgp.NewAsPathParam(2, []uint16{65001, 65000, 65004, 65004, 65005}),
		bgp.NewAsPathParam(1, []uint16{65001, 65000, 65004, 65005}),
	}
	aspath := bgp.NewPathAttributeAsPath(aspathParam)
	nexthop := bgp.NewPathAttributeNextHop("10.0.0.1")
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{origin, aspath, nexthop, med}
	nlri := []bgp.NLRInfo{*bgp.NewNLRInfo(24, "10.10.0.101")}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	table.UpdatePathAttrs4ByteAs(updateMsg.Body.(*bgp.BGPUpdate))
	path := table.ProcessMessage(updateMsg, peer)[0]

	// create policy
	ps := config.PrefixSet{
		PrefixSetName: "ps1",
		PrefixList: []config.Prefix{
			config.Prefix{
				Address:         net.ParseIP("10.10.1.0"),
				Masklength:      16,
				MasklengthRange: "21..24",
			}},
	}
	ns := config.NeighborSet{
		NeighborSetName: "ns1",
		NeighborInfoList: []config.NeighborInfo{
			config.NeighborInfo{
				Address: net.ParseIP("10.0.1.1"),
			}},
	}

	ds := config.DefinedSets{
		PrefixSetList:   []config.PrefixSet{ps},
		NeighborSetList: []config.NeighborSet{ns},
	}

	// create match condition
	asPathLength := config.AsPathLength{
		Operator: "le",
		Value:    10,
	}

	bgpCondition := config.BgpConditions{
		AsPathLength: asPathLength,
	}

	s := config.Statement{
		Name: "statement1",
		Conditions: config.Conditions{
			MatchPrefixSet:   "ps1",
			MatchNeighborSet: "ns1",
			MatchSetOptions:  config.MATCH_SET_OPTIONS_TYPE_ANY,
			BgpConditions:    bgpCondition,
		},
		Actions: config.Actions{
			RejectRoute: true,
		},
	}
	pd := config.PolicyDefinition{"pd1", []config.Statement{s}}
	pl := config.RoutingPolicy{ds, []config.PolicyDefinition{pd}}

	//test
	df := pl.DefinedSets
	p := NewPolicy(pl.PolicyDefinitionList[0], df)
	match, pType, newPath := p.Apply(path)
	assert.Equal(t, true, match)
	assert.Equal(t, ROUTE_TYPE_REJECT, pType)
	assert.Equal(t, nil, newPath)

}

func TestAsPathConditionEvaluate(t *testing.T) {

	// setup
	// create path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("10.0.0.1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam1 := []bgp.AsPathParamInterface{
		bgp.NewAsPathParam(2, []uint16{65001, 65000, 65004, 65005}),
		bgp.NewAsPathParam(1, []uint16{65001, 65010, 65004, 65005}),
	}
	aspath := bgp.NewPathAttributeAsPath(aspathParam1)
	nexthop := bgp.NewPathAttributeNextHop("10.0.0.1")
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{origin, aspath, nexthop, med}
	nlri := []bgp.NLRInfo{*bgp.NewNLRInfo(24, "10.10.0.101")}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg1 := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	table.UpdatePathAttrs4ByteAs(updateMsg1.Body.(*bgp.BGPUpdate))
	path1 := table.ProcessMessage(updateMsg1, peer)[0]

	aspathParam2 := []bgp.AsPathParamInterface{
		bgp.NewAsPathParam(2, []uint16{65010}),
		bgp.NewAsPathParam(1, []uint16{65010}),
	}
	aspath2 := bgp.NewPathAttributeAsPath(aspathParam2)
	pathAttributes = []bgp.PathAttributeInterface{origin, aspath2, nexthop, med}
	updateMsg2 := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	table.UpdatePathAttrs4ByteAs(updateMsg2.Body.(*bgp.BGPUpdate))
	path2 := table.ProcessMessage(updateMsg2, peer)[0]

	// create match condition
	asPathSet1 := config.AsPathSet{
		AsPathSetName:    "asset1",
		AsPathSetMembers: []string{"^65001"},
	}

	asPathSet2 := config.AsPathSet{
		AsPathSetName:    "asset2",
		AsPathSetMembers: []string{"65005$"},
	}

	asPathSet3 := config.AsPathSet{
		AsPathSetName:    "asset3",
		AsPathSetMembers: []string{"65004", "65005$"},
	}

	asPathSet4 := config.AsPathSet{
		AsPathSetName:    "asset4",
		AsPathSetMembers: []string{"65000$"},
	}

	asPathSet5 := config.AsPathSet{
		AsPathSetName:    "asset5",
		AsPathSetMembers: []string{"65010"},
	}

	asPathSet6 := config.AsPathSet{
		AsPathSetName:    "asset6",
		AsPathSetMembers: []string{"^65010$"},
	}

	asPathSetList := []config.AsPathSet{asPathSet1, asPathSet2, asPathSet3,
		asPathSet4, asPathSet5, asPathSet6}

	p1 := NewAsPathCondition("asset1", asPathSetList)
	p2 := NewAsPathCondition("asset2", asPathSetList)
	p3 := NewAsPathCondition("asset3", asPathSetList)
	p4 := NewAsPathCondition("asset4", asPathSetList)
	p5 := NewAsPathCondition("asset5", asPathSetList)
	p6 := NewAsPathCondition("asset6", asPathSetList)

	// test
	assert.Equal(t, true, p1.evaluate(path1))
	assert.Equal(t, true, p2.evaluate(path1))
	assert.Equal(t, true, p3.evaluate(path1))
	assert.Equal(t, false, p4.evaluate(path1))
	assert.Equal(t, false, p5.evaluate(path1))
	assert.Equal(t, false, p6.evaluate(path1))
	assert.Equal(t, true, p6.evaluate(path2))

}

func TestAsPathConditionWithOtherCondition(t *testing.T) {

	// setup
	// create path
	peer := &table.PeerInfo{AS: 65001, Address: net.ParseIP("10.0.0.1")}
	origin := bgp.NewPathAttributeOrigin(0)
	aspathParam := []bgp.AsPathParamInterface{
		bgp.NewAsPathParam(2, []uint16{65001, 65000, 65004, 65004, 65005}),
		bgp.NewAsPathParam(1, []uint16{65001, 65000, 65004, 65005}),
	}
	aspath := bgp.NewPathAttributeAsPath(aspathParam)
	nexthop := bgp.NewPathAttributeNextHop("10.0.0.1")
	med := bgp.NewPathAttributeMultiExitDisc(0)
	pathAttributes := []bgp.PathAttributeInterface{origin, aspath, nexthop, med}
	nlri := []bgp.NLRInfo{*bgp.NewNLRInfo(24, "10.10.0.101")}
	withdrawnRoutes := []bgp.WithdrawnRoute{}
	updateMsg := bgp.NewBGPUpdateMessage(withdrawnRoutes, pathAttributes, nlri)
	table.UpdatePathAttrs4ByteAs(updateMsg.Body.(*bgp.BGPUpdate))
	path := table.ProcessMessage(updateMsg, peer)[0]

	// create policy
	asPathSet := config.AsPathSet{
		AsPathSetName:    "asset1",
		AsPathSetMembers: []string{"65005$"},
	}

	prefixSet := config.PrefixSet{
		PrefixSetName: "ps1",
		PrefixList: []config.Prefix{
			config.Prefix{
				Address:         net.ParseIP("10.11.1.0"),
				Masklength:      16,
				MasklengthRange: "21..24",
			}},
	}

	neighborSet := config.NeighborSet{
		NeighborSetName: "ns1",
		NeighborInfoList: []config.NeighborInfo{
			config.NeighborInfo{
				Address: net.ParseIP("10.2.1.1"),
			}},
	}

	ds := config.DefinedSets{
		PrefixSetList:   []config.PrefixSet{prefixSet},
		NeighborSetList: []config.NeighborSet{neighborSet},
		BgpDefinedSets: config.BgpDefinedSets{
			AsPathSetList: []config.AsPathSet{asPathSet},
		},
	}

	s := config.Statement{
		Name: "statement1",
		Conditions: config.Conditions{
			MatchPrefixSet:   "ps1",
			MatchNeighborSet: "ns1",
			BgpConditions: config.BgpConditions{
				MatchAsPathSet: "asset1",
			},
			MatchSetOptions: config.MATCH_SET_OPTIONS_TYPE_ANY,
		},
		Actions: config.Actions{
			AcceptRoute: false,
			RejectRoute: true,
		},
	}

	pd := config.PolicyDefinition{"pd1", []config.Statement{s}}
	pl := config.RoutingPolicy{
		DefinedSets:          ds,
		PolicyDefinitionList: []config.PolicyDefinition{pd},
	}

	//test
	df := pl.DefinedSets
	p := NewPolicy(pl.PolicyDefinitionList[0], df)
	match, pType, newPath := p.Apply(path)
	assert.Equal(t, true, match)
	assert.Equal(t, ROUTE_TYPE_REJECT, pType)
	assert.Equal(t, nil, newPath)

}