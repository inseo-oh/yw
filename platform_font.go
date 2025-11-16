package main

// #cgo pkg-config: freetype2
// #include <ft2build.h>
// #include FT_FREETYPE_H
import "C"
import "image"

type ft_font struct {
	face C.FT_Face
}

func (fnt ft_font) SetTextSize(size int) {
	// STUB
}
func (fnt ft_font) DrawText(text string, dest *image.Image, offset_x, offset_y float64) (width, height float64) {
	// STUB
	return 10, 10
}
