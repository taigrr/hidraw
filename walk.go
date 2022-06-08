package hidraw

import (
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const hidrawPath = "/sys/class/hidraw"
const ueventFile = "/device/uevent"

type Hidraw struct {
	PathName string
	Path     string
	DRIVER   string
	HID_ID   string
	HID_NAME string
	HID_PHYS string
	HID_UNIQ string
	MODALIAS string
}

func Walk() []Hidraw {
	var hidraw []Hidraw
	filepath.WalkDir(hidrawPath, func(path string, d fs.DirEntry, err error) error {
		dev := Hidraw{}
		dev.Path = d.Name()
		dev.PathName = "/dev/" + d.Name()
		if path == hidrawPath {
			return nil
		}
		uevent, err := ioutil.ReadFile(filepath.Join(path, ueventFile))
		if err != nil {
			return nil
		}
		content := string(uevent)
		for _, line := range strings.Split(content, "\n") {
			split := strings.Split(line, "=")
			if len(split) < 2 {
				continue
			}
			switch split[0] {
			case "DRIVER":
				dev.DRIVER = split[1]
			case "HID_ID":
				dev.HID_ID = split[1]
			case "HID_NAME":
				dev.HID_NAME = split[1]
			case "HID_PHYS":
				dev.HID_PHYS = split[1]
			case "HID_UNIQ":
				dev.HID_UNIQ = split[1]
			case "MODALIAS":
				dev.MODALIAS = split[1]
			}
		}
		hidraw = append(hidraw, dev)
		return nil
	})
	return hidraw
}
