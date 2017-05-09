package console

import (
	"github.com/veandco/go-sdl2/sdl"
)

type mode struct {
	pixelBuffer
}

type pixelBuffer struct {
	textCursor   pos // note print pos in char/line pos not pixel pos
	printColor   Color
	charCols     int
	charRows     int
	pixelSurface *sdl.Surface // offscreen pixel buffer
	psRect       *sdl.Rect    // rect of pixelSurface
}

type pos struct {
	x int
	y int
}

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

	p.textCursor.x = 0
	p.textCursor.y = 0
	p.printColor = WHITE

	p.charCols = int(int32(cfg.ConsoleWidth) / _charWidth)
	p.charRows = int(int32(cfg.ConsoleHeight) / _charHeight)

	return p, nil
}

func (p *pixelBuffer) Render() error {

	// clear offscreen buffer
	p.ClsWithColor(3)

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

// ClsWithColor - fill pixel buffer with a set color
func (p *pixelBuffer) ClsWithColor(colorID Color) {
	_, color := _console.palette.getRGBA(colorID)
	p.pixelSurface.FillRect(p.psRect, color)
}

func (p *pixelBuffer) Color(colorID Color) {
	p.printColor = colorID
}

func (p *pixelBuffer) Cursor(x, y int) {
	p.textCursor.x = x
	p.textCursor.y = y
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

func (p *pixelBuffer) GetCursor() pos {
	return p.textCursor
}

func charToPixel(charPos pos) pos {
	return pos{
		x: charPos.x * _charWidth,
		y: charPos.y * _charHeight,
	}
}

func pixelToChar(pixelPos pos) pos {
	return pos{
		x: pixelPos.x / _charWidth,
		y: pixelPos.y / _charHeight,
	}
}

// ScrolUpLLine - scrolls display up a single line
func (p *pixelBuffer) ScrollUpLine() {
	fromRect := &sdl.Rect{X: 0, Y: _charHeight, W: p.pixelSurface.W, H: p.pixelSurface.H - _charHeight}
	toRect := &sdl.Rect{X: 0, Y: 0, W: p.pixelSurface.W, H: p.pixelSurface.H - _charHeight}
	p.pixelSurface.Blit(fromRect, p.pixelSurface, toRect)
	p.textCursor.y = p.charRows - 2
}

// Print string of characters to the screen
func (p *pixelBuffer) Print(str string) {
	pixelPos := charToPixel(p.textCursor)

	p.PrintAtWithColor(str, 0, int(pixelPos.y), p.printColor)

	// increase printPos by 1 line
	p.textCursor.y++

	if p.textCursor.y > p.charRows-2 {
		p.ScrollUpLine()
	}
}

// PrintAtWithColor a string of characters to the screen at position with color
func (p *pixelBuffer) PrintAtWithColor(str string, x, y int, colorID Color) {
	p.printColor = colorID
	if str != "" {
		rgba, _ := _console.palette.getRGBA(colorID)
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
	// save print pos
	p.textCursor = pixelToChar(pos{x, y})

}

func (p *pixelBuffer) RectFillWithColor(x0, y0, x1, y1 int, colorID Color) {

	_, color := _console.palette.getRGBA(colorID)

	fRect := &sdl.Rect{X: int32(x0), Y: int32(y0), W: int32(x1 - x0), H: int32(y1 - x0)}

	p.pixelSurface.FillRect(fRect, color)

}

// Destroy cleans up any resources at end
func (p *pixelBuffer) Destroy() {
	p.pixelSurface.Free()
}
