package hidraw

import (
	"os"
	"path/filepath"
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

func TestParseDeviceReadsUeventFields(t *testing.T) {
	tempDir := t.TempDir()
	sysPath := filepath.Join(tempDir, "hidraw0")
	if err := os.MkdirAll(filepath.Join(sysPath, "device"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	uevent := "DRIVER=hid-generic\n" +
		"HID_ID=0003:0000046D:0000C52B\n" +
		"HID_NAME=Gaming Receiver\n" +
		"HID_PHYS=usb-0000:00:14.0-3/input2\n" +
		"HID_UNIQ=abc123\n" +
		"MODALIAS=hid:b0003g0001v0000046Dp0000C52B\n"
	if err := os.WriteFile(filepath.Join(sysPath, "device", "uevent"), []byte(uevent), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	dev := parseDevice(sysPath, "hidraw0")
	if dev.Path != "hidraw0" || dev.PathName != "/dev/hidraw0" {
		t.Fatalf("unexpected path fields: %+v", dev)
	}
	if dev.Driver != "hid-generic" {
		t.Errorf("Driver = %q, want %q", dev.Driver, "hid-generic")
	}
	if dev.HidID != "0003:0000046D:0000C52B" {
		t.Errorf("HidID = %q, want %q", dev.HidID, "0003:0000046D:0000C52B")
	}
	if dev.HidName != "Gaming Receiver" {
		t.Errorf("HidName = %q, want %q", dev.HidName, "Gaming Receiver")
	}
	if dev.HidPhys != "usb-0000:00:14.0-3/input2" {
		t.Errorf("HidPhys = %q, want %q", dev.HidPhys, "usb-0000:00:14.0-3/input2")
	}
	if dev.HidUniq != "abc123" {
		t.Errorf("HidUniq = %q, want %q", dev.HidUniq, "abc123")
	}
	if dev.Modalias != "hid:b0003g0001v0000046Dp0000C52B" {
		t.Errorf("Modalias = %q, want %q", dev.Modalias, "hid:b0003g0001v0000046Dp0000C52B")
	}
	if dev.VendorID != usbid.ID(0x046d) || dev.DeviceID != usbid.ID(0xc52b) {
		t.Errorf("parsed IDs = (%#x, %#x), want (%#x, %#x)", dev.VendorID, dev.DeviceID, usbid.ID(0x046d), usbid.ID(0xc52b))
	}
	if dev.VendorName != "Logitech, Inc." {
		t.Errorf("VendorName = %q, want %q", dev.VendorName, "Logitech, Inc.")
	}
	if dev.DeviceName == "" {
		t.Error("DeviceName should not be empty for a known Logitech receiver")
	}
}
