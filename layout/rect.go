// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package layout

// BoxEdges store four values for each edge of a rectangle.
type BoxEdges struct{ Top, Right, Bottom, Left float64 }

// BoxRect represents a rectangular area.
type BoxRect struct{ Left, Top, Width, Height float64 }

// Right returns right of the rectangle.
func (r BoxRect) Right() float64 { return r.Left + r.Width - 1 }

// Bottom returns bottom of the rectangle.
func (r BoxRect) Bottom() float64 { return r.Top + r.Height - 1 }

// AddPadding adds given amount to each edge(top, right, bottom, left), and
// returns resulting rectangle.
func (r BoxRect) AddPadding(edges BoxEdges) BoxRect {
	r.Top += edges.Top
	r.Left += edges.Left
	r.Width -= edges.Left + edges.Right
	r.Height -= edges.Top + edges.Bottom
	return r
}
