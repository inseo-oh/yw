package dom

import "strconv"

type Text interface {
	CharacterData
}
type textImpl struct {
	CharacterData
}

func NewText(doc Document, text string) Text {
	return &textImpl{NewCharacterData(doc, text, TextCharacterData)}
}
func (t textImpl) String() string {
	return strconv.Quote(t.Text())
}
