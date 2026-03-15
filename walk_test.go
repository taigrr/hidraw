package hidraw

import (
	"testing"

	usbid "github.com/taigrr/hidraw/usbids"
)

func TestParseHidID(t *testing.T) {
	tests := []struct {
		input      string
		wantVendor usbid.ID
		wantDevice usbid.ID
	}{
		{"0003:0000046D:0000C52B", usbid.ID(0x046d), usbid.ID(0xc52b)},
		{"0005:0000054C:000005C4", usbid.ID(0x054c), usbid.ID(0x05c4)},
		{"", 0, 0},
		{"invalid", 0, 0},
		{"0003:notahex:0000C52B", 0, 0},
		{"0003:0000046D:notahex", 0, 0},
	}
	for _, tt := range tests {
		vendor, device := parseHidID(tt.input)
		if vendor != tt.wantVendor || device != tt.wantDevice {
			t.Errorf("parseHidID(%q) = (%v, %v), want (%v, %v)",
				tt.input, vendor, device, tt.wantVendor, tt.wantDevice)
		}
	}
}

func TestParseDevice(t *testing.T) {
	// parseDevice with a non-existent path should return a device with
	// just the path fields populated.
	dev := parseDevice("/nonexistent/path", "hidraw99")
	if dev.PathName != "/dev/hidraw99" {
		t.Errorf("PathName = %q, want %q", dev.PathName, "/dev/hidraw99")
	}
	if dev.Path != "hidraw99" {
		t.Errorf("Path = %q, want %q", dev.Path, "hidraw99")
	}
	if dev.Driver != "" {
		t.Errorf("Driver should be empty for nonexistent path, got %q", dev.Driver)
	}
}
