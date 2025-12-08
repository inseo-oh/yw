// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

// Package platform provides abstract platform interface.
package platform

import "github.com/inseo-oh/yw/gfx"

// FontProvider is abstract interface used to provide access to platform's fonts.
type FontProvider interface {
	// OpenFont opens a font with given name.
	OpenFont(name string) gfx.Font
}
