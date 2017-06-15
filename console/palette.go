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
	colorMap       map[Color]rgba
	rgbaMap        map[uint32]Color
	colors         []sdl.Color
	originalColors []sdl.Color
}

func newPalette() *palette {

	p := &palette{}
	// set colours in palette
	p.colors = make([]sdl.Color, TOTAL_COLORS)
	p.originalColors = make([]sdl.Color, TOTAL_COLORS)
	p.originalColors[BLACK] = sdl.Color{R: 0, G: 0, B: 0, A: 255} // black
	p.originalColors[DARK_BLUE] = sdl.Color{29, 43, 83, 255}      // dark-blue
	p.originalColors[DARK_PURPLE] = sdl.Color{126, 37, 83, 255}   // dark-purple
	p.originalColors[DARK_GREEN] = sdl.Color{0, 135, 81, 255}     // dark-green
	p.originalColors[BROWN] = sdl.Color{171, 82, 54, 255}         // brown
	p.originalColors[DARK_GRAY] = sdl.Color{95, 87, 79, 255}      // dark-gray
	p.originalColors[LIGHT_GRAY] = sdl.Color{194, 195, 199, 255}  // light-gray
	p.originalColors[WHITE] = sdl.Color{255, 241, 232, 255}       // white
	p.originalColors[RED] = sdl.Color{255, 0, 77, 255}            // red
	p.originalColors[ORANGE] = sdl.Color{255, 163, 0, 255}        // orange
	p.originalColors[YELLOW] = sdl.Color{255, 236, 39, 255}       // yellow
	p.originalColors[GREEN] = sdl.Color{0, 228, 54, 255}          // green
	p.originalColors[BLUE] = sdl.Color{41, 173, 255, 255}         // blue
	p.originalColors[INDIGO] = sdl.Color{131, 118, 156, 255}      // indigo
	p.originalColors[PINK] = sdl.Color{255, 119, 168, 255}        // pink
	p.originalColors[PEACH] = sdl.Color{255, 204, 170, 255}       //  peach

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
	// default to black
	return BLACK
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
	p2 := newPalette()
	p.colorMap = p2.colorMap
	p.colors = p2.colors
}

func (p *palette) PaletteCopy() Paletter {
	p2 := newPalette()
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
