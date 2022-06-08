package main

import (
	"fmt"

	"github.com/taigrr/hidraw"
)

func main() {
	devs := hidraw.Walk()
	for _, d := range devs {
		fmt.Printf("%v\n", d)
	}
}
