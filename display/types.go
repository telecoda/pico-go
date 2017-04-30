package display

import (
	"github.com/telecoda/pico-go/config"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_gfx"
)

// The display package handles all the output to screen

var title = "pico-go virtual games console"

type Display interface {
	Render() error
	Destroy()
}

type size struct {
	width  uint
	height uint
}

type display struct {
	config       config.Config
	window       *sdl.Window
	pixelSurface *sdl.Surface
	renderer     *sdl.Renderer
	fpsMan       *gfx.FPSmanager
}

func New(cfg config.Config) (Display, error) {
	d := &display{
		config: cfg,
	}

	// init SDL
	sdl.Init(sdl.INIT_EVERYTHING)
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")

	// create window
	window, err := sdl.CreateWindow(
		title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int(d.config.WindowWidth), int(d.config.WindowHeight), sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE)

	if err != nil {
		return nil, err
	}

	d.window = window

	// create offset surface for pixel rendering
	rmask := uint32(0xff000000)
	gmask := uint32(0x00ff0000)
	bmask := uint32(0x0000ff00)
	amask := uint32(0x000000ff)
	ps, err := sdl.CreateRGBSurface(0, int32(d.config.ConsoleWidth), int32(d.config.ConsoleHeight), 32, rmask, gmask, bmask, amask)
	if err != nil {
		return nil, err
	}

	d.pixelSurface = ps

	// create renderer
	r, _ := sdl.CreateRenderer(window, -2, sdl.RENDERER_ACCELERATED)
	if r == nil {
		// revert to software
		r, _ = sdl.CreateRenderer(window, -2, sdl.RENDERER_SOFTWARE)
		if r == nil {
			return nil, err
		}
	}

	d.renderer = r

	return d, nil
}

func (d *display) Render() error {

	// draw to offscreen surface
	rect := sdl.Rect{X: 0, Y: 0, W: 64, H: 64}
	vcRect := sdl.Rect{X: 0, Y: 0, W: d.pixelSurface.W, H: d.pixelSurface.H}
	// clear offscreen buffer
	d.pixelSurface.FillRect(&vcRect, 0x000000ff)
	// draw white rect top corner
	d.pixelSurface.FillRect(&rect, 0xffffffff)

	pixels := d.pixelSurface.Pixels()
	// update specific pixel
	x := 50
	y := 50
	w := 128
	pixels[4*(y*w+x)+0] = 255 // r
	pixels[4*(y*w+x)+1] = 0   // g
	pixels[4*(y*w+x)+2] = 0   // b

	tex, err := d.renderer.CreateTextureFromSurface(d.pixelSurface)
	if err != nil {
		return err
	}

	// clear screen
	d.renderer.Clear()

	// calc how big to render on window

	winW, winH := d.window.GetSize()
	var winRect sdl.Rect
	x1 := int32(0)
	//x2 := int32(0)
	y1 := int32(0)
	//y2 := int32(0)

	// sW, sH - screen width + height
	sW := int32(winW)
	sH := int32(winH)

	// maintain aspect ratio even on resize
	if winW == winH {
		// same dimensions (no padding)
		sW = int32(winW)
		sH = int32(winH)
	}

	if winW > winH {
		// wider (use full height)
		y1 = 0
		//y2 = 0
		sH = int32(winH)
		sW = int32(winH)
		diff := (winW - winH) / 2
		x1 = int32(diff)
		//x2 = int32(diff)

	}

	if winW < winH {
		// taller (use full width)
		x1 = 0
		//x2 = 0
		sH = int32(winW)
		sW = int32(winW)
		diff := (winH - winW) / 2
		y1 = int32(diff)
		//y2 = int32(diff)
	}

	winRect = sdl.Rect{X: x1, Y: y1, W: sW, H: sH}

	// copy and scale offscreen buffer
	d.renderer.Copy(tex, &vcRect, &winRect)

	d.renderer.Present()
	return nil
}

func (d *display) Destroy() {
	if d.window != nil {
		d.window.Destroy()
	}
}
