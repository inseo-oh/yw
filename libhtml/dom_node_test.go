package libhtml

import (
	"fmt"
	"slices"
	"testing"
)

func TestDomIter(t *testing.T) {
	nodes := make([]dom_TestNode, 12)
	for i := range len(nodes) {
		nodes[i].name = fmt.Sprintf("N%d", i)
	}
	dom_node_append_child(&nodes[0], &nodes[1])
	dom_node_append_child(&nodes[0], &nodes[6])
	dom_node_append_child(&nodes[0], &nodes[7])
	dom_node_append_child(&nodes[1], &nodes[2])
	dom_node_append_child(&nodes[1], &nodes[5])
	dom_node_append_child(&nodes[7], &nodes[8])
	dom_node_append_child(&nodes[7], &nodes[11])
	dom_node_append_child(&nodes[2], &nodes[3])
	dom_node_append_child(&nodes[2], &nodes[4])
	dom_node_append_child(&nodes[8], &nodes[9])
	dom_node_append_child(&nodes[8], &nodes[10])

	do_test := func(got_items []dom_Node, iter_name string, expected_items []dom_Node) {
		t.Run(iter_name, func(t *testing.T) {
			if !slices.Equal(got_items, expected_items) {
				t.Errorf("result mismatch! got=%v, expected=%v", got_items, expected_items)
			}
		})
	}
	do_test(dom_node_inclusive_descendants(&nodes[0]), "inclusive descendants", []dom_Node{
		&nodes[0], &nodes[1], &nodes[2], &nodes[3],
		&nodes[4], &nodes[5], &nodes[6], &nodes[7],
		&nodes[8], &nodes[9], &nodes[10], &nodes[11],
	})
	do_test(dom_node_inclusive_ancestors(&nodes[11]), "inclusive ancestors", []dom_Node{
		&nodes[11], &nodes[7], &nodes[0],
	})
}
