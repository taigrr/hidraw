package usbid

import (
	"bytes"
	_ "embed"
	"log"
)

//go:embed usb.ids
var EmbeddedIDs []byte

func init() {
	var err error
	Vendors, err = Parse(bytes.NewReader(EmbeddedIDs))
	if err != nil {
		log.Printf("hidraw: failed to parse embedded USB IDs: %v", err)
	}
}
