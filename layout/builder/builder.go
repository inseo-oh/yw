// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package builder

import (
	"image/color"
	"log"
	"regexp"
	"strings"

	"github.com/inseo-oh/yw/css"
	"github.com/inseo-oh/yw/css/cssom"
	"github.com/inseo-oh/yw/css/display"
	"github.com/inseo-oh/yw/css/float"
	"github.com/inseo-oh/yw/css/fonts"
	"github.com/inseo-oh/yw/css/props"
	"github.com/inseo-oh/yw/css/sizing"
	"github.com/inseo-oh/yw/css/textdecor"
	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/gfx"
	"github.com/inseo-oh/yw/layout"
	"github.com/inseo-oh/yw/platform"
	"github.com/inseo-oh/yw/util"
)

// BuildLayout builds the layout starting from the DOM node root.
func BuildLayout(root dom.Element, viewportWidth, viewportHeight float64, fontProvider platform.FontProvider) layout.Box {
	// https://www.w3.org/TR/css-display-3/#initial-containing-block
	tb := treeBuilder{}
	tb.font = fontProvider.OpenFont("this_is_not_real_filename.ttf")
	tb.font.SetTextSize(32)
	boxRect := layout.LogicalRect{
		LogicalX:      0,
		LogicalY:      0,
		LogicalWidth:  layout.LogicalPos(viewportWidth),
		LogicalHeight: layout.LogicalPos(viewportHeight),
	}
	icb := tb.newBlockContainer(
		nil, nil, nil, nil, nil, boxRect, layout.PhysicalEdges{}, layout.PhysicalEdges{},
		true, true, false, []dom.Node{root}, []gfx.TextDecorOptions{},
	)
	return icb
}

var (
	spacesAndTabsAfterSegmentBreak  = regexp.MustCompile("\n +")
	spacesAndTabsBeforeSegmentBreak = regexp.MustCompile(" +\n")
	multipleSegmentBreaks           = regexp.MustCompile("\n+")
	multipleSpaces                  = regexp.MustCompile(" +")
)

