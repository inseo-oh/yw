// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

package layout

import (
	"image/color"
	"log"
	"regexp"
	"strings"

	"github.com/inseo-oh/yw/css"
	"github.com/inseo-oh/yw/css/cssom"
	"github.com/inseo-oh/yw/css/display"
	"github.com/inseo-oh/yw/css/fonts"
	"github.com/inseo-oh/yw/css/sizing"
	"github.com/inseo-oh/yw/css/textdecor"
	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/gfx"
	"github.com/inseo-oh/yw/util"
)

var (
	spacesAndTabsAfterSegmentBreak  = regexp.MustCompile("\n +")
	spacesAndTabsBeforeSegmentBreak = regexp.MustCompile(" +\n")
	multipleSegmentBreaks           = regexp.MustCompile("\n+")
	multipleSpaces                  = regexp.MustCompile(" +")
)

// https://www.w3.org/TR/css-text-3/#white-space-phase-1
func applyWhitespaceCollapsing(str string, ifc *inlineFormattingContext) string {
	// TODO: Add support for white-space: pre, white-space:pre-wrap, or white-space: break-spaces

	//==========================================================================
	// Ignore collapsible spaces and tabs immediately following/preceding segment break.
	//==========================================================================
	// "foo   \n   bar" --> "foo\nbar"
	str = spacesAndTabsAfterSegmentBreak.ReplaceAllLiteralString(str, "\n")
	str = spacesAndTabsBeforeSegmentBreak.ReplaceAllLiteralString(str, "\n")

	//==========================================================================
	// Transform segment breaks according to segment break transform rules.
	// "foo\n\nbar" -> "foo\nbar"
	//==========================================================================
	str = applySegmentBreakTransform(str)

	//==========================================================================
	// Replace tabs with spaces.
	// "foo\t\tbar" -> "foo  bar"
	//==========================================================================
	str = strings.ReplaceAll(str, "\t", " ")

	//==========================================================================
	// Ignore any space following the another, including the ones outside of
	// current text, as long as it's part of the same IFC.
	// "foo   bar" -> "foo bar"
	//
	// TODO: CSS says these extra sapces don't have zero-advance width, and thus invisible,
	// but still retains its soft wrap opportunity, if any.
	//==========================================================================
	if strings.HasSuffix(ifc.writtenText, " ") {
		str = strings.TrimLeft(str, " ")
	}
	str = multipleSpaces.ReplaceAllLiteralString(str, " ")

	return str
}

