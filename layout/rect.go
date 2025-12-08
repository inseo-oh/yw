// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package layout

type physicalEdges struct{ top, right, bottom, left float64 }

func (e physicalEdges) verticalSum() float64   { return e.top + e.bottom }
func (e physicalEdges) horizontalSum() float64 { return e.left + e.right }

type logicalRect struct{ inlinePos, blockPos, logicalWidth, logicalHeight float64 }

func (r logicalRect) addPadding(edges physicalEdges) logicalRect {
	// TODO: Support vertical writing mode
	r.blockPos += edges.top
	r.inlinePos += edges.left
	r.logicalWidth -= edges.left + edges.right
	r.logicalHeight -= edges.top + edges.bottom
	return r
}
func (r logicalRect) toPhysicalRect() physicalRect {
	// TODO: Support vertical writing mode
	return physicalRect{
		Left:   r.inlinePos,
		Top:    r.blockPos,
		Width:  r.logicalWidth,
		Height: r.logicalHeight,
	}
}

type physicalRect struct{ Left, Top, Width, Height float64 }

func (r physicalRect) right() float64  { return r.Left + r.Width - 1 }
func (r physicalRect) bottom() float64 { return r.Top + r.Height - 1 }
func physicalSizeToLogical(physicalWidth, physicalHeight float64) (logicalWidth, logicalHeight float64) {
	// TODO: Support vertical writing mode
	return physicalWidth, physicalHeight
}