// https://www.w3.org/TR/css-text-3/#white-space-phase-1
func applyWhitespaceCollapsing(str string, ifc *layout.InlineFormattingContext) string {
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
	if strings.HasSuffix(ifc.WrittenText, " ") {
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

func fontSizeOf(styleSetSrc props.ComputedStyleSetSource) float64 {
	parentFontSize := func() css.Num {
		parentSetSrc := styleSetSrc.ParentSource()
		var parentSize float64
		if !util.IsNil(parentSetSrc) {
			parentSize = fontSizeOf(parentSetSrc)
		} else {
			parentSize = props.DescriptorsMap["font-size"].Initial.(fonts.LengthFontSize).CalculateRealFontSize(nil, nil)
		}
		return css.NumFromFloat(parentSize)
	}
	return styleSetSrc.ComputedStyleSet().FontSize().CalculateRealFontSize(parentFontSize, parentFontSize)
}

func closestDomElementForBox(bx layout.Box) dom.Element {
	currBox := bx
	for currBox.BoxElement() == nil {
		parent := currBox.BoxParent()
		if parent == nil {
			break
		}
		currBox = parent
	}
	return currBox.BoxElement()
}
func closestParentBlockContainer(bx layout.Box) *layout.BlockContainerBox {
	currBox := bx
	for {
		if _, ok := currBox.(*layout.BlockContainerBox); ok {
			break
		}
		currBox = currBox.BoxParent()
	}
	return currBox.(*layout.BlockContainerBox)
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
func elementMarginAndPadding(elem dom.Element, boxParent layout.Box) (margin, padding layout.PhysicalEdges) {
	styleSetSrc := cssom.ComputedStyleSetSourceOf(elem)
	styleSet := styleSetSrc.ComputedStyleSet()

	if styleSet.MarginTop().IsAuto() || styleSet.MarginBottom().IsAuto() {
		panic("TODO: Support auto margin")
	}
	if styleSet.MarginLeft().IsAuto() || styleSet.MarginRight().IsAuto() {
		panic("TODO: Support auto margin")
	}

	parentLogicalWidth := func() css.Num { return css.NumFromFloat(float64(boxParent.LogicalWidth())) }
	fontSize := func() css.Num { return css.NumFromFloat(fontSizeOf(styleSetSrc)) }
	margin = layout.PhysicalEdges{
		Top:    layout.PhysicalPos(styleSet.MarginTop().Value.AsLength(parentLogicalWidth).ToPx(fontSize)),
		Right:  layout.PhysicalPos(styleSet.MarginRight().Value.AsLength(parentLogicalWidth).ToPx(fontSize)),
		Bottom: layout.PhysicalPos(styleSet.MarginBottom().Value.AsLength(parentLogicalWidth).ToPx(fontSize)),
		Left:   layout.PhysicalPos(styleSet.MarginLeft().Value.AsLength(parentLogicalWidth).ToPx(fontSize)),
	}
	padding = layout.PhysicalEdges{
		Top:    layout.PhysicalPos(styleSet.PaddingTop().AsLength(parentLogicalWidth).ToPx(fontSize)),
		Right:  layout.PhysicalPos(styleSet.PaddingRight().AsLength(parentLogicalWidth).ToPx(fontSize)),
		Bottom: layout.PhysicalPos(styleSet.PaddingBottom().AsLength(parentLogicalWidth).ToPx(fontSize)),
		Left:   layout.PhysicalPos(styleSet.PaddingLeft().AsLength(parentLogicalWidth).ToPx(fontSize)),
	}
	return margin, padding
}
func computeNextPosition(bfc *layout.BlockFormattingContext, ifc *layout.InlineFormattingContext, parentBcon *layout.BlockContainerBox, isInline bool) (logicalX, logicalY layout.LogicalPos) {
	if isInline {
		baseLogicalY := bfc.ContextOwnerBox().BoxContentRect().LogicalY
		baseLogicalX := bfc.ContextOwnerBox().BoxContentRect().LogicalX
		logicalX = baseLogicalX
		if len(ifc.LineBoxes) != 0 {
			lb := ifc.CurrentLineBox()
			logicalY = lb.InitialLogicalY
			logicalX += ifc.NaturalPos()
		} else {
			logicalY = baseLogicalY + bfc.NaturalPos()
		}
	} else {
		baseLogicalY := bfc.ContextOwnerBox().BoxContentRect().LogicalY
		baseLogicalX := bfc.ContextOwnerBox().BoxContentRect().LogicalX
		logicalY = bfc.NaturalPos() + baseLogicalY
		logicalX = baseLogicalX
	}
	logicalX += layout.LogicalPos(parentBcon.AccumulatedMarginLeft)
	logicalX += layout.LogicalPos(parentBcon.AccumulatedPaddingLeft)
	return logicalX, logicalY
}
func computeBoxRect(
	elem dom.Element, bfc *layout.BlockFormattingContext, ifc *layout.InlineFormattingContext,
	boxParent layout.Box, parentBcon *layout.BlockContainerBox,
	margin, padding layout.PhysicalEdges,
	styleDisplay display.Display,
) (boxRect layout.LogicalRect, physWidthAuto, physHeightAuto bool) {
	styleSetSrc := cssom.ComputedStyleSetSourceOf(elem)
	styleSet := styleSetSrc.ComputedStyleSet()
	isFloat := styleSet.Float() != float.None
	isInline := !isFloat && styleDisplay.Mode == display.OuterInnerMode && styleDisplay.OuterMode == display.Inline
	isInlineFlowRoot := isInline && styleDisplay.InnerMode == display.FlowRoot

	// Calculate left/top position
	logicalX, logicalY := computeNextPosition(bfc, ifc, parentBcon, isInline)

	var boxWidth, boxHeight sizing.Size
	var boxWidthPhysical, boxHeightPhysical layout.PhysicalPos
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
		containerSize := func() css.Num { return css.NumFromFloat(float64(boxParent.BoxContentRect().ToPhysicalRect().Width)) }
		fontSize := func() css.Num { return css.NumFromFloat(fontSizeOf(styleSetSrc)) }
		boxWidthPhysical = layout.PhysicalPos(boxWidth.ComputeUsedValue(containerSize).ToPx(fontSize))
	} else {
		physWidthAuto = true
	}
	boxWidthPhysical += margin.HorizontalSum() + padding.HorizontalSum()
	if boxHeight.Type != sizing.Auto {
		parentSize := func() css.Num { return css.NumFromFloat(float64(boxParent.BoxContentRect().ToPhysicalRect().Height)) }
		fontSize := func() css.Num { return css.NumFromFloat(fontSizeOf(styleSetSrc)) }
		boxHeightPhysical = layout.PhysicalPos(boxHeight.ComputeUsedValue(parentSize).ToPx(fontSize))
	} else {
		physHeightAuto = true
	}
	boxHeightPhysical += margin.VerticalSum() + padding.VerticalSum()
	boxWidthLogical, boxHeightLogical := layout.PhysicalSizeToLogical(boxWidthPhysical, boxHeightPhysical)

	return layout.LogicalRect{LogicalX: logicalX, LogicalY: logicalY, LogicalWidth: boxWidthLogical, LogicalHeight: boxHeightLogical},
		physWidthAuto, physHeightAuto
}

type treeBuilder struct {
	font gfx.Font
}

func (tb treeBuilder) newText(
	txt string,
	rect layout.PhysicalRect,
	color color.Color,
	fontSize float64,
	textDecors []gfx.TextDecorOptions,
) *layout.Text {
	t := layout.Text{}
	t.Text = txt
	t.Rect = rect
	t.Font = tb.font
	t.Color = color
	t.FontSize = fontSize
	t.Decors = textDecors
	return &t
}
func (tb treeBuilder) newInlineBox(
	parentBcon *layout.BlockContainerBox,
	elem dom.Element,
	marginRect layout.LogicalRect,
	margin, padding layout.PhysicalEdges,
	physWidthAuto, physHeightAuto bool,
	children []dom.Node, textDecors []gfx.TextDecorOptions,
) *layout.InlineBox {
	ibox := &layout.InlineBox{}
	ibox.Parent = parentBcon
	ibox.Elem = elem
	ibox.MarginRect = marginRect
	ibox.Margin = margin
	ibox.Padding = padding
	ibox.PhysicalWidthAuto = physWidthAuto
	ibox.PhysicalHeightAuto = physHeightAuto
	ibox.ParentBcon = parentBcon

	for _, childNode := range children {
		nodes := tb.layoutNode(ibox.ParentBcon.Ifc, ibox.ParentBcon.Bfc, ibox.ParentBcon.Ifc, textDecors, ibox, childNode)
		if len(nodes) == 0 {
			continue
		}
		for _, node := range nodes {
			if subBx, ok := node.(layout.Box); ok {
				ibox.AddChildBox(subBx)
			} else if txt, ok := node.(*layout.Text); ok {
				ibox.AddChildText(txt)
			} else {
				log.Panicf("unknown node result %v", node)
			}
		}
	}

	return ibox
}
func (tb treeBuilder) newBlockContainer(
	parentFctx layout.FormattingContext,
	ifc *layout.InlineFormattingContext,
	parentBox layout.Box,
	parentBcon *layout.BlockContainerBox,
	elem dom.Element,
	marginRect layout.LogicalRect,
	margin, padding layout.PhysicalEdges,
	physWidthAuto, physHeightAuto bool,
	isInlineFlowRoot bool,
	children []dom.Node, textDecors []gfx.TextDecorOptions,
) *layout.BlockContainerBox {
	bcon := &layout.BlockContainerBox{}

	// ICBs don't have any formatting context yet -- we have to create one.
	if util.IsNil(parentFctx) {
		bfc := &layout.BlockFormattingContext{}
		bfc.OwnerBox = bcon
		parentFctx = bfc
	}
	// ICBs don't have any IFC yet -- we have to create one.
	if ifc == nil {
		ifc = &layout.InlineFormattingContext{}
		ifc.OwnerBox = bcon
		ifc.BlockContainer = bcon
		ifc.InitialAvailableWidth = layout.LogicalPos(marginRect.ToPhysicalRect().Width)
		ifc.InitialLogicalY = 0
	}

	bcon.Parent = parentBox
	bcon.ParentBcon = parentBcon
	bcon.Elem = elem
	bcon.MarginRect = marginRect
	bcon.Margin = margin
	bcon.Padding = padding
	bcon.PhysicalWidthAuto = physWidthAuto
	bcon.PhysicalHeightAuto = physHeightAuto
	bcon.ParentFctx = parentFctx
	bcon.Ifc = ifc
	bcon.IsInlineFlowRoot = isInlineFlowRoot

	if parentBcon != nil {
		bcon.AccumulatedMarginLeft = parentBcon.AccumulatedMarginLeft + margin.Left
		bcon.AccumulatedMarginRight = parentBcon.AccumulatedMarginRight + margin.Right
		bcon.AccumulatedPaddingLeft = parentBcon.AccumulatedPaddingLeft + padding.Left
		bcon.AccumulatedPaddingRight = parentBcon.AccumulatedPaddingRight + padding.Right
	}
	if _, ok := parentFctx.(*layout.BlockFormattingContext); !ok || isInlineFlowRoot {
		bcon.Bfc = &layout.BlockFormattingContext{}
		bcon.Bfc.OwnerBox = bcon
		bcon.OwnsBfc = true
	} else {
		bcon.Bfc = parentFctx.(*layout.BlockFormattingContext)
	}

	// Check each children's display type.
	hasInline, hasBlock := false, false
	isInline := make([]bool, len(children))
	for i, childNode := range children {
		isBlockLevel := tb.isElementBlockLevel(bcon.ParentFctx, childNode)
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
		var currrentInitialAvailableWidth layout.LogicalPos
		if len(bcon.Ifc.LineBoxes) != 0 {
			currrentInitialAvailableWidth = bcon.Ifc.CurrentLineBox().AvailableWidth
		} else {
			currrentInitialAvailableWidth = bcon.Ifc.InitialAvailableWidth
		}
		// Initialize new IFC --------------------------------------------------
		bcon.Ifc = &layout.InlineFormattingContext{}
		bcon.Ifc.OwnerBox = bcon
		bcon.Ifc.BlockContainer = bcon
		// If display mode is inline flow-root, and width is auto, we inherit initial available width from parent.
		if bcon.IsInlineFlowRoot && bcon.IsWidthAuto() {
			bcon.Ifc.InitialAvailableWidth = currrentInitialAvailableWidth
		} else {
			bcon.Ifc.InitialAvailableWidth = bcon.MarginRect.LogicalWidth
		}
		bcon.OwnsIfc = true
		// Calculate common margin-top -----------------------------------------
		commonMarginTop := layout.PhysicalPos(0.0)
		commonMarginBottom := layout.PhysicalPos(0.0)
		for _, child := range children {
			var margin layout.PhysicalEdges
			if elem, ok := child.(dom.Element); ok {
				styleDisplay := cssom.ComputedStyleSetSourceOf(elem).ComputedStyleSet().Display()
				if styleDisplay.Mode == display.OuterInnerMode && (styleDisplay.OuterMode != display.Inline || styleDisplay.InnerMode == display.FlowRoot) {
					margin, _ = elementMarginAndPadding(elem, bcon)
					commonMarginTop = max(commonMarginTop, margin.Top)
					commonMarginBottom = max(commonMarginBottom, margin.Bottom)
				}
			}
		}
		// Create root inline box ----------------------------------------------
		bcon.Bfc.IncrementNaturalPos(layout.LogicalPos(commonMarginTop))
		bcon.Ifc.InitialLogicalY = bcon.Bfc.OwnerBox.BoxContentRect().LogicalY + bcon.Bfc.CurrentNaturalPos
		ibox := tb.newInlineBox(bcon, nil, bcon.BoxContentRect(), layout.PhysicalEdges{}, layout.PhysicalEdges{}, false, true, children, textDecors)
		bcon.Bfc.IncrementNaturalPos(layout.LogicalPos(commonMarginBottom))
		bcon.AddChildBox(ibox)
		bcon.IncrementSize(0, layout.LogicalPos(commonMarginTop+commonMarginBottom))
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
					logicalX, logicalY := computeNextPosition(bcon.Bfc, bcon.Ifc, bcon, true)
					boxRect := layout.LogicalRect{LogicalX: logicalX, LogicalY: logicalY, LogicalWidth: bcon.MarginRect.LogicalWidth, LogicalHeight: 0}
					anonBcon := tb.newBlockContainer(bcon.ParentFctx, bcon.Ifc, bcon, bcon, nil, boxRect, layout.PhysicalEdges{}, layout.PhysicalEdges{}, false, true, false, anonChildren, textDecors)
					anonBcon.IsAnonymous = true
					bcon.Bfc.IncrementNaturalPos(anonBcon.MarginRect.LogicalHeight)
					anonChildren = []dom.Node{} // Clear children list
					boxes = []any{anonBcon}
				}

			} else {
				// Create layout node normally
				boxes = tb.layoutNode(bcon.ParentFctx, bcon.Bfc, bcon.Ifc, textDecors, bcon, childNode)
			}
			if len(boxes) == 0 {
				continue
			}
			for _, bx := range boxes {
				// NOTE: We should only have boxes at this point
				bcon.AddChildBox(bx.(layout.Box))
			}

		}
	}

	return bcon
}
func (tb treeBuilder) isElementBlockLevel(parentFctx layout.FormattingContext, domNode dom.Node) bool {
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
					case *layout.BlockFormattingContext, *layout.InlineFormattingContext:
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
func (tb treeBuilder) layoutText(txt dom.Text, boxParent layout.Box, bfc *layout.BlockFormattingContext, ifc *layout.InlineFormattingContext, textDecors []gfx.TextDecorOptions) []any {
	parentElem := closestDomElementForBox(boxParent)
	parentBcon := closestParentBlockContainer(boxParent)
	parentStyleSetSrc := cssom.ComputedStyleSetSourceOf(parentElem)
	parentStyleSet := parentStyleSetSrc.ComputedStyleSet()

	str := applyWhitespaceCollapsing(txt.Text(), ifc)
	if str == "" {
		return nil
	}
	ifc.WrittenText += str

	// Apply text-transform
	if v := parentStyleSet.TextTransform(); !util.IsNil(v) {
		str = v.Apply(str)
	}

	// Calculate the font size
	fontSize := fontSizeOf(parentStyleSetSrc)
	tb.font.SetTextSize(int(fontSize)) // NOTE: Size we set here will only be used for measuring
	metrics := tb.font.Metrics()

	fragmentRemaining := str
	textNodes := []any{}

	for 0 < len(fragmentRemaining) {
		// https://www.w3.org/TR/css-text-3/#white-space-phase-2
		// S1.
		fragmentRemaining = strings.TrimLeft(fragmentRemaining, " ")
		if fragmentRemaining == "" {
			break
		}

		// Create line box if needed
		firstLineBoxCreated := false
		if len(ifc.LineBoxes) == 0 {
			ifc.AddLineBox(metrics.LineHeight)
			firstLineBoxCreated = true
		}
		lineBox := ifc.CurrentLineBox()

		var rect layout.PhysicalRect
		var logicalWidth layout.LogicalPos
		strLen := len(fragmentRemaining)

		// Figure out where we should end current fragment, so that we don't
		// overflow the line box.
		// TODO: We should not do this if we are not doing text wrapping(e.g. whitespace: nowrap).
		for {
			// FIXME: This is very brute-force way of fragmenting text.
			//        We need smarter way to handle this.

			// Calculate physWidth/height using dimensions of the text
			physWidth, _ := gfx.MeasureText(tb.font, fragmentRemaining[:strLen])

			rect = layout.PhysicalRect{Left: 0, Top: 0, Width: layout.PhysicalPos(physWidth), Height: layout.PhysicalPos(metrics.LineHeight)}

			// Check if parent's size is auto and we have to grow its size.
			logicalWidth = layout.LogicalPos(rect.Width)
			// Check if we overflow beyond available size
			if lineBox.CurrentNaturalPos+logicalWidth <= lineBox.AvailableWidth {
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

		lineBox.CurrentLineHeight = max(lineBox.CurrentLineHeight, float64(rect.Height))

		// If we just created a line box, we may have to increase the height.
		if firstLineBoxCreated && boxParent.IsHeightAuto() {
			boxParent.IncrementSize(0, layout.LogicalPos(lineBox.CurrentLineHeight))
		}

		// https://www.w3.org/TR/css-text-3/#white-space-phase-2
		// S3.
		fragment = strings.TrimRight(fragment, " ")

		if fragment == "" {
			continue
		}

		// Calculate left/top position -------------------------------------
		left, top := computeNextPosition(bfc, ifc, parentBcon, true)
		rect.Left = layout.PhysicalPos(left)
		rect.Top = layout.PhysicalPos(top)

		// Make text node --------------------------------------------------
		color := parentStyleSet.Color().ToStdColor(parentStyleSetSrc.CurrentColor())
		textNode := tb.newText(fragment, rect, color, fontSize, textDecors)

		if boxParent.IsWidthAuto() {
			boxParent.IncrementSize(layout.LogicalPos(rect.Width), 0)
		}

		ifc.IncrementNaturalPos(logicalWidth)
		textNodes = append(textNodes, textNode)
		if len(fragmentRemaining) != 0 && strings.TrimLeft(fragmentRemaining, " ") != "" {
			// Create next line --------------------------------------------
			ifc.AddLineBox(metrics.LineHeight)
			if boxParent.IsHeightAuto() {
				boxParent.IncrementSize(0, layout.LogicalPos(metrics.LineHeight))
			}
		}
	}

	return textNodes
}
func (tb treeBuilder) layoutElement(elem dom.Element, boxParent layout.Box, parentFctx layout.FormattingContext, bfc *layout.BlockFormattingContext, ifc *layout.InlineFormattingContext, textDecors []gfx.TextDecorOptions) layout.Box {
	parentBcon := closestParentBlockContainer(boxParent)

	styleSetSrc := cssom.ComputedStyleSetSourceOf(elem)
	styleSet := styleSetSrc.ComputedStyleSet()

	textDecors = elementTextDecoration(elem, textDecors)
	margin, padding := elementMarginAndPadding(elem, boxParent)

	styleDisplay := styleSet.Display()
	styleFloat := styleSet.Float()
	switch styleDisplay.Mode {
	case display.DisplayNone:
		return nil
	case display.OuterInnerMode:
		if styleDisplay.OuterMode == display.Inline {
			// Top and bottom margins are handled when creating inline box.
			margin.Top = 0
			margin.Bottom = 0
		}

		boxRect, physWidthAuto, physHeightAuto := computeBoxRect(elem, bfc, ifc, boxParent, parentBcon, margin, padding, styleDisplay)
		isFloat := styleFloat != float.None

		switch styleDisplay.OuterMode {
		case display.Block:
			// Check if we have auto size on a block element. If so, use parent's size and unset auto.
			if physWidthAuto && !isFloat {
				boxRect.LogicalWidth = boxParent.BoxContentRect().LogicalWidth
				physWidthAuto = false
			}
		case display.Inline:
			// Check if we have auto size on a inline element. If so, use current line height and unset auto.
			if physHeightAuto && len(ifc.LineBoxes) != 0 {
				boxRect.LogicalHeight = layout.LogicalPos(ifc.CurrentLineBox().CurrentLineHeight)
				physHeightAuto = false
			}
		}

		// Increment natural position(if it's auto)
		// XXX: Should we increment width/height if the element uses absolute positioning?
		switch styleDisplay.OuterMode {
		case display.Block:
			if boxParent.IsWidthAuto() {
				boxParent.IncrementIfNeeded(layout.LogicalPos(boxRect.ToPhysicalRect().Width), 0)
			}
			if boxParent.IsHeightAuto() {
				boxParent.IncrementSize(0, layout.LogicalPos(boxRect.ToPhysicalRect().Height))
			}
		case display.Inline:
			if boxParent.IsWidthAuto() {
				boxParent.IncrementSize(layout.LogicalPos(boxRect.ToPhysicalRect().Width), 0)
			}
			if boxParent.IsHeightAuto() {
				// TODO
			}
		}

		var bx layout.Box
		var oldLogicalX layout.LogicalPos
		oldLogicalY := bfc.CurrentNaturalPos
		if len(ifc.LineBoxes) != 0 {
			oldLogicalX = ifc.CurrentLineBox().CurrentNaturalPos
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
				case *layout.BlockFormattingContext, *layout.InlineFormattingContext:
					shouldMakeInlineBox = true
				}
			}
			if shouldMakeInlineBox {
				ibox := tb.newInlineBox(parentBcon, elem, boxRect, margin, padding, physWidthAuto, physHeightAuto, elem.Children(), textDecors)
				bx = ibox
			} else {
				bfc.IncrementNaturalPos(layout.LogicalPos(margin.Top + padding.Top)) // Consume top margin+padding first
				bcon := tb.newBlockContainer(
					parentFctx, ifc, boxParent, parentBcon, elem, boxRect, margin, padding, physWidthAuto, physHeightAuto, false, elem.Children(), textDecors)
				bfc.IncrementNaturalPos(layout.LogicalPos(margin.Bottom + padding.Bottom)) // Consume bottom margin+padding
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
		newLogicalY := bfc.CurrentNaturalPos
		var newLogicalX layout.LogicalPos
		if len(ifc.LineBoxes) != 0 {
			newLogicalX = ifc.CurrentLineBox().CurrentNaturalPos
		}

		switch styleFloat {
		case float.None:
			if bcon, ok := bx.(*layout.BlockContainerBox); ok {
				// Increment natural position (but only the amount that hasn't been incremented)
				switch styleDisplay.OuterMode {
				case display.Block:
					logicalHeight := bcon.BoxMarginRect().LogicalHeight
					posDiff := newLogicalY - oldLogicalY
					bfc.IncrementNaturalPos(logicalHeight - posDiff)
				case display.Inline:
					logicalWidth := bcon.BoxMarginRect().LogicalWidth
					posDiff := newLogicalX - oldLogicalX
					if len(ifc.LineBoxes) == 0 {
						ifc.AddLineBox(0)
					}
					ifc.IncrementNaturalPos(logicalWidth - posDiff)

					lb := ifc.CurrentLineBox()
					heightDiff := float64(bcon.BoxMarginRect().ToPhysicalRect().Height) - lb.CurrentLineHeight
					lb.CurrentLineHeight = max(lb.CurrentLineHeight, float64(bcon.BoxMarginRect().ToPhysicalRect().Height))
					if boxParent.IsHeightAuto() {
						boxParent.IncrementSize(0, layout.LogicalPos(heightDiff))
					}
				}

			}
		case float.Left:
			bfc.LeftFloatingBoxes = append(bfc.LeftFloatingBoxes, bx)
		case float.Right:
			bfc.RightFloatingBoxes = append(bfc.RightFloatingBoxes, bx)
		}
		return bx

	default:
		log.Panicf("TODO: Support display: %v", styleDisplay)
	}
	panic("unreachable")
}
func (tb treeBuilder) layoutNode(
	parentFctx layout.FormattingContext,
	bfc *layout.BlockFormattingContext,
	ifc *layout.InlineFormattingContext,
	textDecors []gfx.TextDecorOptions,
	boxParent layout.Box,
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
