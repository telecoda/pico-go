package console

type BaseCartridge struct {
	cfg         Config // holds details of console config
	PixelBuffer        // ref to console display
	running     bool
}

// NewBaseCart - initialise a struct implementing Cartridge interface
func NewBaseCart() *BaseCartridge {
	cart := &BaseCartridge{
		cfg: DefaultConfig(),
	}

	return cart
}

// GetConfig - return config need for Cart to run
func (bc *BaseCartridge) GetConfig() Config {
	return bc.cfg
}

// Init - called once when cart is initialised
func (bc *BaseCartridge) initPb(pb PixelBuffer) {
	// the Init method receives a PixelBuffer reference
	// hold onto this reference, this is the display that
	// your code will be drawing onto each frame
	bc.PixelBuffer = pb
	bc.running = true
}

func (bc *BaseCartridge) IsRunning() bool {
	return bc.running
}

func (bc *BaseCartridge) Stop() {
	bc.running = false
}
