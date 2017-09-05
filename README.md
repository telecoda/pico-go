# pico-go

Written by @telecoda

Pico-go is virtual console inspired by [pico8](https://www.lexaloffle.com/pico-8.php)

<font color="red">**
NOTE:- this branch is an experiment to convert pico-go from using [go-sdl2](https://github.com/veandco/go-sdl2) to the [ebiten](https://github.com/hajimehoshi/ebiten)**

There were several reasons for evaluating ebiten as an alternative.

Installing the native SDL2 dependencies are non-trivial.  Especially on Windows where the gfx-sdl.dll is being used which requires explicit compilation.  These steps would be too much of a deterrent for people to use pico-go.

I switched to ebiten instead for the lower external dependencies, plus the advantage of being able to use GopherJS to generate game code that runs directly in the browser.

The goal of pico-go was to create a simple interactive development environment. As go is a compiled language as opposed to lua being interpreted this throws up a few more challenges.

I considered using ebiten wrapped in a browser app, similar to the go playground. As @jdb has done here [ebiten playground](https://j7b.github.io/playground/) This would allow editing in the browser and generating JS code to run immediately. Unfortunately the browser wrapper alternatives I discovered only appear to support OSX at present (alpha windows support).
https://github.com/alexflint/gallium
https://github.com/murlokswarm/app

Due to the retro nature of pico-go I have tried to reproduce the feel of pico8 as closely as possible. This required writing some simple graphic primatives and sprite rendering that update a Go image.RGBA struct.  The resulting pixels can be copied to an ebiten image for rendering using the `ReplacePixels()` method.  In native code this retains a 60FPS frame rate.  However these same image manipulations using the go std lib in generated javascript cause the performance to drop rapidly. The code is effectively iterating over arrays of pixel data and copying them from one image to another. (Could probably be optimised..)

The best solution would be to adopt the native image handling of ebiten itself for greater performance.  The downside is that ebiten does not support the drawing primatives I desire as it integrates closely with OpenGL.  Going down this path makes me feel that I am adding little value on top of the existing ebiten engine.

Therefore with a heavy heart I'm deciding to suspend development on pico-go for now.  The effort to continue developing a tool with a very limited (maybe non-existent) userbase is not worthwhile, but it was a fun ride.
</font>

