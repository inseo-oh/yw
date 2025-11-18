package libplatform

import (
	"image"
	"yw/libgfx"
)

type Font interface {
	SetTextSize(size int)
	// Note that dest may be nil -- in that case DrawText should perform a dry-run and return resulting size.
	DrawText(text string, dest *image.RGBA, offset_x, offset_y float64) libgfx.Rect
}

func MeasureText(font Font, text string) (width, height float64) {
	rect := font.DrawText(text, nil, 0, 0)
	return rect.Width, rect.Height
}
