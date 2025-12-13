// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package cssom

// AtRule represents CSS at-rule (e.g. @media { ... })
type AtRule struct {
	Name    string
	Prelude any
	Value   any
}
