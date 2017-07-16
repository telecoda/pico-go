package console

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// PICO8 - colors
const (
	PICO8_BLACK Color = iota
	PICO8_DARK_BLUE
	PICO8_DARK_PURPLE
	PICO8_DARK_GREEN
	PICO8_BROWN
	PICO8_DARK_GRAY
	PICO8_LIGHT_GRAY
	PICO8_WHITE
	PICO8_RED
	PICO8_ORANGE
	PICO8_YELLOW
	PICO8_GREEN
	PICO8_BLUE
	PICO8_INDIGO
	PICO8_PINK
	PICO8_PEACH
)

// TIC80 - colors
const (
	TIC80_BLACK Color = iota
	TIC80_DARK_RED
	TIC80_DARK_BLUE
	TIC80_DARK_GRAY
	TIC80_BROWN
	TIC80_DARK_GREEN
	TIC80_RED
	TIC80_LIGHT_GRAY
	TIC80_LIGHT_BLUE
	TIC80_ORANGE
	TIC80_BLUE_GRAY
	TIC80_LIGHT_GREEN
	TIC80_PEACH
	TIC80_CYAN
	TIC80_YELLOW
	TIC80_WHITE
)

// ZX Spectrum - colors
const (
	ZX_BLACK Color = iota
	ZX_BLUE
	ZX_RED
	ZX_MAGENTA
	ZX_GREEN
	ZX_CYAN
	ZX_YELLOW
	ZX_WHITE
	ZX_BRIGHT_BLACK
	ZX_BRIGHT_BLUE
	ZX_BRIGHT_RED
	ZX_BRIGHT_MAGENTA
	ZX_BRIGHT_GREEN
	ZX_BRIGHT_CYAN
	ZX_BRIGHT_YELLOW
	ZX_BRIGHT_WHITE
)

// Commodore 64 - colors
const (
	C64_BLACK Color = iota
	C64_WHITE
	C64_RED
	C64_CYAN
	C64_PURPLE
	C64_GREEN
	C64_BLUE
	C64_YELLOW
	C64_ORANGE
	C64_BROWN
	C64_LIGHT_RED
	C64_DARK_GREY
	C64_GREY
	C64_LIGHT_GREEN
	C64_LIGHT_BLUE
	C64_LIGHT_GREY
)

type rgba struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type palette struct {
	colorMap       map[Color]rgba
	rgbaMap        map[uint32]Color
	colors         []sdl.Color
	originalColors []sdl.Color
}

func newPalette(consoleType ConsoleType) *palette {
	switch consoleType {
	case PICO8:
		return newPico8Palette()
	case TIC80:
		return newTic80Palette()
	case ZX_SPECTRUM:
		return newZXSpectrumPalette()
	case CBM64:
		return newCBM64Palette()
	}
	return newPico8Palette() // always default to PICO8
}

func newPico8Palette() *palette {

	p := &palette{}
	// set colours in palette
	p.colors = make([]sdl.Color, TOTAL_COLORS)
	p.originalColors = make([]sdl.Color, TOTAL_COLORS)
	p.originalColors[PICO8_BLACK] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	p.originalColors[PICO8_DARK_BLUE] = sdl.Color{R: 29, G: 43, B: 83, A: 255}
	p.originalColors[PICO8_DARK_PURPLE] = sdl.Color{R: 126, G: 37, B: 83, A: 255}
	p.originalColors[PICO8_DARK_GREEN] = sdl.Color{R: 0, G: 135, B: 81, A: 255}
	p.originalColors[PICO8_BROWN] = sdl.Color{R: 171, G: 82, B: 54, A: 255}
	p.originalColors[PICO8_DARK_GRAY] = sdl.Color{R: 95, G: 87, B: 79, A: 255}
	p.originalColors[PICO8_LIGHT_GRAY] = sdl.Color{R: 194, G: 195, B: 199, A: 255}
	p.originalColors[PICO8_WHITE] = sdl.Color{R: 255, G: 241, B: 232, A: 255}
	p.originalColors[PICO8_RED] = sdl.Color{R: 255, G: 0, B: 77, A: 255}
	p.originalColors[PICO8_ORANGE] = sdl.Color{R: 255, G: 163, B: 0, A: 255}
	p.originalColors[PICO8_YELLOW] = sdl.Color{R: 255, G: 236, B: 39, A: 255}
	p.originalColors[PICO8_GREEN] = sdl.Color{R: 0, G: 228, B: 54, A: 255}
	p.originalColors[PICO8_BLUE] = sdl.Color{R: 41, G: 173, B: 255, A: 255}
	p.originalColors[PICO8_INDIGO] = sdl.Color{R: 131, G: 118, B: 156, A: 255}
	p.originalColors[PICO8_PINK] = sdl.Color{R: 255, G: 119, B: 168, A: 255}
	p.originalColors[PICO8_PEACH] = sdl.Color{R: 255, G: 204, B: 170, A: 255}

	// copy to working colors
	for i := range p.originalColors {
		p.colors[i] = p.originalColors[i]
	}

	p.updateColorMaps()

	return p
}

