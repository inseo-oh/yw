package main

// #cgo pkg-config: freetype2
// #include <ft2build.h>
// #include FT_FREETYPE_H
import "C"
import (
	"log"
	"unsafe"

	"github.com/inseo-oh/yw/gfx"
)

type platformImpl struct {
	ftLib C.FT_Library
}

func initPlatform() *platformImpl {
	var ftLib C.FT_Library
	if res := C.FT_Init_FreeType(&ftLib); res != C.FT_Err_Ok {
		log.Fatalf("Failed to initialize FreeType (FT Error %d)", res)
	}

	return &platformImpl{
		ftLib: ftLib,
	}
}
func (plat platformImpl) OpenFont(name string) gfx.Font {
	var face C.FT_Face
	fontName := C.CString("res/font/static/NotoSansKR-Regular.ttf")
	if res := C.FT_New_Face(plat.ftLib, fontName, 0, &face); res == C.FT_Err_Unknown_File_Format {
		log.Fatalf("Unrecognized font (FT Error %d)", res)
	} else if res != C.FT_Err_Ok {
		log.Fatalf("Failed to open font %s (FT Error %d)", name, res)
	}
	C.free(unsafe.Pointer(fontName))
	return ftFont{face: face}
}
