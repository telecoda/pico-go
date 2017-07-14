package console

import (
	"fmt"
	"time"

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
	palette      *palette
	charCols     int
	charRows     int
	pixelSurface *sdl.Surface // offscreen pixel buffer
	psRect       *sdl.Rect    // rect of pixelSurface
	renderer     *sdl.Renderer
	fps          int
	timeBudget   int64
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
	p := &pixelBuffer{
		fps: cfg.FPS,
	}

	p.timeBudget = time.Duration(1*time.Second).Nanoseconds() / int64(p.fps)

	ps, err := sdl.CreateRGBSurface(0, int32(cfg.ConsoleWidth), int32(cfg.ConsoleHeight), 8, 0, 0, 0, 0)
	if err != nil {
		return nil, err
	}

	if ps == nil {
		return nil, fmt.Errorf("Surface is nil")
	}

	p.palette = cfg.palette

	if err := setSurfacePalette(p.palette, ps); err != nil {
		return nil, err
	}

	p.pixelSurface = ps

	p.psRect = &sdl.Rect{X: 0, Y: 0, W: p.pixelSurface.W, H: p.pixelSurface.H}

	p.textCursor.x = 0
	p.textCursor.y = 0
	p.fgColor = 7
	p.bgColor = 0

	p.charCols = cfg.ConsoleWidth / _console.Config.fontWidth
	p.charRows = cfg.ConsoleHeight / _console.Config.fontHeight

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

	// this is never called, always locally implemented

	return nil

}

// API

// Cls - clears pixel buffer
func (p *pixelBuffer) Cls() {
	p.pixelSurface.FillRect(p.psRect, uint32(p.bgColor))
}

// ClsWithColor - fill pixel buffer with a set color
func (p *pixelBuffer) ClsWithColor(colorID Color) {
	p.bgColor = colorID
	p.Cls()
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

	// calc how big to render on window
	winW, winH := _console.window.GetSize()

	// clear window
	fullRect := &sdl.Rect{X: 0, Y: 0, W: int32(winW), H: int32(winH)}
	rgba, _ := p.palette.GetRGBA(_console.borderColor)
	_console.renderer.SetDrawColor(rgba.R, rgba.G, rgba.B, rgba.A)
	_console.renderer.FillRect(fullRect)

	var winRect sdl.Rect
	x1 := int32(0)
	y1 := int32(0)

	// sW, sH - screen width + height
	sW := int32(winW)
	sH := int32(winH)

	// aspect ratio
	ratio := float64(_console.ConsoleHeight) / float64(_console.ConsoleWidth)

	// maintain aspect ratio even on resize
	if winW == winH {
		// same dimensions (no padding)
		sW = int32(winW)
		sH = int32(float64(winH) * ratio)
	}

	if winW < winH {
		y1 = 0
		sH = int32(float64(winW) * ratio)
		sW = int32(winW)
		diff := (winH - int(sH)) / 2
		y1 = int32(diff)
		if diff < 0 {
			y1 = 0
			sW = int32(float64(winH) * ratio)
			sH = int32(winH)
			diff := (winW - int(sW)) / 2
			x1 = int32(diff)
		}
	}

	if winW > winH {
		x1 = 0
		sH = int32(winH)
		sW = int32(float64(winH) / ratio)
		diff := (winW - int(sW)) / 2
		x1 = int32(diff)
		if diff < 0 {
			x1 = 0
			sW = int32(winW)
			sH = int32(float64(winW) * ratio)
			diff := (winH - int(sH)) / 2
			y1 = int32(diff)
		}
	}

	x1 += int32(_console.BorderWidth)
	y1 += int32(_console.BorderWidth)
	sH -= int32(_console.BorderWidth * 2)
	sW -= int32(_console.BorderWidth * 2)

	winRect = sdl.Rect{X: x1, Y: y1, W: sW, H: sH}

	// copy and scale offscreen buffer
	_console.renderer.Copy(tex, p.psRect, &winRect)

	_console.renderer.Present()

	p.lockFps()

	// record frame
	_console.recorder.AddFrame(p.GetFrame(), p)

	// at end of frame delay start timing for next one
	startFrame = time.Now()

	// handle events
	_console.handleEvents()

	return nil
}