func newTic80Palette() *palette {

	p := &palette{}
	// set colours in palette
	p.colors = make([]sdl.Color, TOTAL_COLORS)
	p.originalColors = make([]sdl.Color, TOTAL_COLORS)
	p.originalColors[TIC80_BLACK] = sdl.Color{R: 20, G: 12, B: 28, A: 255}
	p.originalColors[TIC80_DARK_RED] = sdl.Color{R: 68, G: 36, B: 52, A: 255}
	p.originalColors[TIC80_DARK_BLUE] = sdl.Color{R: 48, G: 52, B: 109, A: 255}
	p.originalColors[TIC80_DARK_GRAY] = sdl.Color{R: 78, G: 74, B: 78, A: 255}
	p.originalColors[TIC80_BROWN] = sdl.Color{R: 133, G: 76, B: 48, A: 255}
	p.originalColors[TIC80_DARK_GREEN] = sdl.Color{R: 52, G: 101, B: 36, A: 255}
	p.originalColors[TIC80_RED] = sdl.Color{R: 208, G: 70, B: 72, A: 255}
	p.originalColors[TIC80_LIGHT_GRAY] = sdl.Color{R: 117, G: 113, B: 97, A: 255}
	p.originalColors[TIC80_LIGHT_BLUE] = sdl.Color{R: 89, G: 125, B: 206, A: 255}
	p.originalColors[TIC80_ORANGE] = sdl.Color{R: 210, G: 125, B: 44, A: 255}
	p.originalColors[TIC80_BLUE_GRAY] = sdl.Color{R: 133, G: 149, B: 161, A: 255}
	p.originalColors[TIC80_LIGHT_GREEN] = sdl.Color{R: 109, G: 170, B: 44, A: 255}
	p.originalColors[TIC80_PEACH] = sdl.Color{R: 210, G: 170, B: 153, A: 255}
	p.originalColors[TIC80_CYAN] = sdl.Color{R: 109, G: 194, B: 202, A: 255}
	p.originalColors[TIC80_YELLOW] = sdl.Color{R: 218, G: 212, B: 94, A: 255}
	p.originalColors[TIC80_WHITE] = sdl.Color{R: 222, G: 238, B: 214, A: 255}

	// copy to working colors
	for i := range p.originalColors {
		p.colors[i] = p.originalColors[i]
	}

	p.updateColorMaps()

	return p
}

func newZXSpectrumPalette() *palette {

	p := &palette{}
	// set colours in palette
	p.colors = make([]sdl.Color, TOTAL_COLORS)
	p.originalColors = make([]sdl.Color, TOTAL_COLORS)
	p.originalColors[ZX_BLACK] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	p.originalColors[ZX_BLUE] = sdl.Color{R: 0, G: 41, B: 197, A: 255}
	p.originalColors[ZX_RED] = sdl.Color{R: 213, G: 39, B: 30, A: 255}
	p.originalColors[ZX_MAGENTA] = sdl.Color{R: 211, G: 58, B: 199, A: 255}
	p.originalColors[ZX_GREEN] = sdl.Color{R: 0, G: 197, B: 49, A: 255}
	p.originalColors[ZX_CYAN] = sdl.Color{R: 0, G: 200, B: 201, A: 255}
	p.originalColors[ZX_YELLOW] = sdl.Color{R: 205, G: 200, B: 59, A: 255}
	p.originalColors[ZX_WHITE] = sdl.Color{R: 203, G: 203, B: 203, A: 255}
	p.originalColors[ZX_BRIGHT_BLACK] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	p.originalColors[ZX_BRIGHT_BLUE] = sdl.Color{R: 0, G: 54, B: 247, A: 255}
	p.originalColors[ZX_BRIGHT_RED] = sdl.Color{R: 255, G: 52, B: 40, A: 255}
	p.originalColors[ZX_BRIGHT_MAGENTA] = sdl.Color{R: 255, G: 75, B: 250, A: 255}
	p.originalColors[ZX_BRIGHT_GREEN] = sdl.Color{R: 0, G: 247, B: 63, A: 255}
	p.originalColors[ZX_BRIGHT_CYAN] = sdl.Color{R: 0, G: 252, B: 253, A: 255}
	p.originalColors[ZX_BRIGHT_YELLOW] = sdl.Color{R: 255, G: 251, B: 76, A: 255}
	p.originalColors[ZX_BRIGHT_WHITE] = sdl.Color{R: 255, G: 255, B: 255, A: 255}

	// copy to working colors
	for i := range p.originalColors {
		p.colors[i] = p.originalColors[i]
	}

	p.updateColorMaps()

	return p
}

