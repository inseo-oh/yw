package layout

import (
	"fmt"
	"log"

	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/gfx"
)

// https://www.w3.org/TR/css-display-3/#inline-box
type inlineBox struct {
	boxCommon
	parentBcon *blockContainer
}

func (bx inlineBox) String() string {
	leftStr := fmt.Sprintf("%g(%g+%g)", bx.boxContentRect().Left, bx.marginRect.Left, bx.margin.Left)
	topStr := fmt.Sprintf("%g(%g+%g)", bx.boxContentRect().Top, bx.marginRect.Top, bx.margin.Top)
	rightStr := fmt.Sprintf("%g(%g-%g)", bx.boxContentRect().Right(), bx.marginRect.Right(), bx.margin.Right)
	bottomStr := fmt.Sprintf("%g(%g-%g)", bx.boxContentRect().Bottom(), bx.marginRect.Bottom(), bx.margin.Bottom)
	return fmt.Sprintf(
		"inline-box %v at [LTRB %s %s %s %s (%gx%g)]",
		bx.elem, leftStr, topStr, rightStr, bottomStr, bx.marginRect.Width, bx.marginRect.Height,
	)
}
func (bx inlineBox) NodeType() NodeType      { return NodeTypeInlineBox }
func (bx inlineBox) boxMarginRect() gfx.Rect { return bx.marginRect }
func (bx inlineBox) IsBlockLevel() bool      { return false }

// NOTE: This should *only* be called once after making layout node.
func (bx *inlineBox) initChildren(
	tb treeBuilder,
	children []dom.Node,
	writeMode writeMode,
) {
	if len(bx.childBoxes) != 0 || len(bx.childTexts) != 0 {
		panic("this method should be called only once")
	}
	for _, childNode := range children {
		nodes := tb.makeLayoutForNode(bx.parentBcon.ifc, bx.parentBcon.bfc, bx.parentBcon.ifc, writeMode, bx, childNode, false)
		if len(nodes) == 0 {
			continue
		}
		for _, node := range nodes {
			if subBx, ok := node.(box); ok {
				bx.childBoxes = append(bx.childBoxes, subBx)
			} else if txt, ok := node.(*Text); ok {
				bx.childTexts = append(bx.childTexts, txt)
			} else {
				log.Panicf("unknown node result %v", node)
			}
		}
	}
}
