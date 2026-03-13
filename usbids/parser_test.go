package usbid

import (
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
