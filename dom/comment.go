// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

package dom

import (
	"fmt"
)

// Comment represents a [DOM Comment].
//
// [DOM Comment]: https://dom.spec.whatwg.org/#comment
type Comment interface {
	CharacterData
}
type commentImpl struct {
	CharacterData
}

// NewComment constructs a new [Comment] node.
func NewComment(doc Document, text string) Comment {
	return &commentImpl{newCharacterData(doc, text, CommentCharacterData)}
}
func (cm commentImpl) String() string {
	return fmt.Sprintf("<!-- %s -->", cm.Text())
}
