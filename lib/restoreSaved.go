package lib

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"strconv"
)

func RestoreSavedState() {
	fmt.Println("info: read last leases saved state")
	// read file and restore to the routing table
	readConfiguration()

	for _, x := range SavedLeasesList {
		fmt.Println("info: restore saved lease for prefix:", x.Address+"/"+strconv.Itoa(x.PrefixLen))
		netlinkRestoreSavedLease(x)
	}
}

func netlinkRestoreSavedLease(lease savedLeasesStrut) {
	currentIfIndex, err := netlink.LinkByName(lease.IfName)
	if err == nil {
		if lease.IfIndex != currentIfIndex.Attrs().Index {
			lease.IfIndex = currentIfIndex.Attrs().Index
		}
	}

	addRoute(lease.Address, lease.PrefixLen, lease.Gateway, lease.IfIndex)
}
