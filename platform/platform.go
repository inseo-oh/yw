// Package platform provides abstract platform interface.
package platform

import "github.com/inseo-oh/yw/gfx"

// Platform is abstract interface used to isloate core web engine from
// external libs or platform dependent code.
type Platform interface {
	// OpenFont opens a font with given name.
	OpenFont(name string) gfx.Font
}
