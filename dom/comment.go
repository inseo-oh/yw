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
	return &commentImpl{NewCharacterData(doc, text, CommentCharacterData)}
}
func (cm commentImpl) String() string {
	return fmt.Sprintf("<!-- %s -->", cm.Text())
}
