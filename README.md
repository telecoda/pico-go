# pico-go

Written by @telecoda

Pico-go is virtual console inspired by [pico8](https://www.lexaloffle.com/pico-8.php)

The plan is to build something similar to pico8 but where the games can be coded in Go instead of Lua.

<font color="red">**
Note: this project is very much work in progress at the moment, watch this space..**</font>

See [TODO](TODO.md) for current implementation status.


## Building

    go get github.com/telecoda/pico-go

First fetch all the dependencies:

    go get -u -v

You will also need to install the SDL2 dependencies for your platform.  See the [go-sdl README](https://github.com/veandco/go-sdl2/blob/master/README.md) for details

    go install

More docs will be added as I progress..

## Dependencies

This project is dependent on the [github.com/veandco/go-sdl2](https://github.com/veandco/go-sdl2) project for SDL support.

TODO: document installation options thoroughly

## Getting started

This section describes how to get started using pico-go

### Creating a new project
To create a new project

    pico-go new <project-name>

Run this command in the directory where you would like your new project created.  The command will create a new directory matching your project name containing all the required asset.

![image](/docs/images/new-project.png)


### Running the project

Once you have created your pico-go project you can run the code immediately.

    cd <project-name>
    pico-go run

The generated project contains a simple code example that renders animated text on the screen.

### Hot reloading
The `pico-go run` command monitors the filesystem for file changes.  If you amend any files in the code directory your project will be restarted. (Only if the code compiles successfully)

You can also run the project by directly calling the `go run main.go` command but the code will not be automatically restarted.

# References

[PICO8 font](https://drive.google.com/file/d/0B97Um39fHXlcWUFRZlBqUndhbXM/view) by RhythmLynx
