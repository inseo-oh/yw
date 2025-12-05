// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package layout

import (
	"image/color"
	"log"
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
	rect gfx.Rect,
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
	marginRect gfx.Rect,
	margin gfx.Edges,
	widthAuto, heightAuto bool,
) *inlineBox {
	ibox := &inlineBox{}
	ibox.parent = parentBcon
	ibox.elem = elem
	ibox.marginRect = marginRect
	ibox.margin = margin
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
	elem dom.Element,
	marginRect gfx.Rect,
	margin gfx.Edges,
	widthAuto, heightAuto bool,
) *blockContainer {
	bcon := &blockContainer{}
	bcon.parent = parent
	bcon.elem = elem
	bcon.marginRect = marginRect
	bcon.margin = margin
	bcon.widthAuto = widthAuto
	bcon.heightAuto = heightAuto
	bcon.parentFctx = parentFctx
	bcon.ifc = ifc
	if util.IsNil(parentFctx) || parentFctx.formattingContextType() != formattingContextTypeBlock {
		bcon.bfc = makeBfc(bcon)
		bcon.createdBfc = true
	} else {
		bcon.bfc = parentFctx.(*blockFormattingContext)
	}
	return bcon
}

func calcNextPosition(bfc *blockFormattingContext, ifc *inlineFormattingContext, writeMode writeMode, isInline bool) (float64, float64) {
	if isInline {
		var inlinePos, blockBos float64
		if len(ifc.lineBoxes) != 0 {
			inlinePos = ifc.naturalPos()
			blockBos = ifc.currentLineBox().initialBlockPos
		} else {
			inlinePos = 0
			blockBos = bfc.naturalPos()
		}
		baseRect := bfc.contextCreatorBox().boxMarginRect()
		if writeMode == writeModeVertical {
			return baseRect.Left + blockBos, baseRect.Top + inlinePos
		}
		return baseRect.Left + inlinePos, baseRect.Top + blockBos
	} else {
		var x, y float64
		if len(ifc.lineBoxes) != 0 {
			x = ifc.naturalPos()
		} else {
			x = 0
		}
		y = bfc.naturalPos()
		baseRect := bfc.contextCreatorBox().boxMarginRect()
		if writeMode == writeModeVertical {
			return baseRect.Left + y, baseRect.Top + x
		}
		return baseRect.Left + x, baseRect.Top + y
	}

}

