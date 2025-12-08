// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

package dom

import (
	"fmt"
	"slices"
	"testing"
)

func TestDomIter(t *testing.T) {
	nodes := make([]testNode, 12)
	for i := range len(nodes) {
		nodes[i].name = fmt.Sprintf("N%d", i)
	}
	AppendChild(&nodes[0], &nodes[1])
	AppendChild(&nodes[0], &nodes[6])
	AppendChild(&nodes[0], &nodes[7])
	AppendChild(&nodes[1], &nodes[2])
	AppendChild(&nodes[1], &nodes[5])
	AppendChild(&nodes[7], &nodes[8])
	AppendChild(&nodes[7], &nodes[11])
	AppendChild(&nodes[2], &nodes[3])
	AppendChild(&nodes[2], &nodes[4])
	AppendChild(&nodes[8], &nodes[9])
	AppendChild(&nodes[8], &nodes[10])

	doTest := func(gotItems []Node, iterName string, expectedItems []Node) {
		t.Run(iterName, func(t *testing.T) {
			if !slices.Equal(gotItems, expectedItems) {
				t.Errorf("result mismatch! got=%v, expected=%v", gotItems, expectedItems)
			}
		})
	}
	doTest(InclusiveDescendants(&nodes[0]), "inclusive descendants", []Node{
		&nodes[0], &nodes[1], &nodes[2], &nodes[3],
		&nodes[4], &nodes[5], &nodes[6], &nodes[7],
		&nodes[8], &nodes[9], &nodes[10], &nodes[11],
	})
	doTest(InclusiveAncestors(&nodes[11]), "inclusive ancestors", []Node{
		&nodes[11], &nodes[7], &nodes[0],
	})
}
