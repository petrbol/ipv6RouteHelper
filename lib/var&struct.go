package lib

import "net"

var (
	NetlinkRoutePriority   int
	UnknownRouteExpiration int
	ListenPort             int
	ListenAddress          string
	ActiveLeasesList       []leasesStrut
	SavedLeasesList        []savedLeasesStrut
	SavedLeasesListFile    string
)

type leasesStrut struct {
	address   string
	prefixLen int
	leaseTime int
	gateway   net.IP
	ifIndex   int
	ifName    string
}

type savedLeasesStrut struct {
	Address   string
	PrefixLen int
	Gateway   net.IP
	IfIndex   int
	IfName    string
}

type KeaGeneratedStructure struct {
	Arguments struct {
		DeletedLeases []struct {
			Duid         string `json:"duid"`
			Expire       int    `json:"expire"`
			FqdnFwd      bool   `json:"fqdn-fwd"`
			FqdnRev      bool   `json:"fqdn-rev"`
			Hostname     string `json:"hostname"`
			HwAddress    string `json:"hw-address"`
			Iaid         int    `json:"iaid"`
			IPAddress    string `json:"ip-address"`
			PreferredLft int    `json:"preferred-lft"`
			State        int    `json:"state"`
			SubnetID     int    `json:"subnet-id"`
			Type         string `json:"type"`
			ValidLft     int    `json:"valid-lft"`
			PrefixLen    int    `json:"prefix-len,omitempty"`
		} `json:"deleted-leases"`
		Leases []struct {
			Duid         string `json:"duid"`
			Expire       int    `json:"expire"`
			FqdnFwd      bool   `json:"fqdn-fwd"`
			FqdnRev      bool   `json:"fqdn-rev"`
			Hostname     string `json:"hostname"`
			HwAddress    string `json:"hw-address"`
			Iaid         int    `json:"iaid"`
			IPAddress    string `json:"ip-address"`
			PreferredLft int    `json:"preferred-lft"`
			State        int    `json:"state"`
			SubnetID     int    `json:"subnet-id"`
			Type         string `json:"type"`
			ValidLft     int    `json:"valid-lft"`
			PrefixLen    int    `json:"prefix-len,omitempty"`
		} `json:"leases"`
	} `json:"arguments"`
	Command string   `json:"command"`
	Service []string `json:"service"`
}
