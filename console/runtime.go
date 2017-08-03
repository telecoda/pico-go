package console

type runtime struct {
	console *console
	running bool
	PixelBuffer
	Cartridge
}

func newRuntimeMode(c *console) Runtime {
	runtime := &runtime{
		console: c,
	}
	pb, _ := newPixelBuffer(c.Config)

	runtime.PixelBuffer = pb

	return runtime
}

func (r *runtime) HandleEvent(event string) error {
	// TODO
	//	return _console.storeKeyState(event)
	return nil
}

func (r *runtime) LoadCart(cart Cartridge) error {
	r.Cartridge = cart
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
