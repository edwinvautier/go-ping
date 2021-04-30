package services

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func getConstructor(mac string) string {
	url := "https://api.macvendors.com/"
	log.Info("getting constructor for ", mac)
	resp, err := http.Get(url + mac)
	if err != nil || resp.StatusCode != 200 {
		if resp.StatusCode == 429 {
			dur := rand.Intn(5)
			time.Sleep(time.Duration(dur) * time.Second)
			return getConstructor(mac)
		} else {
			return "not-found"
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "not-found"
	}
	
	company := string(body)
	log.Info(company)
	return company
}

func FindAllConstructor(records *[]*ARPRecord) {
	for _, record := range *records {
		record.Constructor = getConstructor(record.MAC)
		time.Sleep(500 * time.Millisecond)
	}
}

type MACLookupBody struct {
	Company	string
}