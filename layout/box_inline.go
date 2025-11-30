package layout

import (
	"fmt"
	"log"
	"yw/dom"
	"yw/gfx"
)

// https://www.w3.org/TR/css-display-3/#inline-box
type inlineBox struct {
	boxCommon
	parentBcon *blockContainer
}

func (bx inlineBox) String() string     { return fmt.Sprintf("inline-box %v at [%v]", bx.elem, bx.rect) }
func (bx inlineBox) NodeType() NodeType { return NodeTypeInlineBox }
func (bx inlineBox) boxRect() gfx.Rect  { return bx.rect }
func (bx inlineBox) IsBlockLevel() bool { return false }

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
