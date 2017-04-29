package console

import (
	"fmt"

	"time"

	"github.com/telecoda/pico-go/config"
	"github.com/telecoda/pico-go/display"
	"github.com/telecoda/pico-go/events"
	"github.com/veandco/go-sdl2/sdl"
)

// This is the main control package for the virtual console
// the console package coordinates with all other console subsystems

type Console interface {
	Run() error
	Destroy()
}

type console struct {
	config.Config
	Display      display.Display
	EventHandler events.EventHandler
}

func New(cfg config.Config) (Console, error) {
	console := &console{
		Config: cfg,
	}

	// initialise display
	if console.Verbose {
		fmt.Printf("Initialising display\n")
	}

	display, err := display.New(cfg)
	if err != nil {
		return nil, err
	}

	console.Display = display

	// initialise event handler
	if console.Verbose {
		fmt.Printf("Initialising event handler\n")
	}

	handler, err := events.New()
	if err != nil {
		return nil, err
	}

	console.EventHandler = handler

	return console, nil
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
	for !c.EventHandler.HasQuit() {
		startFrame = time.Now() // used for framerate timing

		if err := c.EventHandler.PollEvents(); err != nil {
			return err
		}

		if err := c.Display.Render(); err != nil {
			return err
		}

		// lock delay
		lockFps()

	}

	return nil
}

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
