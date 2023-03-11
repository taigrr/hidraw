package usbid

import (
	"os"
)

const systemDB = "/usr/share/hwdata/usb.ids"

var Vendors map[ID]Vendor

func LoadSystemDB() ([]Vendor, error) {
	file, err := os.Open(systemDB)
	if err != nil {
		return []Vendor{}, err
	}
	defer file.Close()
	return []Vendor{}, nil
}
