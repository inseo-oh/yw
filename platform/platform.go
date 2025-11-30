package platform

import "github.com/inseo-oh/yw/gfx"

type Platform interface {
	//==========================================================================
	// Fonts
	//==========================================================================
	OpenFont(name string) gfx.Font
}
