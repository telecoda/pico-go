package main

import (
	"flag"

	"github.com/telecoda/pico-go/console"
	"github.com/telecoda/pico-go/examples/palette/code"
)

/*

	THIS IS GENERATED CODE - DO NOT AMEND

*/

func main() {

	flag.Parse()

	// Create virtual console - based on cart config
	con, err := console.NewConsole(console.PICO8)
	if err != nil {
		panic(err)
	}
	defer con.Destroy()

	cart := code.NewCart()

	if err := con.LoadCart(cart); err != nil {
		panic(err)
	}

	if err := con.Run(); err != nil {
		panic(err)
	}
}
