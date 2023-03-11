package usbid

import (
	_ "embed"
)

//go:embed usb.ids
var EmbeddedIDs []byte

func init() {
}
