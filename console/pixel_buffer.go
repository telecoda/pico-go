package console

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	gfx "github.com/veandco/go-sdl2/sdl_gfx"
)

type mode struct {
	pixelBuffer
}

type pixelBuffer struct {
	textCursor   pos // note print pos in char/line pos not pixel pos
	fgColor      Color
	bgColor      Color
	charCols     int
	charRows     int
	pixelSurface *sdl.Surface // offscreen pixel buffer
	psRect       *sdl.Rect    // rect of pixelSurface
	renderer     *sdl.Renderer
}

type pos struct {
	x int
	y int
}

var rmask = uint32(0xff000000)
var gmask = uint32(0x00ff0000)
var bmask = uint32(0x0000ff00)
var amask = uint32(0x000000ff)

func newPixelBuffer(cfg Config) (PixelBuffer, error) {
	p := &pixelBuffer{}
	ps, err := sdl.CreateRGBSurface(0, int32(cfg.ConsoleWidth), int32(cfg.ConsoleHeight), 32, rmask, gmask, bmask, amask)
	if err != nil {
		return nil, err
	}

	p.pixelSurface = ps

	p.psRect = &sdl.Rect{X: 0, Y: 0, W: p.pixelSurface.W, H: p.pixelSurface.H}

	p.textCursor.x = 0
	p.textCursor.y = 0
	p.fgColor = WHITE
	p.bgColor = BLACK

	p.charCols = int(int32(cfg.ConsoleWidth) / _charWidth)
	p.charRows = int(int32(cfg.ConsoleHeight) / _charHeight)

	p.renderer, err = sdl.CreateSoftwareRenderer(p.pixelSurface)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *pixelBuffer) GetFrame() *sdl.Surface {
	return p.pixelSurface
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
	_, color := _console.palette.getRGBA(p.bgColor)
	p.pixelSurface.FillRect(p.psRect, color)
}

// ClsWithColor - fill pixel buffer with a set color
func (p *pixelBuffer) ClsWithColor(colorID Color) {
	p.bgColor = colorID
	p.Cls()
}

