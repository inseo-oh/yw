// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package layout

import (
	"fmt"
)

// https://www.w3.org/TR/css-display-3/#block-container
type BlockContainerBox struct {
	boxCommon
	Bfc              *BlockFormattingContext
	Ifc              *InlineFormattingContext
	ParentFctx       FormattingContext
	ParentBcon       *BlockContainerBox
	OwnsBfc          bool
	OwnsIfc          bool
	IsAnonymous      bool
	IsInlineFlowRoot bool

	AccumulatedMarginLeft   PhysicalPos
	AccumulatedPaddingLeft  PhysicalPos
	AccumulatedMarginRight  PhysicalPos
	AccumulatedPaddingRight PhysicalPos
}

func (bcon BlockContainerBox) String() string {
	fcStr := ""
	if bcon.OwnsBfc {
		fcStr += "[BFC]"
	}
	if bcon.OwnsIfc {
		fcStr += "[IFC]"
	}
	physMarginRect := bcon.MarginRect.ToPhysicalRect()
	leftStr := fmt.Sprintf("%g+%g+%g", physMarginRect.Left, bcon.Margin.Left, bcon.Padding.Left)
	topStr := fmt.Sprintf("%g+%g+%g", physMarginRect.Top, bcon.Margin.Top, bcon.Padding.Top)
	rightStr := fmt.Sprintf("%g-%g-%g", physMarginRect.right(), bcon.Margin.Right, bcon.Padding.Right)
	bottomStr := fmt.Sprintf("%g-%g-%g", physMarginRect.bottom(), bcon.Margin.Bottom, bcon.Padding.Bottom)
	return fmt.Sprintf(
		"block-container [elem %v] at [LTRB %s %s %s %s (%gx%g)] %s",
		bcon.Elem, leftStr, topStr, rightStr, bottomStr, physMarginRect.Width, physMarginRect.Height, fcStr)
}
func (bcon BlockContainerBox) isBlockLevel() bool { return true }
