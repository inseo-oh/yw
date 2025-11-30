package layout

import (
	"fmt"
	"log"

	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/gfx"
)

// https://www.w3.org/TR/css-display-3/#block-container
type blockContainer struct {
	boxCommon
	bfc        *blockFormattingContext
	ifc        *inlineFormattingContext
	parentFctx formattingContext
	createdBfc bool
	createdIfc bool
}

func (bcon blockContainer) String() string {
	fcStr := ""
	if bcon.createdBfc {
		fcStr += "[BFC]"
	}
	if bcon.createdIfc {
		fcStr += "[IFC]"
	}
	return fmt.Sprintf("block-container [elem %v] at [%v] %s", bcon.elem, bcon.rect, fcStr)
}
func (bcon blockContainer) NodeType() NodeType { return NodeTypeBlockContainer }
func (bcon blockContainer) IsBlockLevel() bool { return true }

// NOTE: This should *only* be called once after making layout node.
func (bcon *blockContainer) initChildren(tb treeBuilder, children []dom.Node, writeMode writeMode) {
	if len(bcon.childBoxes) != 0 || len(bcon.childTexts) != 0 {
		panic("this method should be called only once")
	}

	// Check each children's display type - By running dry-run layout on each of them
	hasInline, hasBlock := false, false
	isInline := make([]bool, len(children))
	for i, childNode := range children {
		nodes := tb.makeLayoutForNode(bcon.parentFctx, bcon.bfc, bcon.ifc, writeMode, bcon, childNode, true)
		isInline[i] = false
		if len(nodes) == 0 {
			continue
		}
		for _, node := range nodes {
			if node.IsBlockLevel() {
				hasBlock = true
			} else {
				hasInline = true
				isInline[i] = true
			}
		}
	}

	// If we have both inline and block-level, we need anonymous block container for inline nodes.
	// (This is actually part of CSS spec)
	needAnonymousBlockContainer := hasInline && hasBlock
	if hasInline && !hasBlock {
		bcon.ifc = &inlineFormattingContext{}
		bcon.ifc.creatorBox = bcon
		bcon.ifc.bcon = bcon
		if bcon.bfc.isDummyContext {
			bcon.ifc.isDummyContext = true
		}
		bcon.createdIfc = true
	}

	// Now we can layout the children for real
	for i, childNode := range children {
		var nodes []Node
		if isInline[i] && needAnonymousBlockContainer {
			// Create anonymous block container
			boxLeft, boxTop := calcNextPosition(bcon.bfc, bcon.ifc, writeMode, false)
			boxRect := gfx.Rect{Left: boxLeft, Top: boxTop, Width: bcon.rect.Width, Height: bcon.rect.Height}
			anonBcon := tb.newBlockContainer(bcon.parentFctx, bcon.ifc, bcon, nil, boxRect, false, false)
			anonBcon.ifc = bcon.ifc
			anonBcon.initChildren(tb, []dom.Node{childNode}, writeMode)
			nodes = []Node{anonBcon}
		} else {
			// Create layout node normally
			nodes = tb.makeLayoutForNode(bcon.parentFctx, bcon.bfc, bcon.ifc, writeMode, bcon, childNode, false)
		}
		if len(nodes) == 0 {
			continue
		}
		for _, node := range nodes {
			if bx, ok := node.(box); ok {
				bcon.childBoxes = append(bcon.childBoxes, bx)
			} else if txt, ok := node.(*Text); ok {
				bcon.childTexts = append(bcon.childTexts, txt)
			} else {
				log.Panicf("unknown node result %v", node)
			}
		}

	}
}