// https://www.w3.org/TR/css-text-3/#line-break-transform
func applySegmentBreakTransform(str string) string {
	//==========================================================================
	// Remove segment breaks immediately following another.
	// "foo\n\nbar" -> "foo\nbar"
	//==========================================================================
	str = multipleSegmentBreaks.ReplaceAllLiteralString(str, "\n")

	//==========================================================================
	// Turn remaining segment breaks into spaces.
	// "foo\nbar\njaz" -> "foo bar jaz"
	//==========================================================================
	str = strings.ReplaceAll(str, "\n", " ")

	return str
}
func closestDomElementForBox(bx box) dom.Element {
	currNode := bx
	for currNode.boxElement() == nil {
		parent := currNode.parentNode()
		if parent == nil {
			break
		}
		currNode = parent.(box)
	}
	return currNode.boxElement()
}
func closestParentBlockContainer(bx box) *blockContainer {
	currBox := bx
	for {
		if _, ok := currBox.(*blockContainer); ok {
			break
		}
		currBox = currBox.parentNode().(box)
	}
	return currBox.(*blockContainer)
}
func elementTextDecoration(elem dom.Element, textDecors []gfx.TextDecorOptions) []gfx.TextDecorOptions {
	styleSetSrc := cssom.ComputedStyleSetSourceOf(elem)
	styleSet := styleSetSrc.ComputedStyleSet()

	// text decoration in CSS is a bit unusual, because it performs box-level
	// inherit, not normal inheritance.
	// This means, for example, if box A's decoration color is currentColor
	// (value of color property), its children B's decoration color will inherit
	// A's currentColor, not B's one.

	if styleSet.TextDecorationLineValue != nil {
		var decorColor color.Color
		if len(textDecors) != 0 {
			decorColor = textDecors[0].Color
		} else {
			decorColor = styleSetSrc.CurrentColor()
		}
		var decorStyle gfx.TextDecorStyle
		if len(textDecors) != 0 {
			decorStyle = textDecors[0].Style
		} else {
			decorStyle = gfx.SolidLine
		}

		textDecors = []gfx.TextDecorOptions{}
		decorLine := styleSet.TextDecorationLine()
		if (decorLine & textdecor.Overline) != 0 {
			textDecors = append(textDecors, gfx.TextDecorOptions{Type: gfx.Overline, Color: decorColor, Style: decorStyle})
		}
		if (decorLine & textdecor.Underline) != 0 {
			textDecors = append(textDecors, gfx.TextDecorOptions{Type: gfx.Underline, Color: decorColor, Style: decorStyle})
		}
		if (decorLine & textdecor.LineThrough) != 0 {
			textDecors = append(textDecors, gfx.TextDecorOptions{Type: gfx.ThroughText, Color: decorColor, Style: decorStyle})
		}
	}
	if styleSet.TextDecorationColorValue != nil {
		decorColor := styleSet.TextDecorationColor().ToStdColor(styleSetSrc.CurrentColor())
		for i := range len(textDecors) {
			textDecors[i].Color = decorColor
		}
	}
	if styleSet.TextDecorationStyleValue != nil {
		var decorStyle gfx.TextDecorStyle
		switch styleSet.TextDecorationStyle() {
		case textdecor.Solid:
			decorStyle = gfx.SolidLine
		case textdecor.Double:
			decorStyle = gfx.DoubleLine
		case textdecor.Dotted:
			decorStyle = gfx.DottedLine
		case textdecor.Dashed:
			decorStyle = gfx.DashedLine
		case textdecor.Wavy:
			decorStyle = gfx.WavyLine
		}
		for i := range len(textDecors) {
			textDecors[i].Style = decorStyle
		}
	}

	return textDecors
}
func elementMarginAndPadding(elem dom.Element, parentNode box) (margin, padding physicalEdges) {
	styleSetSrc := cssom.ComputedStyleSetSourceOf(elem)
	styleSet := styleSetSrc.ComputedStyleSet()

	if styleSet.MarginTop().IsAuto() || styleSet.MarginBottom().IsAuto() {
		panic("TODO: Support auto margin")
	}
	if styleSet.MarginLeft().IsAuto() || styleSet.MarginRight().IsAuto() {
		panic("TODO: Support auto margin")
	}

	parentLogicalWidth := parentNode.logicalWidth()
	parentFontSize := css.NumFromFloat(fonts.PreferredFontSize) // STUB
	fontSize := styleSet.FontSize().CalculateRealFontSize(parentFontSize).ToPx(parentFontSize)
	margin = physicalEdges{
		top:    styleSet.MarginTop().Value.AsLength(css.NumFromFloat(parentLogicalWidth)).ToPx(css.NumFromFloat(fontSize)),
		right:  styleSet.MarginRight().Value.AsLength(css.NumFromFloat(parentLogicalWidth)).ToPx(css.NumFromFloat(fontSize)),
		bottom: styleSet.MarginBottom().Value.AsLength(css.NumFromFloat(parentLogicalWidth)).ToPx(css.NumFromFloat(fontSize)),
		left:   styleSet.MarginLeft().Value.AsLength(css.NumFromFloat(parentLogicalWidth)).ToPx(css.NumFromFloat(fontSize)),
	}
	padding = physicalEdges{
		top:    styleSet.PaddingTop().AsLength(css.NumFromFloat(parentLogicalWidth)).ToPx(css.NumFromFloat(fontSize)),
		right:  styleSet.PaddingRight().AsLength(css.NumFromFloat(parentLogicalWidth)).ToPx(css.NumFromFloat(fontSize)),
		bottom: styleSet.PaddingBottom().AsLength(css.NumFromFloat(parentLogicalWidth)).ToPx(css.NumFromFloat(fontSize)),
		left:   styleSet.PaddingLeft().AsLength(css.NumFromFloat(parentLogicalWidth)).ToPx(css.NumFromFloat(fontSize)),
	}
	return margin, padding
}
func computeNextPosition(bfc *blockFormattingContext, ifc *inlineFormattingContext, parentBcon *blockContainer, isInline bool) (inlinePos, blockPos float64) {
	if isInline {
		baseInlinePos := bfc.contextOwnerBox().boxContentRect().inlinePos
		inlinePos = baseInlinePos
		if len(ifc.lineBoxes) != 0 {
			lb := ifc.currentLineBox()
			blockPos = lb.initialBlockPos
			inlinePos += ifc.naturalPos()
		} else {
			blockPos = bfc.naturalPos()
		}
	} else {
		baseBlockPos := bfc.contextOwnerBox().boxContentRect().blockPos
		baseInlinePos := bfc.contextOwnerBox().boxContentRect().inlinePos
		blockPos = bfc.naturalPos() + baseBlockPos
		inlinePos = baseInlinePos
	}
	inlinePos += parentBcon.accumulatedMarginLeft
	inlinePos += parentBcon.accumulatedPaddingLeft
	return inlinePos, blockPos
}
func computeBoxRect(
	elem dom.Element, bfc *blockFormattingContext, ifc *inlineFormattingContext,
	parentNode box, parentBcon *blockContainer,
	margin, padding physicalEdges,
	styleDisplay display.Display,
) (boxRect logicalRect, physWidthAuto, physHeightAuto bool) {
	styleSetSrc := cssom.ComputedStyleSetSourceOf(elem)
	styleSet := styleSetSrc.ComputedStyleSet()
	isInline := styleDisplay.Mode == display.OuterInnerMode && styleDisplay.OuterMode == display.Inline
	isInlineFlowRoot := isInline && styleDisplay.InnerMode == display.FlowRoot

	// Calculate left/top position
	inlinePos, blockPos := computeNextPosition(bfc, ifc, parentBcon, isInline)

	var boxWidth, boxHeight sizing.Size
	var boxWidthPhysical, boxHeightPhysical float64
	if !isInline || isInlineFlowRoot {
		// Calculate width/height using `width` and `height` property
		boxWidth = styleSet.Width()
		boxHeight = styleSet.Height()
	} else {
		// Inline elemenrs always have auto size
		boxWidth = sizing.Size{Type: sizing.Auto}
		boxHeight = sizing.Size{Type: sizing.Auto}
	}

	// If width or height is auto, we start from 0 and expand it as we layout the children.
	if boxWidth.Type != sizing.Auto {
		parentSize := css.NumFromFloat(parentNode.boxContentRect().toPhysicalRect().Width)
		boxWidthPhysical = boxWidth.ComputeUsedValue(parentSize).ToPx(parentSize)
	} else {
		physWidthAuto = true
	}
	boxWidthPhysical += margin.horizontalSum() + padding.horizontalSum()
	if boxHeight.Type != sizing.Auto {
		parentSize := css.NumFromFloat(parentNode.boxContentRect().toPhysicalRect().Height)
		boxHeightPhysical = boxHeight.ComputeUsedValue(parentSize).ToPx(parentSize)
	} else {
		physHeightAuto = true
	}
	boxHeightPhysical += margin.verticalSum() + padding.verticalSum()
	boxWidthLogical, boxHeightLogical := physicalSizeToLogical(boxWidthPhysical, boxHeightPhysical)

	return logicalRect{inlinePos: inlinePos, blockPos: blockPos, logicalWidth: boxWidthLogical, logicalHeight: boxHeightLogical},
		physWidthAuto, physHeightAuto
}

