package console

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Mode interface {
	Init() error
	Render() error
	Update() error
	HandleEvent(event sdl.Event) error
	PixelBuffer
}

func (c *console) initModes() (map[ModeType]Mode, error) {
	modes := map[ModeType]Mode{
		CLI:         newCLIMode(c),
		CODE_EDITOR: newCodeEditorMode(c),
		//		RUNTIME:     newRuntimeMode(c),
	}

	for _, mode := range modes {
		if err := mode.Init(); err != nil {
			return nil, err
		}
	}

	return modes, nil
}

func (m *mode) OffscreenBuffer() *sdl.Surface {
	return m.pixelSurface
}
