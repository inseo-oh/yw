package layout

import "github.com/inseo-oh/yw/gfx/paint"

type Node interface {
	// MakePaintNode creates a paint node for given node and its children.
	// (So calling this on the root node will generate paint tree for the whole page)
	MakePaintNode() paint.Node

	// String returns description of the node.
	String() string

	nodeType() nodeType
	parentNode() Node
	isBlockLevel() bool
}
type nodeCommon struct {
	parent Node
}

type nodeType uint8

const (
	nodeTypeInlineBox nodeType = iota
	nodeTypeBlockContainer
	nodeTypeText
)

func (n nodeCommon) parentNode() Node {
	return n.parent
}
