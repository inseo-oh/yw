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

// https://www.w3.org/TR/css-display-3/#inline-box
type inlineBox struct {
	boxCommon
	parentBcon *blockContainer
}

func (bx inlineBox) String() string {
	physMarginRect := bx.marginRect.toPhysicalRect()
	leftStr := fmt.Sprintf("%g+%g+%g", physMarginRect.Left, bx.margin.left, bx.padding.left)
	topStr := fmt.Sprintf("%g+%g+%g", physMarginRect.Top, bx.margin.top, bx.padding.top)
	rightStr := fmt.Sprintf("%g-%g-%g", physMarginRect.right(), bx.margin.right, bx.padding.right)
	bottomStr := fmt.Sprintf("%g-%g-%g", physMarginRect.bottom(), bx.margin.bottom, bx.padding.bottom)
	return fmt.Sprintf(
		"inline-box [elem %v] at [LTRB %s %s %s %s (%gx%g)]",
		bx.elem, leftStr, topStr, rightStr, bottomStr, physMarginRect.Width, physMarginRect.Height)
}
func (bx inlineBox) isBlockLevel() bool { return false }

// NOTE: This should *only* be called once after making layout node.
func (bx *inlineBox) initChildren(tb treeBuilder, children []dom.Node, textDecors []gfx.TextDecorOptions) {
	if len(bx.childBoxes) != 0 || len(bx.childTexts) != 0 {
		panic("this method should be called only once")
	}
	for _, childNode := range children {
		nodes := tb.layoutNode(bx.parentBcon.ifc, bx.parentBcon.bfc, bx.parentBcon.ifc, textDecors, bx, childNode)
		if len(nodes) == 0 {
			continue
		}
		for _, node := range nodes {
			if subBx, ok := node.(box); ok {
				bx.childBoxes = append(bx.childBoxes, subBx)
			} else if txt, ok := node.(*text); ok {
				bx.childTexts = append(bx.childTexts, txt)
			} else {
				log.Panicf("unknown node result %v", node)
			}
		}
	}
}
