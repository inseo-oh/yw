// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#include "shadowroot.hh"
#include "documentfragment.hh"
#include "node.hh"
#include <memory>
#include <string>
#include <utility>

namespace yw::dom {

Shadow_Root::Shadow_Root([[maybe_unused]] Constructor_Badge badge,
    std::string debug_name, std::shared_ptr<Document> const& node_document)
    : Document_Fragment(Document_Fragment::_constructor_badge(),
          std::move(debug_name), node_document)
{
}

std::shared_ptr<Node> Shadow_Root::host() const
{
    return _host();
}

std::shared_ptr<Shadow_Root> Shadow_Root::_create(
    std::string debug_name, std::shared_ptr<Document> const& node_document)
{
    std::shared_ptr<Shadow_Root> node = std::make_shared<Shadow_Root>(
        Constructor_Badge {}, debug_name, node_document);
    return node;
}

} // namespace yw::dom
