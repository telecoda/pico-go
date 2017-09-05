package console

import (
	"fmt"
	"image"
	"image/color"
	"time"

	drawx "golang.org/x/image/draw"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten"
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
	pixelSurface *image.RGBA // offscreen pixel buffer
	screen       *ebiten.Image
	gc           *gg.Context     // graphics context
	psRect       image.Rectangle // rect of pixelSurface
	renderRect   image.Rectangle // rect on main window that pixelbuffer is rendered into
	//renderer     *sdl.Renderer
	fps        int
	timeBudget int64
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

	// ps, err := sdl.CreateRGBSurface(0, int32(cfg.ConsoleWidth), int32(cfg.ConsoleHeight), 8, 0, 0, 0, 0)
	// if err != nil {
	// 	return nil, err
	// }
	p.psRect = image.Rect(0, 0, cfg.ConsoleWidth, cfg.ConsoleHeight)
	p.renderRect = image.Rect(0, 0, cfg.ConsoleWidth, cfg.ConsoleHeight)
	ps := image.NewRGBA(p.psRect)

	if ps == nil {
		return nil, fmt.Errorf("Surface is nil")
	}

	p.palette = cfg.palette

	if err := setSurfacePalette(p.palette, ps); err != nil {
		return nil, err
	}

	var err error
	screen, err := ebiten.NewImageFromImage(ps, ebiten.FilterNearest)
	if err != nil {
		return nil, fmt.Errorf("Failed to create ebiten image %s", err)
	}
	p.screen = screen
	p.pixelSurface = ps
	// p.psRect = &sdl.Rect{X: 0, Y: 0, W: p.pixelSurface.W, H: p.pixelSurface.H}
	// p.renderRect = &sdl.Rect{X: 0, Y: 0, W: p.pixelSurface.W, H: p.pixelSurface.H}

	p.textCursor.x = 0
	p.textCursor.y = 0
	p.fgColor = 7
	p.bgColor = 0

	p.charCols = cfg.ConsoleWidth / _console.Config.fontWidth
	p.charRows = cfg.ConsoleHeight / _console.Config.fontHeight

	// p.renderer, err = sdl.CreateSoftwareRenderer(p.pixelSurface)
	// if err != nil {
	// 	return nil, err
	// }
	p.gc = gg.NewContextForRGBA(ps)
	// if err := p.gc.LoadFontFace("/Library/Fonts/Arial.ttf", 6); err != nil {
	// 	panic(err)
	// }
	p.gc.SetFontFace(_console.font)
	p.gc.SetLineWidth(1)
	p.gc.Identity()
	p.gc.SetStrokeStyle(gg.NewSolidPattern(color.White))

	return p, nil
}

func (p *pixelBuffer) GetFrame() *image.RGBA {
	return p.pixelSurface
}

func (p *pixelBuffer) Render() error {

	// this is never called, always locally implemented

	return nil

}

// API

// Cls - clears pixel buffer
func (p *pixelBuffer) Cls() {
	p.gc.SetColor(p.palette.GetColor(p.bgColor))
	p.gc.Clear()
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

	// tex, err := _console.renderer.CreateTextureFromSurface(p.pixelSurface)
	// if err != nil {
	// 	return err
	// }
	// defer tex.Destroy()

	// calc how big to render on window
	// winW, winH := _console.window.GetSize()

	// // clear window
	// fullRect := &sdl.Rect{X: 0, Y: 0, W: int32(winW), H: int32(winH)}
	// rgba, _ := p.palette.GetRGBA(_console.BorderColor)
	// _console.renderer.SetDrawColor(rgba.R, rgba.G, rgba.B, rgba.A)
	// _console.renderer.FillRect(fullRect)

	// //var renderRect sdl.Rect
	// x1 := int32(0)
	// y1 := int32(0)

	// // sW, sH - screen width + height
	// sW := int32(winW)
	// sH := int32(winH)

	// // aspect ratio
	// ratio := float64(_console.ConsoleHeight) / float64(_console.ConsoleWidth)

	// // maintain aspect ratio even on resize
	// if winW == winH {
	// 	// same dimensions (no padding)
	// 	sW = int32(winW)
	// 	sH = int32(float64(winH) * ratio)
	// }

	// if winW < winH {
	// 	y1 = 0
	// 	sH = int32(float64(winW) * ratio)
	// 	sW = int32(winW)
	// 	diff := (winH - int(sH)) / 2
	// 	y1 = int32(diff)
	// 	if diff < 0 {
	// 		y1 = 0
	// 		sW = int32(float64(winH) * ratio)
	// 		sH = int32(winH)
	// 		diff := (winW - int(sW)) / 2
	// 		x1 = int32(diff)
	// 	}
	// }

	// if winW > winH {
	// 	x1 = 0
	// 	sH = int32(winH)
	// 	sW = int32(float64(winH) / ratio)
	// 	diff := (winW - int(sW)) / 2
	// 	x1 = int32(diff)
	// 	if diff < 0 {
	// 		x1 = 0
	// 		sW = int32(winW)
	// 		sH = int32(float64(winW) * ratio)
	// 		diff := (winH - int(sH)) / 2
	// 		y1 = int32(diff)
	// 	}
	// }

	// x1 += int32(_console.BorderWidth)
	// y1 += int32(_console.BorderWidth)
	// sH -= int32(_console.BorderWidth * 2)
	// sW -= int32(_console.BorderWidth * 2)

	// p.renderRect.X = x1
	// p.renderRect.Y = y1
	// p.renderRect.W = sW
	// p.renderRect.H = sH

	// // copy and scale offscreen buffer
	// _console.renderer.Copy(tex, p.psRect, p.renderRect)

	// _console.renderer.Present()

	// TODO update ebiten screen

	// update
	// func update(screen *ebiten.Image) error {

	// 	if ebiten.IsRunningSlowly() {
	// 		return nil
	// 	}
	// 	screen.ReplacePixels(_console.pixelBuffer..Pix)
	// 	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.CurrentFPS()))
	// 	return nil
	// }

	// p.lockFps()

	// record frame
	_console.recorder.AddFrame(p.GetFrame(), p)

	// at end of frame delay start timing for next one
	startFrame = time.Now()

	// handle events
	_console.handleEvents()

	return nil
}

