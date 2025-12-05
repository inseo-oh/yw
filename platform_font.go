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
		// Below appear to be 26.6 fixed point values.
		Ascender:   float64(rawMetrics.ascender) / 64.0,
		Descender:  float64(rawMetrics.descender) / 64.0,
		LineHeight: float64(rawMetrics.height) / 64.0,
		// Below appear to be 12.4 fixed point values? Idk, these seem to work fine for now.
		UnderlinePosition:  float64(fnt.face.underline_position) / 16.0,
		UnderlineThickness: float64(fnt.face.underline_thickness) / 16.0,
	}
}
func (fnt ftFont) DrawText(text string, dest *image.RGBA, offsetX, offsetY int, textColor color.Color) image.Rectangle {
	// 26.6 Fixed Point -> float64
	ft26p6PosToFloat := func(p C.FT_Pos) float64 {
		return float64(p) / 64.0
	}

	textColorR, textColorG, textColorB, _ := textColor.RGBA()
	penX, penY := offsetX, offsetY
	lineHeight := 0
	rect := image.Rect(penX, penY, penX, penY)
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
		rect.Min.Y = min(rect.Min.Y, penY-bitmapTop)
		destX := int(penX) + bitmapLeft
		destY := int(penY) - bitmapTop
		destLeft := destX
		if dest != nil {
			for range bitmap.rows {
				srcIdx := srcLineIdx
				for range bitmap.width {
					val := (uint32(bytes[srcIdx]) * 65535) / 255
					r, g, b, _ := dest.At(destX, destY).RGBA()
					calcChannel := func(old, new, alpha uint32) uint32 {
						return (uint32(old)*(65535-uint32(alpha))/65535 + (uint32(new)*uint32(alpha))/65535)
					}
					r = calcChannel(r, textColorR, val)
					g = calcChannel(g, textColorG, val)
					b = calcChannel(b, textColorB, val)
					dest.Set(destX, destY, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255})
					srcIdx++
					destX++
				}
				srcLineIdx += int(bitmap.width)
				destY++
				destX = destLeft
			}
		}
		penX += int(ft26p6PosToFloat(gslot.advance.x))
		penY += int(ft26p6PosToFloat(gslot.advance.y))

		rect.Max.X += int(ft26p6PosToFloat(gslot.advance.x))
		lineHeight = max(lineHeight, int(bitmap.rows))
	}
	rect.Max.Y = rect.Min.Y + lineHeight
	return rect
}
