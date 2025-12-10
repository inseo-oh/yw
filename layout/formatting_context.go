// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

package layout

type formattingContext interface {
	naturalPos() float64
	incrementNaturalPos(inc float64)
	contextOwnerBox() box
}

type formattingContextCommon struct {
	ownerBox box
}

func (fc formattingContextCommon) contextOwnerBox() box {
	return fc.ownerBox
}
