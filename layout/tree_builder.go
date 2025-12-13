// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

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
func closestDomElementForBox(bx Box) dom.Element {
	currBox := bx
	for currBox.boxElement() == nil {
		parent := currBox.boxParent()
		if parent == nil {
			break
		}
		currBox = parent
	}
	return currBox.boxElement()
}
func closestParentBlockContainer(bx Box) *blockContainer {
	currBox := bx
	for {
		if _, ok := currBox.(*blockContainer); ok {
			break
		}
		currBox = currBox.boxParent()
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
func elementMarginAndPadding(elem dom.Element, boxParent Box) (margin, padding physicalEdges) {
	styleSetSrc := cssom.ComputedStyleSetSourceOf(elem)
	styleSet := styleSetSrc.ComputedStyleSet()

	if styleSet.MarginTop().IsAuto() || styleSet.MarginBottom().IsAuto() {
		panic("TODO: Support auto margin")
	}
	if styleSet.MarginLeft().IsAuto() || styleSet.MarginRight().IsAuto() {
		panic("TODO: Support auto margin")
	}

	parentLogicalWidth := boxParent.logicalWidth()
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
		baseBlockPos := bfc.contextOwnerBox().boxContentRect().blockPos
		baseInlinePos := bfc.contextOwnerBox().boxContentRect().inlinePos
		inlinePos = baseInlinePos
		if len(ifc.lineBoxes) != 0 {
			lb := ifc.currentLineBox()
			blockPos = lb.initialBlockPos
			inlinePos += ifc.naturalPos()
		} else {
			blockPos = baseBlockPos + bfc.naturalPos()
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
	boxParent Box, parentBcon *blockContainer,
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
		parentSize := css.NumFromFloat(boxParent.boxContentRect().toPhysicalRect().Width)
		boxWidthPhysical = boxWidth.ComputeUsedValue(parentSize).ToPx(parentSize)
	} else {
		physWidthAuto = true
	}
	boxWidthPhysical += margin.horizontalSum() + padding.horizontalSum()
	if boxHeight.Type != sizing.Auto {
		parentSize := css.NumFromFloat(boxParent.boxContentRect().toPhysicalRect().Height)
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
	txt string,
	rect physicalRect,
	color color.Color,
	fontSize float64,
	textDecors []gfx.TextDecorOptions,
) *text {
	t := text{}
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
			if subBx, ok := node.(Box); ok {
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
func (tb treeBuilder) newBlockContainer(
	parentFctx formattingContext,
	ifc *inlineFormattingContext,
	parentBox Box,
	parentBcon *blockContainer,
	elem dom.Element,
	marginRect logicalRect,
	margin, padding physicalEdges,
	physWidthAuto, physHeightAuto bool,
	isInlineFlowRoot bool,
	children []dom.Node, textDecors []gfx.TextDecorOptions,
) *blockContainer {
	bcon := &blockContainer{}

	// ICBs don't have any formatting context yet -- we have to create one.
	if util.IsNil(parentFctx) {
		bfc := &blockFormattingContext{}
		bfc.ownerBox = bcon
		parentFctx = bfc
	}
	// ICBs don't have any IFC yet -- we have to create one.
	if ifc == nil {
		ifc = &inlineFormattingContext{}
		ifc.ownerBox = bcon
		ifc.bcon = bcon
		ifc.initialAvailableWidth = marginRect.toPhysicalRect().Width
		ifc.initialBlockPos = 0
	}

	bcon.parent = parentBox
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

	// Check each children's display type.
	hasInline, hasBlock := false, false
	isInline := make([]bool, len(children))
	for i, childNode := range children {
		isBlockLevel := tb.isElementBlockLevel(bcon.parentFctx, childNode)
		isInline[i] = false
		if isBlockLevel {
			hasBlock = true
		} else {
			hasInline = true
			isInline[i] = true
		}
	}

	// If we have both inline and block-level, we need anonymous block container for inline nodes.
	// (This is actually part of CSS spec)
	needAnonymousBlockContainer := hasInline && hasBlock

	if hasInline && !hasBlock {
		//======================================================================
		// We only have inline contents
		//======================================================================

		// Calculate current initial available width ---------------------------
		var currrentInitialAvailableWidth float64
		if len(bcon.ifc.lineBoxes) != 0 {
			currrentInitialAvailableWidth = bcon.ifc.currentLineBox().availableWidth
		} else {
			currrentInitialAvailableWidth = bcon.ifc.initialAvailableWidth
		}
		// Initialize new IFC --------------------------------------------------
		bcon.ifc = &inlineFormattingContext{}
		bcon.ifc.ownerBox = bcon
		bcon.ifc.bcon = bcon
		// If display mode is inline flow-root, and width is auto, we inherit initial available width from parent.
		if bcon.isInlineFlowRoot && bcon.isWidthAuto() {
			bcon.ifc.initialAvailableWidth = currrentInitialAvailableWidth
		} else {
			bcon.ifc.initialAvailableWidth = bcon.marginRect.logicalWidth
		}
		bcon.ownsIfc = true
		// Calculate common margin-top -----------------------------------------
		commonMarginTop := 0.0
		commonMarginBottom := 0.0
		for _, child := range children {
			var margin physicalEdges
			if elem, ok := child.(dom.Element); ok {
				styleDisplay := cssom.ComputedStyleSetSourceOf(elem).ComputedStyleSet().Display()
				if styleDisplay.Mode == display.OuterInnerMode && (styleDisplay.OuterMode != display.Inline || styleDisplay.InnerMode == display.FlowRoot) {
					margin, _ = elementMarginAndPadding(elem, bcon)
					commonMarginTop = max(commonMarginTop, margin.top)
					commonMarginBottom = max(commonMarginBottom, margin.bottom)
				}
			}
		}
		// Create root inline box ----------------------------------------------
		bcon.bfc.incrementNaturalPos(commonMarginTop)
		bcon.ifc.initialBlockPos = bcon.bfc.ownerBox.boxContentRect().blockPos + bcon.bfc.currentNaturalPos
		ibox := tb.newInlineBox(bcon, nil, bcon.boxContentRect(), physicalEdges{}, physicalEdges{}, false, true, children, textDecors)
		bcon.bfc.incrementNaturalPos(commonMarginBottom)
		bcon.childBoxes = append(bcon.childBoxes, ibox)
		bcon.incrementSize(0, commonMarginTop+commonMarginBottom)
	} else {
		//======================================================================
		// We have either only block contents, or mix of inline and block contents.
		// (In the latter case, we create anonymous block container, so we end up having
		// only block contents)
		//======================================================================

		anonChildren := []dom.Node{}
		for i, childNode := range children {
			var boxes []any
			if isInline[i] && needAnonymousBlockContainer {
				anonChildren = append(anonChildren, childNode)
				if i == len(children)-1 || !isInline[i+1] {
					// Create anonymous block container
					inlinePos, blockPos := computeNextPosition(bcon.bfc, bcon.ifc, bcon, true)
					boxRect := logicalRect{inlinePos: inlinePos, blockPos: blockPos, logicalWidth: bcon.marginRect.logicalWidth, logicalHeight: 0}
					anonBcon := tb.newBlockContainer(bcon.parentFctx, bcon.ifc, bcon, bcon, nil, boxRect, physicalEdges{}, physicalEdges{}, false, true, false, anonChildren, textDecors)
					anonBcon.isAnonymous = true
					bcon.bfc.incrementNaturalPos(anonBcon.marginRect.logicalHeight)
					anonChildren = []dom.Node{} // Clear children list
					boxes = []any{anonBcon}
				}

			} else {
				// Create layout node normally
				boxes = tb.layoutNode(bcon.parentFctx, bcon.bfc, bcon.ifc, textDecors, bcon, childNode)
			}
			if len(boxes) == 0 {
				continue
			}
			for _, bx := range boxes {
				// NOTE: We should only have boxes at this point
				bcon.childBoxes = append(bcon.childBoxes, bx.(Box))
			}

		}
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
func (tb treeBuilder) layoutText(txt dom.Text, boxParent Box, bfc *blockFormattingContext, ifc *inlineFormattingContext, textDecors []gfx.TextDecorOptions) []Node {
	parentElem := closestDomElementForBox(boxParent)
	parentBcon := closestParentBlockContainer(boxParent)
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
		if firstLineBoxCreated && boxParent.isHeightAuto() {
			boxParent.incrementSize(0, lineBox.currentLineHeight)
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
		textNode := tb.newText(fragment, rect, color, fontSize, textDecors)

		if boxParent.isWidthAuto() {
			boxParent.incrementSize(rect.Width, 0)
		}

		ifc.incrementNaturalPos(inlineAxisSize)
		textNodes = append(textNodes, textNode)
		if len(fragmentRemaining) != 0 && strings.TrimLeft(fragmentRemaining, " ") != "" {
			// Create next line --------------------------------------------
			ifc.addLineBox(metrics.LineHeight)
			if boxParent.isHeightAuto() {
				boxParent.incrementSize(0, metrics.LineHeight)
			}
		}
	}

	return textNodes
}
func (tb treeBuilder) layoutElement(elem dom.Element, boxParent Box, parentFctx formattingContext, bfc *blockFormattingContext, ifc *inlineFormattingContext, textDecors []gfx.TextDecorOptions) Box {
	parentBcon := closestParentBlockContainer(boxParent)

	styleSetSrc := cssom.ComputedStyleSetSourceOf(elem)
	styleSet := styleSetSrc.ComputedStyleSet()

	textDecors = elementTextDecoration(elem, textDecors)
	margin, padding := elementMarginAndPadding(elem, boxParent)

	styleDisplay := styleSet.Display()
	switch styleDisplay.Mode {
	case display.DisplayNone:
		return nil
	case display.OuterInnerMode:
		if styleDisplay.OuterMode == display.Inline {
			// Top and bottom margins are handled when creating inline box.
			margin.top = 0
			margin.bottom = 0
		}

		boxRect, physWidthAuto, physHeightAuto := computeBoxRect(elem, bfc, ifc, boxParent, parentBcon, margin, padding, styleDisplay)

		switch styleDisplay.OuterMode {
		case display.Block:
			// Check if we have auto size on a block element. If so, use parent's size and unset auto.
			if physWidthAuto {
				// TODO: Support vertical writing mode
				boxRect.logicalWidth = boxParent.boxContentRect().logicalWidth
				physWidthAuto = false
			}
		case display.Inline:
			// Check if we have auto size on a inline element. If so, use current line height and unset auto.
			if physHeightAuto && len(ifc.lineBoxes) != 0 {
				// TODO: Support vertical writing mode
				boxRect.logicalHeight = ifc.currentLineBox().currentLineHeight
				physHeightAuto = false
			}
		}

		// Increment natural position(if it's auto)
		// XXX: Should we increment width/height if the element uses absolute positioning?
		switch styleDisplay.OuterMode {
		case display.Block:
			if boxParent.isWidthAuto() {
				boxParent.incrementIfNeeded(boxRect.toPhysicalRect().Width, 0)
			}
			if boxParent.isHeightAuto() {
				boxParent.incrementSize(0, boxRect.toPhysicalRect().Height)
			}
		case display.Inline:
			if boxParent.isWidthAuto() {
				boxParent.incrementSize(boxRect.toPhysicalRect().Width, 0)
			}
			if boxParent.isHeightAuto() {
				// TODO
			}
		}

		var bx Box
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
					parentFctx, ifc, boxParent, parentBcon, elem, boxRect, margin, padding, physWidthAuto, physHeightAuto, false, elem.Children(), textDecors)
				bfc.incrementNaturalPos(margin.bottom + padding.bottom) // Consume bottom margin+padding
				bx = bcon
			}
		case display.FlowRoot:
			//==================================================================
			// "flow-root" mode (flow-root, inline-block display modes)
			//==================================================================
			// https://www.w3.org/TR/css-display-3/#valdef-display-flow-root
			bcon := tb.newBlockContainer(parentFctx, ifc, boxParent, parentBcon, elem, boxRect, margin, padding, physWidthAuto, physHeightAuto, true, elem.Children(), textDecors)
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
				if boxParent.isHeightAuto() {
					boxParent.incrementSize(0, heightDiff)
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
	boxParent Box,
	domNode dom.Node,
) []any {
	if n, ok := domNode.(dom.CharacterData); ok && n.CharacterDataType() == dom.CommentCharacterData {
		// No layout is needed for comment nodes
		return nil
	}
	if txt, ok := domNode.(dom.CharacterData); ok && txt.CharacterDataType() == dom.TextCharacterData {
		texts := tb.layoutText(txt, boxParent, bfc, ifc, textDecors)
		res := []any{}
		for _, t := range texts {
			res = append(res, t)
		}
		return res
	} else if elem, ok := domNode.(dom.Element); ok {
		elem := tb.layoutElement(elem, boxParent, parentFctx, bfc, ifc, textDecors)
		if util.IsNil(elem) {
			return nil
		}
		return []any{elem}
	}
	panic("unreachable")
}
