package layout

import "yw/gfx"

type Node interface {
	NodeType() NodeType
	ParentNode() Node
	MakePaintNode() gfx.PaintNode
	IsBlockLevel() bool
	String() string
}
type NodeCommon struct {
	parent Node
}

type NodeType uint8

const (
	NodeTypeInlineBox = NodeType(iota)
	NodeTypeBlockContainer
	NodeTypeText
)

func (n NodeCommon) ParentNode() Node {
	return n.parent
}
