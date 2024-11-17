// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#include "element.hh"
#include "../infra/namespaces.hh"
#include "document.hh"
#include "node.hh"

#include <algorithm>
#include <cctype>
#include <memory>
#include <optional>
#include <string>
#include <utility>

namespace yw::dom {

Element::Element([[maybe_unused]] Constructor_Badge badge,
    std::string debug_name, std::optional<std::string> namespace_,
    std::optional<std::string> namespace_prefix, std::string local_name,
    Custom_Element_State custom_element_state, std::optional<std::string> is,
    std::shared_ptr<Document> const& node_document)
    : Node(Node::_constructor_badge(), std::move(debug_name),
          Node::Type::ELEMENT, node_document)
    , m_custom_element_state(custom_element_state)
    , m_local_name(std::move(local_name))
    , m_namespace(std::move(namespace_))
    , m_namespace_prefix(std::move(namespace_prefix))
    , m_is(std::move(is))
{
}

Element::Constructor_Badge Element::_constructor_badge()
{
    return {};
}

std::shared_ptr<Node> Element::shadow_root()
{
    return m_shadow_root;
}

std::shared_ptr<Node const> Element::shadow_root() const
{
    return m_shadow_root;
}

bool Element::is_shadow_host() const
{
    return !!shadow_root();
}

Element::Custom_Element_State Element::custom_element_state() const
{
    return m_custom_element_state;
}

[[nodiscard]] std::string Element::qualified_name() const
{
    if (!m_namespace_prefix) {
        return m_local_name;
    }
    return m_namespace_prefix.value() + ":" + m_local_name;
}

[[nodiscard]] std::string Element::html_uppercased_qualified_name() const
{
    // https://dom.spec.whatwg.org/#element-html-uppercased-qualified-name

    // 1. Let qualifiedName be this’s qualified name.
    std::string my_qualified_name = qualified_name();

    // 2. If this is in the HTML namespace and its node document is an HTML
    // document, then set qualifiedName to qualifiedName in ASCII uppercase.
    if ((m_namespace == infra::HTML_NAMESPACE)
        && node_document()->type() == Document::Type::HTML) {
        std::transform(my_qualified_name.begin(), my_qualified_name.end(),
            my_qualified_name.begin(), ::toupper);
    }

    // 3. Return qualifiedName.
    return my_qualified_name;
}

[[nodiscard]] std::string Element::tag_name() const
{
    return html_uppercased_qualified_name();
}

[[nodiscard]] bool Element::is_custom() const
{
    return m_custom_element_state == Custom_Element_State::CUSTOM;
}

std::shared_ptr<Element> Element::_create(std::string debug_name,
    std::optional<std::string> namespace_,
    std::optional<std::string> namespace_prefix, std::string local_name,
    Custom_Element_State custom_element_state, std::optional<std::string> is,
    std::shared_ptr<Document> node_document)
{
    std::shared_ptr<Element> node = std::make_shared<Element>(
        Constructor_Badge {}, debug_name, namespace_, namespace_prefix,
        local_name, custom_element_state, is, node_document);
    return node;
}

} // namespace yw::dom
