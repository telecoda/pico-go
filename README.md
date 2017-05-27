# pico-go

Written by @telecoda

Pico-go is virtual console inspired by [pico8](https://www.lexaloffle.com/pico-8.php)

The plan is to build something similar to pico8 but where the games can be coded in Go instead of Lua.

<font color="red">**
Note: this project is very much work in progress at the moment, watch this space..**</font>


## Building

    go get github.com/telecoda/pico-go

First fetch all the dependencies:

    go get -u -v

You will also need to install the SDL2 dependencies for your platform.  See the [go-sdl README](https://github.com/telecoda/go-sdl2/blob/master/README.md) for details

    go run main.go

or

    pico-go run

More docs will be added as I progress..

# References

[PICO8 font](https://drive.google.com/file/d/0B97Um39fHXlcWUFRZlBqUndhbXM/view) by RhythmLynx
