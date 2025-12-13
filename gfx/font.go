// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package gfx

import (
	"image"
	"image/color"
)

// FontMetrics provides various font metrics. Here's what some of fields means:
//
//	|             __                     __
//	| Ascender -> | @    @  LineHeight -> |
//	|             |  @  @                 |
//	| baseline ->_|___@@_____             |
//	|                 @    |              |
//	|                @     | <- Descender |
//	|               @      |_             |
//	|                                     |
//	|                                    _|
type FontMetrics struct {
	Ascender           float64 // Ascender (Distance from top of text to baseline)
	Descender          float64 // Descender (Distance from bottom of text to baseline)
	LineHeight         float64 // Line height of the text
	UnderlinePosition  float64 // Position of underline relative to baseline
	UnderlineThickness float64 // Thickness of underline
}

// Font is an abstract interface that is used to access font information and
// draw text.
type Font interface {
	// SetTextSize sets size of text (in pixels).
	SetTextSize(size int)

	// Metrics returns [FontMetrics] for current text size.
	Metrics() FontMetrics

	// DrawText draws the text to (offsetX, offsetY) position of the dest,
	// using textColor as color.
	//
	// Note that offsetY points to the baseline position, not the top of the text.
	// (Use [FontMetrics] to calculate where the top position should be)
	//
	// DrawText can also perform dry-run. To do so, pass nil to dest.
	// Dry-runs can be used to measure dimensions of text.
	DrawText(text string, dest *image.RGBA, offsetX, offsetY int, textColor color.Color) image.Rectangle
}

// MeasureText performs dry-run text drawing, and returns dimensions of the text.
func MeasureText(font Font, text string) (width, height int) {
	rect := font.DrawText(text, nil, 0, 0, color.RGBA{})
	return rect.Dx(), rect.Dy()
}
