package layout

import (
	"image/color"
	"log"
	"strings"
	"yw/css"
	"yw/css/cssom"
	"yw/css/display"
	"yw/css/fonts"
	"yw/css/sizing"
	"yw/dom"
	"yw/gfx"
	"yw/util"
)

type treeBuilder struct {
	font gfx.Font
}

func (tb treeBuilder) newText(
	parent box,
	text string,
	rect gfx.Rect,
	color color.RGBA,
	fontSize float64,
) *Text {
	t := Text{}
	t.parent = parent
	t.text = text
	t.rect = rect
	t.font = tb.font
	t.color = color
	t.fontSize = fontSize
	return &t
}
func (tb treeBuilder) newInlineBox(
	parentBcon *blockContainer,
	elem dom.Element,
	rect gfx.Rect,
	widthAuto, heightAuto bool,
) *inlineBox {
	ibox := &inlineBox{}
	ibox.parent = parentBcon
	ibox.elem = elem
	ibox.rect = rect
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
	rect gfx.Rect,
	widthAuto, heightAuto bool,
) *blockContainer {
	bcon := &blockContainer{}
	bcon.parent = parent
	bcon.elem = elem
	bcon.rect = rect
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
		baseRect := bfc.contextCreatorBox().boxRect()
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
		baseRect := bfc.contextCreatorBox().boxRect()
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
	parentNode box,
	domNode dom.Node,
	dryRun bool,
) []Node {
	var parentElem dom.Element
	{
		currNode := parentNode
		for currNode.boxElement() == nil {
			parent := currNode.ParentNode()
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
		panic("FFC should not be nil at this point")
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
	if text, ok := domNode.(dom.CharacterData); ok && text.CharacterDataType() == dom.TextCharacterData {
		parentStyleSet := cssom.ElementDataOf(parentElem).ComputedStyleSet

		//======================================================================
		// Layout for Text nodes
		//======================================================================
		var textNode *Text
		str := text.Text()

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

			// Figure out where we should end current fragment, so that we don't overflow the line box.
			// TODO: We should not do this if we are not doing text wrapping(e.g. whitespace: nowrap).
			for {
				// FIXME: This is very brute-force way of fragmenting text.
				//        We need smarter way to handle this.

				// Calculate width/height using dimensions of the text
				width, _ := gfx.MeasureText(tb.font, fragmentRemaining[:strLen])

				rect = gfx.Rect{Left: left, Top: top, Width: width, Height: metrics.LineHeight}

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
			color := parentStyleSet.Color().ToRgba()
			textNode = tb.newText(parentNode, fragment, rect, color, fontSize)

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
		styleSet := cssom.ElementDataOf(elem).ComputedStyleSet
		computeBoxRect := func(isInline bool) (r gfx.Rect, widthAuto, heightAuto bool) {
			// Calculate left/top position
			boxLeft, boxTop := calcNextPosition(bfc, ifc, writeMode, isInline)

			// Calculate width/height using `width` and `height` property
			boxWidth := styleSet.Width()
			boxHeight := styleSet.Height()
			boxWidthPx := 0.0
			boxHeightPx := 0.0

			// If width or height is auto, we start from 0 and expand it as we layout the children.
			if boxWidth.Type != sizing.Auto {
				parentSize := css.NumFromFloat(parentNode.boxRect().Width)
				boxWidthPx = boxWidth.ComputeUsedValue(parentSize).ToPx(parentSize)
			} else {
				widthAuto = true
			}
			if boxHeight.Type != sizing.Auto {
				parentSize := css.NumFromFloat(parentNode.boxRect().Height)
				boxHeightPx = boxHeight.ComputeUsedValue(parentSize).ToPx(parentSize)
			} else {
				heightAuto = true
			}
			return gfx.Rect{Left: boxLeft, Top: boxTop, Width: boxWidthPx, Height: boxHeightPx}, widthAuto, heightAuto
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
					boxRect.Width = parentNode.boxRect().Width
					widthAuto = false
				} else if writeMode == writeModeVertical && heightAuto {
					boxRect.Height = parentNode.boxRect().Height
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
					if parentNode.boxRect().Height < lineBox.currentLineHeight {
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
						parentBcon = parentBcon.ParentNode().(box)
					}
					ibox := tb.newInlineBox(parentBcon.(*blockContainer), elem, boxRect, widthAuto, heightAuto)
					if !dryRun {
						ibox.initChildren(tb, elem.Children(), writeMode)
					}
					bx = ibox
				} else {
					bcon := tb.newBlockContainer(parentFctx, ifc, parentNode, elem, boxRect, widthAuto, heightAuto)
					if !dryRun {
						bcon.initChildren(tb, elem.Children(), writeMode)
					}
					bx = bcon
				}

				// Increment natural position
				inlineAxisSize := bx.boxRect().Width
				blockAxisSize := bx.boxRect().Height
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
				bcon := tb.newBlockContainer(parentFctx, ifc, parentNode, elem, boxRect, widthAuto, heightAuto)
				if !dryRun {
					bcon.initChildren(tb, elem.Children(), writeMode)
				}

				// Increment natural position
				inlineAxisSize := bcon.boxRect().Width
				blockAxisSize := bcon.boxRect().Height
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
