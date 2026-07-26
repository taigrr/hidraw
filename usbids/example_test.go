package usbid_test

import (
	"fmt"
	"strings"

	usbid "github.com/taigrr/hidraw/usbids"
)

func ExampleLookupVendor() {
	fmt.Println(usbid.LookupVendor(usbid.ID(0x046d)))

	// Output:
	// Logitech, Inc.
}

func ExampleParse() {
	const db = `046d  Logitech, Inc.
	c52b  Unifying Receiver
		00  Keyboard
`

	vendors, err := usbid.Parse(strings.NewReader(db))
	if err != nil {
		fmt.Println(err)
		return
	}

	vendor := vendors[usbid.ID(0x046d)]
	device := vendor.Devices[usbid.ID(0xc52b)]
	iface := device.Interfaces[usbid.ID(0x00)]

	fmt.Println(vendor.Name)
	fmt.Println(device.Name)
	fmt.Println(iface.Name)

	// Output:
	// Logitech, Inc.
	// Unifying Receiver
	// Keyboard
}
