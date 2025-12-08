// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

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

type treeBuilder struct {
	font gfx.Font
}

func (tb treeBuilder) newText(
	parent box,
	txt string,
	rect BoxRect,
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
	marginRect BoxRect,
	margin, padding BoxEdges,
	widthAuto, heightAuto bool,
) *inlineBox {
	ibox := &inlineBox{}
	ibox.parent = parentBcon
	ibox.elem = elem
	ibox.marginRect = marginRect
	ibox.margin = margin
	ibox.padding = padding
	ibox.widthAuto = widthAuto
	ibox.heightAuto = heightAuto
	ibox.parentBcon = parentBcon
	return ibox
}

// ifc may get overwritten during initChildren()
func (tb treeBuilder) newBlockContainer(
	parentFctx formattingContext,
	ifc *inlineFormattingContext,
	parent Node,
	parentBcon *blockContainer,
	elem dom.Element,
	marginRect BoxRect,
	margin, padding BoxEdges,
	widthAuto, heightAuto bool,
) *blockContainer {
	bcon := &blockContainer{}
	bcon.parent = parent
	bcon.elem = elem
	bcon.marginRect = marginRect
	bcon.margin = margin
	bcon.padding = padding
	bcon.widthAuto = widthAuto
	bcon.heightAuto = heightAuto
	bcon.parentFctx = parentFctx
	bcon.ifc = ifc
	if parentBcon != nil {
		bcon.accumulatedMarginLeft = parentBcon.accumulatedMarginLeft + margin.Left
		bcon.accumulatedMarginRight = parentBcon.accumulatedMarginRight + margin.Right
		bcon.accumulatedPaddingLeft = parentBcon.accumulatedPaddingLeft + padding.Left
		bcon.accumulatedPaddingRight = parentBcon.accumulatedPaddingRight + padding.Right
	}
	if util.IsNil(parentFctx) || parentFctx.formattingContextType() != formattingContextTypeBlock {
		bcon.bfc = makeBfc(bcon)
		bcon.createdBfc = true
	} else {
		bcon.bfc = parentFctx.(*blockFormattingContext)
	}
	return bcon
}

func calcNextPosition(bfc *blockFormattingContext, ifc *inlineFormattingContext, parentBcon *blockContainer, isInline bool) (left, top float64) {
	if isInline {
		baseLeft := ifc.contextCreatorBox().boxContentRect().Left
		left = baseLeft
		if len(ifc.lineBoxes) != 0 {
			lb := ifc.currentLineBox()
			top = lb.initialBlockPos
			left += ifc.naturalPos()
		} else {
			top = bfc.naturalPos()
		}
	} else {
		baseTop := bfc.contextCreatorBox().boxContentRect().Top
		baseLeft := bfc.contextCreatorBox().boxContentRect().Left
		top = bfc.naturalPos() + baseTop
		left = baseLeft
	}
	left += parentBcon.accumulatedMarginLeft
	left += parentBcon.accumulatedPaddingLeft
	return left, top
}

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
					if parentFctx.formattingContextType() == formattingContextTypeBlock ||
						parentFctx.formattingContextType() == formattingContextTypeInline {
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
				return true
			default:
				log.Panicf("TODO: Support display: %v", styleDisplay)
			}

		default:
			log.Panicf("TODO: Support display: %v", styleDisplay)
		}
	}

	panic("unreachable")
}

