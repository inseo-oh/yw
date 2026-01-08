// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package layout

type physicalEdges struct{ top, right, bottom, left PhysicalPos }

func (e physicalEdges) verticalSum() PhysicalPos   { return e.top + e.bottom }
func (e physicalEdges) horizontalSum() PhysicalPos { return e.left + e.right }

type logicalRect struct{ logicalX, logicalY, logicalWidth, logicalHeight LogicalPos }

func (r logicalRect) addPadding(edges physicalEdges) logicalRect {
	// TODO: Support vertical writing mode
	r.logicalY += LogicalPos(edges.top)
	r.logicalX += LogicalPos(edges.left)
	r.logicalWidth -= LogicalPos(edges.left + edges.right)
	r.logicalHeight -= LogicalPos(edges.top + edges.bottom)
	return r
}
func (r logicalRect) toPhysicalRect() physicalRect {
	// TODO: Support vertical writing mode
	return physicalRect{
		Left:   PhysicalPos(r.logicalX),
		Top:    PhysicalPos(r.logicalY),
		Width:  PhysicalPos(r.logicalWidth),
		Height: PhysicalPos(r.logicalHeight),
	}
}

type physicalRect struct{ Left, Top, Width, Height PhysicalPos }

func (r physicalRect) right() PhysicalPos  { return r.Left + r.Width - 1 }
func (r physicalRect) bottom() PhysicalPos { return r.Top + r.Height - 1 }
func physicalSizeToLogical(physicalWidth, physicalHeight PhysicalPos) (logicalWidth, logicalHeight LogicalPos) {
	// TODO: Support vertical writing mode
	return LogicalPos(physicalWidth), LogicalPos(physicalHeight)
}
