// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

package layout

import "github.com/inseo-oh/yw/gfx/paint"

type Node interface {
	// MakePaintNode creates a paint node for given node and its children.
	// (So calling this on the root node will generate paint tree for the whole page)
	MakePaintNode() paint.Node

	// String returns description of the node.
	String() string

	parentNode() Node
	isBlockLevel() bool
}
type nodeCommon struct {
	parent Node
}

func (n nodeCommon) parentNode() Node {
	return n.parent
}
