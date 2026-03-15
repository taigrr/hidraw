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
	var currentVendor *Vendor
	var currentDeviceID ID
	var currentDevice *Device

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
				continue
			}
			hexStr := trimmed[:4]
			id, err := parseHexID(hexStr)
			if err != nil {
				// Non-vendor top-level lines (like "C 00  ...") are
				// device classes or other sections — stop parsing vendors.
				if trimmed[0] == 'C' || trimmed[0] == 'A' || trimmed[0] == 'H' || trimmed[0] == 'L' || trimmed[0] == 'P' || trimmed[0] == 'V' || trimmed[0] == 'R' || trimmed[0] == 'B' {
					break
				}
				continue
			}
			name := strings.TrimSpace(trimmed[4:])
			v := Vendor{
				Name:    name,
				Devices: make(map[ID]Device),
			}
			vendors[id] = v
			currentVendorID = id
			vRef := vendors[currentVendorID]
			currentVendor = &vRef
			currentDevice = nil

		case 1:
			// Device line: "\txxxx  Name"
			if currentVendor == nil || len(trimmed) < 6 {
				continue
			}
			hexStr := trimmed[:4]
			id, err := parseHexID(hexStr)
			if err != nil {
				continue
			}
			name := strings.TrimSpace(trimmed[4:])
			d := Device{
				Name:       name,
				Interfaces: make(map[ID]Interface),
			}
			currentVendor.Devices[id] = d
			currentDeviceID = id
			dRef := currentVendor.Devices[currentDeviceID]
			currentDevice = &dRef
			vendors[currentVendorID] = *currentVendor

		case 2:
			// Interface line: "\t\txx  Name" (2-digit hex ID)
			if currentDevice == nil || len(trimmed) < 4 {
				continue
			}
			hexStr := trimmed[:2]
			id, err := parseHexID(hexStr)
			if err != nil {
				continue
			}
			name := strings.TrimSpace(trimmed[2:])
			currentDevice.Interfaces[id] = Interface{Name: name}
			currentVendor.Devices[currentDeviceID] = *currentDevice
			vendors[currentVendorID] = *currentVendor
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
