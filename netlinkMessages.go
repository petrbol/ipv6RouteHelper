package main

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"log"
	"net"
	"strconv"
	"strings"
	"syscall"
)

func delRoute(address string, prefixLen int, systemRoutes []netlink.Route) {
	prefixToDel := address + "/" + strconv.Itoa(prefixLen)

	for _, systemRoute := range systemRoutes {
		if systemRoute.Dst.String() == prefixToDel {
			route := netlink.Route{
				Scope:     netlink.SCOPE_UNIVERSE,
				Dst:       systemRoute.Dst,
				Gw:        systemRoute.Gw,
				LinkIndex: systemRoute.LinkIndex,
				Priority:  netlinkRoutePriority,
			}
			if err := netlink.RouteDel(&route); err == syscall.EEXIST {
				//Ignore this error
			} else if err != nil {
				fmt.Println("err: api remove route failed", err)
			}
			fmt.Println("netlink: del route", prefixToDel, "via", systemRoute.Gw, "index", systemRoute.LinkIndex)
		}
	}
}

func addRoute(address string, prefixLen int, gateway net.IP, linkIndex int) {
	prefixToAdd := address + "/" + strconv.Itoa(prefixLen)
	_, destinationNet, _ := net.ParseCIDR(prefixToAdd)

	route := netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		Dst:       destinationNet,
		Gw:        gateway,
		LinkIndex: linkIndex,
		Priority:  netlinkRoutePriority,
	}

	if err := netlink.RouteAdd(&route); err == syscall.EEXIST {
		//Ignore this error
	} else if err != nil {
		fmt.Println("err: rctKeaRouteHelper add route failed", err)
	}

	fmt.Println("netlink: add route", prefixToAdd, "via", gateway, "index", linkIndex)
}

func netlinkGetSystemRoutes() {
	fmt.Println("info: read system DHCP routes with metric", netlinkRoutePriority)
	// get system routes list
	systemRoutes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		log.Fatalln("Failed to get routes")
	}

	for _, systemRoute := range systemRoutes {
		if systemRoute.Priority == netlinkRoutePriority {
			ipPrefix := strings.Split(systemRoute.Dst.String(), "/")

			if len(ipPrefix) == 2 {
				if cidrMask, err := strconv.Atoi(ipPrefix[1]); err == nil {
					fmt.Println("info: add system route", systemRoute.Dst.String(), "to the dynamic lease list")
					ifName, _ := netlink.LinkByIndex(systemRoute.LinkIndex)
					activeLeasesList = append(activeLeasesList, leasesStrut{
						address:   ipPrefix[0],
						prefixLen: cidrMask,
						leaseTime: unknownRouteExpiration,
						gateway:   systemRoute.Gw,
						ifIndex:   systemRoute.LinkIndex,
						ifName:    ifName.Attrs().Name,
					})
				}
			}
		}
	}
}
