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
	// Unknown vendor returns empty.
	name := LookupDevice(ID(0x0000), ID(0x0001))
	if name != "" {
		t.Errorf("expected empty name for unknown vendor, got %q", name)
	}

	// Known vendor, unknown device returns empty.
	name = LookupDevice(ID(0x046d), ID(0xfffe))
	if name != "" {
		t.Errorf("expected empty name for unknown device, got %q", name)
	}
}
