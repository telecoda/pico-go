# pico-go API

So pico-po is "inspired" by the pico8 virtual console.  Therefore I have tried to create an API as close as possible to the original pico8 API.

When developing in pico8 you use `Lua` whereas pico-go uses `Go`. The Lua language supports optional/default parameters but Go does not.  Therefore instead of creating overloaded methods there are separate methods with different signatures.

Watch this space as I may amend this in the future to use a different style (I'm not entirely happy with too many params for a function).

##API Comparison
This is a side by side comparison of pico8 functions vs the equivalent pico-go functions


### Drawing API

These are the drawing functions.

| pico8 | pico-go | description |
| ------------ | ------------- | ------------ |
| cls() | Cls()  | Clear screen with current color |
| cls(col) | ClsWithColor(colorID Color)  | Clear screen with specifc color |
| color(col) | Color(colorID Color)  | Set the default drawing color |
| circ(x,y,r) | Circle(x, y, r int)| Draw circle with default color |
| circ(x,y,r,col) | CircleWithColor(x, y, r int, colorID Color)| Draw circle with specific color |
| circfill(x,y,r) | CircleFill(x, y, r int)| Draw filled circle with default color |
| circfill(x,y,r,col) | CircleFillWithColor(x, y, r int, colorID Color)| Draw filled circle with specific color |
| line(x0,y0,x1,y1) | Line(x0, y0, x1, y1 int)| Draw line with default color |
| line(x0,y0,x1,y1,col) | LineWithColor(x0, y0, x1, y1 int, colorID Color)| Draw line with specific color |
| pget(x,y) | PGet(x, y int) Color| Get color of a pixel |
| pset(x,y) | 	PSet(x, y int) | Set color of a pixel with default color|
| pset(x,y,col) | 	PSetWithColor(x, y int, colorID Color) | Set color of a pixel with specifc color|
| rect(x0,y0,x1,y1) | 	Rect(x0, y0, x1, y1 int) | Draw rectangle with default color|
| rect(x0,y0,x1,y1,col) | 	RectWithColor(x0, y0, x1, y1 int, colorID Color) | Draw rectangle with specifc color|
| rectfill(x0,y0,x1,y1) | 	RectFill(x0, y0, x1, y1 int) | Draw rectangle with default color|
| rectfill(x0,y0,x1,y1,col) | 	RectFillWithColor(x0, y0, x1, y1 int, colorID Color) | Draw rectangle with specifc color|

### Palette API

These are the palette functions.

| pico8 | pico-go | description |
| ------------ | ------------- | ------------ |
| pal() | PaletteReset()  | Reset palette colors to original settings |
| pal(c0,c1) | MapColor(fromColor Color, toColor Color) error  | Map from one color to another |
| palt(c0,true) | SetTransparent(color Color, enabled bool) error  | Enable/Disable transparency on a color |


### Printing API

These are the text printing functions.

| pico8 | pico-go | description |
| ------------ | ------------- | ------------ |
| print(str) | Print(str string)  | Print text at current cursor position |
| print(str,x,y) | PrintAt(str string, x, y int) | Print text at x,y |
| print(str,x,y,col) | PrintAtWithColor(str string, x, y int, colorID Color) | Print text at x,y with a specific color |
| cursor(x,y) | Cursor(x, y int)  | Set cursor position |


### Sprite API

These are the sprite rendering functions.

| pico8 | pico-go | description |
| ------------ | ------------- | ------------ |
| spr(n,x,y,w,h,flip_x,flip_y,| Sprite(n, x, y, w, h, dw, dh int, rot float64, flipX, flipY bool)  | Render sprite |

TODO:- Sprite function contains waaaayyy too many params and should be simplifed.

	