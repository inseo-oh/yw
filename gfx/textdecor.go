// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

package gfx

import "image/color"

// TextDecorOptions represents text decoration options.
type TextDecorOptions struct {
	Type  TextDecorType  // Type of decoration
	Color color.Color    // Decoration color
	Style TextDecorStyle // Decoration style
}

// Type of text decoration
type TextDecorType uint8

const (
	Underline   TextDecorType = iota // Line under text
	Overline                         // Line over text (the opposite of underline)
	ThroughText                      // Line through text
)

// Style of text decoration
type TextDecorStyle uint8

const (
	SolidLine  TextDecorStyle = iota // Solid line
	DoubleLine                       // Double line
	DottedLine                       // Dotted line
	DashedLine                       // Dashed line
	WavyLine                         // Wavy line
)
