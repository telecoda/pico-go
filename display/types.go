package display

import (
	"github.com/telecoda/pico-go/api"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_gfx"
)

// The display package handles all the output to screen

var title = "pico-go virtual games console"

type Display interface {
	Render() error
	Destroy()

	Cls()
	ClsColor(color api.Color)
}

type size struct {
	width  uint
	height uint
}

type rgba struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type palette map[int]rgba

type display struct {
	config       api.Config
	window       *sdl.Window
	pixelSurface *sdl.Surface
	renderer     *sdl.Renderer
	fpsMan       *gfx.FPSmanager
	palette
}

func initPico8Palette() palette {
	return map[int]rgba{
		0:  rgba{0, 0, 0, 255},       // black
		1:  rgba{29, 43, 83, 255},    // dark-blue
		2:  rgba{126, 37, 83, 255},   // dark-purple
		3:  rgba{0, 135, 81, 255},    // dark-green
		4:  rgba{171, 82, 54, 255},   // brown
		5:  rgba{95, 87, 79, 255},    // dark-gray
		6:  rgba{194, 195, 199, 255}, // light-gray
		7:  rgba{255, 241, 232, 255}, // white
		8:  rgba{255, 0, 77, 255},    // red
		9:  rgba{255, 163, 0, 255},   // orange
		10: rgba{255, 236, 39, 255},  // yellow
		11: rgba{0, 228, 54, 255},    // green
		12: rgba{41, 173, 255, 255},  // blue
		13: rgba{131, 118, 156, 255}, // indigo
		14: rgba{255, 119, 168, 255}, // pink
		15: rgba{255, 204, 170, 255}, //  peach
	}
}
