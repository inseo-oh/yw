// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package layout

type formattingContext interface {
	formattingContextType() formattingContextType
	naturalPos() float64
	incrementNaturalPos(inc float64)
	contextCreatorBox() box
}

type formattingContextCommon struct {
	creatorBox box
}
type formattingContextType uint8

const (
	formattingContextTypeBlock formattingContextType = iota
	formattingContextTypeInline
)

func (fc formattingContextCommon) contextCreatorBox() box {
	return fc.creatorBox
}