// lockFps - locks rendering to a steady framerate
func (p *pixelBuffer) lockFps() float64 {

	var timeBudget = p.timeBudget
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

func (p *pixelBuffer) GetCursor() pos {
	return p.textCursor
}

func charToPixel(charPos pos) pos {
	return pos{
		x: charPos.x * _console.Config.fontWidth,
		y: charPos.y * _console.Config.fontHeight,
	}
}

func pixelToChar(pixelPos pos) pos {
	return pos{
		x: pixelPos.x / _console.Config.fontWidth,
		y: pixelPos.y / _console.Config.fontHeight,
	}
}

// ScrolUpLLine - scrolls display up a single line
func (p *pixelBuffer) ScrollUpLine() {
	fromRect := &sdl.Rect{X: 0, Y: int32(_console.Config.fontHeight), W: p.pixelSurface.W, H: p.pixelSurface.H - int32(_console.Config.fontHeight)}
	toRect := &sdl.Rect{X: 0, Y: 0, W: p.pixelSurface.W, H: p.pixelSurface.H - int32(_console.Config.fontHeight)}
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
		rgbaFg, _ := p.palette.GetRGBA(colorID)
		fgColor := sdl.Color{R: rgbaFg.R, G: rgbaFg.G, B: rgbaFg.B, A: rgbaFg.A}
		textSurface, err := _console.font.RenderUTF8_Solid(str, fgColor)
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
	rgba, _ := p.palette.GetRGBA(p.fgColor)
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
	rgba, _ := p.palette.GetRGBA(p.fgColor)
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
	rgba, _ := p.palette.GetRGBA(p.fgColor)
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

	return p.palette.GetColorID(rgba)
}

// PSet - pixel set in drawing color
func (p *pixelBuffer) PSet(x, y int) {
	p.PSetWithColor(x, y, p.fgColor)
}

// PSetWithColor - pixel set with color
func (p *pixelBuffer) PSetWithColor(x0, y0 int, colorID Color) {
	p.fgColor = colorID
	rgba, _ := p.palette.GetRGBA(p.fgColor)
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
	rgba, _ := p.palette.GetRGBA(p.fgColor)
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
	fRect := &sdl.Rect{X: int32(x0), Y: int32(y0), W: int32(x1 - x0), H: int32(y1 - y0)}
	p.pixelSurface.FillRect(fRect, uint32(colorID))
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
	ss1, err := sdl.CreateRGBSurface(0, int32(sw), int32(sh), 32, 0, 0, 0, 0)
	if err != nil {
		fmt.Printf("Failed to create surface1: %s\n", err)
		return
	}
	defer ss1.Free()

	// convert sprite number into x,y pos
	xCell := n % _spritesPerLine
	yCell := (n - xCell) / _spritesPerLine

	xPos := int32(xCell * _spriteWidth)
	yPos := int32(yCell * _spriteHeight)

	// this is the rect to copy from sprite sheet
	spriteSrcRect := &sdl.Rect{X: xPos, Y: yPos, W: sw, H: sh}
	// this rect represents the size of the resulting sprite
	ss1Rect := &sdl.Rect{X: 0, Y: 0, W: ss1.W, H: ss1.H}

	// set palette for sprites based on current palette
	if err := setSurfacePalette(p.palette, _console.sprites); err != nil {
		fmt.Printf("Failed to update sprite surface palette: %s\n", err)
		return
	}

	// copy sprite data from sprite sheet onto sprite surface
	err = _console.sprites.Blit(spriteSrcRect, ss1, ss1Rect)
	if err != nil {
		fmt.Printf("Failed to blit surface1: %s\n", err)
		return
	}

	// create 2nd sprite for blitscaling
	ss2, err := sdl.CreateRGBSurface(0, int32(dw), int32(dh), 32, 0, 0, 0, 0)
	if err != nil {
		fmt.Printf("Failed to create surface2: %s\n", err)
		return
	}
	defer ss2.Free()

	ss2Rect := &sdl.Rect{X: 0, Y: 0, W: ss2.W, H: ss2.H}

	// copy sprite data from sprite sheet onto sprite surface
	err = ss1.BlitScaled(ss1Rect, ss2, ss2Rect)
	if err != nil {
		fmt.Printf("Failed to blit surface2: %s\n", err)
		return
	}

	texture, err := p.renderer.CreateTextureFromSurface(ss2)
	if err != nil {
		fmt.Printf("Failed to create texture: %s\n", err)
		return
	}
	defer texture.Destroy()

	centre := &sdl.Point{X: int32(dw / 2), Y: int32(dh / 2)}

	screenRect := &sdl.Rect{X: int32(x), Y: int32(y), W: int32(dw), H: int32(dh)}

	err = p.renderer.CopyEx(texture, ss2Rect, screenRect, rot, centre, flip)
	if err != nil {
		fmt.Printf("Error: 1 %s\n", err)
	}

}

// Color - Set current drawing color
func (p *pixelBuffer) Color(colorID Color) {
	p.fgColor = colorID
}

// Paletter methods

// getRGBA - returns color as Color and uint32
func (p *pixelBuffer) GetRGBA(color Color) (rgba, uint32) {
	return p.palette.GetRGBA(color)

}

// GetColorID - find color from rgba
func (p *pixelBuffer) GetColorID(rgba rgba) Color {
	return p.palette.GetColorID(rgba)
}

func (p *pixelBuffer) PaletteReset() {
	p.palette.PaletteReset()
	// reset palette on sprites
	setSurfacePalette(p.palette, _console.sprites)
	// reset palette on pixel buffer
	setSurfacePalette(p.palette, p.pixelSurface)
}

func (p *pixelBuffer) PaletteCopy() Paletter {
	return p.palette.PaletteCopy()
}

func (p *pixelBuffer) GetSDLColors() []sdl.Color {
	return p.palette.GetSDLColors()
}

func (p *pixelBuffer) MapColor(fromColor Color, toColor Color) error {
	if err := p.palette.MapColor(fromColor, toColor); err != nil {
		return err
	}
	// update palette for surface
	return setSurfacePalette(p.palette, p.pixelSurface)
}

func (p *pixelBuffer) SetTransparent(color Color, enabled bool) error {
	if err := p.palette.SetTransparent(color, enabled); err != nil {
		return err
	}
	// update palette for surface
	return setSurfacePalette(p.palette, p.pixelSurface)
}

// Destroy cleans up any resources at end
func (p *pixelBuffer) Destroy() {
	p.pixelSurface.Free()
}
