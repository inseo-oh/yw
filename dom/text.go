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
	return &textImpl{NewCharacterData(doc, text, TextCharacterData)}
}
func (t textImpl) String() string {
	return strconv.Quote(t.Text())
}
