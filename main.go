package main

import (
	"flag"

	"github.com/telecoda/pico-go/config"
	"github.com/telecoda/pico-go/console"
)

var optVerbose bool

func init() {
	flag.BoolVar(&optVerbose, "v", false, "verbose logging")
}

func main() {

	flag.Parse()

	config := config.Config{
		ConsoleWidth:  128,
		ConsoleHeight: 128,
		WindowWidth:   800,
		WindowHeight:  800,
		FPS:           60,
		Verbose:       optVerbose,
	}

	// Create virtual console
	c, err := console.New(config)
	if err != nil {
		panic(err)
	}

	defer c.Destroy()

	if err := c.Run(); err != nil {
		panic(err)
	}
}