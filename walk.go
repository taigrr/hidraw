package hidraw

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	usbid "github.com/taigrr/hidraw/usbids"
)

const (
	hidrawPath = "/sys/class/hidraw"
	ueventFile = "/device/uevent"
)

// Hidraw represents a HID raw device discovered via sysfs.
type Hidraw struct {
	PathName   string
	Path       string
	Driver     string
	HidID      string
	HidName    string
	HidPhys    string
	HidUniq    string
	Modalias   string
	VendorID   usbid.ID
	DeviceID   usbid.ID
	VendorName string
	DeviceName string
}

// Walk discovers all hidraw devices and returns them. Errors reading
// individual device uevent files are silently ignored.
func Walk() []Hidraw {
	devices, _ := WalkErr()
	return devices
}

// WalkErr discovers all hidraw devices. It returns an error only if the
// hidraw sysfs directory cannot be walked at all. Errors reading individual
// device uevent files are silently skipped.
func WalkErr() ([]Hidraw, error) {
	var devices []Hidraw
	err := filepath.WalkDir(hidrawPath, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == hidrawPath {
			return nil
		}
		dev := parseDevice(path, d.Name())
		devices = append(devices, dev)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking %s: %w", hidrawPath, err)
	}
	return devices, nil
}

func parseDevice(sysPath, name string) Hidraw {
	dev := Hidraw{
		Path:     name,
		PathName: "/dev/" + name,
	}

	uevent, err := os.ReadFile(filepath.Join(sysPath, ueventFile))
	if err != nil {
		return dev
	}

	for _, line := range strings.Split(string(uevent), "\n") {
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch key {
		case "DRIVER":
			dev.Driver = value
		case "HID_ID":
			dev.HidID = value
		case "HID_NAME":
			dev.HidName = value
		case "HID_PHYS":
			dev.HidPhys = value
		case "HID_UNIQ":
			dev.HidUniq = value
		case "MODALIAS":
			dev.Modalias = value
		}
	}

	dev.VendorID, dev.DeviceID = parseHidID(dev.HidID)
	dev.VendorName = usbid.LookupVendor(dev.VendorID)
	dev.DeviceName = usbid.LookupDevice(dev.VendorID, dev.DeviceID)

	return dev
}

// parseHidID extracts vendor and device IDs from a HID_ID string.
// Format: "BBBB:VVVVVVVV:DDDDDDDD" (bus:vendor:device, hex).
func parseHidID(hidID string) (vendor, device usbid.ID) {
	parts := strings.Split(hidID, ":")
	if len(parts) != 3 {
		return 0, 0
	}
	v, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		return 0, 0
	}
	d, err := strconv.ParseUint(parts[2], 16, 16)
	if err != nil {
		return 0, 0
	}
	return usbid.ID(v), usbid.ID(d)
}
