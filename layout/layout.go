// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

// Package layout implements the layout engine.
package layout

import (
	"fmt"
	"strings"

	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/gfx"
	"github.com/inseo-oh/yw/platform"
	"github.com/inseo-oh/yw/util"
)

// Build builds the layout starting from the DOM node root.
func Build(root dom.Element, viewportWidth, viewportHeight float64, fontProvider platform.FontProvider) Node {
	// https://www.w3.org/TR/css-display-3/#initial-containing-block
	tb := treeBuilder{}
	tb.font = fontProvider.OpenFont("this_is_not_real_filename.ttf")
	tb.font.SetTextSize(32)
	bfc := &blockFormattingContext{}
	ifc := &inlineFormattingContext{}
	boxRect := logicalRect{inlinePos: 0, blockPos: 0, logicalWidth: viewportWidth, logicalHeight: viewportHeight}
	icb := tb.newBlockContainer(bfc, ifc, nil, nil, nil, boxRect, physicalEdges{}, physicalEdges{}, true, true)
	bfc.ownerBox = icb
	ifc.ownerBox = icb
	ifc.bcon = icb
	icb.initChildren(tb, []dom.Node{root}, []gfx.TextDecorOptions{})
	return icb
}

// PrintTree prints the layout tree.
func PrintTree(node Node) {
	currNode := node
	count := 0
	if !util.IsNil(currNode.parentNode()) {
		for n := currNode.parentNode(); !util.IsNil(n); n = n.parentNode() {
			count += 4
		}
	}
	indent := strings.Repeat(" ", count)
	fmt.Printf("%s%v\n", indent, node)
	if bx, ok := currNode.(box); ok {
		for _, child := range bx.ChildBoxes() {
			PrintTree(child)
		}
		for _, child := range bx.ChildTexts() {
			PrintTree(child)
		}
	}
}
