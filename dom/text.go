// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package dom

import "strconv"

// Comment represents a [DOM Text].
//
// [DOM Text]: https://dom.spec.whatwg.org/#text
type Text interface {
	CharacterData
}
type textImpl struct {
	CharacterData
}

// NewText constructs a new [Text] node.
func NewText(doc Document, text string) Text {
	return &textImpl{newCharacterData(doc, text, TextCharacterData)}
}
func (t textImpl) String() string {
	return strconv.Quote(t.Text())
}
