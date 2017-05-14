package console

import (
	"flag"
)

var optVerbose bool

func init() {
	flag.BoolVar(&optVerbose, "v", false, "verbose logging")
}

func DefaultConfig() Config {
	flag.Parse()

	config := Config{
		ConsoleWidth:  128,
		ConsoleHeight: 128,
		WindowWidth:   400,
		WindowHeight:  400,
		FPS:           60,
		Verbose:       optVerbose,
	}

	return config
}
