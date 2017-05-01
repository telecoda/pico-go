package console

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (c *console) initModes() map[ModeType]Mode {
	return map[ModeType]Mode{
		CLI:         newCLIMode(c),
		CODE_EDITOR: newCodeEditorMode(c),
		RUNTIME:     newRuntimeMode(c),
	}
}

func (m *mode) OffscreenBuffer() *sdl.Surface {
	return m.pixelSurface
}
