package lib

import (
	"fmt"
	"time"
)

func CompareAndSaveCurrentState() {
	for {
		time.Sleep(15 * time.Second)
		var activeLeasesListToCompare []savedLeasesStrut
		for _, x := range ActiveLeasesList {
			activeLeasesListToCompare = append(activeLeasesListToCompare, savedLeasesStrut{
				Address:   x.address,
				PrefixLen: x.prefixLen,
				Gateway:   x.gateway,
				IfIndex:   x.ifIndex,
				IfName:    x.ifName,
			})
		}

		if equal(activeLeasesListToCompare, SavedLeasesList) == false {
			fmt.Println("info: startup lease file not match current state, new version saved")
			SavedLeasesList = activeLeasesListToCompare
			saveActiveLeasesListToFile(SavedLeasesList)
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
