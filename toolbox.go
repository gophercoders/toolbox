// Copyright 2016 Kulawe Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package toolbox is a package that simplifies the handling graphics, windows and
// key presses offered by the underlying SDL library.
// The package is not intended as a generic game library or as a simplification
// of the SDL library. The package is intended solely for use in the Code Club
// to support very simple games.
//
// Warning: This package is not safe in the presence of multiple go routines.
// This is by design.  This package contains an internal, unguarded,
// unexported global variables.
// Given the use cases of this package, this is an acceptable trade off.
package toolbox

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

// Graphic is the type used for an image or graphic.
type Graphic *sdl.Texture

// Window is the type of the window that the game takes place within.
type Window *sdl.Window

// Key is the type that represents the key or window button that has been pressed.
type Key int

// Colour is the tye that represens a solout, held as a Read, Green Blue, Alpha tuple.
// Note: the English spelling.
type Colour sdl.Color

// Possible constants returned from GetKey()
const (
	KeyNone     = iota // no key was pressed
	KeyUp              // the up cursor key was pressed
	KeyDown            // the down cursor key was pressed
	KeyPause           // the Pause key was pressed
	ButtonClose        // the windows close button was clicked on
)

const (
	notInitialisedMessage       = "The toolbox has not been initialised. Have you called toolbox.Initialise()?"
	nilGraphicMessage           = "The Graphic is invalid. Have you loaded the Graphic using toolbox.LoadGraphic(filename)?"
	windowNotInitialisedMessage = "The window has not been created. Have you called toolbox.CreateWindow(...)?"
	noRenderer                  = "No Renderer. Have you called toolbox.CreateWindow(...)?"
	couldNotDrawPoint           = "The point could not be plotted. Something very bad has happened."
)

// This the abstraction to the graphics hardware inside the computer
// that actually does the drawing
var renderer *sdl.Renderer

// This is the intialised flag. It is true only if the Initialise function has
// been called.
var initialised bool

// Initialise prepares the toolbox for use. It must be called before any other
// functions in the tool box are used.
func Initialise() {
	sdl.Init(sdl.INIT_EVERYTHING)
	initialised = true
}

// Close closes the toolbox, releasing any underlying resources. It must
// be called before the program exits to prevent a resource leak. The toolbox
// cannot be used again until the Initialise function is called again.
func Close() {
	sdl.Quit()
	initialised = true
}

// LoadGraphic loads a graphic from disk or panics trying. The image can be in
// BMP, GIF, JPEG, LBM, PCX, PNG, PNM, TGA, TIFF, WEBP, XCF, XPM, or XV format.
// The user must supply the filename of the graphic file to load, in the
// parameter filename.
//
// If the function succeeds then a variable of type Graphic will be returned
// back to the calling function. It is the programmers responsibility to
// store this in a variable of type Graphic.
//
// If the function fails it will panic and crash the program.
// The reasons for a panic are:
//
// 1. The toolbox has not been initalised
//
// 2. The filename does not exist, or is otherwise inaccessable. The specific
// reason will be contained in the panic message itself. This message will
// be prefixed with "Failed to load file: ".
//
// 3. The file could not be converted into a Graphic type. Again the specific
// reason will be contained in the panic message itself. This message will be
// prefixed with "Failed to create Graphic: "
func LoadGraphic(filename string) Graphic {
	if !initialised {
		// this stops execution here, so ne need for an else after the if
		panic(notInitialisedMessage)
	}

	var err error
	var image *sdl.Surface
	image, err = img.Load(filename)
	if err != nil {
		fmt.Print("Failed to load file: ")
		fmt.Println(err)
		panic(err)
	}
	defer image.Free()
	var graphic *sdl.Texture
	graphic, err = renderer.CreateTextureFromSurface(image)
	if err != nil {
		fmt.Print("Failed to create Graphic: ")
		fmt.Println(err)
		panic(err)
	}
	return Graphic(graphic)
}

// DestroyGraphic destroys a Graphic that has been previously loaded by
// LoadGraphic. It must be called once for each graphic that has been loaded
// before the program exists.
//
// DestroyGraphic will panic if:
//
// 1. The Graphic, g, does not contain a Graphic type.
//
// 2. The toolbox has not been initialised.
func DestroyGraphic(g Graphic) {
	if !initialised {
		// this stops execution here, so ne need for an else after the if
		panic(notInitialisedMessage)
	}
	if g == nil {
		panic(nilGraphicMessage)
	}
	var t *sdl.Texture
	t = g
	t.Destroy()
}

