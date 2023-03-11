package usbid

type (
	ID     uint16
	Vendor struct {
		Name    string
		Devices map[ID]Device
	}
	Device struct {
		Name       string
		Interfaces map[ID]Interface
	}
	Interface struct {
		Name string
	}
)
