// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

// Package props provides various CSS properties, and [ComputedStyleSet]
// containing computed values for properties.
package props

import "image/color"

//go:generate go run ./gen

// Source of [ComputedStyleSet]. [ComputedStyleSet] is stored inside DOM element,
// and this type is used to avoid props package from depending on dom package.
type ComputedStyleSetSource interface {
	ComputedStyleSet() *ComputedStyleSet
	ParentSource() ComputedStyleSetSource
	CurrentColor() color.Color
}

// PropertyValue represents an property value.
type PropertyValue interface {
	String() string
}

// Descriptor represents information about each property.
type Descriptor struct {
	Initial   PropertyValue
	ApplyFunc func(dest *ComputedStyleSet, value any)
}
