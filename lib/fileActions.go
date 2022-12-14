package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func readConfiguration() {
	data := readFile(SavedLeasesListFile)
	err := json.Unmarshal(data, &SavedLeasesList)
	if err != nil {
		fmt.Println("err: failed to unmarshal json read file", err)
	}
}

func readFile(fileName string) []byte {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("err: failed to read file", err)
		return nil
	}
	return content
}

func saveActiveLeasesListToFile(activeLeaseList []savedLeasesStrut) {
	jsonData, err := json.MarshalIndent(activeLeaseList, "", " ")
	if err != nil {
		fmt.Println("err: failed to convert class map to json structure", err)
	}

	errWrite := ioutil.WriteFile(SavedLeasesListFile, jsonData, 0666)
	if errWrite != nil {
		fmt.Println("err: failed to write json structure to file", errWrite)
	}
}
