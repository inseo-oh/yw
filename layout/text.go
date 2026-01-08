// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package layout

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/inseo-oh/yw/gfx"
	"github.com/inseo-oh/yw/gfx/paint"
)

type Text struct {
	Rect     PhysicalRect
	Text     string
	Font     gfx.Font
	FontSize float64
	Color    color.Color
	Decors   []gfx.TextDecorOptions
}

func (txt Text) String() string {
	return fmt.Sprintf("text %s at [%v]", strconv.Quote(txt.Text), txt.Rect)
}
func (txt Text) MakePaintNode() paint.Node {
	return paint.TextPaint{
		Left:   int(txt.Rect.Left),
		Top:    int(txt.Rect.Top),
		Text:   txt.Text,
		Font:   txt.Font,
		Size:   txt.FontSize,
		Color:  txt.Color,
		Decors: txt.Decors,
	}
}
func (txt Text) isBlockLevel() bool { return false }
