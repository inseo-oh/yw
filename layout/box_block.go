// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

package layout

import (
	"fmt"
	"log"

	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/gfx"
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

	accumulatedMarginLeft   float64
	accumulatedPaddingLeft  float64
	accumulatedMarginRight  float64
	accumulatedPaddingRight float64
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

// NOTE: This should *only* be called once after making layout node.
func (bcon *blockContainer) initChildren(tb treeBuilder, children []dom.Node, textDecors []gfx.TextDecorOptions) {
	if len(bcon.childBoxes) != 0 || len(bcon.childTexts) != 0 {
		panic("this method should be called only once")
	}

	// Check each children's display type - By running dry-run layout on each of them
	hasInline, hasBlock := false, false
	isInline := make([]bool, len(children))
	for i, childNode := range children {
		isBlockLevel := tb.isElementBlockLevel(bcon.parentFctx, childNode)
		isInline[i] = false
		if isBlockLevel {
			hasBlock = true
		} else {
			hasInline = true
			isInline[i] = true
		}
	}

	// If we have both inline and block-level, we need anonymous block container for inline nodes.
	// (This is actually part of CSS spec)
	needAnonymousBlockContainer := hasInline && hasBlock
	if hasInline && !hasBlock {
		var initialAvailableWidth float64
		if len(bcon.ifc.lineBoxes) != 0 {
			initialAvailableWidth = bcon.ifc.currentLineBox().availableWidth
		} else {
			initialAvailableWidth = bcon.ifc.initialAvailableWidth
		}
		bcon.ifc = &inlineFormattingContext{}
		bcon.ifc.ownerBox = bcon
		bcon.ifc.bcon = bcon
		if bcon.isInlineFlowRoot && bcon.isWidthAuto() {
			bcon.ifc.initialAvailableWidth = initialAvailableWidth
		} else {
			bcon.ifc.initialAvailableWidth = bcon.marginRect.logicalWidth
		}
		bcon.ownsIfc = true
	}

	// Now we can layout the children for real
	anonChildren := []dom.Node{}
	for i, childNode := range children {
		var nodes []Node
		if isInline[i] && needAnonymousBlockContainer {
			anonChildren = append(anonChildren, childNode)
			if i == len(children)-1 || !isInline[i+1] {
				// Create anonymous block container
				inlinePos, blockPos := computeNextPosition(bcon.bfc, bcon.ifc, bcon, true)
				boxRect := logicalRect{inlinePos: inlinePos, blockPos: blockPos, logicalWidth: bcon.marginRect.logicalWidth, logicalHeight: 0}
				anonBcon := tb.newBlockContainer(bcon.parentFctx, bcon.ifc, bcon, bcon, nil, boxRect, physicalEdges{}, physicalEdges{}, false, true, false)
				anonBcon.isAnonymous = true
				anonBcon.initChildren(tb, anonChildren, textDecors)
				bcon.bfc.incrementNaturalPos(anonBcon.marginRect.logicalHeight)
				anonChildren = []dom.Node{} // Clear children list
				nodes = []Node{anonBcon}
			}

		} else {
			// Create layout node normally
			nodes = tb.layoutNode(bcon.parentFctx, bcon.bfc, bcon.ifc, textDecors, bcon, childNode)
		}
		if len(nodes) == 0 {
			continue
		}
		for _, node := range nodes {
			if bx, ok := node.(box); ok {
				bcon.childBoxes = append(bcon.childBoxes, bx)
			} else if txt, ok := node.(*text); ok {
				bcon.childTexts = append(bcon.childTexts, txt)
			} else {
				log.Panicf("unknown node result %v", node)
			}
		}

	}
}
