package console

import (
	"fmt"

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

type palette struct {
	colorMap map[Color]rgba
	colors   []sdl.Color
}

func NewPalette() *palette {

	p := &palette{}
	// set colours in palette
	p.colors = make([]sdl.Color, TOTAL_COLORS)
	p.colors[BLACK] = sdl.Color{R: 0, G: 0, B: 0, A: 255} // black
	p.colors[DARK_BLUE] = sdl.Color{29, 43, 83, 255}      // dark-blue
	p.colors[DARK_PURPLE] = sdl.Color{126, 37, 83, 255}   // dark-purple
	p.colors[DARK_GREEN] = sdl.Color{0, 135, 81, 255}     // dark-green
	p.colors[BROWN] = sdl.Color{171, 82, 54, 255}         // brown
	p.colors[DARK_GRAY] = sdl.Color{95, 87, 79, 255}      // dark-gray
	p.colors[LIGHT_GRAY] = sdl.Color{194, 195, 199, 255}  // light-gray
	p.colors[WHITE] = sdl.Color{255, 241, 232, 255}       // white
	p.colors[RED] = sdl.Color{255, 0, 77, 255}            // red
	p.colors[ORANGE] = sdl.Color{255, 163, 0, 255}        // orange
	p.colors[YELLOW] = sdl.Color{255, 236, 39, 255}       // yellow
	p.colors[GREEN] = sdl.Color{0, 228, 54, 255}          // green
	p.colors[BLUE] = sdl.Color{41, 173, 255, 255}         // blue
	p.colors[INDIGO] = sdl.Color{131, 118, 156, 255}      // indigo
	p.colors[PINK] = sdl.Color{255, 119, 168, 255}        // pink
	p.colors[PEACH] = sdl.Color{255, 204, 170, 255}       //  peach

	p.updateColorMap()

	return p
}

func (p *palette) updateColorMap() {
	// create a mpa of the colors
	p.colorMap = make(map[Color]rgba, len(p.colors))

	for i, c := range p.colors {
		p.colorMap[Color(i)] = rgba{R: c.R, G: c.G, B: c.B, A: c.A}
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
	rgbaCombined := uint32(c.R)<<24 | uint32(c.G)<<16 | uint32(c.B)<<8 | uint32(c.A)
	return c, rgbaCombined

}

// GetColorID - find color from rgba
func (p *palette) GetColorID(rgba rgba) Color {
	// lookup color using rgba
	for color, attrs := range p.colorMap {
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

	palette.SetColors(_console.palette.GetSDLColors())
	return surface.SetPalette(palette)
}

func (p *palette) PaletteReset() {
	p = NewPalette()
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
	p.colors[int(fromIdx)] = p.colors[int(toIdx)]
	p.updateColorMap()
	return nil
}
