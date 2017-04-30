package console

import "github.com/telecoda/pico-go/api"

func (c *console) Cls() {
	c.Display.Cls()
}

func (c *console) ClsColor(colorId api.Color) {
	c.Display.ClsColor(colorId)
}
