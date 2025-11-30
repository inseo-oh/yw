package gfx

import (
	"image"
	"image/color"
)

type FontMetrics struct {
	//
	// Ascender -> | @    @  LineHeight -> |
	//             |  @  @                 |
	// baseline -> |___@@___|              |
	//                 @    |              |
	//                @     | <- Descender |
	//               @      |              |
	//                                     |
	//                                     |

	Ascender   float64
	Descender  float64
	LineHeight float64
}

type Font interface {
	SetTextSize(size int)
	Metrics() FontMetrics
	// Note that dest may be nil -- in that case DrawText should perform a dry-run and return resulting size.
	DrawText(text string, dest *image.RGBA, offsetX, offsetY float64, textColor color.RGBA) Rect
}

func MeasureText(font Font, text string) (width, height float64) {
	rect := font.DrawText(text, nil, 0, 0, color.RGBA{})
	return rect.Width, rect.Height
}
