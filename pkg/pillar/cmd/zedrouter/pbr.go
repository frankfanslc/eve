// Copyright (c) 2017 Zededa, Inc.
// SPDX-License-Identifier: Apache-2.0

// Create ip rules and ip routing tables for each ifindex for the bridges used
// for network instances.

package zedrouter

import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"syscall"

	"github.com/lf-edge/eve/pkg/pillar/types"
	"github.com/vishvananda/netlink"
)

var baseTableIndex = 500 // Number tables from here + ifindex

// Call before setting up routeChanges, addrChanges, and linkChanges
func PbrInit(ctx *zedrouterContext) {

	log.Tracef("PbrInit()\n")
}

// PbrRouteAddAll adds all the routes for the bridgeName table to the specific port
// Separately we handle changes in PbrRouteChange
// XXX used by networkinstance only
func PbrRouteAddAll(bridgeName string, port string) error {
	log.Functionf("PbrRouteAddAll(%s, %s)\n", bridgeName, port)

	// for airgap internal switch case
	if port == "" {
		log.Functionf("PbrRouteAddAll: for internal switch, skip for ACL and Route installation\n")
		return nil
	}

	ifindex, err := IfnameToIndex(log, port)
	if err != nil {
		errStr := fmt.Sprintf("IfnameToIndex(%s) failed: %s",
			port, err)
		log.Errorln(errStr)
		return errors.New(errStr)
	}
	link, err := netlink.LinkByName(bridgeName)
	if err != nil {
		errStr := fmt.Sprintf("LinkByName(%s) failed: %s",
			bridgeName, err)
		log.Errorln(errStr)
		return errors.New(errStr)
	}
	index := link.Attrs().Index
	// Add the lowest-prio default-drop route.
	// The route is used to drop all packets otherwise not matched by any route
	// and prevent them from escaping the NI-specific routing table.
	err = AddDefaultDropRoute(index, true)
	if err != nil {
		errStr := fmt.Sprintf("Failed to add default-drop route: %s", err)
		log.Errorln(errStr)
	}
	routes := getAllIPv4Routes(ifindex)
	if routes == nil {
		log.Warnf("PbrRouteAddAll(%s, %s) no routes",
			bridgeName, port)
		return nil
	}
	// Add to ifindex specific table
	ifindex, err = IfnameToIndex(log, bridgeName)
	if err != nil {
		errStr := fmt.Sprintf("IfnameToIndex(%s) failed: %s",
			bridgeName, err)
		log.Errorln(errStr)
		return errors.New(errStr)
	}
	// XXX do they differ? Yes
	if index != ifindex {
		log.Warnf("XXX Different ifindex vs index %d vs %x",
			ifindex, index)
		ifindex = index
	}
	MyTable := baseTableIndex + ifindex
	for _, rt := range routes {
		myrt := rt
		myrt.Table = MyTable
		// Clear any RTNH_F_LINKDOWN etc flags since add doesn't like them
		if rt.Flags != 0 {
			myrt.Flags = 0
		}
		log.Functionf("PbrRouteAddAll(%s, %s) adding %v\n",
			bridgeName, port, myrt)
		if err := netlink.RouteAdd(&myrt); err != nil {
			errStr := fmt.Sprintf("Failed to add %v to %d: %s",
				myrt, myrt.Table, err)
			log.Errorln(errStr)
			return errors.New(errStr)
		}
	}
	return nil
}

// PbrRouteDeleteAll deletes all the routes for the bridgeName table to the specific port
// Separately we handle changes in PbrRouteChange
// XXX used by networkinstance only
// XXX can't we flush the table?
func PbrRouteDeleteAll(bridgeName string, port string) error {
	log.Functionf("PbrRouteDeleteAll(%s, %s)\n", bridgeName, port)

	// for airgap internal switch case
	if port == "" {
		log.Functionf("PbrRouteDeleteAll: for internal switch, skip for ACL and Route deletion\n")
		return nil
	}

	ifindex, err := IfnameToIndex(log, port)
	if err != nil {
		errStr := fmt.Sprintf("IfnameToIndex(%s) failed: %s",
			port, err)
		log.Errorln(errStr)
		return errors.New(errStr)
	}
	routes := getAllIPv4Routes(ifindex)
	if routes == nil {
		log.Warnf("PbrRouteDeleteAll(%s, %s) no routes",
			bridgeName, port)
		return nil
	}
	// Remove from ifindex specific table
	ifindex, err = IfnameToIndex(log, bridgeName)
	if err != nil {
		errStr := fmt.Sprintf("IfnameToIndex(%s) failed: %s",
			bridgeName, err)
		log.Errorln(errStr)
		return errors.New(errStr)
	}
	MyTable := baseTableIndex + ifindex
	for _, rt := range routes {
		myrt := rt
		myrt.Table = MyTable
		// Clear any RTNH_F_LINKDOWN etc flags since del might not like them
		if rt.Flags != 0 {
			myrt.Flags = 0
		}
		log.Functionf("PbrRouteDeleteAll(%s, %s) deleting %v\n",
			bridgeName, port, myrt)
		if err := netlink.RouteDel(&myrt); err != nil {
			errStr := fmt.Sprintf("Failed to delete %v from %d: %s",
				myrt, myrt.Table, err)
			log.Errorln(errStr)
			// We continue to try to delete all
		}
	}
	// Delete the lowest-prio default-drop route.
	err = DelDefaultDropRoute(ifindex, true)
	if err != nil {
		errStr := fmt.Sprintf("Failed to delete default-drop route: %s", err)
		log.Errorln(errStr)
	}
	return nil
}

