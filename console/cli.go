package console

import (
	"fmt"

	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

type cli struct {
	console *console
	PixelBuffer
	*cursor
	maxLineLen int
	lastLine   int
	cmd        string
	cmdPos     pos
}

var event sdl.Event
var joysticks [16]*sdl.Joystick

func newCLIMode(c *console) Mode {
	cli := &cli{

		console: c,
	}
	pb, err := newPixelBuffer(c.Config)
	if err != nil {
		panic(err)
	}
	cursor := newCursor(pb, c.Config.cursorColor)

	cli.PixelBuffer = pb
	cli.cursor = cursor
	cursor.x = 2
	cursor.y = 8
	cli.cmdPos = cursor.pos
	// calc max line width
	cli.maxLineLen = (int(c.Config.ConsoleWidth) / _console.Config.fontWidth) - 2
	cli.lastLine = c.Config.ConsoleWidth / _console.Config.fontHeight

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
		case sdl.K_RETURN:
			c.cmdEnter()
		case sdl.K_q:
			if t.Keysym.Mod == sdl.KMOD_CTRL {
				c.console.Quit()
			}
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
	t = strings.ToUpper(t)
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

// cmdEnter - parse and execute command entered
func (c *cli) cmdEnter() {
	// redraw without cursor
	c.cursor.on = false
	lines := c.getCmdLines()
	c.clearCmd(len(lines))
	c.renderCmd(lines, false)
	err := c.cmdExec(c.cmd)
	if err != nil {
		c.Color(c.console.errColor)
		c.Print("")
		c.Print(err.Error())
	}
	c.initCmd()
}

func (c *cli) cmdExec(statement string) error {
	// parse command entered
	cmd, err := c.cmdParse(statement)
	if err != nil {
		return err
	}
	c.Print("")
	return cmd.Exec(c.PixelBuffer, statement)
}

func (c *cli) cmdParse(statement string) (Command, error) {

	// split statement into tokens
	tokens := strings.Split(statement, " ")
	if len(tokens) == 0 {
		return nil, nil
	}
	if cmd, ok := commands[tokens[0]]; ok {
		return cmd, nil
	}
	return nil, fmt.Errorf("SYNTAX ERROR")
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
		// if cursor off screen we need to scroll
		if c.cursor.y >= c.lastLine-1 {
			// we're gonna scroll
			lines := c.getCmdLines()
			c.clearCmd(len(lines) + 1)
			c.renderCmd(lines, false)

			c.ScrollUpLine()
			c.cursor.y--
			c.cmdPos.y--
		}
	}
}

func (c *cli) Init() error {
	// get native pixel buffer
	c.PixelBuffer.ClsWithColor(c.console.BgColor)
	pb := c.PixelBuffer.(*pixelBuffer)

	logoRect := &sdl.Rect{X: 0, Y: 0, W: _logoWidth, H: _logoHeight}
	screenRect := &sdl.Rect{X: 0, Y: 0, W: _logoWidth, H: _logoHeight}
	_console.logo.Blit(logoRect, pb.pixelSurface, screenRect)

	title := fmt.Sprintf("PICO-GO %s", _version)
	c.Cursor(0, 4)
	c.Color(c.console.FgColor)
	c.PixelBuffer.Print(title)

	c.PixelBuffer.Print("(C) 2017 @TELECODA")
	c.PixelBuffer.Print("TYPE HELP FOR HELP")

	c.initCmd()
	return nil
}

func (c *cli) initCmd() {
	currPos := c.GetCursor()
	// cmdPos is location of first char of command on screen
	// cursor is current location of cursor on screen
	c.cursor.pos = currPos
	c.cursor.pos.x = 2
	c.cmdPos = c.cursor.pos
	c.cmd = ""
}

func (c *cli) Update() error {
	return nil
}

func (c *cli) Render() error {
	// render text
	lines := c.getCmdLines()
	c.clearCmd(len(lines))
	c.renderCmd(lines, true)
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
	x0 := 0
	y0 := c.cmdPos.y * int(_console.Config.fontHeight)
	x1 := c.console.ConsoleWidth
	y1 := y0 + (count)*int(_console.Config.fontHeight)
	c.RectFillWithColor(x0, y0, x1, y1, c.console.BgColor)
}

// renderCmd - renders command string across multiple lines
func (c *cli) renderCmd(lines []string, resetCursor bool) {
	c.PixelBuffer.Color(c.console.FgColor)

	// set print color
	currentPos := c.cmdPos
	c.Color(c.console.FgColor)
	c.Cursor(0, currentPos.y)
	c.PrintAtWithColor(">", 0, currentPos.y*_console.Config.fontHeight, c.console.FgColor)
	c.Cursor(2, currentPos.y)
	for i := range lines {
		if currentPos.y < c.lastLine {
			c.PrintAtWithColor(lines[i], 2*_console.Config.fontWidth, (currentPos.y+i)*_console.Config.fontHeight, c.console.FgColor)
		} else {
			c.PrintAtWithColor(lines[i], 2*_console.Config.fontWidth, c.lastLine*_console.Config.fontHeight, c.console.FgColor)
		}
	}
	if resetCursor {
		c.Cursor(0, currentPos.y)
	}
}
