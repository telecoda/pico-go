package console

import (
	"github.com/veandco/go-sdl2/sdl"
)

func newPixelBuffer(cfg Config) (PixelBuffer, error) {
	p := &pixelBuffer{}
	rmask := uint32(0xff000000)
	gmask := uint32(0x00ff0000)
	bmask := uint32(0x0000ff00)
	amask := uint32(0x000000ff)
	ps, err := sdl.CreateRGBSurface(0, int32(cfg.ConsoleWidth), int32(cfg.ConsoleHeight), 32, rmask, gmask, bmask, amask)
	if err != nil {
		return nil, err
	}

	p.pixelSurface = ps

	p.psRect = &sdl.Rect{X: 0, Y: 0, W: p.pixelSurface.W, H: p.pixelSurface.H}

	return p, nil
}

func (p *pixelBuffer) Render() error {

	// clear offscreen buffer
	p.ClsColor(3)

	// draw to offscreen surface
	rect := sdl.Rect{X: 0, Y: 0, W: 64, H: 64}

	// draw white rect top corner
	p.pixelSurface.FillRect(&rect, 0xffffffff)

	pixels := p.pixelSurface.Pixels()
	// update specific pixel
	x := 50
	y := 50
	w := 128
	pixels[4*(y*w+x)+0] = 255 // r
	pixels[4*(y*w+x)+1] = 0   // g
	pixels[4*(y*w+x)+2] = 0   // b

	return p.Flip()

}

// API

// Cls - clears pixel buffer
func (p *pixelBuffer) Cls() {
	_, color := _console.palette.getRGBA(0)
	p.pixelSurface.FillRect(p.psRect, color)
}

// Cls - fill pixel buffer with a set color
func (p *pixelBuffer) ClsColor(colorId Color) {
	_, color := _console.palette.getRGBA(colorId)
	p.pixelSurface.FillRect(p.psRect, color)
}

// Flip - copy offscreen buffer to onscreen buffer
func (p *pixelBuffer) Flip() error {
	tex, err := _console.renderer.CreateTextureFromSurface(p.pixelSurface)
	if err != nil {
		return err
	}
	defer tex.Destroy()

	// clear window
	_console.renderer.Clear()
	// calc how big to render on window
	winW, winH := _console.window.GetSize()
	var winRect sdl.Rect
	x1 := int32(0)
	y1 := int32(0)

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
		sH = int32(winH)
		sW = int32(winH)
		diff := (winW - winH) / 2
		x1 = int32(diff)
	}

	if winW < winH {
		// taller (use full width)
		x1 = 0
		sH = int32(winW)
		sW = int32(winW)
		diff := (winH - winW) / 2
		y1 = int32(diff)
	}

	winRect = sdl.Rect{X: x1, Y: y1, W: sW, H: sH}

	// copy and scale offscreen buffer
	_console.renderer.Copy(tex, p.psRect, &winRect)

	_console.renderer.Present()

	return nil
}

// PrintColorAt a string of characters to the screen at position with color
func (p *pixelBuffer) PrintColorAt(str string, x, y int, colorId Color) {
	rgba, _ := _console.palette.getRGBA(colorId)
	sColor := sdl.Color{R: rgba.R, G: rgba.G, B: rgba.B, A: rgba.A}
	textSurface, err := _console.font.RenderUTF8_Blended(str, sColor)
	if err != nil {
		panic(err)
	}
	defer textSurface.Free()

	// copy text surface to offscreen buffer

	tRect := &sdl.Rect{X: 0, Y: 0, W: textSurface.W, H: textSurface.H}
	posRect := &sdl.Rect{X: int32(x), Y: int32(y), W: textSurface.W, H: textSurface.H}

	textSurface.Blit(tRect, p.pixelSurface, posRect)
}

// Destroy cleans up any resources at end
func (p *pixelBuffer) Destroy() {
	p.pixelSurface.Free()
}
