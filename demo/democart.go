package demo

import (
	"fmt"

	"github.com/telecoda/pico-go/console"
)

/*
This is the code that the pico-go developer creates
They will not need to touch anything outside of this package.
Lets try and keeo them as isolated from the rest of the
code as possible
*/

// Code must implement Cartridge interface

type democart struct {
	pb console.PixelBuffer

	counter int
}

func NewCart() console.Cartridge {
	return &democart{}
}

func (d *democart) Init(pb console.PixelBuffer) {
	d.pb = pb
	d.pb.ClsWithColor(console.BLUE)
	d.pb.PrintAtWithColor("Runtime Print Test", 10, 10, console.WHITE)

	d.counter = 0
}

func (d *democart) Update() {
	d.counter++
}

func (d *democart) Render() {
	d.pb.Print(fmt.Sprintf("c:%d", d.counter))

}
