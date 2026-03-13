package usbid

import "fmt"

// String returns the hex representation of an ID (e.g. "04b3").
func (id ID) String() string {
	return fmt.Sprintf("%04x", uint16(id))
}

// LookupVendor returns the vendor name for the given ID, or an empty string
// if not found.
func LookupVendor(vendorID ID) string {
	if v, ok := Vendors[vendorID]; ok {
		return v.Name
	}
	return ""
}

// LookupDevice returns the device name for the given vendor and device IDs,
// or an empty string if not found.
func LookupDevice(vendorID, deviceID ID) string {
	v, ok := Vendors[vendorID]
	if !ok {
		return ""
	}
	d, ok := v.Devices[deviceID]
	if !ok {
		return ""
	}
	return d.Name
}
