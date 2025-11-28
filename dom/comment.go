package dom

import "fmt"

type Comment interface {
	CharacterData
}
type commentImpl struct {
	CharacterData
}

func NewComment(doc Document, text string) Comment {
	return &commentImpl{NewCharacterData(doc, text, CommentCharacterData)}
}
func (cm commentImpl) String() string {
	return fmt.Sprintf("<!-- %s -->", cm.Text())
}
