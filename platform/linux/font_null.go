// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package linux

import (
	"image"
	"image/color"

	"github.com/inseo-oh/yw/gfx"
	"github.com/inseo-oh/yw/platform"
)

type nullFontProvider struct{}

// Returns new [platform.FontProvider] that doesn't do anything.
func NewNullFontProvider() platform.FontProvider {
	return &nullFontProvider{}
}

func (prv nullFontProvider) OpenFont(name string) gfx.Font {
	return nullFont{}
}

type nullFont struct{}

func (fnt nullFont) SetTextSize(size int) {}
func (fnt nullFont) Metrics() gfx.FontMetrics {
	return gfx.FontMetrics{}
}
func (fnt nullFont) DrawText(text string, dest *image.RGBA, offsetX, offsetY int, textColor color.Color) image.Rectangle {
	return image.Rectangle{}
}
