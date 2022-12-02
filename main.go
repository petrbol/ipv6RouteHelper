package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"precioz.net/kea6RouteHelperPublicVersion/lib"
	"strconv"
)

func main() {
	flag.IntVar(&lib.NetlinkRoutePriority, "priority", 4096, "default DHCP router helper priority")
	flag.IntVar(&lib.UnknownRouteExpiration, "expiration", 86400, "automatic remove not renewed route after timeout")
	flag.IntVar(&lib.ListenPort, "port", 8082, "listen TCP port for Kea messages")
	flag.StringVar(&lib.ListenAddress, "address", "127.0.0.1", "default listen ip address")
	flag.StringVar(&lib.SavedLeasesListFile, "leasesFile", "/etc/rct/savedLeaseRouteHelper.json", "where to store backup leases file")
	flag.Parse()

	if lib.ListenAddress != net.ParseIP(lib.ListenAddress).String() {
		log.Fatalln("invalid listen address")
	}

	lib.RestoreSavedState()
	lib.NetlinkGetSystemRoutes()

	go lib.MaintainLeasesList()
	go lib.CompareAndSaveCurrentState()

	http.HandleFunc("/", lib.ReceiveKeaPostJson)
	log.Fatal(http.ListenAndServe(lib.ListenAddress+":"+strconv.Itoa(lib.ListenPort), nil))
}
