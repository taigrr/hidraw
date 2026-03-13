package usbid

import (
	"os"
)

const systemDB = "/usr/share/hwdata/usb.ids"

// Vendors holds the parsed vendor database. It is populated from the embedded
// usb.ids file at init time. Call LoadSystemDB to replace it with the
// system-provided database if available.
var Vendors map[ID]Vendor

// LoadSystemDB reads and parses the system USB ID database, replacing the
// package-level Vendors map. Returns the parsed vendors or an error if the
// system database cannot be read or parsed.
func LoadSystemDB() (map[ID]Vendor, error) {
	file, err := os.Open(systemDB)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	vendors, err := Parse(file)
	if err != nil {
		return nil, err
	}
	Vendors = vendors
	return vendors, nil
}
