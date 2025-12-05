// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

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
	bfc         *blockFormattingContext
	ifc         *inlineFormattingContext
	parentFctx  formattingContext
	createdBfc  bool
	createdIfc  bool
	isAnonymous bool
}

func (bcon blockContainer) String() string {
	fcStr := ""
	if bcon.createdBfc {
		fcStr += "[BFC]"
	}
	if bcon.createdIfc {
		fcStr += "[IFC]"
	}
	leftStr := fmt.Sprintf("%g(%g+%g)", bcon.boxContentRect().Left, bcon.marginRect.Left, bcon.margin.Left)
	topStr := fmt.Sprintf("%g(%g+%g)", bcon.boxContentRect().Top, bcon.marginRect.Top, bcon.margin.Top)
	rightStr := fmt.Sprintf("%g(%g-%g)", bcon.boxContentRect().Right(), bcon.marginRect.Right(), bcon.margin.Right)
	bottomStr := fmt.Sprintf("%g(%g-%g)", bcon.boxContentRect().Bottom(), bcon.marginRect.Bottom(), bcon.margin.Bottom)
	return fmt.Sprintf(
		"block-container [elem %v] at [LTRB %s %s %s %s (%gx%g)] %s",
		bcon.elem, leftStr, topStr, rightStr, bottomStr, bcon.marginRect.Width, bcon.marginRect.Height, fcStr)
}
func (bcon blockContainer) nodeType() nodeType { return nodeTypeBlockContainer }
func (bcon blockContainer) isBlockLevel() bool { return true }

// NOTE: This should *only* be called once after making layout node.
func (bcon *blockContainer) initChildren(tb treeBuilder, children []dom.Node, writeMode writeMode, textDecors []gfx.TextDecorOptions) {
	if len(bcon.childBoxes) != 0 || len(bcon.childTexts) != 0 {
		panic("this method should be called only once")
	}

	// Check each children's display type - By running dry-run layout on each of them
	hasInline, hasBlock := false, false
	isInline := make([]bool, len(children))
	for i, childNode := range children {
		nodes := tb.makeLayoutForNode(bcon.parentFctx, bcon.bfc, bcon.ifc, writeMode, textDecors, bcon, childNode, true)
		isInline[i] = false
		if len(nodes) == 0 {
			continue
		}
		for _, node := range nodes {
			if node.isBlockLevel() {
				hasBlock = true
			} else {
				hasInline = true
				isInline[i] = true
			}
		}
	}

	// If we have both inline and block-level, we need anonymous block container for inline nodes.
	// (This is actually part of CSS spec)
	needAnonymousBlockContainer := hasInline && hasBlock
	if hasInline && !hasBlock && !bcon.isAnonymous {
		bcon.ifc = &inlineFormattingContext{}
		bcon.ifc.creatorBox = bcon
		bcon.ifc.bcon = bcon
		if bcon.bfc.isDummyContext {
			bcon.ifc.isDummyContext = true
		}
		bcon.createdIfc = true
	}

	// Now we can layout the children for real
	anonChildren := []dom.Node{}
	for i, childNode := range children {
		var nodes []Node
		if isInline[i] && needAnonymousBlockContainer {
			anonChildren = append(anonChildren, childNode)
			if i == len(children)-1 || isInline[i+1] {
				// Create anonymous block container
				boxLeft, boxTop := calcNextPosition(bcon.bfc, bcon.ifc, writeMode, false)
				boxRect := BoxRect{Left: boxLeft, Top: boxTop, Width: bcon.marginRect.Width, Height: bcon.marginRect.Height}
				anonBcon := tb.newBlockContainer(bcon.parentFctx, bcon.ifc, bcon, nil, boxRect, BoxEdges{}, true, true)
				anonBcon.isAnonymous = true
				anonBcon.initChildren(tb, anonChildren, writeMode, textDecors)
				anonChildren = []dom.Node{} // Clear children list
				nodes = []Node{anonBcon}
			}

		} else {
			// Create layout node normally
			nodes = tb.makeLayoutForNode(bcon.parentFctx, bcon.bfc, bcon.ifc, writeMode, textDecors, bcon, childNode, false)
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