// func update(screen *ebiten.Image) error {

// 	if ebiten.IsRunningSlowly() {
// 		return nil
// 	}
// 	screen.ReplacePixels(_console.pixelBuffer..Pix)
// 	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.CurrentFPS()))
// 	return nil
// }

// lockFps - locks rendering to a steady framerate
func (p *pixelBuffer) lockFps() float64 {

	var timeBudget = p.timeBudget
	now := time.Now()
	// calc time to process frame so since start
	procTime := now.Sub(startFrame)

	// delay for remainder of time budget (based on fps)
	delay := time.Duration(timeBudget - procTime.Nanoseconds())
	if delay > 0 {
		// TODO
		//sdl.Delay(uint32(delay / 1000000))
	}

	// calc actual fps being achieved
	endFrame = time.Now()
	frameTime := endFrame.Sub(startFrame)

	endFrame = time.Now()

	return float64(time.Second) / float64(frameTime.Nanoseconds())
}

func (p *pixelBuffer) getPixelBuffer() *pixelBuffer {
	return p
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
	// TODO
	// fromRect := &sdl.Rect{X: 0, Y: int32(_console.Config.fontHeight), W: p.pixelSurface.W, H: p.pixelSurface.H - int32(_console.Config.fontHeight)}
	// toRect := &sdl.Rect{X: 0, Y: 0, W: p.pixelSurface.W, H: p.pixelSurface.H - int32(_console.Config.fontHeight)}
	// p.pixelSurface.Blit(fromRect, p.pixelSurface, toRect)
	// p.textCursor.y = p.charRows - 2
}

