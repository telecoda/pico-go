package display

import "github.com/telecoda/pico-go/api"

// getRGBA - returns color as api.Color and uint32
func (p palette) getRGBA(color api.Color) (rgba, uint32) {
	// lookup colour
	var c rgba
	var ok bool
	if c, ok = p[int(color)]; ok {
	} else {
		// if not found default to color 0
		c = p[0]
	}
	rgbaCombined := uint32(c.R)<<24 | uint32(c.G)<<16 | uint32(c.B)<<8 | uint32(c.A)
	return c, rgbaCombined

}