// GetSizeOfGraphic returns the width and heigth of the Graphic, g.
// GetSizeOfGraphic returns two numbers of type int. The first numnber,
// marked as 'width', is the width of the Graphic, g, in pixels. The second number,
// marked as 'height' is the height of the Graphic, g, in pixels.
//
// GetSizeOfGraphic will panic if:
//
// 1. The Graphic, g, does not contain a Graphic type.
//
// 2. The toolbox has not been initialised.
func GetSizeOfGraphic(g Graphic) (width, height int) {
	if !initialised {
		// this stops execution here, so ne need for an else after the if
		panic(notInitialisedMessage)
	}
	if g == nil {
		panic(nilGraphicMessage)
	}
	var w, h int32
	var err error
	var t *sdl.Texture
	t = g
	_, _, w, h, err = t.Query()
	if err != nil {
		fmt.Print("Failed to query texture: ")
		fmt.Println(err)
		panic(err)
	}
	width = int(w)
	height = int(h)
	// return is implicit - using named return paramaters
	return
}

// RenderGraphic draws the Graphic, g, on the screen at postion (x,y).
// 'width' and 'height' specifiy the width and height of the Graphic, g.
//
// RenderGraphic will panic if:
// 1. The Graphic, g, does not contain a Graphic type.
//
// 2. The toolbox has not been initialised.
//
// 3. Any one of x, y, width or height are negative
func RenderGraphic(g Graphic, x, y, width, height int) {
	if !initialised {
		// this stops execution here, so ne need for an else after the if
		panic(notInitialisedMessage)
	}
	if g == nil {
		panic(nilGraphicMessage)
	}
	if x < 0 || y < 0 || width < 0 || height < 0 {
		panic("One of x, y, width or height is negative.")
	}
	var src, dst sdl.Rect

	src.X = 0
	src.Y = 0
	src.W = int32(width)
	src.H = int32(height)

	dst.X = int32(x)
	dst.Y = int32(y)
	dst.W = int32(width)
	dst.H = int32(height)

	renderer.Copy(g, &src, &dst)
}

// GetKey returns the key or button that has been pressed by the user.
// Possible return values are the constants:
//
//  KeyNone - indicating that no key, or an unrecognised key as been pressed
//  KeyUp	- indicating that the up cursor key has been pressed
//  KeyDown - indicating that the down cursor key has been pressed
//  KeyPause - indicating that the pause key has been presses
//  ButtonClose - indicating that the windows close button has been pressed
//
// GetKey will panic if:
//
// 1. The toolbox has not been initialised.
//
// 2. If an internal check fails. In this case the panic message is "KeyDownEvent type assertion failed!"
// This is highly unlikely to occur and indcates a problem with the underlying
// graphics llibrary.
func GetKey() Key {
	if !initialised {
		// this stops execution here, so ne need for an else after the if
		panic(notInitialisedMessage)
	}

	var event sdl.Event
	event = sdl.PollEvent()
	if event != nil {
		if isQuitEvent(event) {
			return ButtonClose
		}
		if isKeyDownEvent(event) {
			if isKeyUp(event) {
				return KeyUp
			}
			if isKeyDown(event) {
				return KeyDown
			}
			// We must always respond to the paused key being pressed - if the
			// game is not over.
			// If the game is running the pause key pauses the game.
			// But if the game is paused, we must still respond to the paused key.
			// This is the only way to unpause the game.
			if isKeyPause(event) {
				return KeyPause
			}
		}
	}
	return KeyNone
}

