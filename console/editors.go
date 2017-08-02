package console

/* this file contains common editor related functionality */

func renderModeHeader(pb PixelBuffer) {
	// render menu bar
	pb.RectFillWithColor(0, 0, _console.ConsoleWidth, 8, PICO8_RED)
	x := _console.ConsoleWidth - 8*len(_console.modes)
	for k, _ := range _console.modes {
		index := int(k)
		if _console.currentMode == k {
			// current mode
			pb.systemSprite(index, x+index*8, -1, 1, 1, 8, 8, 0, false, false)
		} else {
			// other modes
			pb.systemSprite(index+6, x+index*8, -1, 1, 1, 8, 8, 0, false, false)
		}
	}
}
