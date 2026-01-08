// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package layout

type PhysicalEdges struct{ Top, Right, Bottom, Left PhysicalPos }

func (e PhysicalEdges) VerticalSum() PhysicalPos   { return e.Top + e.Bottom }
func (e PhysicalEdges) HorizontalSum() PhysicalPos { return e.Left + e.Right }

type LogicalRect struct{ LogicalX, LogicalY, LogicalWidth, LogicalHeight LogicalPos }

func (r LogicalRect) addPadding(edges PhysicalEdges) LogicalRect {
	// TODO: Support vertical writing mode
	r.LogicalY += LogicalPos(edges.Top)
	r.LogicalX += LogicalPos(edges.Left)
	r.LogicalWidth -= LogicalPos(edges.Left + edges.Right)
	r.LogicalHeight -= LogicalPos(edges.Top + edges.Bottom)
	return r
}
func (r LogicalRect) ToPhysicalRect() PhysicalRect {
	// TODO: Support vertical writing mode
	return PhysicalRect{
		Left:   PhysicalPos(r.LogicalX),
		Top:    PhysicalPos(r.LogicalY),
		Width:  PhysicalPos(r.LogicalWidth),
		Height: PhysicalPos(r.LogicalHeight),
	}
}

type PhysicalRect struct{ Left, Top, Width, Height PhysicalPos }

func (r PhysicalRect) right() PhysicalPos  { return r.Left + r.Width - 1 }
func (r PhysicalRect) bottom() PhysicalPos { return r.Top + r.Height - 1 }
func PhysicalSizeToLogical(physicalWidth, physicalHeight PhysicalPos) (logicalWidth, logicalHeight LogicalPos) {
	// TODO: Support vertical writing mode
	return LogicalPos(physicalWidth), LogicalPos(physicalHeight)
}
