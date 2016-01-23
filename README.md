## Toolbox  - a simple, minimal wrapper over the go SDL bindings

Package Toolbox is a package that simplifies the handling graphics, windows and
key presses offered by the underlying SDL library.

The package is not intended as a generic game library or as a simplification
of the SDL library. The package is intended solely for use in the Code Club
to support very simple games.

_Warning:_ This package is not safe in the presence of multiple go routines.
This is by design.  This package contains an internal, unguarded,
 unexported global variables.
 Given the use cases of this package, this is an acceptable trade off.

### Dependencies

Toolbox relies on the `go-sdl2` package. See [here](https://github.com/veandco/go-sdl2)
for the installation instructions.

Please see the [COPYRIGHT](https://github.com/gophercoders/codeclub/blob/master/COPYRIGHT)
file for copyright information.

Please see the [LICENSE](https://github.com/gophercoders/codeclub/blob/master/LICENSE)
file for license information.