type treeBuilder struct {
	font gfx.Font
}

func (tb treeBuilder) newText(
	parent box,
	txt string,
	rect physicalRect,
	color color.Color,
	fontSize float64,
	textDecors []gfx.TextDecorOptions,
) *text {
	t := text{}
	t.parent = parent
	t.text = txt
	t.rect = rect
	t.font = tb.font
	t.color = color
	t.fontSize = fontSize
	t.decors = textDecors
	return &t
}
func (tb treeBuilder) newInlineBox(
	parentBcon *blockContainer,
	elem dom.Element,
	marginRect logicalRect,
	margin, padding physicalEdges,
	physWidthAuto, physHeightAuto bool,
	children []dom.Node, textDecors []gfx.TextDecorOptions,
) *inlineBox {
	ibox := &inlineBox{}
	ibox.parent = parentBcon
	ibox.elem = elem
	ibox.marginRect = marginRect
	ibox.margin = margin
	ibox.padding = padding
	ibox.physicalWidthAuto = physWidthAuto
	ibox.physicalHeightAuto = physHeightAuto
	ibox.parentBcon = parentBcon

	for _, childNode := range children {
		nodes := tb.layoutNode(ibox.parentBcon.ifc, ibox.parentBcon.bfc, ibox.parentBcon.ifc, textDecors, ibox, childNode)
		if len(nodes) == 0 {
			continue
		}
		for _, node := range nodes {
			if subBx, ok := node.(box); ok {
				ibox.childBoxes = append(ibox.childBoxes, subBx)
			} else if txt, ok := node.(*text); ok {
				ibox.childTexts = append(ibox.childTexts, txt)
			} else {
				log.Panicf("unknown node result %v", node)
			}
		}
	}

	return ibox
}

