package console

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type cli struct {
	console *console
	PixelBuffer
	*cursor
	maxLineLen int
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
	cli.maxLineLen = (int(c.Config.ConsoleWidth) / _charWidth) - 2
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
		//fmt.Printf("Some event: %#v \n", event)
	}

	return nil
}

// getCmdCharIndex - returns relative position in command string that cursor currently is
func (c *cli) getCmdCharIndex() int {
	return int(c.cursor.x-c.cmdPos.x) + int(((c.cursor.y - c.cmdPos.y) * c.maxLineLen))
}

func (c *cli) cmdInsert(t string) {
	if len(c.cmd) >= _maxCmdLen {
		return
	}
	if c.cmd == "" {
		c.cmd = t
	} else {
		// should insert text relative to cursor
		charIndex := c.getCmdCharIndex()
		// check if cursor if at end of cmd string
		if int(charIndex) < len(c.cmd) {
			newCmd := c.cmd[0:charIndex]
			newCmd += t
			newCmd += c.cmd[charIndex:]
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
	charIndex := int(c.cursor.x-c.cmdPos.x) + int(((c.cursor.y - c.cmdPos.y) * c.maxLineLen))

	newCmd := c.cmd[0:charIndex]
	if int(charIndex) < len(c.cmd) {
		newCmd += c.cmd[charIndex+1:]
	}
	c.cmd = newCmd
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
	// if at beginning of command - can't move left
	if c.getCmdCharIndex() == 0 {
		return
	}
	c.cursor.x--
	if c.cursor.x < c.cmdPos.x {
		c.cursor.x = c.maxLineLen + c.cmdPos.x - 1
		c.cursor.y--
	}
	if c.cursor.y < c.cmdPos.y {
		c.cursor.y = c.cmdPos.y
	}
}

func (c *cli) cursorRight() {
	// if at end of command - can't move right
	if c.getCmdCharIndex() >= len(c.cmd) {
		return
	}
	c.cursor.x++
	if c.cursor.x >= c.maxLineLen+c.cmdPos.x {
		c.cursor.x = c.cmdPos.x
		c.cursor.y++
	}
}

func (c *cli) Init() error {
	// get native pixel buffer
	c.PixelBuffer.ClsWithColor(BLACK)
	pb := c.PixelBuffer.(*pixelBuffer)

	logoRect := &sdl.Rect{X: 0, Y: 0, W: _logoWidth, H: _logoHeight}
	screenRect := &sdl.Rect{X: 0, Y: 0, W: _logoWidth, H: _logoHeight}
	_console.logo.Blit(logoRect, pb.pixelSurface, screenRect)

	title := fmt.Sprintf("PICO-GO %s", _version)
	c.PixelBuffer.PrintAtWithColor(title, 0, 24, LIGHT_GRAY)

	c.PixelBuffer.Print("(C) 2017 @TELECODA")
	c.PixelBuffer.Print("TYPE HELP FOR HELP")

	c.initCursor()
	return nil
}

func (c *cli) initCursor() {
	currPos := c.GetCursor()
	// cmdPos is location of first char of command on screen
	// cursor is current location of cursor on screen
	c.cursor.pos = currPos
	c.cursor.pos.x = 2
	c.cmdPos = c.cursor.pos
}

func (c *cli) Update() error {
	return nil
}

func (c *cli) Render() error {
	// render text
	lines := c.getCmdLines()
	c.clearCmd(len(lines))
	c.renderCmd(lines)
	c.cursor.render()
	return nil
}

// getCmdLines - splits a command string into slices
func (c *cli) getCmdLines() []string {

	cmdLen := len(c.cmd)

	if cmdLen > c.maxLineLen {
		// split command across multiple lines
		wholeLines := (cmdLen / c.maxLineLen)
		totalLines := wholeLines
		rem := cmdLen % c.maxLineLen
		if rem > 0 {
			totalLines++
		}
		lines := make([]string, totalLines)
		pos := 0
		for i := 0; i < wholeLines; i++ {
			lines[i] = c.cmd[pos : pos+c.maxLineLen]
			pos += c.maxLineLen
		}
		// append remainder
		if rem > 0 {
			lines[len(lines)-1] = c.cmd[pos:]
		}
		return lines

	}
	return []string{c.cmd}

}

// clearCmd - clears screen where command will be rendered
func (c *cli) clearCmd(count int) {
	x0 := c.cmdPos.x * int(_charWidth)
	y0 := c.cmdPos.y * int(_charHeight)
	x1 := x0 + c.maxLineLen*int(_charWidth)
	y1 := y0 + (count+1)*int(_charHeight)
	c.RectFillWithColor(x0, y0, x1, y1, BLACK)
}

// renderCmd - renders command string across multiple lines
func (c *cli) renderCmd(lines []string) {
	c.PixelBuffer.Color(WHITE)

	// set print color
	currentPos := c.GetCursor()
	c.Color(WHITE)
	c.Cursor(0, currentPos.y)
	c.Print(">")
	c.Cursor(2, currentPos.y)
	for i := range lines {
		c.Print(lines[i])
	}
	c.Cursor(0, currentPos.y)
}
