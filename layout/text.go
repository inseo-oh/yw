// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package layout

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/inseo-oh/yw/gfx"
	"github.com/inseo-oh/yw/gfx/paint"
)

type text struct {
	nodeCommon
	rect     physicalRect
	text     string
	font     gfx.Font
	fontSize float64
	color    color.Color
	decors   []gfx.TextDecorOptions
}

func (txt text) String() string {
	return fmt.Sprintf("text %s at [%v]", strconv.Quote(txt.text), txt.rect)
}
func (txt text) nodeType() nodeType {
	return nodeTypeText
}
func (txt text) MakePaintNode() paint.Node {
	return paint.TextPaint{
		Left:   int(txt.rect.Left),
		Top:    int(txt.rect.Top),
		Text:   txt.text,
		Font:   txt.font,
		Size:   txt.fontSize,
		Color:  txt.color,
		Decors: txt.decors,
	}
}
func (txt text) isBlockLevel() bool { return false }
