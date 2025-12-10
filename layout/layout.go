// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

// Package layout implements the layout engine.
package layout

import (
	"fmt"
	"strings"

	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/gfx"
	"github.com/inseo-oh/yw/platform"
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
	icb := tb.newBlockContainer(bfc, ifc, nil, nil, nil, boxRect, physicalEdges{}, physicalEdges{}, true, true, false)
	bfc.ownerBox = icb
	ifc.ownerBox = icb
	ifc.bcon = icb
	ifc.initialAvailableWidth = viewportWidth
	icb.initChildren(tb, []dom.Node{root}, []gfx.TextDecorOptions{})
	return icb
}

// PrintTree prints the layout tree to standard output.
func PrintTree(node Node, indentLevel int) {
	indent := strings.Repeat(" ", indentLevel*4)
	fmt.Printf("%s%v\n", indent, node)
	if bx, ok := node.(box); ok {
		for _, child := range bx.ChildBoxes() {
			PrintTree(child, indentLevel+1)
		}
		for _, child := range bx.ChildTexts() {
			PrintTree(child, indentLevel+1)
		}
	}
}
