// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#include "documentfragment.hh"
#include "node.hh"

#include <memory>
#include <string>
#include <utility>

namespace yw::dom {

Document_Fragment::Document_Fragment([[maybe_unused]] Constructor_Badge badge,
    std::string debug_name, std::shared_ptr<Document> const& node_document)
    : Node(Node::_constructor_badge(), std::move(debug_name),
          Node::Type::DOCUMENT_FRAGMENT, node_document)
{
}

Document_Fragment::Constructor_Badge Document_Fragment::_constructor_badge()
{
    return {};
}

std::shared_ptr<Node> Document_Fragment::_host() const
{
    return m_host.lock();
}

void Document_Fragment::_set_host(std::shared_ptr<Node>& host)
{
    m_host = host;
}

std::shared_ptr<Document_Fragment> Document_Fragment::_create(
    std::string debug_name, std::shared_ptr<Document> const& node_document)
{
    std::shared_ptr<Document_Fragment> node
        = std::make_shared<Document_Fragment>(
            Constructor_Badge {}, debug_name, node_document);
    return node;
}

} // namespace yw::dom
