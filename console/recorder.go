package console

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type recorder struct {
	frames      []*sdl.Surface
	current     int
	fps         int
	seconds     int
	totalFrames int
}

type Recorder interface {
	AddFrame(surface *sdl.Surface)
	SaveVideo(filename string, scale int) error
	SaveScreenshot(filename string, scale int) error
}

func NewRecorder(fps, secondsToRecord int) Recorder {
	r := &recorder{
		fps:         fps,
		seconds:     secondsToRecord,
		totalFrames: fps * secondsToRecord,
	}

	r.frames = make([]*sdl.Surface, r.totalFrames, r.totalFrames)

	return r
}

func (r *recorder) AddFrame(surface *sdl.Surface) {
	r.frames[r.current] = surface

	r.current++
	if r.current > r.totalFrames-1 {
		r.current = 0
	}
}

func (r *recorder) SaveVideo(filename string, scale int) error {
	var palette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0x00, 0x00, 0xff, 0xff},
		color.RGBA{0x00, 0xff, 0x00, 0xff},
		color.RGBA{0x00, 0xff, 0xff, 0xff},
		color.RGBA{0xff, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0x00, 0xff, 0xff},
		color.RGBA{0xff, 0xff, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
	}

	// count used frames
	totalFrames := 0
	for i := 0; i < len(r.frames); i++ {
		if r.frames[i] != nil {
			totalFrames++
		}
	}

	images := make([]*image.Paletted, totalFrames, totalFrames)
	delays := make([]int, totalFrames, totalFrames)

	delay := 1000 / r.fps

	for i := 0; i < totalFrames; i++ {
		delays[i] = delay
		frame := r.frames[i]
		img := image.NewPaletted(image.Rect(0, 0, int(frame.W), int(frame.H)), palette)
		images[i] = img

		// drawing code (too long)
	}

	f, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	return gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}

func (r *recorder) SaveScreenshot(filename string, scale int) error {

	// screenshot uses last frame
	lastFrame := r.current - 1
	if lastFrame < 0 {
		lastFrame = r.totalFrames - 1
	}

	sourceSurface := r.frames[lastFrame]

	srcRect := &sdl.Rect{X: 0, Y: 0, W: sourceSurface.W, H: sourceSurface.H}
	targetRect := &sdl.Rect{X: 0, Y: 0, W: sourceSurface.W * int32(scale), H: sourceSurface.H * int32(scale)}

	targetSurface, err := sdl.CreateRGBSurface(0, sourceSurface.W*int32(scale), sourceSurface.H*int32(scale), 4, 0, 0, 0, 0)
	if err != nil {
		return err
	}

	err = sourceSurface.BlitScaled(srcRect, targetSurface, targetRect)
	if err != nil {
		return err
	}

	imageRect := image.Rect(0, 0, int(targetRect.W), int(targetRect.H))

	rgbaImage := image.NewRGBA(imageRect)

	pixels := targetSurface.Pixels()

	w := int(sourceSurface.W) * scale
	// convert SDL surface to RGBA image
	// process 4 bytes at a time
	for i := 0; i < len(pixels); i += 4 {

		// convert index i to x,y coords
		x := (i % (w * 4)) / 4
		y := (i - (x * 4)) / (w * 4)
		a := pixels[i]
		b := pixels[i+1]
		g := pixels[i+2]
		r := pixels[i+3]
		c := color.RGBA{R: r, G: g, B: b, A: a}
		rgbaImage.Set(x, y, c)
	}

	// save to out.gif
	f, _ := os.OpenFile("out.png", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	fmt.Printf("Saving screenshot to %s\n", f.Name())

	return png.Encode(f, rgbaImage)

}
