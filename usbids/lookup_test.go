package usbid

import "testing"

func TestIDString(t *testing.T) {
	tests := []struct {
		id   ID
		want string
	}{
		{ID(0x0000), "0000"},
		{ID(0x046d), "046d"},
		{ID(0xffff), "ffff"},
		{ID(0x0001), "0001"},
	}
	for _, tt := range tests {
		got := tt.id.String()
		if got != tt.want {
			t.Errorf("ID(%d).String() = %q, want %q", tt.id, got, tt.want)
		}
	}
}

func TestLookupVendor(t *testing.T) {
	name := LookupVendor(ID(0x046d))
	if name == "" {
		t.Error("expected non-empty name for Logitech (046d)")
	}

	name = LookupVendor(ID(0x0000))
	if name != "" {
		t.Errorf("expected empty name for unknown vendor, got %q", name)
	}
}

func TestLookupDevice(t *testing.T) {
	name := LookupDevice(ID(0x046d), ID(0xc52b))
	if name == "" {
		t.Error("expected non-empty name for Logitech device 046d:c52b")
	}

	// Unknown vendor returns empty.
	name = LookupDevice(ID(0x0000), ID(0x0001))
	if name != "" {
		t.Errorf("expected empty name for unknown vendor, got %q", name)
	}

	// Known vendor, unknown device returns empty.
	name = LookupDevice(ID(0x046d), ID(0xfffe))
	if name != "" {
		t.Errorf("expected empty name for unknown device, got %q", name)
	}
}

func TestLookupInterface(t *testing.T) {
	originalVendors := Vendors
	t.Cleanup(func() {
		Vendors = originalVendors
	})

	Vendors = map[ID]Vendor{
		ID(0x0001): {
			Name: "Test Vendor",
			Devices: map[ID]Device{
				ID(0x0002): {
					Name: "Test Device",
					Interfaces: map[ID]Interface{
						ID(0x0003): {Name: "Test Interface"},
					},
				},
			},
		},
	}

	name := LookupInterface(ID(0x0001), ID(0x0002), ID(0x0003))
	if name != "Test Interface" {
		t.Errorf("LookupInterface() = %q, want %q", name, "Test Interface")
	}

	tests := []struct {
		name        string
		vendorID    ID
		deviceID    ID
		interfaceID ID
	}{
		{"unknown vendor", ID(0xffff), ID(0x0002), ID(0x0003)},
		{"unknown device", ID(0x0001), ID(0xffff), ID(0x0003)},
		{"unknown interface", ID(0x0001), ID(0x0002), ID(0xffff)},
	}
	for _, tt := range tests {
		got := LookupInterface(tt.vendorID, tt.deviceID, tt.interfaceID)
		if got != "" {
			t.Errorf("%s: LookupInterface() = %q, want empty string", tt.name, got)
		}
	}
}
