package main

import (
	"fmt"
	"time"
)

func compareAndSaveCurrentState() {
	for {
		time.Sleep(15 * time.Second)
		var activeLeasesListToCompare []savedLeasesStrut
		for _, x := range activeLeasesList {
			activeLeasesListToCompare = append(activeLeasesListToCompare, savedLeasesStrut{
				Address:   x.address,
				PrefixLen: x.prefixLen,
				Gateway:   x.gateway,
				IfIndex:   x.ifIndex,
				IfName:    x.ifName,
			})
		}

		if equal(activeLeasesListToCompare, savedLeasesList) == false {
			fmt.Println("info: startup lease file not match current state, new version saved")
			savedLeasesList = activeLeasesListToCompare
			saveActiveLeasesListToFile(savedLeasesList)
		}
	}
}

func equal(a []savedLeasesStrut, b []savedLeasesStrut) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if fmt.Sprint(v) != fmt.Sprint(b[i]) {
			return false
		}
	}
	return true
}
