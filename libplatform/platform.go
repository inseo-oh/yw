package libplatform

import "yw/libgfx"

type Platform interface {
	//==========================================================================
	// Fonts
	//==========================================================================
	OpenFont(name string) libgfx.Font
}
