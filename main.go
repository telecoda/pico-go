package main

import (
	"flag"

	"github.com/telecoda/pico-go/console"
	"github.com/telecoda/pico-go/demo"
)

var optVerbose bool

func init() {
	flag.BoolVar(&optVerbose, "v", false, "verbose logging")
}

func main() {

	flag.Parse()

	config := console.Config{
		ConsoleWidth:  128,
		ConsoleHeight: 128,
		WindowWidth:   400,
		WindowHeight:  400,
		FPS:           60,
		Verbose:       optVerbose,
	}

	// Create virtual console
	con, err := console.NewConsole(config)
	if err != nil {
		panic(err)
	}

	defer con.Destroy()

	cart := demo.NewCart()

	if err := con.LoadCart(cart); err != nil {
		panic(err)
	}

	if err := con.Run(); err != nil {
		panic(err)
	}
}
