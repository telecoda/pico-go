package console

import (
	"image"
)

type Mode interface {
	Init() error
	Render() error
	Update() error
	HandleEvent(event string) error
	PixelBuffer
}

func (c *console) initModes() (map[ModeType]Mode, error) {
	modes := map[ModeType]Mode{
		CLI:           newCLIMode(c),
		CODE_EDITOR:   newCodeEditorMode(c),
		SPRITE_EDITOR: newSpriteEditorMode(c),
		//		RUNTIME:     newRuntimeMode(c),
	}

	for _, mode := range modes {
		if err := mode.Init(); err != nil {
			return nil, err
		}
	}

	return modes, nil
}

func (m *mode) OffscreenBuffer() *image.RGBA {
	return m.pixelSurface
}