// Print - prints string of characters to the screen with drawing color
func (p *pixelBuffer) Print(str string) {
	pixelPos := charToPixel(p.textCursor)

	p.PrintAtWithColor(str, int(pixelPos.x), int(pixelPos.y), p.fgColor)

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
		p.gc.SetColor(_console.palette.GetColor(colorID))
		p.gc.DrawString(str, float64(x), float64(y+_console.fontHeight/2))
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
func (p *pixelBuffer) CircleWithColor(x0, y0, r int, colorID Color) {
	p.fgColor = colorID

	x := 0
	y := r
	p0 := (5 - r*4) / 4

	/* circle calcs from this blog
	http://xiaohuiliucuriosity.blogspot.co.uk/2015/03/draw-circle-using-integer-arithmetic.html
	*/

	p.circlePoints(x0, y0, x, y)
	for x < y {
		x++
		if p0 < 0 {
			p0 += 2*x + 1
		} else {
			y--
			p0 += 2*(x-y) + 1
		}
		p.circlePoints(x0, y0, x, y)
	}

}

func (p *pixelBuffer) circlePoints(cx, cy, x, y int) {

	if x == 0 {
		p.PSet(cx, cy+y)
		p.PSet(cx, cy-y)
		p.PSet(cx+y, cy)
		p.PSet(cx-y, cy)
	} else if x == y {
		p.PSet(cx+x, cy+y)
		p.PSet(cx-x, cy+y)
		p.PSet(cx+x, cy-y)
		p.PSet(cx-x, cy-y)
	} else if x < y {
		p.PSet(cx+x, cy+y)
		p.PSet(cx-x, cy+y)
		p.PSet(cx+x, cy-y)
		p.PSet(cx-x, cy-y)
		p.PSet(cx+y, cy+x)
		p.PSet(cx-y, cy+x)
		p.PSet(cx+y, cy-x)
		p.PSet(cx-y, cy-x)
	}
}

// CircleFill - fill circle with drawing color
func (p *pixelBuffer) CircleFill(x, y, r int) {
	p.CircleFillWithColor(x, y, r, p.fgColor)
}

// CircleFillWithColor - fill circle with color
func (p *pixelBuffer) CircleFillWithColor(x0, y0, r int, colorID Color) {
	p.fgColor = colorID
	// rgba, _ := p.palette.GetRGBA(p.fgColor)
	// p.gc.SetRGBA255(int(rgba.R), int(rgba.G), int(rgba.B), int(rgba.A))
	// p.gc.DrawCircle(float64(x), float64(y), float64(r))
	// p.gc.Fill()

	x := 0
	y := r
	p0 := (5 - r*4) / 4

	/* circle calcs from this blog
	http://groups.csail.mit.edu/graphics/classes/6.837/F98/Lecture6/circle.html
	*/

	p.circleLines(x0, y0, x, y)
	for x < y {
		x++
		if p0 < 0 {
			p0 += 2*x + 1
		} else {
			y--
			p0 += 2*(x-y) + 1
		}
		p.circleLines(x0, y0, x, y)
	}
}

func (p *pixelBuffer) circleLines(cx, cy, x, y int) {
	p.Line(cx-x, cy+y, cx+x, cy+y)
	p.Line(cx-x, cy-y, cx+x, cy-y)
	p.Line(cx-y, cy+x, cx+y, cy+x)
	p.Line(cx-y, cy-x, cx+y, cy-x)
}

// Line - line in drawing color
func (p *pixelBuffer) Line(x0, y0, x1, y1 int) {
	p.LineWithColor(x0, y0, x1, y1, p.fgColor)
}

// LineWithColor - line with color
func (p *pixelBuffer) LineWithColor(x1, y1, x2, y2 int, colorID Color) {
	p.setFGColor(colorID)

	col := p.palette.GetColor(colorID)

	/* Code from
	https://github.com/StephaneBunel/bresenham/blob/master/drawline.go#L12-L22
	*/
	// draw line

	var dx, dy, e, slope int

	// Because drawing p1 -> p2 is equivalent to draw p2 -> p1,
	// I sort points in x-axis order to handle only half of possible cases.
	if x1 > x2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}

	dx, dy = x2-x1, y2-y1
	// Because point is x-axis ordered, dx cannot be negative
	if dy < 0 {
		dy = -dy
	}

	switch {

	// Is line a point ?
	case x1 == x2 && y1 == y2:
		p.pixelSurface.Set(x1, y1, col)

	// Is line an horizontal ?
	case y1 == y2:
		for ; dx != 0; dx-- {
			p.pixelSurface.Set(x1, y1, col)
			x1++
		}
		p.pixelSurface.Set(x1, y1, col)

	// Is line a vertical ?
	case x1 == x2:
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		for ; dy != 0; dy-- {
			p.pixelSurface.Set(x1, y1, col)
			y1++
		}
		p.pixelSurface.Set(x1, y1, col)

	// Is line a diagonal ?
	case dx == dy:
		if y1 < y2 {
			for ; dx != 0; dx-- {
				p.pixelSurface.Set(x1, y1, col)
				x1++
				y1++
			}
		} else {
			for ; dx != 0; dx-- {
				p.pixelSurface.Set(x1, y1, col)
				x1++
				y1--
			}
		}
		p.pixelSurface.Set(x1, y1, col)

	// wider than high ?
	case dx > dy:
		if y1 < y2 {
			// BresenhamDxXRYD(img, x1, y1, x2, y2, col)
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				p.pixelSurface.Set(x1, y1, col)
				x1++
				e -= dy
				if e < 0 {
					y1++
					e += slope
				}
			}
		} else {
			// BresenhamDxXRYU(img, x1, y1, x2, y2, col)
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				p.pixelSurface.Set(x1, y1, col)
				x1++
				e -= dy
				if e < 0 {
					y1--
					e += slope
				}
			}
		}
		p.pixelSurface.Set(x2, y2, col)

	// higher than wide.
	default:
		if y1 < y2 {
			// BresenhamDyXRYD(img, x1, y1, x2, y2, col)
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				p.pixelSurface.Set(x1, y1, col)
				y1++
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		} else {
			// BresenhamDyXRYU(img, x1, y1, x2, y2, col)
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				p.pixelSurface.Set(x1, y1, col)
				y1--
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		}
		p.pixelSurface.Set(x2, y2, col)
	}
}

