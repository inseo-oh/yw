package main

// #cgo pkg-config: freetype2
// #include <ft2build.h>
// #include FT_FREETYPE_H
import "C"
import (
	"log"
	"unsafe"
	"yw/libhtml/platform"
)

type platform_impl struct {
	ft_lib C.FT_Library
}

func init_platform() *platform_impl {
	var ft_lib C.FT_Library
	if res := C.FT_Init_FreeType(&ft_lib); res != C.FT_Err_Ok {
		log.Fatalf("Failed to initialize FreeType (FT Error %d)", res)
	}

	return &platform_impl{
		ft_lib: ft_lib,
	}
}
func (plat platform_impl) OpenFont(name string) platform.Font {
	var face C.FT_Face
	font_name := C.CString("res/font.ttf")
	if res := C.FT_New_Face(plat.ft_lib, font_name, 0, &face); res == C.FT_Err_Unknown_File_Format {
		log.Fatalf("Unrecognized font (FT Error %d)", res)
	} else if res != C.FT_Err_Ok {
		log.Fatalf("Failed to open font %s (FT Error %d)", name, res)
	}
	C.free(unsafe.Pointer(font_name))
	return ft_font{face: face}
}
