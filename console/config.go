package console

import (
	"flag"
)

type Config struct {
	ConsoleWidth    int
	ConsoleHeight   int
	WindowWidth     int
	WindowHeight    int
	FPS             int
	Verbose         bool
	ScreenshotScale int
	GifScale        int
	GifLength       int
}

var optVerbose bool
var screenshotScale int
var gifScale int
var gifLength int

func init() {
	flag.BoolVar(&optVerbose, "v", false, "verbose logging")
	flag.IntVar(&screenshotScale, "screenshot_scale", 3, "scale of screenshots")
	flag.IntVar(&gifScale, "gif_scale", 2, "scale of gif captures.")
	flag.IntVar(&gifLength, "gif_len", 10, "set the maximum gif length in seconds (1..120)")
}

func DefaultConfig() Config {
	flag.Parse()

	config := Config{
		ConsoleWidth:    128,
		ConsoleHeight:   128,
		WindowWidth:     400,
		WindowHeight:    400,
		FPS:             60,
		Verbose:         optVerbose,
		ScreenshotScale: screenshotScale,
		GifScale:        gifScale,
		GifLength:       gifLength,
	}

	return config
}
