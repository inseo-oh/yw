// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

// Package layout implements the layout engine.
package layout

import (
	"fmt"
	"strings"
)

// LogicalPos represents logical position(i.e. inline/block position)
type LogicalPos float64

// PhysicalPos represents physical position(i.e. location on physical screen)
type PhysicalPos float64

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
