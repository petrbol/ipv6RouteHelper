package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func ReceiveKeaPostJson(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var keaMessage KeaGeneratedStructure
	err = json.Unmarshal(body, &keaMessage)
	if err != nil {
		log.Println("receive non json message\t", string(body))
	}

	NetlinkReconfigureByKea(keaMessage)
}
