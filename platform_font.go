// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

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

	"github.com/inseo-oh/yw/gfx"
)

type ftFont struct {
	face C.FT_Face
}

func (fnt ftFont) SetTextSize(size int) {
	if res := C.FT_Set_Pixel_Sizes(fnt.face, 0, C.FT_UInt(size)); res != C.FT_Err_Ok {
		log.Printf("Failed to set font size (FT_Set_Pixel_Sizes error %d)", res)
	}
}
func (fnt ftFont) Metrics() gfx.FontMetrics {
	rawMetrics := fnt.face.size.metrics
	return gfx.FontMetrics{
		Ascender:   float64(rawMetrics.ascender) / 64.0,
		Descender:  float64(rawMetrics.descender) / 64.0,
		LineHeight: float64(rawMetrics.height) / 64.0,
	}
}
func (fnt ftFont) DrawText(text string, dest *image.RGBA, offsetX, offsetY float64, textColor color.RGBA) gfx.Rect {
	// 26.6 Fixed Point -> float64
	ft26p6PosToFloat := func(p C.FT_Pos) float64 {
		return float64(p) / 64.0
	}

	penX, penY := offsetX, offsetY
	lineHeight := 0
	rect := gfx.Rect{Left: penX, Top: penY, Width: 0, Height: 0}
	for _, char := range text {
		glyphIndex := C.FT_Get_Char_Index(fnt.face, C.FT_ULong(char))
		if res := C.FT_Load_Glyph(fnt.face, glyphIndex, C.FT_LOAD_DEFAULT); res != C.FT_Err_Ok {
			log.Printf("Failed to load glyph for char %c (FT_Load_Glyph error %d)", char, res)
		}
		if fnt.face.glyph.format != C.FT_GLYPH_FORMAT_BITMAP {
			if res := C.FT_Render_Glyph(fnt.face.glyph, C.FT_RENDER_MODE_NORMAL); res != C.FT_Err_Ok {
				log.Printf("Failed to render glyph for char %c (FT_Render_Glyph error %d)", char, res)
			}
		}
		gslot := fnt.face.glyph
		bitmap := gslot.bitmap
		bitmapLeft := int(gslot.bitmap_left)
		bitmapTop := int(gslot.bitmap_top)
		bytes := C.GoBytes(unsafe.Pointer(bitmap.buffer), C.int(bitmap.rows*bitmap.width))
		srcLineIdx := 0
		rect.Top = min(rect.Top, penY-float64(bitmapTop))
		destX := int(penX) + bitmapLeft
		destY := int(penY) - bitmapTop
		destLeft := destX
		if dest != nil {
			for range bitmap.rows {
				srcIdx := srcLineIdx
				for range bitmap.width {
					val := bytes[srcIdx]
					rgba := dest.RGBAAt(destX, destY)
					calcChannel := func(old, new, alpha uint8) uint8 {
						return uint8((int(old)*(255-int(alpha)))/255) + uint8((int(new)*int(alpha))/255)
					}
					rgba.R = calcChannel(rgba.R, textColor.R, val)
					rgba.G = calcChannel(rgba.G, textColor.G, val)
					rgba.B = calcChannel(rgba.B, textColor.B, val)
					rgba.A = 255 // Just make sure it's fully opaque
					dest.SetRGBA(destX, destY, rgba)
					srcIdx++
					destX++
				}
				srcLineIdx += int(bitmap.width)
				destY++
				destX = destLeft
			}
		}
		penX += ft26p6PosToFloat(gslot.advance.x)
		penY += ft26p6PosToFloat(gslot.advance.y)

		rect.Width += ft26p6PosToFloat(gslot.advance.x)
		lineHeight = max(lineHeight, int(bitmap.rows))
	}
	rect.Height = float64(lineHeight)
	return rect
}
