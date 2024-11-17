// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#pragma once
#include "../../dom/node.hh"
#include <memory>

namespace yw::dom {

void dump_node(std::shared_ptr<Node const> const& node, int indent = 0);

} // namespace yw::dom
