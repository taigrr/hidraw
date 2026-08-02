package usbid

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testDB = `# Test USB IDs
0001  Test Vendor
	0001  Test Device One
		00  Test Interface
	0002  Test Device Two
0002  Another Vendor
	abcd  Some Device
`

func TestParse(t *testing.T) {
	vendors, err := Parse(strings.NewReader(testDB))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(vendors) != 2 {
		t.Fatalf("expected 2 vendors, got %d", len(vendors))
	}

	v1 := vendors[ID(0x0001)]
	if v1.Name != "Test Vendor" {
		t.Errorf("vendor 0001 name = %q, want %q", v1.Name, "Test Vendor")
	}
	if len(v1.Devices) != 2 {
		t.Fatalf("vendor 0001: expected 2 devices, got %d", len(v1.Devices))
	}

	d1 := v1.Devices[ID(0x0001)]
	if d1.Name != "Test Device One" {
		t.Errorf("device 0001:0001 name = %q, want %q", d1.Name, "Test Device One")
	}
	if len(d1.Interfaces) != 1 {
		t.Fatalf("device 0001:0001: expected 1 interface, got %d", len(d1.Interfaces))
	}
	iface := d1.Interfaces[ID(0x00)]
	if iface.Name != "Test Interface" {
		t.Errorf("interface name = %q, want %q", iface.Name, "Test Interface")
	}

	d2 := v1.Devices[ID(0x0002)]
	if d2.Name != "Test Device Two" {
		t.Errorf("device 0001:0002 name = %q, want %q", d2.Name, "Test Device Two")
	}

	v2 := vendors[ID(0x0002)]
	if v2.Name != "Another Vendor" {
		t.Errorf("vendor 0002 name = %q, want %q", v2.Name, "Another Vendor")
	}
}

func TestParseEmpty(t *testing.T) {
	vendors, err := Parse(strings.NewReader(""))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(vendors) != 0 {
		t.Errorf("expected 0 vendors, got %d", len(vendors))
	}
}

func TestParseCommentsOnly(t *testing.T) {
	input := "# This is a comment\n# Another comment\n"
	vendors, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(vendors) != 0 {
		t.Errorf("expected 0 vendors, got %d", len(vendors))
	}
}

func TestParseDoesNotAttachClassDevicesToPreviousVendor(t *testing.T) {
	input := `0001  Test Vendor
	0001  Test Device
C 00  Device
	03  Human Interface Device
0002  Another Vendor
	0002  Another Device
`

	vendors, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	vendor := vendors[ID(0x0001)]
	if _, ok := vendor.Devices[ID(0x03)]; ok {
		t.Fatal("class entry 03 was parsed as a device under vendor 0001")
	}
	if _, ok := vendors[ID(0x0002)]; !ok {
		t.Fatal("vendor 0002 was not parsed after class section")
	}
}

func TestParseDoesNotAttachInterfaceAfterMalformedDevice(t *testing.T) {
	input := `0001  Test Vendor
	0001  Test Device
	bad   Malformed Device
		02  Wrong Interface
	0002  Next Device
		03  Right Interface
`

	vendors, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	vendor := vendors[ID(0x0001)]
	firstDevice := vendor.Devices[ID(0x0001)]
	if _, ok := firstDevice.Interfaces[ID(0x02)]; ok {
		t.Fatal("interface after malformed device was attached to the previous valid device")
	}

	nextDevice := vendor.Devices[ID(0x0002)]
	if _, ok := nextDevice.Interfaces[ID(0x03)]; !ok {
		t.Fatal("interface after valid device was not parsed")
	}
}

func TestParseEmbeddedIDs(t *testing.T) {
	// Verify that the embedded DB was parsed successfully at init.
	if Vendors == nil {
		t.Fatal("Vendors is nil after init")
	}
	if len(Vendors) < 100 {
		t.Errorf("expected at least 100 vendors from embedded DB, got %d", len(Vendors))
	}

	// Spot check a well-known vendor.
	logitech := Vendors[ID(0x046d)]
	if logitech.Name == "" {
		t.Error("expected Logitech (046d) to be present in embedded DB")
	}
}

func TestLoadDBReplacesVendors(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "usb.ids")
	const customDB = `1234  Custom Vendor
	5678  Custom Device
`
	if err := os.WriteFile(dbPath, []byte(customDB), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	original := Vendors
	t.Cleanup(func() {
		Vendors = original
	})

	vendors, err := loadDB(dbPath)
	if err != nil {
		t.Fatalf("loadDB() error = %v", err)
	}
	if len(vendors) != 1 {
		t.Fatalf("len(vendors) = %d, want 1", len(vendors))
	}
	if Vendors[ID(0x1234)].Name != "Custom Vendor" {
		t.Fatalf("loaded vendor name = %q, want %q", Vendors[ID(0x1234)].Name, "Custom Vendor")
	}
	if LookupDevice(ID(0x1234), ID(0x5678)) != "Custom Device" {
		t.Fatalf("LookupDevice() = %q, want %q", LookupDevice(ID(0x1234), ID(0x5678)), "Custom Device")
	}
}

func TestLoadDBReturnsErrorForMissingFile(t *testing.T) {
	original := Vendors
	originalLogitech := LookupVendor(ID(0x046d))
	t.Cleanup(func() {
		Vendors = original
	})

	if _, err := loadDB(filepath.Join(t.TempDir(), "missing.ids")); err == nil {
		t.Fatal("loadDB() error = nil, want error")
	}
	if LookupVendor(ID(0x046d)) != originalLogitech {
		t.Fatal("Vendors changed after loadDB() error")
	}
}
