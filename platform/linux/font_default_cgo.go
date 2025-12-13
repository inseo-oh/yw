// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

//go:build cgo

package linux

import "github.com/inseo-oh/yw/platform"

// Returns new default [platform.FontProvider]. (In current build configuration, it is the same as [NewFreetypeFontProvider])
func NewDefaultFontProvider() platform.FontProvider {
	return NewFreetypeFontProvider()
}
