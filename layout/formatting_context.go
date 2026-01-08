// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package layout

type FormattingContext interface {
	NaturalPos() LogicalPos
	IncrementNaturalPos(inc LogicalPos)
	ContextOwnerBox() Box
}

type formattingContextCommon struct {
	OwnerBox Box
}

func (fc formattingContextCommon) ContextOwnerBox() Box {
	return fc.OwnerBox
}
