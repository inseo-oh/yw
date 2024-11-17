// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#include "../../dom/node.hh"
#include <iostream>
#include <memory>
#include <string>

namespace yw::dom {

namespace {

    std::string node_type_string(Node::Type type)
    {
        switch (type) {
        case Node::Type::ELEMENT:
            return "Element";
        case Node::Type::ATTRIBUTE:
            return "Attr";
        case Node::Type::TEXT:
            return "Text";
        case Node::Type::CDATA_SECTION:
            return "CDATASection";
        case Node::Type::PROCESSING_INSTRUCTION:
            return "ProcessingInstruction";
        case Node::Type::COMMENT:
            return "Comment";
        case Node::Type::DOCUMENT:
            return "Document";
        case Node::Type::DOCUMENT_TYPE:
            return "DocumentType";
        case Node::Type::DOCUMENT_FRAGMENT:
            return "DocumentFragment";
        }
    }

} // namespace

// NOLINTNEXTLINE(misc-no-recursion)
void dump_node(std::shared_ptr<Node const> const& node, int indent)
{
    for (int i = 0; i < indent; i++) {
        std::cout << ' ';
    }
    std::cout << node->_debug_name() << "("
              << node_type_string(node->node_type()) << ")" << '\n';
    std::shared_ptr<Node> child = node->first_child();
    while (child) {
        dump_node(child, indent + 2);
        child = child->next_sibling();
    }
}

} // namespace yw::dom
