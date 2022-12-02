package lib

import (
	"github.com/vishvananda/netlink"
	"log"
	"time"
)

func MaintainLeasesList() {
	for {
		//fmt.Println("leases list:", ActiveLeasesList)
		for i, x := range ActiveLeasesList {
			if x.leaseTime >= 1 {
				ActiveLeasesList[i].leaseTime = x.leaseTime - 1
			}
		}

		for _, x := range ActiveLeasesList {
			if x.leaseTime == 0 {
				netlinkDelExpiredLease(x)
			}
		}

		var tmpList []leasesStrut
		for _, x := range ActiveLeasesList {
			if x.leaseTime != 0 {
				tmpList = append(tmpList, x)
			}
		}
		ActiveLeasesList = tmpList

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
