# hidraw

A Go library for discovering HID (Human Interface Device) raw devices on Linux via sysfs, with built-in USB vendor/device name resolution.

## Overview

`hidraw` reads device information from `/sys/class/hidraw` and returns structured data for each HID raw device, including driver info, HID IDs, and human-readable vendor/device names looked up from the USB ID database.

## Features

- Discover all HID raw devices on the system
- Parse HID ID strings into vendor and device IDs
- Look up vendor and device names from an embedded USB ID database
- Optionally load the system USB ID database (`/usr/share/hwdata/usb.ids`)

## Installation

```bash
go get github.com/taigrr/hidraw
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/taigrr/hidraw"
)

func main() {
	devs := hidraw.Walk()
	for _, d := range devs {
		fmt.Printf("%s: %s", d.PathName, d.HidName)
		if d.VendorName != "" {
			fmt.Printf(" [%s", d.VendorName)
			if d.DeviceName != "" {
				fmt.Printf(" — %s", d.DeviceName)
			}
			fmt.Print("]")
		}
		fmt.Println()
	}
}
```

## API

### `hidraw.Walk() []Hidraw`

Discovers all hidraw devices. Errors reading individual devices are silently ignored.

### `hidraw.WalkErr() ([]Hidraw, error)`

Like `Walk`, but returns an error if the hidraw sysfs directory cannot be walked.

### `Hidraw` struct

| Field        | Type      | Description                              |
| ------------ | --------- | ---------------------------------------- |
| `PathName`   | `string`  | Full device path (e.g. `/dev/hidraw0`)   |
| `Path`       | `string`  | Device name (e.g. `hidraw0`)             |
| `Driver`     | `string`  | Kernel driver name                       |
| `HidID`      | `string`  | Raw HID ID (`bus:vendor:device` hex)     |
| `HidName`    | `string`  | Human-readable device name from uevent   |
| `HidPhys`    | `string`  | Physical location of the device          |
| `HidUniq`    | `string`  | Unique device identifier                 |
| `Modalias`   | `string`  | Device modalias string                   |
| `VendorID`   | `usbid.ID`| Parsed vendor ID                         |
| `DeviceID`   | `usbid.ID`| Parsed device ID                         |
| `VendorName` | `string`  | Looked-up vendor name (e.g. `Logitech`)  |
| `DeviceName` | `string`  | Looked-up device name                    |

### `usbids` subpackage

```go
import usbid "github.com/taigrr/hidraw/usbids"

// Look up by ID
name := usbid.LookupVendor(usbid.ID(0x046d))   // "Logitech, Inc."
name  = usbid.LookupDevice(usbid.ID(0x046d), usbid.ID(0xc52b))
name  = usbid.LookupInterface(usbid.ID(0x046d), usbid.ID(0xc52b), usbid.ID(0x00))

// Load system DB instead of embedded
usbid.LoadSystemDB()

// Parse a custom usb.ids file
vendors, err := usbid.Parse(reader)
```

## Requirements

- Linux operating system
- Go 1.26 or later
- Access to `/sys/class/hidraw` directory

## License

0BSD — see [LICENSE](LICENSE).
