package console

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	BLACK Color = iota
	DARK_BLUE
	DARK_PURPLE
	DARK_GREEN
	BROWN
	DARK_GRAY
	LIGHT_GRAY
	WHITE
	RED
	ORANGE
	YELLOW
	GREEN
	BLUE
	INDIGO
	PINK
	PEACH
)

type rgba struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type palette map[Color]rgba

var colors []sdl.Color

func initPico8Palette() palette {
	// set colours in palette
	colors = make([]sdl.Color, TOTAL_COLORS)
	colors[BLACK] = sdl.Color{R: 0, G: 0, B: 0, A: 255} // black
	colors[DARK_BLUE] = sdl.Color{29, 43, 83, 255}      // dark-blue
	colors[DARK_PURPLE] = sdl.Color{126, 37, 83, 255}   // dark-purple
	colors[DARK_GREEN] = sdl.Color{0, 135, 81, 255}     // dark-green
	colors[BROWN] = sdl.Color{171, 82, 54, 255}         // brown
	colors[DARK_GRAY] = sdl.Color{95, 87, 79, 255}      // dark-gray
	colors[LIGHT_GRAY] = sdl.Color{194, 195, 199, 255}  // light-gray
	colors[WHITE] = sdl.Color{255, 241, 232, 255}       // white
	colors[RED] = sdl.Color{255, 0, 77, 255}            // red
	colors[ORANGE] = sdl.Color{255, 163, 0, 255}        // orange
	colors[YELLOW] = sdl.Color{255, 236, 39, 255}       // yellow
	colors[GREEN] = sdl.Color{0, 228, 54, 255}          // green
	colors[BLUE] = sdl.Color{41, 173, 255, 255}         // blue
	colors[INDIGO] = sdl.Color{131, 118, 156, 255}      // indigo
	colors[PINK] = sdl.Color{255, 119, 168, 255}        // pink
	colors[PEACH] = sdl.Color{255, 204, 170, 255}       //  peach

	return map[Color]rgba{
		BLACK: rgba{0, 0, 0, 255}, // black
		//DARK_BLUE:   rgba{29, 43, 83, 255},    // dark-blue
		DARK_BLUE:   rgba{255, 0, 0, 255},     // dark-blue
		DARK_PURPLE: rgba{126, 37, 83, 255},   // dark-purple
		DARK_GREEN:  rgba{0, 135, 81, 255},    // dark-green
		BROWN:       rgba{171, 82, 54, 255},   // brown
		DARK_GRAY:   rgba{95, 87, 79, 255},    // dark-gray
		LIGHT_GRAY:  rgba{194, 195, 199, 255}, // light-gray
		WHITE:       rgba{255, 241, 232, 255}, // white
		RED:         rgba{255, 0, 77, 255},    // red
		ORANGE:      rgba{255, 163, 0, 255},   // orange
		YELLOW:      rgba{255, 236, 39, 255},  // yellow
		GREEN:       rgba{0, 228, 54, 255},    // green
		BLUE:        rgba{41, 173, 255, 255},  // blue
		INDIGO:      rgba{131, 118, 156, 255}, // indigo
		PINK:        rgba{255, 119, 168, 255}, // pink
		PEACH:       rgba{255, 204, 170, 255}, //  peach
	}
}

// getRGBA - returns color as Color and uint32
func (p palette) getRGBA(color Color) (rgba, uint32) {
	// lookup color
	var c rgba
	var ok bool
	if c, ok = p[color]; ok {
	} else {
		// if not found default to color 0
		c = p[0]
	}
	rgbaCombined := uint32(c.R)<<24 | uint32(c.G)<<16 | uint32(c.B)<<8 | uint32(c.A)
	return c, rgbaCombined

}

// getColor - find color from rgba
func (p palette) getColor(rgba rgba) Color {
	// lookup color using rgba
	for color, attrs := range p {
		if rgba.R == attrs.R && rgba.G == attrs.G && rgba.B == attrs.B && rgba.A == attrs.A {
			return color
		}
	}
	// default to black
	return BLACK
}

func setSurfacePalette(surface *sdl.Surface) error {
	palette, err := sdl.AllocPalette(TOTAL_COLORS)
	if err != nil {
		return err
	}

	palette.SetColors(colors)
	return surface.SetPalette(palette)
}
