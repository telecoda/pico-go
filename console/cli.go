package console

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type cli struct {
	console *console
	PixelBuffer
	*cursor
	maxLineLen int32
	cmd        string
	cmdPos     pos
}

var event sdl.Event
var joysticks [16]*sdl.Joystick

func newCLIMode(c *console) Mode {
	cli := &cli{

		console: c,
	}
	pb, _ := newPixelBuffer(c.Config)
	cursor := newCursor(pb, RED)

	cli.PixelBuffer = pb
	cli.cursor = cursor
	cursor.x = 2
	cursor.y = 8
	cli.cmdPos = cursor.pos
	// calc max line width
	cli.maxLineLen = (int32(c.Config.ConsoleWidth) / _charWidth) - 2
	return cli
}

func (c *cli) HandleEvent(event sdl.Event) error {
	switch t := event.(type) {
	case *sdl.TextInputEvent:
		c.cmdInsert(string(t.Text[0]))
	case *sdl.TextEditingEvent:
		fmt.Printf("TEMP: text editing %s\n", t.Text)
	case *sdl.KeyDownEvent:
		switch t.Keysym.Sym {
		case sdl.K_DELETE:
			c.cmdDelete()
		case sdl.K_BACKSPACE:
			c.cmdBackspace()
		case sdl.K_LEFT:
			c.cursorLeft()
		case sdl.K_RIGHT:
			c.cursorRight()
			// fmt.Printf("fSwitching to code editor\n")
			// c.console.SetMode(CODE_EDITOR)
		}
	default:
		fmt.Printf("Some event: %#v \n", event)
	}

	return nil
}

func (c *cli) cmdInsert(t string) {
	if c.cmd == "" {
		c.cmd = t
	} else {
		// should insert text relative to cursor
		curpos := int(c.cursor.x-c.cmdPos.x) + int(((c.cursor.y - c.cmdPos.y) * c.maxLineLen))
		// check if cursor if at end of cmd string
		if int(curpos) < len(c.cmd) {
			newCmd := c.cmd[0:curpos]
			newCmd += t
			newCmd += c.cmd[curpos:]
			c.cmd = newCmd
		} else {
			// append to end
			c.cmd += t
		}
	}
	// then move cursor
	c.cursorRight()
}

func (c *cli) cmdDelete() {
	if c.cmd == "" {
		return
	}
	// should delete text relative to cursor
	curpos := int(c.cursor.x-c.cmdPos.x) + int(((c.cursor.y - c.cmdPos.y) * c.maxLineLen))

	newCmd := c.cmd[0:curpos]
	if int(curpos) < len(c.cmd) {
		newCmd += c.cmd[curpos+1:]
	}
	c.cmd = newCmd
	// don't move cursor
}

func (c *cli) cmdBackspace() {
	if c.cmd == "" {
		return
	}

	cmdLen := len(c.cmd)

	if cmdLen == 1 {
		c.cmd = ""
		c.cursorLeft()
		return
	}
	// should delete text relative to cursor
	curpos := int(c.cursor.x-c.cmdPos.x) + int(((c.cursor.y - c.cmdPos.y) * c.maxLineLen))

	newCmd := c.cmd[0 : curpos-1]
	if int(curpos) < len(c.cmd) {
		newCmd += c.cmd[curpos:]
	}
	c.cmd = newCmd
	// then move cursor
	c.cursorLeft()
}

func (c *cli) cursorLeft() {
	c.cursor.x--
	if int32(c.cursor.x) < c.cmdPos.x {
		c.cursor.x = c.maxLineLen + c.cmdPos.x - 1
		c.cursor.y--
	}
	if c.cursor.y < c.cmdPos.y {
		c.cursor.y = c.cmdPos.y
	}
}

func (c *cli) cursorRight() {
	c.cursor.x++
	if int32(c.cursor.x) >= c.maxLineLen+c.cmdPos.x {
		c.cursor.x = c.cmdPos.x
		c.cursor.y++
	}
}

func (c *cli) Init() error {
	// get native pixel buffer
	c.PixelBuffer.ClsColor(BLACK)
	pb := c.PixelBuffer.(*pixelBuffer)

	logoRect := &sdl.Rect{X: 0, Y: 0, W: _logoWidth, H: _logoHeight}
	screenRect := &sdl.Rect{X: 0, Y: 0, W: _logoWidth, H: _logoHeight}
	_console.logo.Blit(logoRect, pb.pixelSurface, screenRect)

	title := fmt.Sprintf("PICO-GO %s", _version)
	c.PixelBuffer.PrintColorAt(title, 0, 24, LIGHT_GRAY)
	c.PixelBuffer.PrintColorAt("(C) 2017 @TELECODA", 0, 32, LIGHT_GRAY)
	c.PixelBuffer.PrintColorAt("TYPE HELP FOR HELP", 0, 48, LIGHT_GRAY)

	c.PixelBuffer.PrintColorAt(">", 0, 64, WHITE)

	return nil
}

func (c *cli) Update() error {
	return nil
}

func (c *cli) Render() error {
	//c.Cls()
	// render text
	lines := c.getCmdLines()
	c.clearCmd(len(lines))
	c.renderCmd(lines)
	c.cursor.render()
	return nil
}

// getCmdLines - splits a command string into slices
func (c *cli) getCmdLines() []string {

	cmdLen := int32(len(c.cmd))

	if cmdLen > c.maxLineLen {
		// split command across multiple lines
		wholeLines := (cmdLen / c.maxLineLen)
		totalLines := wholeLines
		rem := cmdLen % c.maxLineLen
		if rem > 0 {
			totalLines++
		}
		lines := make([]string, totalLines)
		pos := int32(0)
		for i := int32(0); i < wholeLines; i++ {
			lines[i] = c.cmd[pos : pos+c.maxLineLen]
			pos += c.maxLineLen
		}
		// append remainder
		if rem > 0 {
			lines[len(lines)-1] = c.cmd[pos:]
		}
		return lines

	} else {
		return []string{c.cmd}
	}
}

// clearCmd - clears screen where command will be rendered
func (c *cli) clearCmd(count int) {
	x0 := int(c.cmdPos.x * _charWidth)
	y0 := int(c.cmdPos.y * _charHeight)
	x1 := x0 + int(c.maxLineLen*_charWidth)
	y1 := y0 + (count+1)*int(_charHeight)
	c.RectFillWithColor(x0, y0, x1, y1, BLACK)
}

// renderCmd - renders command string across multiple lines
func (c *cli) renderCmd(lines []string) {

	// render prompt
	c.PixelBuffer.PrintColorAt(">", 0, int(c.cmdPos.y*_charHeight), WHITE)

	for i := range lines {
		c.PrintColorAt(lines[i], int(c.cmdPos.x*_charWidth), int((c.cmdPos.y+int32(i))*_charHeight), WHITE)
	}

}
