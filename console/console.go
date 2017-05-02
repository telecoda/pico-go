package console

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

// Global var

var _console *console

const (
	_version     = "v0.1"
	_logoWidth   = 57
	_logoHeight  = 24
	_charWidth   = int32(4)
	_charHeight  = int32(8)
	_cursorFlash = time.Duration(500 * time.Millisecond)
)

func NewConsole(cfg Config) (Console, error) {
	_console = &console{
		Config:        cfg,
		currentMode:   CLI,
		secondaryMode: CODE_EDITOR,
		hasQuit:       false,
	}

	// init SDL
	sdl.Init(sdl.INIT_EVERYTHING)
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")
	ttf.Init()

	// initialise window
	window, err := sdl.CreateWindow(
		title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int(cfg.WindowWidth), int(cfg.WindowHeight), sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE)

	if err != nil {
		return nil, err
	}
	_console.window = window

	// initialse renderer
	// create renderer
	r, _ := sdl.CreateRenderer(window, -2, sdl.RENDERER_ACCELERATED)
	if r == nil {
		// revert to software
		r, _ = sdl.CreateRenderer(window, -2, sdl.RENDERER_SOFTWARE)
		if r == nil {
			return nil, err
		}
	}

	_console.renderer = r

	// init palette
	_console.palette = initPico8Palette()

	// init font
	font, err := ttf.OpenFont("./fonts/PICO-8.ttf", 4)
	if err != nil {
		return nil, fmt.Errorf("Error loading font:%s", err)
	}

	_console.font = font

	// init logo
	logo, err := img.Load("./images/pico-go-logo.png")
	if err != nil {
		return nil, fmt.Errorf("Error loading image: %s\n", err)
	}

	_console.logo = logo

	// initialise modes
	modes, err := _console.initModes()
	if err != nil {
		return nil, err
	}
	_console.modes = modes

	return _console, nil
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

			// poll all events
			for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				// this is the common event handling code
				switch t := event.(type) {
				case *sdl.QuitEvent:
					c.hasQuit = true
				case *sdl.KeyDownEvent:
					switch t.Keysym.Sym {
					case sdl.K_ESCAPE:
						c.toggleCLI()
					default:
						// if not handled pass event to mode event handler
						if err := mode.HandleEvent(event); err != nil {
							return err
						}
					}
				}
			}

			if err := mode.Update(); err != nil {
				return err
			}

			if err := mode.Render(); err != nil {
				return err
			}

			mode.Flip()

			// lock framerate
			lockFps()

		} else {
			return fmt.Errorf("Mode :%d not found in console.modes", c.currentMode)
		}

	}

	return nil
}

// toggleCLI - toggle between CLI and secondary mode
func (c *console) toggleCLI() {
	if c.currentMode == CLI {
		c.SetMode(c.secondaryMode)
	} else {
		c.secondaryMode = c.currentMode
		c.SetMode(CLI)
	}
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
	c.window.Destroy()
}
