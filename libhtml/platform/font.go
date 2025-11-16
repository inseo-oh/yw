package platform

import "image"

type Font interface {
	SetTextSize(size int)
	DrawText(text string, dest *image.Image, offset_x, offset_y float64) (width, height float64)
}

func MeasureText(font Font, text string) (width, height float64) {
	return font.DrawText(text, nil, 0, 0)
}
