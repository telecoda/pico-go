package console

import (
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"time"

	"golang.org/x/image/font"

	"sync"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Global var

var _console *console

const (
	_version        = "v0.1"
	_logoWidth      = 57
	_logoHeight     = 24
	_spriteWidth    = 8
	_spriteHeight   = 8
	_spritesPerLine = 16
	//_charWidth      = 4
	//_charHeight     = 6
	_maxCmdLen   = 254
	_cursorFlash = time.Duration(500 * time.Millisecond)
)

const (
	systemSpriteBank = 0
	userSpriteBank1  = 1
)

type Console interface {
	LoadCart(cart Cartridge) error
	Run() error
	Destroy()
	//GetWindow() *sdl.Window
	GetBounds() image.Rectangle
	SetMode(newMode ModeType)
	Inputter
}

type console struct {
	sync.Mutex
	Config

	currentMode   ModeType
	secondaryMode ModeType
	modes         map[ModeType]Mode
	hasQuit       bool

	// files
	baseDir    string
	currentDir string

	cart Cartridge

	screen *ebiten.Image
	pImage *image.Paletted

	font              font.Face
	logo              *image.RGBA
	sprites           []*image.RGBA
	currentSpriteBank int
	originalPalette   *palette

	state    Persister
	recorder Recorder
	Inputter
}

func NewConsole(consoleType ConsoleType) (Console, error) {

	// validate type
	if _, ok := ConsoleTypes[consoleType]; !ok {
		return nil, fmt.Errorf("Console type: %s not supported", consoleType)
	}

	cfg := NewConfig(consoleType)

	_console = &console{
		Config:        cfg,
		currentMode:   CLI,
		secondaryMode: CODE_EDITOR,
		hasQuit:       false,
	}

	// TODO
	// persister, err := NewStateManager()
	// if err != nil {
	// 	return nil, err
	// }

	// _console.state = persister

	// TODO
	// state, err := persister.LoadState()
	// if err != nil {
	// 	fmt.Printf("Warning loading state: %s - using default\n", err)
	// }
	// // init SDL
	// sdl.Init(sdl.INIT_EVERYTHING)
	// sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")
	// ttf.Init()

	//if state == nil {
	// use defaults
	// state = &ConsoleState{
	// 	X: sdl.WINDOWPOS_CENTERED,
	// 	Y: sdl.WINDOWPOS_CENTERED,
	// 	W: cfg.WindowWidth,
	// 	H: cfg.WindowHeight,
	// }
	//}

	_console.recorder = NewRecorder(cfg.FPS, cfg.GifLength)
	_console.Inputter = NewInputter()

	// initialise window
	// window, err := sdl.CreateWindow(
	// 	title, state.X, state.Y,
	// 	state.W, state.H, sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE)

	// if err != nil {
	// 	return nil, err
	// }
	// _console.window = window

	// // initialse renderer
	// // create renderer
	// r, _ := sdl.CreateRenderer(window, -2, sdl.RENDERER_ACCELERATED)
	// if r == nil {
	// 	// revert to software
	// 	r, _ = sdl.CreateRenderer(window, -2, sdl.RENDERER_SOFTWARE)
	// 	if r == nil {
	// 		return nil, err
	// 	}
	// }

	// _console.renderer = r

	//goPath := build.Default.GOPATH

	// init font

	//fontPath := fmt.Sprintf("%s/src/github.com/telecoda/pico-go/consoles/%s/font.ttf", goPath, _console.consoleType)
	fontPath := fmt.Sprintf("./font.ttf")
	// font, err := ttf.OpenFont(fontPath, _console.Config.fontWidth)
	// if err != nil {
	// 	return nil, fmt.Errorf("Error loading font:%s", err)
	// }

	// INIT font
	// TODO move elsewhere..
	f, err := ebitenutil.OpenFile(fontPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	tt, err := truetype.Parse(b)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 48
	mplusNormalFont := truetype.NewFace(tt, &truetype.Options{
		Size:    6,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	// mplusBigFont = truetype.NewFace(tt, &truetype.Options{
	// 	Size:    48,
	// 	DPI:     dpi,
	// 	Hinting: font.HintingFull,
	// })

	_console.font = mplusNormalFont

	// init logo
	// logoPath := fmt.Sprintf("%s/src/github.com/telecoda/pico-go/consoles/%s/logo.png", goPath, _console.consoleType)
	// logo, err := img.Load(logoPath)
	// if err != nil {
	// 	return nil, fmt.Errorf("Error loading image: %s\n", err)
	// }

	//	logoPath := fmt.Sprintf("%s/src/github.com/telecoda/pico-go/consoles/%s/logo.png", goPath, _console.consoleType)
	_, logo, err := ebitenutil.NewImageFromFile("./logo.png", ebiten.FilterNearest)
	if err != nil {
		return nil, fmt.Errorf("Error loading image: %s\n", err)
	}
	_console.logo = logo.(*image.RGBA)

	// init sprites
	// There are 2 sprite banks
	// 0 = System sprites
	// 1 = User sprite bank 1
	_console.sprites = make([]*image.RGBA, 2)

	// tempSurface, err := sdl.CreateRGBSurface(0, int32(cfg.ConsoleWidth), int32(cfg.ConsoleHeight), 8, 0, 0, 0, 0)
	// if err != nil {
	// 	return nil, err
	// }

	_console.palette = newPalette(cfg.consoleType)
	_console.originalPalette = newPalette(cfg.consoleType)

	// create paletted image
	rect := image.Rect(0, 0, cfg.ConsoleWidth, cfg.ConsoleHeight)
	pImage := image.NewPaletted(rect, _console.palette.colors)
	_console.pImage = pImage

	// init icons
	//	iconsPath := fmt.Sprintf("%s/src/github.com/telecoda/pico-go/consoles/%s/icons.png", goPath, _console.consoleType)
	_, icons, err := ebitenutil.NewImageFromFile("./icons.png", ebiten.FilterNearest)
	//icons, err := img.Load(iconsPath)
	if err != nil {
		return nil, fmt.Errorf("Error loading image: %s\n", err)
	}

	// if err := setSurfacePalette(_console.originalPalette, tempSurface); err != nil {
	// 	return nil, err
	// }

	// iconsSurface, err := icons.Convert(tempSurface.Format, tempSurface.Flags)
	// if err != nil {
	// 	return nil, err
	// }

	// _console.sprites[systemSpriteBank] = iconsSurface
	_console.sprites[systemSpriteBank] = icons.(*image.RGBA)

	// if err := setSurfacePalette(_console.originalPalette, tempSurface); err != nil {
	// 	return nil, err
	// }

	//sprites, err := img.Load("./sprites/sprites.png")
	_, sprites, err := ebitenutil.NewImageFromFile("./sprites/sprites.png", ebiten.FilterNearest)
	if err != nil {
		return nil, fmt.Errorf("Error loading image: %s\n", err)
	}

	// spritesSurface, err := sprites.Convert(tempSurface.Format, tempSurface.Flags)
	// if err != nil {
	// 	return nil, err
	// }

	// _console.sprites[userSpriteBank1] = spritesSurface
	_console.sprites[userSpriteBank1] = sprites.(*image.RGBA)

	// initialise modes
	modes, err := _console.initModes()
	if err != nil {
		return nil, err
	}
	_console.modes = modes

	// init files
	// _console.currentDir = "/"
	// _console.baseDir, err = os.Getwd()
	// if err != nil {
	// 	return nil, err
	// }

	// text input
	//sdl.StartTextInput()

	return _console, nil
}

// func (c *console) GetWindow() *sdl.Window {
// 	return c.window
// }

func (c *console) GetBounds() image.Rectangle {
	// TODO
	//	return c.screen.Bounds()
	return image.Rect(0, 0, 0, 0)
}

func (c *console) SetMode(newMode ModeType) {
	c.Lock()
	defer c.Unlock()
	c.currentMode = newMode
}

func (c *console) LoadCart(cart Cartridge) error {
	c.cart = cart
	// init runtime mode
	runtime := newRuntimeMode(c)
	if err := runtime.LoadCart(cart); err != nil {
		return err
	}

	if err := runtime.Init(); err != nil {
		return err
	}

	c.modes[RUNTIME] = runtime
	return nil
}

var lastFrame time.Time
var startFrame time.Time
var endFrame time.Time

// Run is the main run loop
func (c *console) Run() error {

	// // default to runtime
	//c.SetMode(CLI)
	c.SetMode(RUNTIME)

	//go c.saveState()

	// poll events
	endFrame = time.Now() // init end frame
	startFrame = time.Now()
	//for !c.hasQuit {

	// if mode, ok := c.modes[c.currentMode]; ok {

	// 	if err := mode.Update(); err != nil {
	// 		return err
	// 	}

	// 	if err := mode.Render(); err != nil {
	// 		return err
	// 	}

	// 	// record frame
	// 	//			c.recorder.AddFrame(mode.GetFrame(), mode)

	// 	mode.Flip()

	// } else {
	// 	return fmt.Errorf("Mode :%d not found in console.modes", c.currentMode)
	// }

	//}

	// c.testImage = image.NewRGBA(image.Rect(0, 0, c.ConsoleWidth, c.ConsoleHeight))
	// rect := image.Rect(0, 0, c.ConsoleWidth, c.ConsoleHeight)
	// palette := color.Palette{
	// 	color.RGBA{0, 0, 0, 0},
	// 	color.RGBA{255, 0, 0, 255},
	// }
	// paletteImage := image.NewPaletted(rect, palette)

	// gc := gg.NewContextForImage(c.testImage)
	// gc.SetRGB(0, 0, 0)
	// gc.SetRGBA(1, 0, 0, 1)
	// gc.SetLineWidth(1)
	// gc.DrawLine(20, 20, 80, 60)
	// gc.Stroke()

	// draw onto palette
	//	draw.Draw(paletteImage, rect, gc.Image(), rect.Bounds().Min, draw.Over)
	// // draw back again
	//	draw.Draw(c.testImage, rect, paletteImage, rect.Bounds().Min, draw.Over)

	return ebiten.Run(c.update, c.Config.ConsoleWidth, c.Config.ConsoleHeight, 6, "pico-go")
}

func (c *console) update(screen *ebiten.Image) error {
	if ebiten.IsRunningSlowly() {
		return nil
	}

	if mode, ok := c.modes[c.currentMode]; ok {

		if err := mode.Update(); err != nil {
			return err
		}

		if err := mode.Render(); err != nil {
			return err
		}

		// record frame
		//			c.recorder.AddFrame(mode.GetFrame(), mode)

		mode.Flip()

		pb := mode.getPixelBuffer()
		//screen.ReplacePixels(pb.pixelSurface.Pix)

		// gc = gg.NewContextForImage(drawImage)
		// gc.SetRGB(0, 0, 0)
		// gc.SetRGBA(1, 0, 0, 1)
		// gc.SetLineWidth(1)
		// gc.DrawLine(20, 20, 80, 60)
		// gc.Stroke()

		// pb.gc.Stroke()
		// pb.gc.SetLineWidth(1)
		// pb.gc.SetLineCapSquare()
		// red := c.palette.GetColor(PICO8_RED)
		// pattern := gg.NewSolidPattern(red)
		// pb.gc.SetStrokeStyle(pattern)
		// //col := pb.palette.GetColor(PICO8_RED)
		// //r, g, b, _ := col.RGBA()
		// //pb.gc.SetRGB255(int(r), int(g), int(b))
		// //		pb.gc.SetRGB(0, 0, 0)
		// //		pb.gc.SetRGBA(1, 0, 0, 1)
		// //pb.gc.SetRGB(255, 0, 0)
		// pb.gc.DrawLine(20, 85, 50, 85)
		// pb.gc.Stroke()

		// set pixels
		// for x := 20; x < 50; x++ {
		// 	pb.pixelSurface.Set(x, 90, col)
		// }

		// draw onto palette
		draw.Draw(c.pImage, c.pImage.Bounds(), pb.gc.Image(), c.pImage.Bounds().Min, draw.Over)
		// draw back again
		draw.Draw(pb.pixelSurface, c.pImage.Bounds(), c.pImage, c.pImage.Bounds().Min, draw.Over)
		screen.ReplacePixels(pb.pixelSurface.Pix)

		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.CurrentFPS()))

	} else {
		return fmt.Errorf("Mode :%d not found in console.modes", c.currentMode)
	}

	return nil
}

func (c *console) handleEvents() error {
	//	if mode, ok := c.modes[c.currentMode]; ok {

	// TODO
	// poll all events
	// for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

	// 	// this is the common event handling code

	// 	// TODO keys to implement
	// 	// F7 Capture cartridge label image
	// 	// F8 Start recording a video
	// 	// F9 Save GIF video to desktop (max: 8 seconds by default)

	// 	switch t := event.(type) {
	// 	case *sdl.QuitEvent:
	// 		fmt.Printf("Quit event..\n")
	// 		c.state.SaveState(c)
	// 		c.hasQuit = true
	// 	case *sdl.KeyDownEvent:
	// 		switch t.Keysym.Sym {
	// 		case sdl.K_ESCAPE:
	// 			c.toggleCLI()
	// 		case sdl.K_F6:
	// 			if err := c.saveScreenshot(); err != nil {
	// 				return err
	// 			}
	// 		case sdl.K_F9:
	// 			if err := c.saveVideo(); err != nil {
	// 				return err
	// 			}
	// 		default:
	// 			// pass keydown events to mode handle
	// 			if err := mode.HandleEvent(event); err != nil {
	// 				return err
	// 			}
	// 		}
	// 	case *sdl.MouseButtonEvent:
	// 		// we only care about mouse clicks
	// 		if t.State == 1 && t.Button == 1 {
	// 			c.mouseClicked(t.X, t.Y)
	// 		}
	// 	default:
	// 		// if not handled pass event to mode event handler
	// 		if err := mode.HandleEvent(event); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	//	}
	return nil
}

func (c *console) mouseClicked(x, y int32) {
	// transform window x,y coords to pixel buffer coords

	fmt.Printf("Mouse clicked at x: %d y: %d\n", x, y)

	// get mode
	if mode, ok := c.modes[c.currentMode]; ok {
		pb := mode.getPixelBuffer()
		fmt.Printf("RenderRect: %#v\n", pb.renderRect)
		fmt.Printf("pixelBuffer: %#v\n", pb.psRect)
		// subtract top left offset
		x -= int32(pb.renderRect.Min.X)
		y -= int32(pb.renderRect.Min.Y)
		fmt.Printf("[adjusted] Mouse clicked at x: %d y: %d\n", x, y)
		// scale to match pixelbuffer
		scale := float32(pb.renderRect.Max.X) / float32(pb.pixelSurface.Bounds().Max.X)
		scaledX := float32(x) / scale
		scaledY := float32(y) / scale
		x = int32(scaledX)
		y = int32(scaledY)
		fmt.Printf("[scaled] Mouse clicked at x: %d y: %d\n", x, y)
	}

}

// saveState - saves console state periodically
func (c *console) saveState() {

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			// save state
			c.state.SaveState(c)

		}
	}
}

func (c *console) Quit() {
	c.hasQuit = true
}

// toggleCLI - toggle between CLI and secondary mode
func (c *console) toggleCLI() {
	switch c.currentMode {
	case CLI:
		c.SetMode(c.secondaryMode)
	case RUNTIME:
		if mode, ok := c.modes[c.currentMode]; ok {
			runtime := mode.(*runtime)
			runtime.Stop()
			c.SetMode(CLI)
		}

	default:
		c.secondaryMode = c.currentMode
		c.SetMode(CLI)
	}
}

// Destroy cleans up any resources at end
func (c *console) Destroy() {
	//c.window.Destroy()
}

// saveScreenshot - saves a screenshot of current frame
func (c *console) saveScreenshot() error {

	return c.recorder.SaveScreenshot("out.png", c.Config.ScreenshotScale)

}

// saveVideo - saves a video of last x seconds
func (c *console) saveVideo() error {
	return c.recorder.SaveVideo("out.gif", c.Config.GifScale)
}
