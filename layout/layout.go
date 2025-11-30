package layout

import (
	"fmt"
	"strings"

	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/gfx"
	"github.com/inseo-oh/yw/platform"
	"github.com/inseo-oh/yw/util"
)

type writeMode uint8

const (
	writeModeHorizontal = writeMode(iota)
	writeModeVertical
)

// https://www.w3.org/TR/css-display-3/#initial-containing-block
func Build(root dom.Element, viewportWidth, viewportHeight float64, plat platform.Platform) Node {
	tb := treeBuilder{}
	tb.font = plat.OpenFont("this_is_not_real_filename.ttf")
	tb.font.SetTextSize(32)
	bfc := &blockFormattingContext{}
	ifc := &inlineFormattingContext{}
	boxRect := gfx.Rect{Left: 0, Top: 0, Width: viewportWidth, Height: viewportHeight}
	icb := tb.newBlockContainer(bfc, ifc, nil, nil, boxRect, gfx.Edges{}, false, false)
	bfc.creatorBox = icb
	ifc.creatorBox = icb
	ifc.bcon = icb
	icb.initChildren(tb, []dom.Node{root}, writeModeHorizontal)
	return icb
}

func PrintTree(node Node) {
	currNode := node
	count := 0
	if !util.IsNil(currNode.ParentNode()) {
		for n := currNode.ParentNode(); !util.IsNil(n); n = n.ParentNode() {
			count += 4
		}
	}
	indent := strings.Repeat(" ", count)
	fmt.Printf("%s%v\n", indent, node)
	if bx, ok := currNode.(box); ok {
		for _, child := range bx.ChildBoxes() {
			PrintTree(child)
		}
		for _, child := range bx.ChildTexts() {
			PrintTree(child)
		}
	}
}