// ifc may get overwritten during initChildren()
func (tb treeBuilder) newBlockContainer(
	parentFctx formattingContext,
	ifc *inlineFormattingContext,
	parent Node,
	parentBcon *blockContainer,
	elem dom.Element,
	marginRect logicalRect,
	margin, padding physicalEdges,
	physWidthAuto, physHeightAuto bool,
	isInlineFlowRoot bool,
) *blockContainer {
	bcon := &blockContainer{}
	bcon.parent = parent
	bcon.parentBcon = parentBcon
	bcon.elem = elem
	bcon.marginRect = marginRect
	bcon.margin = margin
	bcon.padding = padding
	bcon.physicalWidthAuto = physWidthAuto
	bcon.physicalHeightAuto = physHeightAuto
	bcon.parentFctx = parentFctx
	bcon.ifc = ifc
	bcon.isInlineFlowRoot = isInlineFlowRoot
	if parentBcon != nil {
		bcon.accumulatedMarginLeft = parentBcon.accumulatedMarginLeft + margin.left
		bcon.accumulatedMarginRight = parentBcon.accumulatedMarginRight + margin.right
		bcon.accumulatedPaddingLeft = parentBcon.accumulatedPaddingLeft + padding.left
		bcon.accumulatedPaddingRight = parentBcon.accumulatedPaddingRight + padding.right
	}
	if _, ok := parentFctx.(*blockFormattingContext); !ok || isInlineFlowRoot {
		bcon.bfc = &blockFormattingContext{}
		bcon.bfc.ownerBox = bcon
		bcon.ownsBfc = true
	} else {
		bcon.bfc = parentFctx.(*blockFormattingContext)
	}
	return bcon
}

