package main

import (
	"fmt"
	"path/filepath"
	"regexp"

	evdev "github.com/gvalkov/golang-evdev"
)

func main() {
	events, err := filepath.Glob("/dev/input/event*")
	if err != nil {
		panic(err)
	}

	Rtouchpad := regexp.MustCompile("(?i)Touchpad")
	RKeyboard := regexp.MustCompile("(?i)keyboard")

	var tpath string
	var kpath string
	for _, eventFile := range events {
		device, err := evdev.Open(eventFile)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		if Rtouchpad.MatchString(device.Name) {
			tpath = eventFile
		}
		if RKeyboard.MatchString(device.Name) {
			kpath = eventFile
		}
		device.File.Close()
	}

	if tpath == "" || kpath == "" {
		panic("not found touchpad or keyboard device")
	}

	m, err := NewManager(tpath, kpath)
	if err != nil {
		panic(err)
	}
	err = m.worker()
	if err != nil {
		m.Close()
		panic(err)
	}
}
