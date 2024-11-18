// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#include "nodeiterator.hh"
#include <vector>

namespace yw::dom {

std::vector<Node_Iterator>& Node_Iterator::_node_iterators()
{
    static std::vector<Node_Iterator> node_iterators {};

    return node_iterators;
}

} // namespace yw::dom