// Handle a route change
func PbrRouteChange(ctx *zedrouterContext,
	deviceNetworkStatus *types.DeviceNetworkStatus,
	change netlink.RouteUpdate) {

	rt := change.Route
	if rt.Table != getDefaultRouteTable() {
		// Ignore since we will not add to other table
		return
	}
	op := "NONE"
	if change.Type == getRouteUpdateTypeDELROUTE() {
		op = "DELROUTE"
	} else if change.Type == getRouteUpdateTypeNEWROUTE() {
		op = "NEWROUTE"
	}
	ifname, linkType, err := IfindexToName(log, rt.LinkIndex)
	if err != nil {
		log.Errorf("PbrRouteChange IfindexToName failed for %d: %s: route %v\n",
			rt.LinkIndex, err, rt)
		return
	}
	if linkType != "bridge" && !types.IsL3Port(*deviceNetworkStatus, ifname) {
		// Ignore
		log.Functionf("PbrRouteChange ignore %s: neither bridge nor port. route %v\n",
			ifname, rt)
		return
	}
	log.Tracef("RouteChange(%d/%s) %s %+v", rt.LinkIndex, ifname, op, rt)

	// Add to ifindex specific table and to any bridges used by network instances
	myrt := rt
	myrt.Table = baseTableIndex + rt.LinkIndex
	// Clear any RTNH_F_LINKDOWN etc flags since add doesn't like them
	if myrt.Flags != 0 {
		myrt.Flags = 0
	}
	if change.Type == getRouteUpdateTypeDELROUTE() {
		log.Functionf("Received route del %v\n", rt)
		if linkType == "bridge" {
			log.Functionf("Apply route del to bridge %s", ifname)
			if err := netlink.RouteDel(&myrt); err != nil {
				log.Errorf("Failed to remove %v from %d: %s\n",
					myrt, myrt.Table, err)
			}
		}
		// find all bridges for network instances and del for them
		indicies := getAllNIindices(ctx, ifname)
		if len(indicies) != 0 {
			log.Functionf("Apply route del to %v", indicies)
		}
		for _, ifindex := range indicies {
			myrt.Table = baseTableIndex + ifindex
			if err := netlink.RouteDel(&myrt); err != nil {
				log.Errorf("Failed to remove %v from %d: %s\n",
					myrt, myrt.Table, err)
			}
		}
	} else if change.Type == getRouteUpdateTypeNEWROUTE() {
		log.Functionf("Received route add %v\n", rt)
		if linkType == "bridge" {
			log.Functionf("Apply route add to bridge %s", ifname)
			if err := netlink.RouteAdd(&myrt); err != nil {
				// XXX ditto for ENXIO?? for del?
				if isErrno(err, syscall.EEXIST) {
					log.Functionf("Failed to add %v to %d: %s\n",
						myrt, myrt.Table, err)
				} else {
					log.Errorf("Failed to add %v to %d: %s\n",
						myrt, myrt.Table, err)
				}
			}
		}
		// find all bridges for network instances and add for them
		indicies := getAllNIindices(ctx, ifname)
		if len(indicies) != 0 {
			log.Functionf("Apply route add to %v", indicies)
		}
		for _, ifindex := range indicies {
			myrt.Table = baseTableIndex + ifindex
			if err := netlink.RouteAdd(&myrt); err != nil {
				log.Errorf("Failed to add %v to %d: %s\n",
					myrt, myrt.Table, err)
			}
		}
	}
}

