package services

import "os"

func WriteRecord(record ARPRecord) error {
	f, err := os.OpenFile("arp_list.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}
	if _, err := f.WriteString("|IP " + record.IP + " |MAC " + record.MAC + " |CONSTRUCTOR " + record.Constructor + "\n"); err != nil {
		return err
	}

	return nil
}
func WriteList(records []*ARPRecord) {
	for _, record := range records {
		WriteRecord(*record)
	}
}