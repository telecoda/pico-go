package console

import (
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

type runtime struct {
	sync.Mutex
	console *console
	running bool
	PixelBuffer
	PicoInputAPI
	Cartridge
	keydownMap map[int]bool
}

const (
	P1_BUTT_RIGHT = 0x4f
	P1_BUTT_LEFT  = 0x50
	P1_BUTT_DOWN  = 0x51
	P1_BUTT_UP    = 0x52
	//Z,X / C,V / N,M
	P1_BUTT_01 = 0x1d // Z
	P1_BUTT_02 = 0x1b // X
	P1_BUTT_03 = 0x6  // C
	P1_BUTT_04 = 0x19 // V
	P1_BUTT_05 = 0x11 // N
	P1_BUTT_06 = 0x10 // M
)

func newRuntimeMode(c *console) Runtime {
	runtime := &runtime{
		console: c,
	}
	pb, _ := newPixelBuffer(c.Config)

	runtime.Lock()
	defer runtime.Unlock()

	// init key map
	runtime.keydownMap = make(map[int]bool)

	runtime.keydownMap[P1_BUTT_RIGHT] = false
	runtime.keydownMap[P1_BUTT_LEFT] = false
	runtime.keydownMap[P1_BUTT_UP] = false
	runtime.keydownMap[P1_BUTT_DOWN] = false
	runtime.keydownMap[P1_BUTT_01] = false
	runtime.keydownMap[P1_BUTT_02] = false
	runtime.keydownMap[P1_BUTT_03] = false
	runtime.keydownMap[P1_BUTT_04] = false
	runtime.keydownMap[P1_BUTT_05] = false
	runtime.keydownMap[P1_BUTT_06] = false

	runtime.PixelBuffer = pb

	return runtime
}

func (r *runtime) LoadCart(cart Cartridge) error {
	r.Cartridge = cart
	return nil
}

func (r *runtime) HandleEvent(event sdl.Event) error {

	switch event.(type) {
	case *sdl.KeyDownEvent:
		// update keydown map
		keyEvent := event.(*sdl.KeyDownEvent)
		r.Lock()
		defer r.Unlock()
		r.keydownMap[int(keyEvent.Keysym.Scancode)] = true
	case *sdl.KeyUpEvent:
		// update keydown map
		keyEvent := event.(*sdl.KeyUpEvent)
		r.Lock()
		defer r.Unlock()
		r.keydownMap[int(keyEvent.Keysym.Scancode)] = true
	}
	return nil
}

func (r *runtime) Init() error {
	r.Cartridge.initPb(r.PixelBuffer)
	r.Cartridge.Init()
	return nil
}

func (r *runtime) Update() error {
	r.Cartridge.Update()
	return nil
}

func (r *runtime) Render() error {
	r.Cartridge.Render()
	return nil
}