func isErrno(err error, errno syscall.Errno) bool {
	e1, ok := err.(syscall.Errno)
	if !ok {
		log.Warnf("XXX not Errno: %T", err)
		return false
	}
	return e1 == errno
}

func AddOverlayRuleAndRoute(bridgeName string, iifIndex int,
	oifIndex int, ipnet *net.IPNet) error {
	log.Tracef("AddOverlayRuleAndRoute: IIF index %d, Prefix %s, OIF index %d",
		iifIndex, ipnet.String(), oifIndex)

	r := netlink.NewRule()
	myTable := baseTableIndex + iifIndex
	r.Table = myTable
	r.IifName = bridgeName
	r.Priority = 10000
	if ipnet.IP.To4() != nil {
		r.Family = syscall.AF_INET
	} else {
		r.Family = syscall.AF_INET6
	}

	// Avoid duplicate rules
	_ = netlink.RuleDel(r)

	// Add rule
	if err := netlink.RuleAdd(r); err != nil {
		errStr := fmt.Sprintf("AddOverlayRuleAndRoute: RuleAdd %v failed with %s", r, err)
		log.Errorln(errStr)
		return errors.New(errStr)
	}

	// Add a the required route to new table that we created above.

	// Setup a route for the current network's subnet to point out of the given oifIndex
	rt := netlink.Route{Dst: ipnet, LinkIndex: oifIndex, Table: myTable, Flags: 0}
	if err := netlink.RouteAdd(&rt); err != nil {
		errStr := fmt.Sprintf("AddOverlayRuleAndRoute: RouteAdd %s failed: %s",
			ipnet.String(), err)
		log.Errorln(errStr)
		return errors.New(errStr)
	}
	return nil
}

// AddFwMarkRuleToDummy : Create an ip rule that sends packets marked by a Drop ACE
// out of interface with given index.
func AddFwMarkRuleToDummy(iifIndex int) error {

	r := netlink.NewRule()
	myTable := baseTableIndex + iifIndex
	r.Table = myTable
	r.Mark = aceDropAction
	r.Mask = aceActionMask
	// This rule gets added during the starting steps of service.
	// Other ip rules corresponding to network instances get added after this
	// and take higher priority. We want this ip rule to match before anything else.
	// Hence we make the priority of this 1000 and the other rules to have 10000.
	r.Priority = 1000

	// Avoid duplicate rules
	_ = netlink.RuleDel(r)

	// Add rule
	if err := netlink.RuleAdd(r); err != nil {
		errStr := fmt.Sprintf("AddFwMarkRuleToDummy: RuleAdd %v failed with %s", r, err)
		log.Errorln(errStr)
		return errors.New(errStr)
	}

	// Add default route that points to dummy interface.
	err := AddDefaultDropRoute(iifIndex, false)
	if err != nil {
		errStr := fmt.Sprintf("AddFwMarkRuleToDummy: AddDefaultDropRoute failed: %s", err)
		log.Errorln(errStr)
		return errors.New(errStr)
	}
	return nil
}

// AddDefaultDropRoute : Add default route dropping packets either by sending them
// into the dummy interface or by using an unreachable destination.
func AddDefaultDropRoute(ifIndex int, unreachable bool) error {
	route, err := makeDefaultDropRoute(ifIndex, unreachable)
	if err != nil {
		return err
	}
	return netlink.RouteAdd(route)
}

// DelDefaultDropRoute : Delete previously added default route dropping packets.
func DelDefaultDropRoute(ifIndex int, unreachable bool) error {
	route, err := makeDefaultDropRoute(ifIndex, unreachable)
	if err != nil {
		return err
	}
	return netlink.RouteDel(route)
}

func makeDefaultDropRoute(ifIndex int, unreachable bool) (*netlink.Route, error) {
	var (
		routeType    int
		outLinkIndex int
	)
	if unreachable {
		routeType = unix.RTN_UNREACHABLE
	} else {
		link, err := netlink.LinkByName(dummyIntfName)
		if err != nil {
			return nil, fmt.Errorf("failed to get dummy interface: %w", err)
		}
		outLinkIndex = link.Attrs().Index
	}

	_, dst, err := net.ParseCIDR("0.0.0.0/0")
	if err != nil {
		return nil, fmt.Errorf("failed to parse dst for default route: %w", err)
	}

	var prio int
	if unreachable {
		// Do not override any actual default route.
		prio = int(^uint32(0))
	}
	return &netlink.Route{
		LinkIndex: outLinkIndex,
		Dst:       dst,
		Priority:  prio,
		Table:     baseTableIndex + ifIndex,
		Type:      routeType,
	}, nil
}
