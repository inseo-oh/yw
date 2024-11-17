// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#pragma once
#include "element.hh"
#include "node.hh"
#include <memory>
#include <optional>
#include <string>

namespace yw::dom {

// https://dom.spec.whatwg.org/#interface-document
//
// https://developer.mozilla.org/en-US/docs/Web/API/Document
class Document : public Node {
    struct Constructor_Badge { };

public:
    enum class Type { XML, HTML };
    enum class Mode { NO_QUIRKS, QUIRKS, LIMITED_QUIRKS };

    Document(Constructor_Badge, std::string const& debug_name, Type type,
        Mode mode, std::string content_type);

    // https://dom.spec.whatwg.org/#concept-document-type
    Type type() const;

    // https://dom.spec.whatwg.org/#concept-document-mode
    Mode mode() const;

    // https://dom.spec.whatwg.org/#concept-create-element
    std::shared_ptr<Element> create_element(std::string const &local_name,
        std::optional<std::string> namespace_,
        std::optional<std::string> prefix = {},
        std::optional<std::string> is = {},
        bool synchronous_custom_elements = false);

    //--------------------------------------------------------------------------
    // Functions with associated MDN docs
    //--------------------------------------------------------------------------

    // https://dom.spec.whatwg.org/#dom-document-createelement
    //
    // https://developer.mozilla.org/en-US/docs/Web/API/Document/createElement
    std::shared_ptr<Element> create_element(
        std::string local_name /* TODO: { is }*/);

    //--------------------------------------------------------------------------
    // YW Internal functions
    //--------------------------------------------------------------------------

    static std::shared_ptr<Document> _create(std::string debug_name,
        Type type = Type::XML, Mode mode = Mode::NO_QUIRKS,
        std::string content_type = "application/xml");

protected:
    static Constructor_Badge _constructor_badge();

private:
    Type m_type;
    Mode m_mode;
    std::string m_content_type;
};

} // namespace yw::dom
