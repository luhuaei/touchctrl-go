package main

import "flag"

func main() {
	tpath := flag.String("touchpad", "/dev/input/event14", "touchpad device path")
	kpath := flag.String("keyboard", "/dev/input/event0", "keyboard device path")

	flag.Parse()
	m, err := NewManager(*tpath, *kpath)
	if err != nil {
		panic(err)
	}
	err = m.worker()
	if err != nil {
		m.Close()
		panic(err)
	}
}
