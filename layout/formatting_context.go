package layout

type formattingContext interface {
	formattingContextType() FormattingContextType
	naturalPos() float64
	incrementNaturalPos(inc float64)
	contextCreatorBox() box
}

type formattingContextCommon struct {
	isDummyContext bool
	creatorBox     box
}
type FormattingContextType uint8

const (
	formattingContextTypeBlock FormattingContextType = iota
	formattingContextTypeInline
)

func (fc formattingContextCommon) contextCreatorBox() box {
	return fc.creatorBox
}
