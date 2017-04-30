package console

import (
	"fmt"
	"time"

	"github.com/telecoda/pico-go/api"
	"github.com/telecoda/pico-go/display"
	"github.com/veandco/go-sdl2/sdl"
)

func NewConsole(cfg api.Config) (Console, error) {
	console := &console{
		Config:      cfg,
		currentMode: CLI,
		hasQuit:     false,
	}

	// initialise display
	if console.Verbose {
		fmt.Printf("Initialising display\n")
	}

	display, err := display.NewDisplay(cfg)
	if err != nil {
		return nil, err
	}

	console.Display = display

	// initialise modes
	modes := console.initModes()

	console.modes = modes

	return console, nil
}

func (c *console) SetMode(newMode ModeType) {
	c.currentMode = newMode
}

func (c *console) LoadCart(path string) error {
	return nil
}

var timeBudget int64
var lastFrame time.Time
var startFrame time.Time
var endFrame time.Time

// Run is the main run loop
func (c *console) Run() error {
	// poll events

	endFrame = time.Now() // init end frame
	timeBudget = time.Duration(1*time.Second).Nanoseconds() / int64(c.Config.FPS)
	for !c.hasQuit {
		startFrame = time.Now() // used for framerate timing

		if mode, ok := c.modes[c.currentMode]; ok {

			if err := mode.PollEvents(); err != nil {
				return err
			}
			if err := mode.Update(); err != nil {
				return err
			}

			if err := mode.Render(); err != nil {
				return err
			}
			c.Flip()

			// lock framerate
			lockFps()

		} else {
			return fmt.Errorf("Mode :%d not found in console.modes", c.currentMode)
		}

	}

	return nil
}

// lockFps - locks rendering to a steady framerate
func lockFps() float64 {
	now := time.Now()
	// calc time to process frame so since start
	procTime := now.Sub(startFrame)

	// delay for remainder of time budget (based on fps)
	delay := time.Duration(timeBudget - procTime.Nanoseconds())
	if delay > 0 {
		sdl.Delay(uint32(delay / 1000000))
	}

	// calc actual fps being achieved
	endFrame = time.Now()
	frameTime := endFrame.Sub(startFrame)

	endFrame = time.Now()

	return float64(time.Second) / float64(frameTime.Nanoseconds())
}

// Destroy cleans up any resources at end
func (c *console) Destroy() {
	c.Display.Destroy()
}
