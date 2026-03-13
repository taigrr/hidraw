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
