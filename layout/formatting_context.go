package layout

type formattingContext interface {
	formattingContextType() formattingContextType
	naturalPos() float64
	incrementNaturalPos(inc float64)
	contextCreatorBox() box
}

type formattingContextCommon struct {
	isDummyContext bool
	creatorBox     box
}
type formattingContextType uint8

const (
	formattingContextTypeBlock formattingContextType = iota
	formattingContextTypeInline
)

func (fc formattingContextCommon) contextCreatorBox() box {
	return fc.creatorBox
}
