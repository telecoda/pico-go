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
	SaveVideo(filename string, scale int, palette Paletter) error
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
	// save copy of current frame

	newSurface, _ := surface.Convert(surface.Format, 0)

	r.frames[r.current] = newSurface

	r.current++
	if r.current > r.totalFrames-1 {
		r.current = 0
	}
}

func (r *recorder) SaveVideo(filename string, scale int, palette Paletter) error {

	// count used frames
	totalFrames := 0
	for i := 0; i < len(r.frames); i++ {
		if r.frames[i] != nil {
			totalFrames++
		}
	}

	if totalFrames == 0 {
		return fmt.Errorf("No frames to save to video")
	}

	sampledFrames := totalFrames / 3

	colors := make([]color.Color, len(palette.GetSDLColors()))
	for i, c := range palette.GetSDLColors() {
		colors[i] = color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	}

	images := make([]*image.Paletted, sampledFrames, sampledFrames)
	delays := make([]int, sampledFrames, sampledFrames)

	//Delay Time (1/100ths of a second)

	// Sampling at 60 fps
	// - one frame every 16.666 milliseconds
	// - one frame every 1.6666 1/100th of a second
	// all good... except delays can only be an int in 1/100ths of a second..
	// so we either delay to 1/100th and be too quick or 2/100th and be too slow

	// Convert to a 20 fps GIF?
	// Select every 3rd frame
	// Delay = 5/100th :)

	// TODO - handle different frame rates..

	delay := 5

	srcRect := &sdl.Rect{X: 0, Y: 0, W: r.frames[0].W, H: r.frames[0].H}
	targetRect := &sdl.Rect{X: 0, Y: 0, W: r.frames[0].W * int32(scale), H: r.frames[0].H * int32(scale)}

	target32Surface, err := sdl.CreateRGBSurface(0, srcRect.W, srcRect.H, 32, rmask, gmask, bmask, amask)
	if err != nil {
		return err
	}

	target32ScaledSurface, err := sdl.CreateRGBSurface(0, targetRect.W, targetRect.H, 32, rmask, gmask, bmask, amask)
	if err != nil {
		return err
	}

	for i := 0; i < sampledFrames; i++ {
		delays[i] = delay
		frame := r.frames[i*3]

		// scale original frame to video size
		err = frame.Blit(srcRect, target32Surface, srcRect)
		if err != nil {
			return err
		}

		err = target32Surface.BlitScaled(srcRect, target32ScaledSurface, targetRect)
		if err != nil {
			return err
		}

		img := image.NewPaletted(image.Rect(0, 0, int(target32ScaledSurface.W), int(target32ScaledSurface.H)), colors)
		// copy framePixels to image pixels
		pixels := target32ScaledSurface.Pixels()
		w := int(target32ScaledSurface.W)
		for i := 0; i < len(pixels); i += 4 {

			// convert index i to x,y coords
			x := (i % (w * 4)) / 4
			y := (i - (x * 4)) / (w * 4)
			a := pixels[i]
			b := pixels[i+1]
			g := pixels[i+2]
			r := pixels[i+3]

			// lookup color index
			pixelColor := rgba{R: r, G: g, B: b, A: a}
			colorID := palette.GetColorID(pixelColor)

			img.SetColorIndex(x, y, uint8(colorID))
		}

		images[i] = img

	}

	f, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	fmt.Printf("Saving video to %s\n", f.Name())
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
	target32Surface, err := sdl.CreateRGBSurface(0, sourceSurface.W, sourceSurface.H, 32, rmask, gmask, bmask, amask)
	if err != nil {
		return err
	}
	defer target32Surface.Free()

	// convert 8 bit palette image to 32 bit RGBA
	err = sourceSurface.Blit(srcRect, target32Surface, srcRect)
	if err != nil {
		return err
	}

	// convert 32 bit RGBA to scaled 32 bit RGBA
	target32ScaledRect := &sdl.Rect{X: 0, Y: 0, W: sourceSurface.W * int32(scale), H: sourceSurface.H * int32(scale)}
	target32ScaledSurface, err := sdl.CreateRGBSurface(0, sourceSurface.W*int32(scale), sourceSurface.H*int32(scale), 32, rmask, gmask, bmask, amask)
	if err != nil {
		return err
	}

	err = target32Surface.BlitScaled(srcRect, target32ScaledSurface, target32ScaledRect)
	if err != nil {
		return err
	}

	imageRect := image.Rect(0, 0, int(target32ScaledRect.W), int(target32ScaledRect.H))

	rgbaImage := image.NewRGBA(imageRect)

	pixels := target32ScaledSurface.Pixels()

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