func newCBM64Palette() *palette {

	p := &palette{}
	// set colours in palette
	p.colors = make([]sdl.Color, TOTAL_COLORS)
	p.originalColors = make([]sdl.Color, TOTAL_COLORS)
	p.originalColors[C64_BLACK] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	p.originalColors[C64_WHITE] = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	p.originalColors[C64_RED] = sdl.Color{R: 160, G: 78, B: 70, A: 255}
	p.originalColors[C64_CYAN] = sdl.Color{R: 110, G: 193, B: 199, A: 255}
	p.originalColors[C64_PURPLE] = sdl.Color{R: 161, G: 89, B: 163, A: 255}
	p.originalColors[C64_GREEN] = sdl.Color{R: 95, G: 171, B: 98, A: 255}
	p.originalColors[C64_BLUE] = sdl.Color{R: 80, G: 71, B: 154, A: 255}
	p.originalColors[C64_YELLOW] = sdl.Color{R: 202, G: 212, B: 140, A: 255}
	p.originalColors[C64_ORANGE] = sdl.Color{R: 161, G: 104, B: 63, A: 255}
	p.originalColors[C64_BROWN] = sdl.Color{R: 109, G: 83, B: 21, A: 255}
	p.originalColors[C64_LIGHT_RED] = sdl.Color{R: 202, G: 127, B: 120, A: 255}
	p.originalColors[C64_DARK_GREY] = sdl.Color{R: 99, G: 99, B: 99, A: 255}
	p.originalColors[C64_GREY] = sdl.Color{R: 139, G: 139, B: 139, A: 255}
	p.originalColors[C64_LIGHT_GREEN] = sdl.Color{R: 157, G: 226, B: 160, A: 255}
	p.originalColors[C64_LIGHT_BLUE] = sdl.Color{R: 138, G: 129, B: 202, A: 255}
	p.originalColors[C64_LIGHT_GREY] = sdl.Color{R: 174, G: 174, B: 174, A: 255}

	// copy to working colors
	for i := range p.originalColors {
		p.colors[i] = p.originalColors[i]
	}

	p.updateColorMaps()

	return p
}

func (r rgba) toIndex() uint32 {
	return uint32(uint32(r.R)<<24 | uint32(r.G)<<16 | uint32(r.B)<<8 | uint32(r.A))
}

func (p *palette) updateColorMaps() {
	// create a mpa of the colors
	p.colorMap = make(map[Color]rgba, len(p.colors))
	p.rgbaMap = make(map[uint32]Color, len(p.colors))
	for i, c := range p.colors {
		color := rgba{R: c.R, G: c.G, B: c.B, A: c.A}
		p.colorMap[Color(i)] = color
		p.rgbaMap[color.toIndex()] = Color(i)
	}

}

// getRGBA - returns color as Color and uint32
func (p *palette) GetRGBA(color Color) (rgba, uint32) {
	// lookup color
	var c rgba
	var ok bool
	if c, ok = p.colorMap[color]; ok {
	} else {
		// if not found default to color 0
		c = p.colorMap[0]
	}
	rgbaCombined := c.toIndex()
	return c, rgbaCombined

}

// GetColorID - find color from rgba
func (p *palette) GetColorID(rgba rgba) Color {
	// lookup color using rgba
	if colorID, ok := p.rgbaMap[rgba.toIndex()]; ok {
		return colorID
	}
	// default to color 0
	return 0
}

func setSurfacePalette(palette Paletter, surface *sdl.Surface) error {
	p, err := sdl.AllocPalette(TOTAL_COLORS)
	if err != nil {
		return err
	}
	p.SetColors(palette.GetSDLColors())
	return surface.SetPalette(p)
}

func (p *palette) PaletteReset() {
	p2 := _console.originalPalette
	p.colorMap = p2.colorMap
	p.colors = p2.colors
}

func (p *palette) PaletteCopy() Paletter {
	p2 := _console.originalPalette
	for i, c := range p.colors {
		p2.colors[i] = c
	}
	p2.updateColorMaps()
	return p2
}

func (p *palette) GetSDLColors() []sdl.Color {
	return p.colors
}

func (p *palette) MapColor(fromColor Color, toColor Color) error {
	// valid request
	if fromColor < 0 || int(fromColor) > len(p.colors) {
		return fmt.Errorf("Error mapping color - fromColour outside range: %d", fromColor)
	}
	if toColor < 0 || int(toColor) > len(p.colors) {
		return fmt.Errorf("Error mapping color - toColour outside range: %d", toColor)
	}

	// update color
	fromIdx := int(fromColor)
	toIdx := int(toColor)
	p.colors[int(fromIdx)] = p.originalColors[int(toIdx)]
	p.updateColorMaps()
	return nil
}

func (p *palette) SetTransparent(color Color, enabled bool) error {

	fromIdx := int(color)

	if enabled {
		if err := p.MapColor(color, 0); err != nil {
			return err
		}
	} else {
		p.colors[int(fromIdx)] = p.originalColors[int(fromIdx)]
	}

	p.updateColorMaps()
	return nil
}
