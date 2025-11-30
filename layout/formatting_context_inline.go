package layout

// Inline Formatting Contexts(IFC for short) are responsible for tracking X-axis,
// or more accurately, the primary axis of writing mode.
// (English uses X-axis for writing text, so IFC's position grows X-axis)
//
// or can be also thought as "The opposite axis of BFC", if you really want :D
//
// https://www.w3.org/TR/CSS2/visuren.html#inline-formatting
// https://www.w3.org/TR/css-inline-3/#inline-formatting-context
type inlineFormattingContext struct {
	formattingContextCommon

	bcon      *blockContainer // Block container containing this inline node
	lineBoxes []lineBox       // List of line boxes
}

func (ifc inlineFormattingContext) formattingContextType() FormattingContextType {
	return formattingContextTypeInline
}
func (ifc *inlineFormattingContext) addLineBox(bfc *blockFormattingContext) {
	lb := lineBox{}
	lb.currentLineHeight = 0
	if len(ifc.lineBoxes) != 0 {
		lastLb := ifc.currentLineBox()
		lb.initialBlockPos = lastLb.initialBlockPos + lastLb.currentLineHeight
	} else {
		lb.initialBlockPos = bfc.naturalPos()
	}
	lb.availableWidth = ifc.bcon.marginRect.Width
	ifc.lineBoxes = append(ifc.lineBoxes, lb)
}
func (ifc *inlineFormattingContext) currentLineBox() *lineBox {
	return &ifc.lineBoxes[len(ifc.lineBoxes)-1]
}
func (ifc inlineFormattingContext) naturalPos() float64 {
	return ifc.currentLineBox().currentNaturalPos
}
func (ifc *inlineFormattingContext) incrementNaturalPos(pos float64) {
	if len(ifc.lineBoxes) == 0 {
		ifc.addLineBox(ifc.bcon.bfc)
	}
	lb := ifc.currentLineBox()
	if lb.availableWidth < lb.currentNaturalPos+pos && !ifc.isDummyContext {
		panic("content overflow")
	}
	lb.currentNaturalPos += pos
}

// Line box holds state needed for placing inline contents, such as next inline
// position(which gets reset when entering new line).
//
// https://www.w3.org/TR/css-inline-3/#line-box
type lineBox struct {
	availableWidth    float64
	currentNaturalPos float64
	currentLineHeight float64
	initialBlockPos   float64
}
