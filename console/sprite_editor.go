package console

import (
	"github.com/veandco/go-sdl2/sdl"
)

type spriteEditor struct {
	console *console
	PixelBuffer
}

func newSpriteEditorMode(c *console) Mode {
	spriteEditor := &spriteEditor{
		console: c,
	}
	pb, _ := newPixelBuffer(c.Config)

	spriteEditor.PixelBuffer = pb

	return spriteEditor
}

func (s *spriteEditor) HandleEvent(event sdl.Event) error {
	switch t := event.(type) {
	case *sdl.KeyDownEvent:
		switch t.Keysym.Sym {
		// case sdl.K_RIGHT:
		// 	fmt.Printf("Switching to runtime\n")
		// 	c.console.SetMode(RUNTIME)
		}
	default:
		//fmt.Printf("Some event: %#v \n", event)
	}

	return nil
}

func (s *spriteEditor) Init() error {
	return nil
}

func (s *spriteEditor) Update() error {
	return nil
}

func (s *spriteEditor) Render() error {
	s.PixelBuffer.ClsWithColor(5)
	s.PixelBuffer.PrintAtWithColor("Code editor Print Test", 10, 10, 7)
	renderModeHeader(s)
	return nil
}
