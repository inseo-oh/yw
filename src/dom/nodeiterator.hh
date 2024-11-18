// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#pragma once

#include <vector>
namespace yw::dom {

class Node_Iterator {
public:
    static std::vector<Node_Iterator> &_node_iterators();
private:
};

} // namespace yw::dom