func (p *pixelBuffer) Color(colorID Color) {
	p.fgColor = colorID
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

// Print - prints string of characters to the screen with drawing color
func (p *pixelBuffer) Print(str string) {
	pixelPos := charToPixel(p.textCursor)

	p.PrintAtWithColor(str, 0, int(pixelPos.y), p.fgColor)

	// increase printPos by 1 line
	p.textCursor.y++

	if p.textCursor.y > p.charRows-2 {
		p.ScrollUpLine()
	}
}

// PrintAt - prints a string of characters to the screen at position with drawing color
func (p *pixelBuffer) PrintAt(str string, x, y int) {
	p.PrintAtWithColor(str, x, y, p.fgColor)
}

// PrintAtWithColor - prints a string of characters to the screen at position with color
func (p *pixelBuffer) PrintAtWithColor(str string, x, y int, colorID Color) {
	p.fgColor = colorID
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

// Drawer methods

// Circle - draw circle with drawing color
func (p *pixelBuffer) Circle(x, y, r int) {
	p.CircleWithColor(x, y, r, p.fgColor)
}

// CircleWithColor - draw circle with color
func (p *pixelBuffer) CircleWithColor(x, y, r int, colorID Color) {
	p.fgColor = colorID
	rgba, _ := _console.palette.getRGBA(p.fgColor)
	//	p.renderer.SetDrawColor(rgba.R, rgba.G, rgba.B, rgba.A)
	sColor := sdl.Color{R: rgba.R, G: rgba.G, B: rgba.B, A: rgba.A}
	gfx.CircleColor(p.renderer, x, y, r, sColor)
}

// CircleFill - fill circle with drawing color
func (p *pixelBuffer) CircleFill(x, y, r int) {
	p.CircleFillWithColor(x, y, r, p.fgColor)
}

// CircleFillWithColor - fill circle with color
func (p *pixelBuffer) CircleFillWithColor(x, y, r int, colorID Color) {
	p.fgColor = colorID
	rgba, _ := _console.palette.getRGBA(p.fgColor)
	//	p.renderer.SetDrawColor(rgba.R, rgba.G, rgba.B, rgba.A)
	sColor := sdl.Color{R: rgba.R, G: rgba.G, B: rgba.B, A: rgba.A}
	gfx.FilledCircleColor(p.renderer, x, y, r, sColor)
}

// Line - line in drawing color
func (p *pixelBuffer) Line(x0, y0, x1, y1 int) {
	p.LineWithColor(x0, y0, x1, y1, p.fgColor)
}

// LineWithColor - line with color
func (p *pixelBuffer) LineWithColor(x0, y0, x1, y1 int, colorID Color) {
	p.fgColor = colorID
	rgba, _ := _console.palette.getRGBA(p.fgColor)
	p.renderer.SetDrawColor(rgba.R, rgba.G, rgba.B, rgba.A)
	p.renderer.DrawLine(x0, y0, x1, y1)
}

// PGet - pixel get
func (p *pixelBuffer) PGet(x, y int) Color {
	pixels := p.pixelSurface.Pixels()
	// get specific pixel
	w := _console.ConsoleWidth
	offset := 4 * (y*w + x)
	r := pixels[offset+3]
	g := pixels[offset+2]
	b := pixels[offset+1]
	a := pixels[offset+0]

	rgba := rgba{R: r, G: g, B: b, A: a}

	return _console.palette.getColor(rgba)
}

// PSet - pixel set in drawing color
func (p *pixelBuffer) PSet(x, y int) {
	p.PSetWithColor(x, y, p.fgColor)
}

// PSetWithColor - pixel set with color
func (p *pixelBuffer) PSetWithColor(x0, y0 int, colorID Color) {
	p.fgColor = colorID
	rgba, _ := _console.palette.getRGBA(p.fgColor)
	p.renderer.SetDrawColor(rgba.R, rgba.G, rgba.B, rgba.A)
	p.renderer.DrawPoint(x0, y0)
}

// Rect - draw rectangle with drawing color
func (p *pixelBuffer) Rect(x0, y0, x1, y1 int) {
	p.RectWithColor(x0, y0, x1, y1, p.fgColor)
}

// RectWithColor - draw rectangle with color
func (p *pixelBuffer) RectWithColor(x0, y0, x1, y1 int, colorID Color) {
	p.fgColor = colorID
	rgba, _ := _console.palette.getRGBA(p.fgColor)
	rect := &sdl.Rect{X: int32(x0), Y: int32(y0), W: int32(x1 - x0), H: int32(y1 - y0)}
	p.renderer.SetDrawColor(rgba.R, rgba.G, rgba.B, rgba.A)
	p.renderer.DrawRect(rect)
}

// RectFill - fill rectangle with drawing color
func (p *pixelBuffer) RectFill(x0, y0, x1, y1 int) {
	p.RectFillWithColor(x0, y0, x1, y1, p.fgColor)
}

// RectFillWithColor - fill rectangle with color
func (p *pixelBuffer) RectFillWithColor(x0, y0, x1, y1 int, colorID Color) {
	p.fgColor = colorID
	_, color := _console.palette.getRGBA(colorID)
	fRect := &sdl.Rect{X: int32(x0), Y: int32(y0), W: int32(x1 - x0), H: int32(y1 - y0)}
	p.pixelSurface.FillRect(fRect, color)
}

// Spriter methods

func (p *pixelBuffer) Sprite(n, x, y, w, h, dw, dh int, rot float64, flipX, flipY bool) {
	sw := int32(w) * _spriteWidth
	sh := int32(h) * _spriteHeight

	var flip sdl.RendererFlip
	if flipX {
		flip = flip | sdl.FLIP_HORIZONTAL
	}
	if flipY {
		flip = flip | sdl.FLIP_VERTICAL
	}

	if flip == 0 {
		flip = sdl.FLIP_NONE
	}

	// create sprite surface, to copy a single sprite onto
	ss, err := sdl.CreateRGBSurface(0, sw, sh, 32, rmask, gmask, bmask, amask)
	if err != nil {
		fmt.Printf("Failed to create surface: %s\n", err)
		return
	}
	defer ss.Free()

	// convert sprite number into x,y pos
	xCell := n % _spritesPerLine
	yCell := (n - xCell) / _spritesPerLine

	xPos := int32(xCell * _spriteWidth)
	yPos := int32(yCell * _spriteHeight)

	// this is the rect to copy from sprite sheet
	spriteSrcRect := &sdl.Rect{X: xPos, Y: yPos, W: sw, H: sh}
	// this rect represents the size of the resulting sprite
	spriteRect := &sdl.Rect{X: 0, Y: 0, W: sw, H: sh}

	// copy sprite data from sprite sheet onto sprite surface
	err = _console.sprites.Blit(spriteSrcRect, ss, spriteRect)
	if err != nil {
		fmt.Printf("Failed to blit surface: %s\n", err)
		return
	}

	texture, err := p.renderer.CreateTextureFromSurface(ss)
	if err != nil {
		fmt.Printf("Failed to create texture: %s\n", err)
		return
	}
	defer texture.Destroy()

	centre := &sdl.Point{X: int32(dw / 2), Y: int32(dh / 2)}

	screenRect := &sdl.Rect{X: int32(x), Y: int32(y), W: int32(dw), H: int32(dh)}
	err = p.renderer.CopyEx(texture, spriteRect, screenRect, rot, centre, flip)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

}

// Destroy cleans up any resources at end
func (p *pixelBuffer) Destroy() {
	p.pixelSurface.Free()
}
