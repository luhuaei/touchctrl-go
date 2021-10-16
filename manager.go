package main

import (
	"context"

	"github.com/bendahl/uinput"
	evdev "github.com/gvalkov/golang-evdev"
)

type manager struct {
	touchpad *evdev.InputDevice
	keyboard *evdev.InputDevice

	output uinput.Keyboard

	isPress bool
}

func NewManager(tpath, kpath string) (*manager, error) {
	touchpad, err := evdev.Open(tpath)
	if err != nil {
		return nil, err
	}

	keyboard, err := evdev.Open(kpath)
	if err != nil {
		return nil, err
	}

	output, err := uinput.CreateKeyboard("/dev/uinput", []byte("touchctrl-go"))
	if err != nil {
		return nil, err
	}

	return &manager{
		keyboard: keyboard,
		touchpad: touchpad,
		output:   output,
	}, nil
}

func (m *manager) Close() error {
	err := m.touchpad.File.Close()
	if err != nil {
		return err
	}

	err = m.keyboard.File.Close()
	if err != nil {
		return err
	}

	return m.output.Close()
}

func (m *manager) worker() error {
	ctx, cancel := context.WithCancel(context.Background())

	terrC := make(chan error, 1)
	go func() {
		terrC <- m.touchpadWorker(ctx)
		close(terrC)
		cancel()
	}()

	err := m.keyboardWorker(ctx)
	if err != nil {
		cancel()
		return err
	}
	return <-terrC
}

func (m *manager) touchpadWorker(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			es, err := m.touchpad.Read()
			if err != nil {
				return err
			}

			// if the ctrl already pressed.
			if m.isPress {
				continue
			}

			var touch bool
			var finger bool
			var found bool
			for _, e := range es {
				if e.Type != evdev.EV_KEY {
					continue
				}

				switch e.Code {
				case evdev.BTN_TOUCH:
					found = true
					touch = e.Value == int32(1)
				case evdev.BTN_TOOL_FINGER:
					found = true
					finger = e.Value == int32(1)
				}
			}

			if found {
				if touch && finger {
					err = m.output.KeyDown(uinput.KeyLeftctrl)
				} else {
					err = m.output.KeyUp(uinput.KeyLeftctrl)
				}
				if err != nil {
					return err
				}
			}
		}
	}
}

func (m *manager) keyboardWorker(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			es, err := m.keyboard.Read()
			if err != nil {
				return err
			}

			for _, e := range es {
				if e.Type == evdev.EV_KEY && e.Code == evdev.KEY_LEFTCTRL {
					m.isPress = e.Value != 0
				}
			}
		}
	}
}
