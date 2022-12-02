package main

import (
	"github.com/vishvananda/netlink"
	"log"
	"time"
)

func maintainLeasesList() {
	for {
		//fmt.Println("leases list:", activeLeasesList)
		for i, x := range activeLeasesList {
			if x.leaseTime >= 1 {
				activeLeasesList[i].leaseTime = x.leaseTime - 1
			}
		}

		for _, x := range activeLeasesList {
			if x.leaseTime == 0 {
				netlinkDelExpiredLease(x)
			}
		}

		var tmpList []leasesStrut
		for _, x := range activeLeasesList {
			if x.leaseTime != 0 {
				tmpList = append(tmpList, x)
			}
		}
		activeLeasesList = tmpList

		time.Sleep(1 * time.Second)
	}
}

func netlinkDelExpiredLease(lease leasesStrut) {
	// Lock the OS Thread, so we don't accidentally switch namespaces
	// get system routes list
	systemRoutes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		log.Fatalln("Failed to get routes")
	}

	delRoute(lease.address, lease.prefixLen, systemRoutes)
}