// CreateWindow creates a window with the specified title, width and height.
// The resulting window will be created with a title bar, a title, a close button and is moveable.
// The window cannot be resized.
//
// Create window is designed to create a single window. You cannot use
// CreateWindow more than once without first calling DestroyWindow.
//
// If the function succeeds then a variable of type Window will be returned
// back to the calling function. It is the programmers responsibility to
// store this in a variable of type Window.
//
// CreateWindow will panic if:
//
// 1. The toolbox has not been initialised.
//
// 2. Either of the width or height are negative
//
// 3. CreateWindow has been called more than once.
func CreateWindow(title string, width, height int) Window {
	if !initialised {
		// this stops execution here, so no need for an else after the if
		panic(notInitialisedMessage)
	}
	if width < 0 || height < 0 {
		panic("Requested window width or height is negative.")
	}
	if renderer != nil {
		// this stops execution here, so ne need for an else after the if
		panic("CreateWindow() has already been called. Did you call DestroyWindow(...)?")
	}

	var w *sdl.Window
	w = createWindow(width, height, title)
	return Window(w)
}

// DestroyWindow closes the window, freeing any resources as it does so.
// It must be called for each window that has been created before the program
// exists.
//
// DestroyWindow will panic if:
//
// 1. The toolbox has not been initialised.
//
// 2. The window is invalid
//
// 3. CreateWindow has not been called.
func DestroyWindow(window Window) {
	if !initialised {
		// this stops execution here, so ne need for an else after the if
		panic(notInitialisedMessage)
	}
	if renderer == nil {
		// this stops execution here, so ne need for an else after the if
		panic(windowNotInitialisedMessage)
	}
	if window == nil {
		panic("The Window is invalid. Have you used toolbox.CreateWindow(...)?")
	}
	var w *sdl.Window
	w = window
	if renderer != nil { // should always be true - see createWindow - but be defensive
		renderer.Destroy()
	}
	w.Destroy()
}

// ShowWindow redraws the window, displaying any changes that have been made.
//
// ShowWindow panics if:
//
// 1. The toolbox has not been initialised.
//
// 2. CreateWindow has not been called.
func ShowWindow() {
	if !initialised {
		// this stops execution here, so ne need for an else after the if
		panic(notInitialisedMessage)
	}
	if renderer != nil {
		renderer.Present()
	} else {
		panic(noRenderer)
	}
}

// SetBackgroundColour sets the background colour of the window.
// This is the colour that will be used to fill the window when ClearBackground()
// is called.
// The effect will not be seen until ShowWindow is called.
//
// The colour is specified as 3 integers. One for red, one for green
// and one for blue. The range of each of these numbers is 0..255,
// including 0 and 255.
//
// SetBackgroundColour panics if:
//
// 1. The toolbox has not been initialised.
//
// 2. CreateWindow has not been called.
//
// 3. If any of r, g, or b lies outside the range.
func SetBackgroundColour(r, g, b int) {
	if !initialised {
		// this stops execution here, so ne need for an else after the if
		panic(notInitialisedMessage)
	}
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		panic("One of r, g, or b is less than zero or greater than 255.")
	}
	if renderer != nil {
		renderer.SetDrawColor(0, 0, 0, 0)
	} else {
		panic(noRenderer)
	}
}

// SetDrawColour sets the colour that will be used when a point is plotted.
// The effect will not be seen until after DrawPoint and ShowWindow
// have been called.
//
// The colour is specified as a Colour type.
//
// SetDrawColour panics if:
//
// 1. The toolbox has not been initialised.
//
// 2. CreateWindow has not been called.
//
// 3. If any of red, green, blue, or alpha components of the colour
// are outside of the range [0..255].
func SetDrawColour(c Colour) {
	if !initialised {
		// this stops execution here, so ne need for an else after the if
		panic(notInitialisedMessage)
	}
	if c.R < 0 || c.R > 255 || c.G < 0 || c.G > 255 || c.B < 0 || c.B > 255 || c.A < 0 || c.A > 255 {
		panic("One of the r, g, b or alpha values in the colour is less than zero or greater than 255.")
	}
	if renderer != nil {
		renderer.SetDrawColor(c.R, c.G, c.B, c.A)
	} else {
		panic(noRenderer)
	}
}

// ClearBackground clears the window using the background colour set with
// SetBackgroundColour
// The effect will not be seen until ShowWindow is called.
//
// ClearBackgroundColour panics if:
//
// 1. The toolbox has not been initialised.
//
// 2. CreateWindow has not been called.
//
func ClearBackground() {
	if !initialised {
		// this stops execution here, so ne need for an else after the if
		panic(notInitialisedMessage)
	}
	if renderer != nil {
		renderer.Clear()
	} else {
		panic(noRenderer)
	}
}

