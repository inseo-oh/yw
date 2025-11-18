package main

// #cgo pkg-config: freetype2
// #include <ft2build.h>
// #include FT_FREETYPE_H
import "C"
import (
	"image"
	"image/color"
	"log"
	"unsafe"
	"yw/libgfx"
)

type ft_font struct {
	face C.FT_Face
}

func (fnt ft_font) SetTextSize(size int) {
	if res := C.FT_Set_Pixel_Sizes(fnt.face, 0, C.FT_UInt(size)); res != C.FT_Err_Ok {
		log.Printf("Failed to set font size (FT_Set_Pixel_Sizes error %d)", res)
	}
}
func (fnt ft_font) DrawText(text string, dest *image.RGBA, offset_x, offset_y float64, text_color color.RGBA) libgfx.Rect {
	// 26.6 Fixed Point -> float64
	ft_26_6_pos_to_float := func(p C.FT_Pos) float64 {
		return float64(p) / 64.0
	}

	pen_x, pen_y := offset_x, offset_y
	line_height := 0
	rect := libgfx.Rect{Left: pen_x, Top: pen_y, Width: 0, Height: 0}
	for _, char := range text {
		glyph_index := C.FT_Get_Char_Index(fnt.face, C.FT_ULong(char))
		if res := C.FT_Load_Glyph(fnt.face, glyph_index, C.FT_LOAD_DEFAULT); res != C.FT_Err_Ok {
			log.Printf("Failed to load glyph for char %c (FT_Load_Glyph error %d)", char, res)
		}
		if fnt.face.glyph.format != C.FT_GLYPH_FORMAT_BITMAP {
			if res := C.FT_Render_Glyph(fnt.face.glyph, C.FT_RENDER_MODE_NORMAL); res != C.FT_Err_Ok {
				log.Printf("Failed to render glyph for char %c (FT_Render_Glyph error %d)", char, res)
			}
		}
		gslot := fnt.face.glyph
		bitmap := gslot.bitmap
		bitmap_left := int(gslot.bitmap_left)
		bitmap_top := int(gslot.bitmap_top)
		rect.Left = float64(bitmap_left)
		rect.Top = -float64(bitmap_top)
		bytes := C.GoBytes(unsafe.Pointer(bitmap.buffer), C.int(bitmap.rows*bitmap.width))
		src_line_idx := 0
		dest_x := int(pen_x) + bitmap_left
		dest_y := int(pen_y) - bitmap_top
		dest_left := dest_x
		if dest != nil {
			for range bitmap.rows {
				src_idx := src_line_idx
				for range bitmap.width {
					val := bytes[src_idx]
					rgba := dest.RGBAAt(dest_x, dest_y)
					calc_channel := func(old, new, alpha uint8) uint8 {
						return uint8((int(old)*(255-int(alpha)))/255) + uint8((int(new)*int(alpha))/255)
					}
					rgba.R = calc_channel(rgba.R, text_color.R, val)
					rgba.G = calc_channel(rgba.G, text_color.G, val)
					rgba.B = calc_channel(rgba.B, text_color.B, val)
					rgba.A = 255 // Just make sure it's fully opaque
					dest.SetRGBA(dest_x, dest_y, rgba)
					src_idx++
					dest_x++
				}
				src_line_idx += int(bitmap.width)
				dest_y++
				dest_x = dest_left
			}
		}
		pen_x += ft_26_6_pos_to_float(gslot.advance.x)
		pen_y += ft_26_6_pos_to_float(gslot.advance.y)

		rect.Width += ft_26_6_pos_to_float(gslot.advance.x)
		line_height = max(line_height, int(bitmap.rows))
	}
	rect.Height = float64(line_height)
	return rect
}
