package services

import (
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func GetARPTable() []*ARPRecord {
	output, err := exec.Command("arp", "-a").Output()
	if err != nil {
		log.Fatal(err)
	}
	outputString := string(output)
	lines := strings.Split(outputString, "\n")
	arpList := make([]*ARPRecord, 0)
	for _, line := range lines {
		if strings.Contains(line, "incomplet") || len(line) < 4 {
			continue
		}
		elements := strings.Split(line, " ")

		ip := elements[1]
		mac := elements[3]
		arpList = append(arpList, &ARPRecord{IP: ip, MAC: mac})
	}

	return arpList
}

type ARPRecord struct {
	IP          string
	MAC         string
	Constructor string
}