// DrawPoint plots a single point (pixel) on the screen.
// The colour is set via SetDrawColour.
// The effect will not be seen until after ShowWindow has been called.
//
// The colour is specified as a Colour type.
//
// DrawPoint panics if:
//
// 1. The toolbox has not been initialised.
//
// 2. CreateWindow has not been called.
//
// 3. Plotting the point itself fails. This would indicate in internal invarient failure.
func DrawPoint(x, y int) {
	if !initialised {
		// this stops execution here, so ne need for an else after the if
		panic(notInitialisedMessage)
	}
	if renderer != nil {
		var err error
		err = renderer.DrawPoint(x, y)
		if err != nil {
			panic(couldNotDrawPoint)
		}
	} else {
		panic(noRenderer)
	}
}

// GetTickCount returns the amount of time that has passed since the toolbox
// was initialised.
//
// The toolbox counts time in 'ticks'. One 'tick' is 1/1000th of a second.
// Time starts when Initialise() is called. The number for ticks always
// increases. The number of ticks cannot be reset, and time does cannot
// run backwards.
//
// If the function succeeds then a variable of type int64 will be returned
// back to the calling function. It is the programmers responsibility to
// store this in a variable of type int64.
// GetTickCount returns an int64 type, not an int because the maximum number
// of ticks is to large to store in an int.
//
// GetTickCount will panic if:
//
// 1. The toolbox has not been initialised.
func GetTickCount() int64 {
	if !initialised {
		// this stops execution here, so ne need for an else after the if
		panic(notInitialisedMessage)
	}

	var ticks uint32
	ticks = sdl.GetTicks()
	return int64(ticks)
}

// Pause pauses execution of the game for the specified number of ticks.
// During the pause period nothing will happen. All input will be ignored
// and the window will not redraw.
//
// The duration of the pause must be specified in ticks in a variable
// of type int64
//
// Pause will panic if:
//
// 1. The toolbox has not been initialised.
//
// 2. The numnber of ticks is negative.
func Pause(numberOfTicks int64) {
	if !initialised {
		// this stops execution here, so ne need for an else after the if
		panic(notInitialisedMessage)
	}
	if numberOfTicks < 0 {
		panic("Cannot pause. The numberOfTicks is negative.")
	}
	sdl.Delay(uint32(numberOfTicks))
}

// Create the graphics window using the SDl library or crash trying
func createWindow(w, h int, title string) *sdl.Window {
	var window *sdl.Window
	var err error

	window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	// we're all good to create the renderer....
	renderer, err = createRenderer(window)
	if err != nil {
		// Oops, we failed to create the renderer failed
		// cleanup the window before we panic
		window.Destroy()
		panic(err)
	}
	// at this point we have both a good renderer and a good window
	return window
}

// Create the graphics renderer
func createRenderer(w *sdl.Window) (*sdl.Renderer, error) {
	var r *sdl.Renderer
	var err error
	r, err = sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED)
	return r, err
}

func isQuitEvent(event sdl.Event) bool {
	var ok bool
	_, ok = event.(*sdl.QuitEvent)
	return ok
}

func isKeyDownEvent(event sdl.Event) bool {
	var ok bool
	_, ok = event.(*sdl.KeyDownEvent)
	return ok
}

func isKeyUp(event sdl.Event) bool {
	var keyDownEvt *sdl.KeyDownEvent
	var ok bool
	keyDownEvt, ok = event.(*sdl.KeyDownEvent)
	if !ok {
		panic("KeyDownEvent type assertion failed!")
	}
	return (keyDownEvt.Keysym.Sym == sdl.K_UP)
}

func isKeyDown(event sdl.Event) bool {
	var keyDownEvt *sdl.KeyDownEvent
	var ok bool
	keyDownEvt, ok = event.(*sdl.KeyDownEvent)
	if !ok {
		panic("KeyDownEvent type assertion failed!")
	}
	return (keyDownEvt.Keysym.Sym == sdl.K_DOWN)
}

func isKeyPause(event sdl.Event) bool {
	var keyDownEvt *sdl.KeyDownEvent
	var ok bool
	keyDownEvt, ok = event.(*sdl.KeyDownEvent)
	if !ok {
		panic("KeyDownEvent type assertion failed!")
	}
	return (keyDownEvt.Keysym.Sym == sdl.K_PAUSE)
}