// This function can be seen as heart of layout process.
func (tb treeBuilder) makeLayoutForNode(
	parentFctx formattingContext,
	bfc *blockFormattingContext,
	ifc *inlineFormattingContext,
	writeMode writeMode,
	textDecors []gfx.TextDecorOptions,
	parentNode box,
	domNode dom.Node,
) []Node {
	var parentElem dom.Element
	{
		currNode := parentNode
		for currNode.boxElement() == nil {
			parent := currNode.parentNode()
			if parent == nil {
				break
			}
			currNode = parent.(box)
		}
		parentElem = currNode.boxElement()
	}

	parentBconBox := parentNode
	for {
		if _, ok := parentBconBox.(*blockContainer); ok {
			break
		}
		parentBconBox = parentBconBox.parentNode().(box)
	}
	parentBcon := parentBconBox.(*blockContainer)

	if n, ok := domNode.(dom.CharacterData); ok && n.CharacterDataType() == dom.CommentCharacterData {
		//======================================================================
		// Layout for Comment nodes
		//======================================================================

		// If you can see comments on screen without devtools, congraturations!
		// You are a very rare person with built-in devtools inside your brain.
		return nil
	}
	if txt, ok := domNode.(dom.CharacterData); ok && txt.CharacterDataType() == dom.TextCharacterData {
		parentStyleSetSrc := cssom.ComputedStyleSetSourceOf(parentElem)
		parentStyleSet := parentStyleSetSrc.ComputedStyleSet()

		//======================================================================
		// Layout for Text nodes
		//======================================================================
		var textNode *text

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
				ifc.addLineBox(bfc)
				firstLineBoxCreated = true
			}
			lineBox := ifc.currentLineBox()

			var rect BoxRect
			var inlineAxisSize float64
			strLen := len(fragmentRemaining)

			// Figure out where we should end current fragment, so that we don't
			// overflow the line box.
			// TODO: We should not do this if we are not doing text wrapping(e.g. whitespace: nowrap).
			for {
				// FIXME: This is very brute-force way of fragmenting text.
				//        We need smarter way to handle this.

				// Calculate width/height using dimensions of the text
				width, _ := gfx.MeasureText(tb.font, fragmentRemaining[:strLen])

				rect = BoxRect{Left: 0, Top: 0, Width: float64(width), Height: metrics.LineHeight}

				// Check if parent's size is auto and we have to grow its size.
				inlineAxisSize = rect.Width
				if writeMode == writeModeVertical {
					inlineAxisSize = rect.Height
				}
				// Check if we overflow beyond available size
				if lineBox.currentNaturalPos+inlineAxisSize <= lineBox.availableWidth {
					// If not, we don't have to fragment text further.
					break
				}
				strLen-- // Decrement length and try again
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
			left, top := calcNextPosition(bfc, ifc, parentBcon, true)
			rect.Left = left
			rect.Top = top

			// Make text node --------------------------------------------------
			color := parentStyleSet.Color().ToStdColor(parentStyleSetSrc.CurrentColor())
			textNode = tb.newText(parentNode, fragment, rect, color, fontSize, textDecors)

			if parentNode.isWidthAuto() {
				parentNode.incrementSize(rect.Width, 0)
			}

			ifc.incrementNaturalPos(inlineAxisSize)
			textNodes = append(textNodes, textNode)
			if len(fragmentRemaining) != 0 && strings.TrimLeft(fragmentRemaining, " ") != "" {
				// Create next line --------------------------------------------
				ifc.addLineBox(bfc)
				if parentNode.isHeightAuto() {
					parentNode.incrementSize(0, lineBox.currentLineHeight)
				}
			}
		}

		return textNodes
	} else if elem, ok := domNode.(dom.Element); ok {
		//======================================================================
		// Layout for Element nodes
		//======================================================================
		styleSetSrc := cssom.ComputedStyleSetSourceOf(elem)
		styleSet := styleSetSrc.ComputedStyleSet()

		// Calculate text-decoration values ------------------------------------

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
		if styleSet.MarginTop().IsAuto() || styleSet.MarginBottom().IsAuto() {
			panic("TODO: Support auto margin")
		}
		if styleSet.MarginLeft().IsAuto() || styleSet.MarginRight().IsAuto() {
			panic("TODO: Support auto margin")
		}
		parentLogicalWidth := parentNode.logicalWidth(writeMode)
		parentFontSize := css.NumFromFloat(fonts.PreferredFontSize) // STUB
		fontSize := styleSet.FontSize().CalculateRealFontSize(parentFontSize).ToPx(parentFontSize)
		margin := BoxEdges{
			Top:    styleSet.MarginTop().Value.AsLength(css.NumFromFloat(parentLogicalWidth)).ToPx(css.NumFromFloat(fontSize)),
			Right:  styleSet.MarginRight().Value.AsLength(css.NumFromFloat(parentLogicalWidth)).ToPx(css.NumFromFloat(fontSize)),
			Bottom: styleSet.MarginBottom().Value.AsLength(css.NumFromFloat(parentLogicalWidth)).ToPx(css.NumFromFloat(fontSize)),
			Left:   styleSet.MarginLeft().Value.AsLength(css.NumFromFloat(parentLogicalWidth)).ToPx(css.NumFromFloat(fontSize)),
		}
		padding := BoxEdges{
			Top:    styleSet.PaddingTop().AsLength(css.NumFromFloat(parentLogicalWidth)).ToPx(css.NumFromFloat(fontSize)),
			Right:  styleSet.PaddingRight().AsLength(css.NumFromFloat(parentLogicalWidth)).ToPx(css.NumFromFloat(fontSize)),
			Bottom: styleSet.PaddingBottom().AsLength(css.NumFromFloat(parentLogicalWidth)).ToPx(css.NumFromFloat(fontSize)),
			Left:   styleSet.PaddingLeft().AsLength(css.NumFromFloat(parentLogicalWidth)).ToPx(css.NumFromFloat(fontSize)),
		}
		computeBoxRect := func(isInline bool) (boxRect BoxRect, widthAuto, heightAuto bool) {
			// Calculate left/top position
			boxLeft, boxTop := calcNextPosition(bfc, ifc, parentBcon, isInline)

			// Calculate width/height using `width` and `height` property
			boxWidth := styleSet.Width()
			boxHeight := styleSet.Height()
			boxWidthPx := 0.0
			boxHeightPx := 0.0

			// If width or height is auto, we start from 0 and expand it as we layout the children.
			if boxWidth.Type != sizing.Auto {
				parentSize := css.NumFromFloat(parentNode.boxContentRect().Width)
				boxWidthPx = boxWidth.ComputeUsedValue(parentSize).ToPx(parentSize)
			} else {
				widthAuto = true
			}
			boxWidthPx += margin.HorizontalSum() + padding.HorizontalSum()
			if boxHeight.Type != sizing.Auto {
				parentSize := css.NumFromFloat(parentNode.boxContentRect().Height)
				boxHeightPx = boxHeight.ComputeUsedValue(parentSize).ToPx(parentSize)
			} else {
				heightAuto = true
			}
			boxHeightPx += margin.VerticalSum() + padding.VerticalSum()

			return BoxRect{Left: boxLeft, Top: boxTop, Width: boxWidthPx, Height: boxHeightPx},
				widthAuto, heightAuto
		}

		styleDisplay := styleSet.Display()
		switch styleDisplay.Mode {
		case display.DisplayNone:
			return nil
		case display.OuterInnerMode:
			isInline := styleDisplay.OuterMode == display.Inline
			boxRect, widthAuto, heightAuto := computeBoxRect(isInline)

			// Check if we have auto size on a block element. If so, use parent's size and unset auto.
			if styleDisplay.OuterMode == display.Block {
				if writeMode == writeModeHorizontal && widthAuto {
					boxRect.Width = parentNode.boxContentRect().Width
					widthAuto = false
				} else if writeMode == writeModeVertical && heightAuto {
					boxRect.Height = parentNode.boxContentRect().Height
					heightAuto = false
				}
			}

			// Increment natural position(if it's auto)
			// XXX: Should we increment width/height if the element uses absolute positioning?
			switch styleDisplay.OuterMode {
			case display.Block:
				if parentNode.isWidthAuto() {
					parentNode.incrementIfNeeded(boxRect.Width, 0)
				}
				if parentNode.isHeightAuto() {
					parentNode.incrementSize(0, boxRect.Height)
				}
			case display.Inline:
				if parentNode.isWidthAuto() {
					parentNode.incrementSize(boxRect.Width, 0)
					if len(ifc.lineBoxes) == 0 {
						ifc.addLineBox(bfc)
					}
				}
				if parentNode.isHeightAuto() {
					if len(ifc.lineBoxes) == 0 {
						ifc.addLineBox(bfc)
					}
					lineBox := ifc.currentLineBox()
					lineBox.currentLineHeight = max(lineBox.currentLineHeight, boxRect.Height)

					// Increment parent's height if needed.
					if parentNode.boxMarginRect().Height < lineBox.currentLineHeight {
						parentNode.incrementIfNeeded(0, lineBox.currentLineHeight)
					}
				}
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
					if parentFctx.formattingContextType() == formattingContextTypeBlock ||
						parentFctx.formattingContextType() == formattingContextTypeInline {
						shouldMakeInlineBox = true
					}
				}
				var bx box
				if shouldMakeInlineBox {
					ibox := tb.newInlineBox(parentBcon, elem, boxRect, margin, padding, widthAuto, heightAuto)
					ibox.initChildren(tb, elem.Children(), writeMode, textDecors)
					bx = ibox
				} else {
					oldBlockPos := bfc.currentNaturalPos
					bfc.incrementNaturalPos(margin.Top + padding.Top) // Consume top margin+padding first
					bcon := tb.newBlockContainer(parentFctx, ifc, parentNode, parentBcon, elem, boxRect, margin, padding, widthAuto, heightAuto)
					bcon.initChildren(tb, elem.Children(), writeMode, textDecors)
					bfc.incrementNaturalPos(margin.Bottom + padding.Bottom) // Consume bottom margin+padding
					newBlockPos := bfc.currentNaturalPos
					bx = bcon

					// Increment natural position (but only the amount that hasn't been incremented)
					inlineAxisSize := bx.boxMarginRect().Width
					blockAxisSize := bx.boxMarginRect().Height
					if writeMode == writeModeVertical {
						inlineAxisSize, blockAxisSize = blockAxisSize, inlineAxisSize
					}
					_ = inlineAxisSize
					posDiff := newBlockPos - oldBlockPos
					bfc.incrementNaturalPos(blockAxisSize - posDiff)
				}
				return []Node{bx}
			case display.FlowRoot:
				//==================================================================
				// "flow-root" mode (flow-root, inline-block display modes)
				//==================================================================
				// https://www.w3.org/TR/css-display-3/#valdef-display-flow-root
				bcon := tb.newBlockContainer(parentFctx, ifc, parentNode, parentBcon, elem, boxRect, margin, padding, widthAuto, heightAuto)
				bcon.initChildren(tb, elem.Children(), writeMode, textDecors)

				// Increment natural position
				inlineAxisSize := bcon.boxMarginRect().Width
				blockAxisSize := bcon.boxMarginRect().Height
				if writeMode == writeModeVertical {
					inlineAxisSize, blockAxisSize = blockAxisSize, inlineAxisSize
				}
				switch styleDisplay.OuterMode {
				case display.Inline:
					ifc.incrementNaturalPos(inlineAxisSize)
				case display.Block:
					bfc.incrementNaturalPos(blockAxisSize)
				}

				return []Node{bcon}
			default:
				log.Panicf("TODO: Support display: %v", styleDisplay)
			}

		default:
			log.Panicf("TODO: Support display: %v", styleDisplay)
		}
	}

	panic("unreachable")
}
