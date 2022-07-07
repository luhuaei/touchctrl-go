package main

import (
	"context"
	"fmt"

	"github.com/bendahl/uinput"
	evdev "github.com/gvalkov/golang-evdev"
)

type Point struct {
	X int32
	Y int32
}

func (p Point) String() string {
	return fmt.Sprintf("Point X: %d Y: %d\n", p.X, p.Y)
}

type Rect struct {
	TopLeft     Point
	RightBottom Point
}

func (r Rect) Contains(p Point) bool {
	return p.X >= r.TopLeft.X && p.X <= r.RightBottom.X && p.Y >= r.TopLeft.Y && p.Y <= r.RightBottom.Y
}

var (
	// touch area absolute position
	LeftRect = Rect{
		TopLeft:     Point{X: 1, Y: 1},
		RightBottom: Point{X: 250, Y: 250},
	}
	RightRect = Rect{
		TopLeft:     Point{X: 1000, Y: 1},
		RightBottom: Point{X: 1500, Y: 250},
	}
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
	var press bool
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

			var action bool
			var touch bool
			var finger bool
			var point Point
			for _, e := range es {
				if e.Type != evdev.EV_KEY && e.Type != evdev.EV_ABS {
					continue
				}

				switch e.Code {
				case evdev.BTN_TOUCH:
					action = true
					touch = e.Value == int32(1)
				case evdev.BTN_TOOL_FINGER:
					action = true
					finger = e.Value == int32(1)
				case evdev.ABS_X:
					point.X = e.Value
				case evdev.ABS_Y:
					point.Y = e.Value
				}
			}

			// 如果没有发生 press 和 up 事件则跳过
			if !action {
				continue
			}

			// 如果发生的事件处于热区内，则判断事件是 press 还是 up
			if LeftRect.Contains(point) || RightRect.Contains(point) {
				if touch && finger {
					err = m.output.KeyDown(uinput.KeyLeftctrl)
					press = true
				} else {
					err = m.output.KeyUp(uinput.KeyLeftctrl)
					press = false
				}
			} else if press {
				// 如果处于 press 状态，则释放
				err = m.output.KeyUp(uinput.KeyLeftctrl)
				press = false
			}

			if err != nil {
				return err
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
