package console

func NewCartridge() (Cartridge, error) {

	c := &cartridge{}
	return c, nil
}

func (c *cartridge) Init() {

}

func (c *cartridge) Render() {

}

func (c *cartridge) Update() {

}