// This function can be seen as heart of layout process.
//
// dryRun flag is intended for determine resulting box type. If dryRun is true:
//   - parentFctx, parentNode will be internally replaced by dummy ones,
//     so that they don't affect actual parent context.
//   - will not layout its children, and so returned box will have empty children.
//   - New dummy formatting context will have its isDummyContext set to true.
//     (As of writing this comment, this is mostly for debug prints.
//     outputs with dummy contexts can be confusing when mixed with real ones)
func (tb treeBuilder) makeLayoutForNode(
	parentFctx formattingContext,
	bfc *blockFormattingContext,
	ifc *inlineFormattingContext,
	writeMode writeMode,
	textDecors []gfx.TextDecorOptions,
	parentNode box,
	domNode dom.Node,
	dryRun bool,
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

	if dryRun {
		dummyBcon := &blockContainer{}
		dummyBcon.elem = parentNode.boxElement()
		dummyBcon.bfc = &blockFormattingContext{
			formattingContextCommon: formattingContextCommon{
				isDummyContext: true,
				creatorBox:     dummyBcon,
			},
		}
		bfc = dummyBcon.bfc
		ifc = &inlineFormattingContext{
			formattingContextCommon: formattingContextCommon{
				isDummyContext: true,
				creatorBox:     dummyBcon,
			},
			bcon: dummyBcon,
		}
		if parentFctx.formattingContextType() == formattingContextTypeBlock {
			parentFctx = bfc
		} else {
			parentFctx = ifc
		}
		parentNode = dummyBcon
	}
	if bfc == nil {
		panic("BFC should not be nil at this point")
	}
	if ifc == nil {
		panic("IFC should not be nil at this point")
	}

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
		str := txt.Text()

		// Apply text-transform
		if v := parentStyleSet.TextTransform(); !util.IsNil(v) {
			str = v.Apply(str)
		}

		if str == "" {
			// Nothing to display
			return nil
		}
		str = strings.TrimSpace(str)
		if str == "" {
			str = " "
		}

		// Create line box if needed
		if len(ifc.lineBoxes) == 0 {
			ifc.addLineBox(bfc)
		}

		// Calculate the font size
		parentFontSize := css.NumFromFloat(fonts.PreferredFontSize) // STUB
		fontSize := parentStyleSet.FontSize().CalculateRealFontSize(parentFontSize).ToPx(parentFontSize)
		tb.font.SetTextSize(int(fontSize)) // NOTE: Size we set here will only be used for measuring
		metrics := tb.font.Metrics()

		fragmentRemaining := str
		textNodes := []Node{}

		for 0 < len(fragmentRemaining) {
			lineBox := ifc.currentLineBox()

			var rect gfx.Rect
			var inlineAxisSize float64
			strLen := len(fragmentRemaining)

			// Calculate left/top position
			left, top := calcNextPosition(bfc, ifc, writeMode, true)
			left += parentNode.boxMargin().Left
			top += parentNode.boxMargin().Top

			// Figure out where we should end current fragment, so that we don't overflow the line box.
			// TODO: We should not do this if we are not doing text wrapping(e.g. whitespace: nowrap).
			for {
				// FIXME: This is very brute-force way of fragmenting text.
				//        We need smarter way to handle this.

				// Calculate width/height using dimensions of the text
				width, _ := gfx.MeasureText(tb.font, fragmentRemaining[:strLen])

				rect = gfx.Rect{Left: left, Top: top, Width: float64(width), Height: metrics.LineHeight}

				// Check if parent's size is auto and we have to grow its size.
				inlineAxisSize = rect.Width
				if writeMode == writeModeVertical {
					inlineAxisSize = rect.Height
				}
				// Check if we overflow beyond available size
				if ifc.isDummyContext || lineBox.currentNaturalPos+inlineAxisSize <= lineBox.availableWidth {
					// If not, we don't have to fragment text further.
					break
				}
				strLen-- // Decrement length and try again
			}
			fragment := fragmentRemaining[:strLen]
			fragmentRemaining = fragmentRemaining[strLen:]

			// Make text node
			color := parentStyleSet.Color().ToStdColor(parentStyleSetSrc.CurrentColor())
			textNode = tb.newText(parentNode, fragment, rect, color, fontSize, textDecors)

			if parentNode.isWidthAuto() {
				parentNode.incrementSize(rect.Width, 0)
			}

			lineBox.currentLineHeight = max(lineBox.currentLineHeight, rect.Height)
			if parentNode.isHeightAuto() {
				// Increment parent's height if needed.
				parentNode.incrementIfNeeded(0, lineBox.currentLineHeight)
			}
			ifc.incrementNaturalPos(inlineAxisSize)
			textNodes = append(textNodes, textNode)
			if len(fragmentRemaining) != 0 {
				// Make a new line
				ifc.addLineBox(bfc)
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
		marginParentSize := parentNode.logicalWidth(writeMode)
		parentFontSize := css.NumFromFloat(fonts.PreferredFontSize) // STUB
		fontSize := styleSet.FontSize().CalculateRealFontSize(parentFontSize).ToPx(parentFontSize)
		margin := gfx.Edges{
			Top:    styleSet.MarginTop().Value.AsLength(css.NumFromFloat(marginParentSize)).ToPx(css.NumFromFloat(fontSize)),
			Right:  styleSet.MarginRight().Value.AsLength(css.NumFromFloat(marginParentSize)).ToPx(css.NumFromFloat(fontSize)),
			Bottom: styleSet.MarginBottom().Value.AsLength(css.NumFromFloat(marginParentSize)).ToPx(css.NumFromFloat(fontSize)),
			Left:   styleSet.MarginLeft().Value.AsLength(css.NumFromFloat(marginParentSize)).ToPx(css.NumFromFloat(fontSize)),
		}
		computeBoxRect := func(isInline bool) (boxRect gfx.Rect, widthAuto, heightAuto bool) {
			// Calculate left/top position
			boxLeft, boxTop := calcNextPosition(bfc, ifc, writeMode, isInline)
			boxLeft += parentNode.boxMargin().Top
			boxTop += parentNode.boxMargin().Left

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
			boxWidthPx += margin.Left + margin.Right
			if boxHeight.Type != sizing.Auto {
				parentSize := css.NumFromFloat(parentNode.boxContentRect().Height)
				boxHeightPx = boxHeight.ComputeUsedValue(parentSize).ToPx(parentSize)
			} else {
				heightAuto = true
			}
			boxHeightPx += margin.Top + margin.Bottom

			return gfx.Rect{Left: boxLeft, Top: boxTop, Width: boxWidthPx, Height: boxHeightPx},
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
					boxRect.Width = parentNode.boxMarginRect().Width
					widthAuto = false
				} else if writeMode == writeModeVertical && heightAuto {
					boxRect.Height = parentNode.boxMarginRect().Height
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
					parentBcon := parentNode
					for {
						if _, ok := parentBcon.(*blockContainer); ok {
							break
						}
						parentBcon = parentBcon.parentNode().(box)
					}
					ibox := tb.newInlineBox(parentBcon.(*blockContainer), elem, boxRect, margin, widthAuto, heightAuto)
					if !dryRun {
						ibox.initChildren(tb, elem.Children(), writeMode, textDecors)
					}
					bx = ibox
				} else {
					bcon := tb.newBlockContainer(parentFctx, ifc, parentNode, elem, boxRect, margin, widthAuto, heightAuto)
					if !dryRun {
						bcon.initChildren(tb, elem.Children(), writeMode, textDecors)
					}
					bx = bcon
				}

				// Increment natural position
				inlineAxisSize := bx.boxMarginRect().Width
				blockAxisSize := bx.boxMarginRect().Height
				if writeMode == writeModeVertical {
					inlineAxisSize, blockAxisSize = blockAxisSize, inlineAxisSize
				}
				switch styleDisplay.OuterMode {
				case display.Inline:
					_ = inlineAxisSize
					// ifc.incrementNaturalPos(inlineAxisSize)
				case display.Block:
					bfc.incrementNaturalPos(blockAxisSize)
				}

				return []Node{bx}
			case display.FlowRoot:
				//==================================================================
				// "flow-root" mode (flow-root, inline-block display modes)
				//==================================================================
				// https://www.w3.org/TR/css-display-3/#valdef-display-flow-root
				bcon := tb.newBlockContainer(parentFctx, ifc, parentNode, elem, boxRect, margin, widthAuto, heightAuto)
				if !dryRun {
					bcon.initChildren(tb, elem.Children(), writeMode, textDecors)
				}

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
