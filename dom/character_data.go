package dom

type CharacterDataType uint8

const (
	TextCharacterData = CharacterDataType(iota)
	CommentCharacterData
)

type CharacterData interface {
	Node
	CharacterDataType() CharacterDataType
	Text() string
	AppendText(s string)
}
type characterDataImpl struct {
	Node
	tp   CharacterDataType
	text string
}

func NewCharacterData(doc Document, text string, tp CharacterDataType) CharacterData {
	return &characterDataImpl{NewNode(doc), tp, text}
}

func (c characterDataImpl) CharacterDataType() CharacterDataType {
	return c.tp
}
func (c characterDataImpl) Text() string {
	return c.text
}
func (c *characterDataImpl) AppendText(s string) {
	c.text += s
}