func (tb treeBuilder) isElementBlockLevel(parentFctx formattingContext, domNode dom.Node) bool {
	if n, ok := domNode.(dom.CharacterData); ok && n.CharacterDataType() == dom.CommentCharacterData {
		return false
	}
	if txt, ok := domNode.(dom.CharacterData); ok && txt.CharacterDataType() == dom.TextCharacterData {
		return false
	} else if elem, ok := domNode.(dom.Element); ok {
		styleSetSrc := cssom.ComputedStyleSetSourceOf(elem)
		styleSet := styleSetSrc.ComputedStyleSet()
		styleDisplay := styleSet.Display()
		switch styleDisplay.Mode {
		case display.DisplayNone:
			return false
		case display.OuterInnerMode:
			switch styleDisplay.InnerMode {
			case display.Flow:
				//==================================================================
				// "flow" mode (block, inline, run-in, list-item, inline list-item display modes)
				//==================================================================

				// https://www.w3.org/TR/css-display-3/#valdef-display-flow

				shouldMakeInlineBox := false
				if styleDisplay.OuterMode == display.Inline ||
					styleDisplay.OuterMode == display.RunIn {
					switch parentFctx.(type) {
					case *blockFormattingContext, *inlineFormattingContext:
						shouldMakeInlineBox = true
					}
				}
				if shouldMakeInlineBox {
					return false
				}
				return true
			case display.FlowRoot:
				//==================================================================
				// "flow-root" mode (flow-root, inline-block display modes)
				//==================================================================

				// https://www.w3.org/TR/css-display-3/#valdef-display-flow-root
				return false
			default:
				log.Panicf("TODO: Support display: %v", styleDisplay)
			}

		default:
			log.Panicf("TODO: Support display: %v", styleDisplay)
		}
	}

	panic("unreachable")
}
func (tb treeBuilder) layoutText(txt dom.Text, parentNode box, bfc *blockFormattingContext, ifc *inlineFormattingContext, textDecors []gfx.TextDecorOptions) []Node {
	parentElem := closestDomElementForBox(parentNode)
	parentBcon := closestParentBlockContainer(parentNode)
	parentStyleSetSrc := cssom.ComputedStyleSetSourceOf(parentElem)
	parentStyleSet := parentStyleSetSrc.ComputedStyleSet()

	str := applyWhitespaceCollapsing(txt.Text(), ifc)
	if str == "" {
		return nil
	}
	ifc.writtenText += str

	// Apply text-transform
	if v := parentStyleSet.TextTransform(); !util.IsNil(v) {
		str = v.Apply(str)
	}

	// Calculate the font size
	parentFontSize := css.NumFromFloat(fonts.PreferredFontSize) // STUB
	fontSize := parentStyleSet.FontSize().CalculateRealFontSize(parentFontSize).ToPx(parentFontSize)
	tb.font.SetTextSize(int(fontSize)) // NOTE: Size we set here will only be used for measuring
	metrics := tb.font.Metrics()

	fragmentRemaining := str
	textNodes := []Node{}

	for 0 < len(fragmentRemaining) {
		// https://www.w3.org/TR/css-text-3/#white-space-phase-2
		// S1.
		fragmentRemaining = strings.TrimLeft(fragmentRemaining, " ")
		if fragmentRemaining == "" {
			break
		}

		// Create line box if needed
		firstLineBoxCreated := false
		if len(ifc.lineBoxes) == 0 {
			ifc.addLineBox(metrics.LineHeight)
			firstLineBoxCreated = true
		}
		lineBox := ifc.currentLineBox()

		var rect physicalRect
		var inlineAxisSize float64
		strLen := len(fragmentRemaining)

		// Figure out where we should end current fragment, so that we don't
		// overflow the line box.
		// TODO: We should not do this if we are not doing text wrapping(e.g. whitespace: nowrap).
		for {
			// FIXME: This is very brute-force way of fragmenting text.
			//        We need smarter way to handle this.

			// Calculate physWidth/height using dimensions of the text
			physWidth, _ := gfx.MeasureText(tb.font, fragmentRemaining[:strLen])

			rect = physicalRect{Left: 0, Top: 0, Width: float64(physWidth), Height: metrics.LineHeight}

			// Check if parent's size is auto and we have to grow its size.
			inlineAxisSize = rect.Width
			// Check if we overflow beyond available size
			if lineBox.currentNaturalPos+inlineAxisSize <= lineBox.availableWidth {
				// If not, we don't have to fragment text further.
				break
			}
			strLen-- // Decrement length and try again
		}
		if strLen == 0 {
			// Display at least one chracter per line
			strLen = 1
		}
		fragment := fragmentRemaining[:strLen]
		fragmentRemaining = fragmentRemaining[strLen:]

		lineBox.currentLineHeight = max(lineBox.currentLineHeight, rect.Height)

		// If we just created a line box, we may have to increase the height.
		if firstLineBoxCreated && parentNode.isHeightAuto() {
			parentNode.incrementSize(0, lineBox.currentLineHeight)
		}

		// https://www.w3.org/TR/css-text-3/#white-space-phase-2
		// S3.
		fragment = strings.TrimRight(fragment, " ")

		if fragment == "" {
			continue
		}

		// Calculate left/top position -------------------------------------
		left, top := computeNextPosition(bfc, ifc, parentBcon, true)
		rect.Left = left
		rect.Top = top

		// Make text node --------------------------------------------------
		color := parentStyleSet.Color().ToStdColor(parentStyleSetSrc.CurrentColor())
		textNode := tb.newText(parentNode, fragment, rect, color, fontSize, textDecors)

		if parentNode.isWidthAuto() {
			parentNode.incrementSize(rect.Width, 0)
		}

		ifc.incrementNaturalPos(inlineAxisSize)
		textNodes = append(textNodes, textNode)
		if len(fragmentRemaining) != 0 && strings.TrimLeft(fragmentRemaining, " ") != "" {
			// Create next line --------------------------------------------
			ifc.addLineBox(metrics.LineHeight)
			if parentNode.isHeightAuto() {
				parentNode.incrementSize(0, metrics.LineHeight)
			}
		}
	}

	return textNodes
}
func (tb treeBuilder) layoutElement(elem dom.Element, parentNode box, parentFctx formattingContext, bfc *blockFormattingContext, ifc *inlineFormattingContext, textDecors []gfx.TextDecorOptions) Node {
	parentBcon := closestParentBlockContainer(parentNode)

	styleSetSrc := cssom.ComputedStyleSetSourceOf(elem)
	styleSet := styleSetSrc.ComputedStyleSet()

	textDecors = elementTextDecoration(elem, textDecors)
	margin, padding := elementMarginAndPadding(elem, parentNode)

	styleDisplay := styleSet.Display()
	switch styleDisplay.Mode {
	case display.DisplayNone:
		return nil
	case display.OuterInnerMode:
		boxRect, physWidthAuto, physHeightAuto := computeBoxRect(elem, bfc, ifc, parentNode, parentBcon, margin, padding, styleDisplay)
		isLogicalWidthAuto := func() bool { return physWidthAuto } // STUB
		setLogicalWidthAuto := func(v bool) { physWidthAuto = v }  // STUB

		// Check if we have auto size on a block element. If so, use parent's size and unset auto.
		if styleDisplay.OuterMode == display.Block {
			if isLogicalWidthAuto() {
				boxRect.logicalWidth = parentNode.boxContentRect().logicalWidth
				setLogicalWidthAuto(false)
			}
		}

		// Increment natural position(if it's auto)
		// XXX: Should we increment width/height if the element uses absolute positioning?
		switch styleDisplay.OuterMode {
		case display.Block:
			if parentNode.isWidthAuto() {
				parentNode.incrementIfNeeded(boxRect.toPhysicalRect().Width, 0)
			}
			if parentNode.isHeightAuto() {
				parentNode.incrementSize(0, boxRect.toPhysicalRect().Height)
			}
		case display.Inline:
			if parentNode.isWidthAuto() {
				parentNode.incrementSize(boxRect.toPhysicalRect().Width, 0)
			}
			if parentNode.isHeightAuto() {
				// TODO
			}
		}

		var bx box
		var oldInlinePos float64
		oldBlockPos := bfc.currentNaturalPos
		if len(ifc.lineBoxes) != 0 {
			oldInlinePos = ifc.currentLineBox().currentNaturalPos
		}

		switch styleDisplay.InnerMode {
		case display.Flow:
			//==================================================================
			// "flow" mode (block, inline, run-in, list-item, inline list-item display modes)
			//==================================================================

			// https://www.w3.org/TR/css-display-3/#valdef-display-flow

			shouldMakeInlineBox := false
			if styleDisplay.OuterMode == display.Inline ||
				styleDisplay.OuterMode == display.RunIn {
				switch parentFctx.(type) {
				case *blockFormattingContext, *inlineFormattingContext:
					shouldMakeInlineBox = true
				}
			}
			if shouldMakeInlineBox {
				ibox := tb.newInlineBox(parentBcon, elem, boxRect, margin, padding, physWidthAuto, physHeightAuto, elem.Children(), textDecors)
				bx = ibox
			} else {
				bfc.incrementNaturalPos(margin.top + padding.top) // Consume top margin+padding first
				bcon := tb.newBlockContainer(
					parentFctx, ifc, parentNode, parentBcon, elem, boxRect, margin, padding, physWidthAuto, physHeightAuto, false)
				bcon.initChildren(tb, elem.Children(), textDecors)
				bfc.incrementNaturalPos(margin.bottom + padding.bottom) // Consume bottom margin+padding
				bx = bcon
			}
		case display.FlowRoot:
			//==================================================================
			// "flow-root" mode (flow-root, inline-block display modes)
			//==================================================================
			// https://www.w3.org/TR/css-display-3/#valdef-display-flow-root
			bcon := tb.newBlockContainer(parentFctx, ifc, parentNode, parentBcon, elem, boxRect, margin, padding, physWidthAuto, physHeightAuto, true)
			bcon.initChildren(tb, elem.Children(), textDecors)
			bx = bcon
		default:
			log.Panicf("TODO: Support display: %v", styleDisplay)
		}
		newBlockPos := bfc.currentNaturalPos
		var newInlinePos float64
		if len(ifc.lineBoxes) != 0 {
			newInlinePos = ifc.currentLineBox().currentNaturalPos
		}

		if bcon, ok := bx.(*blockContainer); ok {
			// Increment natural position (but only the amount that hasn't been incremented)
			switch styleDisplay.OuterMode {
			case display.Block:
				logicalHeight := bcon.boxMarginRect().logicalHeight
				posDiff := newBlockPos - oldBlockPos
				bfc.incrementNaturalPos(logicalHeight - posDiff)
			case display.Inline:
				logicalWidth := bcon.boxMarginRect().logicalWidth
				posDiff := newInlinePos - oldInlinePos
				if len(ifc.lineBoxes) == 0 {
					ifc.addLineBox(0)
				}
				ifc.incrementNaturalPos(logicalWidth - posDiff)

				lb := ifc.currentLineBox()
				heightDiff := bcon.boxMarginRect().toPhysicalRect().Height - lb.currentLineHeight
				lb.currentLineHeight = max(lb.currentLineHeight, bcon.boxMarginRect().toPhysicalRect().Height)
				if parentNode.isHeightAuto() {
					parentNode.incrementSize(0, heightDiff)
				}
			}

		}
		return bx

	default:
		log.Panicf("TODO: Support display: %v", styleDisplay)
	}
	panic("unreachable")
}
func (tb treeBuilder) layoutNode(
	parentFctx formattingContext,
	bfc *blockFormattingContext,
	ifc *inlineFormattingContext,
	textDecors []gfx.TextDecorOptions,
	parentNode box,
	domNode dom.Node,
) []Node {
	if n, ok := domNode.(dom.CharacterData); ok && n.CharacterDataType() == dom.CommentCharacterData {
		// No layout is needed for comment nodes
		return nil
	}
	if txt, ok := domNode.(dom.CharacterData); ok && txt.CharacterDataType() == dom.TextCharacterData {
		return tb.layoutText(txt, parentNode, bfc, ifc, textDecors)
	} else if elem, ok := domNode.(dom.Element); ok {
		elem := tb.layoutElement(elem, parentNode, parentFctx, bfc, ifc, textDecors)
		if util.IsNil(elem) {
			return nil
		}
		return []Node{elem}
	}
	panic("unreachable")
}
