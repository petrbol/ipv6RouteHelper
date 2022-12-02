package main

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"log"
	"net"
	"strconv"
	"strings"
)

func netlinkReconfigureByKea(keaMessage KeaGeneratedStructure) {
	// get neighbor list
	netlinkNeighbors, err := netlink.NeighList(0, netlink.FAMILY_ALL)
	if err != nil {
		log.Fatalln("Failed to get neighbors")
	}

	// get system routes list
	systemRoutes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		log.Fatalln("Failed to get routes")
	}

	for _, x := range keaMessage.Arguments.Leases {
		for _, neighbor := range netlinkNeighbors {
			if neighbor.HardwareAddr.String() == x.HwAddress && strings.Contains(x.IPAddress, ":") == true {
				if strings.Contains(neighbor.IP.String(), "fe80::") == true {
					if x.PrefixLen == 0 && strings.Contains(x.IPAddress, ":") == true {
						x.PrefixLen = 128
					}

					if routeIsDeclared(systemRoutes, x.IPAddress, x.PrefixLen, neighbor.IP, neighbor.LinkIndex) == true {
						fmt.Println("dhcp: route", x.IPAddress+"/"+strconv.Itoa(x.PrefixLen), "via", neighbor.IP, "index", neighbor.LinkIndex, "renewed")
						if prefixInLeaseList(x.IPAddress, x.PrefixLen, neighbor.IP, neighbor.LinkIndex) == true {
							renewLeaseTime(x.IPAddress, x.PrefixLen, x.ValidLft)
						} else {
							addToLeaseList(x.IPAddress, x.PrefixLen, x.ValidLft, neighbor.IP, neighbor.LinkIndex)
						}
					} else {
						fmt.Println("dhcp: route", x.IPAddress+"/"+strconv.Itoa(x.PrefixLen), "not declared or not declared correctly")
						delRoute(x.IPAddress, x.PrefixLen, systemRoutes)
						addRoute(x.IPAddress, x.PrefixLen, neighbor.IP, neighbor.LinkIndex)
						delFromLeaseList(x.IPAddress, x.PrefixLen)
						addToLeaseList(x.IPAddress, x.PrefixLen, x.ValidLft, neighbor.IP, neighbor.LinkIndex)
					}
				}
			}
		}
	}

	for _, x := range keaMessage.Arguments.DeletedLeases {
		if x.PrefixLen == 0 && strings.Contains(x.IPAddress, ":") == true {
			x.PrefixLen = 128
		}

		delRoute(x.IPAddress, x.PrefixLen, systemRoutes)
		delFromLeaseList(x.IPAddress, x.PrefixLen)
	}
}

func routeIsDeclared(systemRoutes []netlink.Route, address string, prefixLen int, gateway net.IP, linkIndex int) bool {
	for _, x := range systemRoutes {
		prefixToFind := address + "/" + strconv.Itoa(prefixLen)
		if x.LinkIndex == linkIndex && x.Gw.String() == gateway.String() && x.Dst.String() == prefixToFind && x.Priority == netlinkRoutePriority {
			return true
		}
	}
	return false
}

func prefixInLeaseList(address string, prefixLen int, gateway net.IP, linkIndex int) bool {
	for _, x := range activeLeasesList {
		if x.address == address && x.prefixLen == prefixLen && x.gateway.String() == gateway.String() && x.ifIndex == linkIndex {
			return true
		}
	}
	return false
}

func renewLeaseTime(address string, prefixLen int, validLft int) {
	for i, x := range activeLeasesList {
		if x.address == address && x.prefixLen == prefixLen {
			activeLeasesList[i].leaseTime = validLft
		}
	}
}

func delFromLeaseList(address string, prefixLen int) {
	var tmpList []leasesStrut
	for _, x := range activeLeasesList {
		if x.address != address && x.prefixLen != prefixLen {
			tmpList = append(tmpList, x)
		}
	}
	activeLeasesList = tmpList
}

func addToLeaseList(address string, prefixLen int, validLft int, gateway net.IP, linkIndex int) {
	ifName, err := netlink.LinkByIndex(linkIndex)
	if err == nil {
		activeLeasesList = append(activeLeasesList, leasesStrut{
			address:   address,
			prefixLen: prefixLen,
			leaseTime: validLft,
			gateway:   gateway,
			ifIndex:   linkIndex,
			ifName:    ifName.Attrs().Name,
		})
	}
}
