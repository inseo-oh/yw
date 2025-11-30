package platform

import "yw/gfx"

type Platform interface {
	//==========================================================================
	// Fonts
	//==========================================================================
	OpenFont(name string) gfx.Font
}
