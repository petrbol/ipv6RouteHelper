package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
)

func main() {
	flag.IntVar(&netlinkRoutePriority, "priority", 4096, "default DHCP router helper priority")
	flag.IntVar(&unknownRouteExpiration, "expiration", 86400, "automatic remove not renewed route after timeout")
	flag.IntVar(&listenPort, "port", 8082, "listen TCP port for Kea messages")
	flag.StringVar(&listenAddress, "address", "127.0.0.1", "default listen ip address")
	flag.StringVar(&savedLeasesListFile, "leasesFile", "/etc/rct/savedLeaseRouteHelper.json", "where to store backup leases file")
	flag.Parse()

	if listenAddress != net.ParseIP(listenAddress).String() {
		log.Fatalln("invalid listen address")
	}

	restoreSavedState()
	netlinkGetSystemRoutes()

	go maintainLeasesList()
	go compareAndSaveCurrentState()

	http.HandleFunc("/", receiveKeaPostJson)
	log.Fatal(http.ListenAndServe(listenAddress+":"+strconv.Itoa(listenPort), nil))
}

func receiveKeaPostJson(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var keaMessage KeaGeneratedStructure
	err = json.Unmarshal(body, &keaMessage)
	if err != nil {
		log.Println("receive non json message\t", string(body))
	}

	netlinkReconfigureByKea(keaMessage)
}
