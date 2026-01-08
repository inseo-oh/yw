// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package layout

import (
	"fmt"
)

// https://www.w3.org/TR/css-display-3/#block-container
type blockContainer struct {
	boxCommon
	bfc              *blockFormattingContext
	ifc              *inlineFormattingContext
	parentFctx       formattingContext
	parentBcon       *blockContainer
	ownsBfc          bool
	ownsIfc          bool
	isAnonymous      bool
	isInlineFlowRoot bool

	accumulatedMarginLeft   PhysicalPos
	accumulatedPaddingLeft  PhysicalPos
	accumulatedMarginRight  PhysicalPos
	accumulatedPaddingRight PhysicalPos
}

func (bcon blockContainer) String() string {
	fcStr := ""
	if bcon.ownsBfc {
		fcStr += "[BFC]"
	}
	if bcon.ownsIfc {
		fcStr += "[IFC]"
	}
	physMarginRect := bcon.marginRect.toPhysicalRect()
	leftStr := fmt.Sprintf("%g+%g+%g", physMarginRect.Left, bcon.margin.left, bcon.padding.left)
	topStr := fmt.Sprintf("%g+%g+%g", physMarginRect.Top, bcon.margin.top, bcon.padding.top)
	rightStr := fmt.Sprintf("%g-%g-%g", physMarginRect.right(), bcon.margin.right, bcon.padding.right)
	bottomStr := fmt.Sprintf("%g-%g-%g", physMarginRect.bottom(), bcon.margin.bottom, bcon.padding.bottom)
	return fmt.Sprintf(
		"block-container [elem %v] at [LTRB %s %s %s %s (%gx%g)] %s",
		bcon.elem, leftStr, topStr, rightStr, bottomStr, physMarginRect.Width, physMarginRect.Height, fcStr)
}
func (bcon blockContainer) isBlockLevel() bool { return true }
