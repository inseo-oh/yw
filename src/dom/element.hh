// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#pragma once
#include "node.hh"
#include <memory>
#include <optional>
#include <string>

namespace yw::dom {

// https://dom.spec.whatwg.org/#interface-element
//
// https://developer.mozilla.org/en-US/docs/Web/API/Element
class Element : public Node {
    struct Constructor_Badge { };

public:
    enum class Custom_Element_State {
        UNDEFINED,
        FAILED,
        UNCUSTOMIZED,
        PRECUSTOMIZED,
        CUSTOM
    };

    Element(Constructor_Badge, std::string debug_name,
        std::optional<std::string> namespace_,
        std::optional<std::string> namespace_prefix, std::string local_name,
        Custom_Element_State custom_element_state,
        std::optional<std::string> is,
        std::shared_ptr<Document> const& node_document);

    // https://dom.spec.whatwg.org/#concept-element-shadow-root
    [[nodiscard]] std::shared_ptr<Node> shadow_root();
    [[nodiscard]] std::shared_ptr<Node const> shadow_root() const;

    // https://dom.spec.whatwg.org/#element-shadow-host
    [[nodiscard]] bool is_shadow_host() const;

    // https://dom.spec.whatwg.org/#concept-element-custom-element-state
    [[nodiscard]] Custom_Element_State custom_element_state() const;

    // https://dom.spec.whatwg.org/#concept-element-qualified-name
    [[nodiscard]] std::string qualified_name() const;

    // https://dom.spec.whatwg.org/#element-html-uppercased-qualified-name
    [[nodiscard]] std::string html_uppercased_qualified_name() const;

    // https://dom.spec.whatwg.org/#dom-element-tagname
    [[nodiscard]] std::string tag_name() const;

    // https://dom.spec.whatwg.org/#concept-element-custom
    [[nodiscard]] bool is_custom() const;

    //--------------------------------------------------------------------------
    // Functions with associated MDN docs
    //--------------------------------------------------------------------------

    //  https://dom.spec.whatwg.org/#dom-element-tagname
    //
    // https://developer.mozilla.org/en-US/docs/Web/API/Element/tagName

    //--------------------------------------------------------------------------
    // YW Internal functions
    //--------------------------------------------------------------------------

    static std::shared_ptr<Element> _create(std::string debug_name,
        std::optional<std::string> namespace_,
        std::optional<std::string> namespace_prefix, std::string local_name,
        Custom_Element_State custom_element_state,
        std::optional<std::string> is, std::shared_ptr<Document> node_document);

protected:
    static Constructor_Badge _constructor_badge();

private:
    std::shared_ptr<Node> m_shadow_root;
    Custom_Element_State m_custom_element_state;
    std::string m_local_name;
    std::optional<std::string> m_namespace;
    std::optional<std::string> m_namespace_prefix;
    std::optional<std::string> m_is;
};

} // namespace yw::dom
