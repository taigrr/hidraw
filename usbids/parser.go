package usbid

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Parse reads a usb.ids formatted file and returns a map of vendors keyed by
// their hex ID. The file format uses tabs for hierarchy:
//
//	vendor  vendor_name
//		device  device_name
//			interface  interface_name
func Parse(r io.Reader) (map[ID]Vendor, error) {
	vendors := make(map[ID]Vendor)
	scanner := bufio.NewScanner(r)

	var currentVendorID ID
	var currentVendorOK bool
	var currentDeviceID ID
	var currentDeviceOK bool

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 || line[0] == '#' {
			continue
		}

		// Count leading tabs to determine hierarchy level.
		tabs := 0
		for _, ch := range line {
			if ch == '\t' {
				tabs++
			} else {
				break
			}
		}
		trimmed := strings.TrimLeft(line, "\t")

		switch tabs {
		case 0:
			// Vendor line: "xxxx  Name"
			if len(trimmed) < 6 {
				currentVendorOK = false
				currentDeviceOK = false
				continue
			}
			hexStr := trimmed[:4]
			id, err := parseHexID(hexStr)
			if err != nil {
				currentVendorOK = false
				currentDeviceOK = false
				continue
			}
			name := strings.TrimSpace(trimmed[4:])
			vendors[id] = Vendor{
				Name:    name,
				Devices: make(map[ID]Device),
			}
			currentVendorID = id
			currentVendorOK = true
			currentDeviceOK = false

		case 1:
			// Device line: "\txxxx  Name"
			if !currentVendorOK || len(trimmed) < 6 {
				currentDeviceOK = false
				continue
			}
			hexStr := trimmed[:4]
			id, err := parseHexID(hexStr)
			if err != nil {
				currentDeviceOK = false
				continue
			}
			name := strings.TrimSpace(trimmed[4:])
			currentVendor := vendors[currentVendorID]
			currentVendor.Devices[id] = Device{
				Name:       name,
				Interfaces: make(map[ID]Interface),
			}
			vendors[currentVendorID] = currentVendor
			currentDeviceID = id
			currentDeviceOK = true

		case 2:
			// Interface line: "\t\txx  Name" (2-digit hex ID)
			if !currentDeviceOK || len(trimmed) < 4 {
				continue
			}
			hexStr := trimmed[:2]
			id, err := parseHexID(hexStr)
			if err != nil {
				continue
			}
			name := strings.TrimSpace(trimmed[2:])
			currentVendor := vendors[currentVendorID]
			currentDevice := currentVendor.Devices[currentDeviceID]
			currentDevice.Interfaces[id] = Interface{Name: name}
			currentVendor.Devices[currentDeviceID] = currentDevice
			vendors[currentVendorID] = currentVendor
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning usb.ids: %w", err)
	}

	return vendors, nil
}

func parseHexID(s string) (ID, error) {
	val, err := strconv.ParseUint(s, 16, 16)
	if err != nil {
		return 0, err
	}
	return ID(val), nil
}