func (p *pixelBuffer) setFGColor(colorID Color) {
	p.fgColor = colorID
	//	p.gc.SetColor(p.palette.GetColor(colorID))
	c := p.palette.GetColor(colorID)
	r, g, b, _ := c.RGBA()
	p.gc.SetRGB255(int(r), int(g), int(b))
}

// PGet - pixel get
func (p *pixelBuffer) PGet(x, y int) Color {

	c := p.pixelSurface.At(x, y)
	r, g, b, a := c.RGBA()
	color := rgba{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}

	return p.palette.GetColorID(color)
}

// PSet - pixel set in drawing color
func (p *pixelBuffer) PSet(x, y int) {
	p.PSetWithColor(x, y, p.fgColor)
}

// PSetWithColor - pixel set with color
func (p *pixelBuffer) PSetWithColor(x0, y0 int, colorID Color) {
	p.setFGColor(colorID)
	p.pixelSurface.Set(x0, y0, p.palette.GetColor(colorID))
}

// Rect - draw rectangle with drawing color
func (p *pixelBuffer) Rect(x0, y0, x1, y1 int) {
	p.RectWithColor(x0, y0, x1, y1, p.fgColor)
}

// RectWithColor - draw rectangle with color
func (p *pixelBuffer) RectWithColor(x0, y0, x1, y1 int, colorID Color) {
	p.fgColor = colorID
	p.Line(x0, y0, x1, y0)
	p.Line(x1, y0, x1, y1)
	p.Line(x1, y1, x0, y1)
	p.Line(x0, y1, x0, y0)
}

// RectFill - fill rectangle with drawing color
func (p *pixelBuffer) RectFill(x0, y0, x1, y1 int) {
	p.RectFillWithColor(x0, y0, x1, y1, p.fgColor)
}

// RectFillWithColor - fill rectangle with color
func (p *pixelBuffer) RectFillWithColor(x0, y0, x1, y1 int, colorID Color) {
	p.fgColor = colorID
	for x := x0; x < x1; x++ {
		p.Line(x, y0, x, y1)
	}
}

// Spriter methods

func (p *pixelBuffer) systemSprite(n, x, y, w, h, dw, dh int, rot float64, flipX, flipY bool) {
	_console.currentSpriteBank = systemSpriteBank
	p.sprite(n, x, y, w, h, dw, dh, rot, flipX, flipY)
}

func (p *pixelBuffer) Sprite(n, x, y, w, h, dw, dh int, rot float64, flipX, flipY bool) {
	_console.currentSpriteBank = userSpriteBank1
	p.sprite(n, x, y, w, h, dw, dh, rot, flipX, flipY)

}

func (p *pixelBuffer) sprite(n, x, y, w, h, dw, dh int, rot float64, flipX, flipY bool) {

	sw := w * _spriteWidth
	sh := h * _spriteHeight

	// var flip sdl.RendererFlip
	// if flipX {
	// 	flip = flip | sdl.FLIP_HORIZONTAL
	// }
	// if flipY {
	// 	flip = flip | sdl.FLIP_VERTICAL
	// }

	// if flip == 0 {
	// 	flip = sdl.FLIP_NONE
	// }

	// TOOD rotation and flipping not supported yet

	// convert sprite number into x,y pos
	xCell := n % _spritesPerLine
	yCell := (n - xCell) / _spritesPerLine

	xPos := xCell * _spriteWidth
	yPos := yCell * _spriteHeight

	// this is the rect to copy from sprite sheet
	spriteSrcRect := image.Rect(xPos, yPos, xPos+sw, yPos+sh)
	// this rect is where the sprite will be copied to
	screenRect := image.Rect(x, y, x+dw, y+dh)

	options := &drawx.Options{
		SrcMask:  _console.sprites[userSpriteMask1],
		SrcMaskP: image.Point{0, 0},
	}
	fmt.Printf("TEMP: options: %#v\n", options)
	drawx.NearestNeighbor.Scale(p.pixelSurface, screenRect, _console.sprites[userSpriteBank1], spriteSrcRect, drawx.Over, options)

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

// GetColor - get color by id
func (p *pixelBuffer) GetColor(colorID Color) color.Color {
	return p.palette.GetColor(colorID)
}

// GetColorID - find color from rgba
func (p *pixelBuffer) GetColorID(rgba rgba) Color {
	return p.palette.GetColorID(rgba)
}

func (p *pixelBuffer) PaletteReset() {
	p.palette.PaletteReset()
}

func (p *pixelBuffer) PaletteCopy() Paletter {
	return p.palette.PaletteCopy()
}

func (p *pixelBuffer) GetColors() []color.Color {
	return p.palette.GetColors()
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
	p.pixelSurface = nil
}
