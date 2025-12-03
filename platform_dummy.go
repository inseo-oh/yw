// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

//go:build !cgo

package main

import (
	"image"
	"image/color"

	"github.com/inseo-oh/yw/gfx"
)

type dummyPlatformImpl struct{}

func initPlatform() *dummyPlatformImpl {
	return &dummyPlatformImpl{}
}

func (plat dummyPlatformImpl) OpenFont(name string) gfx.Font {
	return dummyFont{}
}

type dummyFont struct{}

func (fnt dummyFont) SetTextSize(size int) {}
func (fnt dummyFont) Metrics() gfx.FontMetrics {
	return gfx.FontMetrics{}
}
func (fnt dummyFont) DrawText(text string, dest *image.RGBA, offsetX, offsetY float64, textColor color.RGBA) gfx.Rect {
	return gfx.Rect{}
}
