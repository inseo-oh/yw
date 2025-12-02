package layout

import "github.com/inseo-oh/yw/gfx/paint"

type Node interface {
	NodeType() NodeType
	ParentNode() Node
	MakePaintNode() paint.PaintNode
	IsBlockLevel() bool
	String() string
}
type NodeCommon struct {
	parent Node
}

type NodeType uint8

const (
	NodeTypeInlineBox NodeType = iota
	NodeTypeBlockContainer
	NodeTypeText
)

func (n NodeCommon) ParentNode() Node {
	return n.parent
}
