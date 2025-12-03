// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package dom

// CharacterData represents a [DOM CharacterData], and holds text data.
// CharacterData is abstract type in DOM, and should not be constructed directly.
// See [Text] or [Comment] for that.
//
// [DOM CharacterData]: https://dom.spec.whatwg.org/#characterdata
type CharacterData interface {
	Node
	CharacterDataType() CharacterDataType
	Text() string
	AppendText(s string)
}

// CharacterDataType is type of [CharacterData].
type CharacterDataType uint8

const (
	TextCharacterData    CharacterDataType = iota // Text Node
	CommentCharacterData                          // CommentCharacter node
)

type characterDataImpl struct {
	Node
	tp   CharacterDataType
	text string
}

// newCharacterData constructs a new [CharacterData] node.
func newCharacterData(doc Document, text string, tp CharacterDataType) CharacterData {
	return &characterDataImpl{NewNode(doc), tp, text}
}

// CharacterDataType returns the type.
func (c characterDataImpl) CharacterDataType() CharacterDataType {
	return c.tp
}

// Text returns the text.
func (c characterDataImpl) Text() string {
	return c.text
}

// AppendText appends the given s to c's text.
func (c *characterDataImpl) AppendText(s string) {
	c.text += s
}
