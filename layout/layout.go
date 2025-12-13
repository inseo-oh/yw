// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

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
func Build(root dom.Element, viewportWidth, viewportHeight float64, fontProvider platform.FontProvider) Box {
	// https://www.w3.org/TR/css-display-3/#initial-containing-block
	tb := treeBuilder{}
	tb.font = fontProvider.OpenFont("this_is_not_real_filename.ttf")
	tb.font.SetTextSize(32)
	boxRect := logicalRect{inlinePos: 0, blockPos: 0, logicalWidth: viewportWidth, logicalHeight: viewportHeight}
	icb := tb.newBlockContainer(nil, nil, nil, nil, nil, boxRect, physicalEdges{}, physicalEdges{}, true, true, false, []dom.Node{root}, []gfx.TextDecorOptions{})
	return icb
}

// PrintTree prints the layout tree to standard output.
func PrintTree(bx Box, indentLevel int) {
	indent := strings.Repeat(" ", indentLevel*4)
	fmt.Printf("%s%v\n", indent, bx)
	for _, child := range bx.ChildBoxes() {
		PrintTree(child, indentLevel+1)
	}

	indent = strings.Repeat(" ", (indentLevel+1)*4)
	for _, child := range bx.ChildTexts() {
		fmt.Printf("%s%v\n", indent, child)
	}
}
