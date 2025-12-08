// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

//go:build !cgo

package stdplatform

import "github.com/inseo-oh/yw/platform"

func NewDefaultFontProvider() platform.FontProvider {
	return NewNullFontProvider()
}
