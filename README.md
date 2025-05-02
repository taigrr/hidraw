# hidraw

A Go library for discovering and interacting with HID (Human Interface Device) raw devices on Linux systems.

## Overview

This library provides a simple interface to discover and access HID raw devices on Linux systems. It reads device information from `/sys/class/hidraw` and provides structured access to device properties.

## Features

- Discover all HID raw devices on the system
- Access device properties including:
  - Driver information
  - HID ID
  - Device name
  - Physical location
  - Unique identifier
  - Modalias

## Usage

```go
package main

import (
    "fmt"
    "github.com/taigrr/hidraw"
)

func main() {
    // Get all HID raw devices
    devices := hidraw.Walk()
    
    // Print device information
    for _, device := range devices {
        fmt.Printf("Device: %s\n", device.PathName)
        fmt.Printf("  Driver: %s\n", device.DRIVER)
        fmt.Printf("  Name: %s\n", device.HID_NAME)
        fmt.Printf("  Physical: %s\n", device.HID_PHYS)
        fmt.Printf("  Unique ID: %s\n", device.HID_UNIQ)
    }
}
```

## Device Information

The `Hidraw` struct provides the following fields:

- `PathName`: Full device path (e.g., "/dev/hidraw0")
- `Path`: Device name (e.g., "hidraw0")
- `DRIVER`: Device driver name
- `HID_ID`: HID device identifier
- `HID_NAME`: Human-readable device name
- `HID_PHYS`: Physical location of the device
- `HID_UNIQ`: Unique device identifier
- `MODALIAS`: Device modalias string

## Requirements

- Linux operating system
- Go 1.16 or later
- Access to `/sys/class/hidraw` directory
