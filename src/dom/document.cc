// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#include "document.hh"
#include "../_utility/logging.hh"
#include "../infra/namespaces.hh"
#include "element.hh"
#include "node.hh"

#include <algorithm>
#include <cassert>
#include <cctype>
#include <memory>
#include <optional>
#include <string>
#include <utility>

namespace yw::dom {

Document::Document([[maybe_unused]] Constructor_Badge badge,
    std::string const& debug_name, Type type, Mode mode,
    std::string content_type)
    : Node(Node::_constructor_badge(), debug_name, Node::Type::DOCUMENT, {})
    , m_type(type)
    , m_mode(mode)
    , m_content_type(std::move(content_type))
{
}

Document::Constructor_Badge Document::_constructor_badge()
{
    return {};
}

Document::Type Document::type() const
{
    return m_type;
}

Document::Mode Document::mode() const
{
    return m_mode;
}

std::shared_ptr<Element> Document::create_element(std::string const& local_name,
    std::optional<std::string> namespace_, std::optional<std::string> prefix,
    std::optional<std::string> is, bool synchronous_custom_elements)
{
    (void)synchronous_custom_elements;
    // https://dom.spec.whatwg.org/#concept-create-element

    // 1. If prefix was not given, let prefix be null.
    // 2. If is was not given, let is be null.
    // NOTE: We've done 1 and 2 using C++ default parameters.

    // 3. Let result be null.
    std::shared_ptr<Element> result;

    // 4. Let definition be the result of looking up a custom element definition
    // given document, namespace, localName, and is.
    std::shared_ptr<int /* STUB */> definition;
    LOG_TODO << "Let definition be the result of looking up a custom element "
                "definition given document, namespace, localName, and is.\n";

    // 5. If definition is non-null, and definition’s name is not equal to its
    // local name (i.e., definition represents a customized built-in element),
    // then:
    if (definition) {
        // 1. Let interface be the element interface for localName and the HTML
        // namespace.
        assert(!"TODO");

        // 2. Set result to a new element that implements interface, with no
        // attributes, namespace set to the HTML namespace, namespace prefix set
        // to prefix, local name set to localName, custom element state set to
        // "undefined", custom element definition set to null, is value set to
        // is, and node document set to document.

        // 3. If the synchronous custom elements flag is set, then run this step
        // while catching any exceptions:
        // - 1. Upgrade result using definition.
        // If this step threw an exception exception:
        // 1. Report exception for definition’s constructor’s corresponding
        // JavaScript object’s associated realm’s global object.
        // 2. Set result’s custom element state to "failed".

        // 4. Otherwise, enqueue a custom element upgrade reaction given result
        // and definition.
    }
    // 6. Otherwise, if definition is non-null:
    else if (definition) {
        // 1. If the synchronous custom elements flag is set, then run these
        // steps while catching any exceptions:
        // - 1. Let C be definition’s constructor.
        // - 2. Set result to the result of constructing C, with no arguments.
        // - 3. Assert: result’s custom element state and custom element
        // definition are initialized.
        // - 4. Assert: result’s namespace is the HTML namespace.
        // - 5. If result’s attribute list is not empty, then throw a
        // "NotSupportedError" DOMException.
        // - 6. If result has children, then throw a "NotSupportedError"
        // DOMException.
        // - 7. If result’s parent is not null, then throw a "NotSupportedError"
        // DOMException.
        // - 8. If result’s node document is not document, then throw a
        // "NotSupportedError" DOMException.
        // - 9. If result’s local name is not equal to localName, then throw a
        // "NotSupportedError" DOMException.
        // - 10. Set result’s namespace prefix to prefix.
        // - 11. Set result’s is value to null.
        // - If any of these steps threw an exception exception:
        // -- 1. Report exception for definition’s constructor’s corresponding
        // JavaScript object’s associated realm’s global object.
        // -- 2. Set result to a new element that implements the
        // HTMLUnknownElement interface, with no attributes, namespace set to
        // the HTML namespace, namespace prefix set to prefix, local name set to
        // localName, custom element state set to "failed", custom element
        // definition set to null, is value set to null, and node document set
        // to document.
        // 2. Otherwise:
        // - 1. Set result to a new element that implements the HTMLElement
        // interface, with no attributes, namespace set to the HTML namespace,
        // namespace prefix set to prefix, local name set to localName, custom
        // element state set to "undefined", custom element definition set to
        // null, is value set to null, and node document set to document.
        // - 2. Enqueue a custom element upgrade reaction given result and
        // definition.
    }
    // 7. Otherwise:
    else {
        // 1. Let interface be the element interface for localName and
        // namespace.
        // 2. Set result to a new element that implements interface, with no
        // attributes, namespace set to namespace, namespace prefix set to
        // prefix, local name set to localName, custom element state set to
        // "uncustomized", custom element definition set to null, is value set
        // to is, and node document set to document.
        result = Element::_create("element[" + local_name + "]",
            std::move(namespace_), std::move(prefix), local_name,
            Element::Custom_Element_State::UNCUSTOMIZED,
            /* TODO: custom element definition */ std::move(is),
            node_document());
        // - 3. If namespace is the HTML namespace, and either localName is a
        // valid custom element name or is is non-null, then set result’s custom
        // element state to "undefined".
        LOG_TODO << "If namespace is the HTML namespace, and either localName "
                    "is a valid custom element name or is is non-null, then "
                    "set result’s custom element state to 'undefined'.";
    }
    // 8. Return result.
    assert(result);
    return result;
}

std::shared_ptr<Element> Document::create_element(
    std::string local_name /* TODO: { is }*/)
{
    // https://dom.spec.whatwg.org/#dom-document-createelement

    // 1. If localName does not match the Name production, then throw an
    // "InvalidCharacterError" DOMException.
    // TODO

    // 2. If this is an HTML document, then set localName to localName in ASCII
    // lowercase.
    if (m_type == Type::HTML) {
        std::transform(local_name.begin(), local_name.end(), local_name.begin(),
            ::tolower);
    }

    // 3. Let is be null.
    std::optional<std::string> is = {};

    // 4. If options is a dictionary and options["is"] exists, then set is to
    // it.
    // TODO

    // 5. Let namespace be the HTML namespace, if this is an HTML document or
    // this’s content type is "application/xhtml+xml"; otherwise null.
    std::optional<std::string> namespace_;
    if ((m_type == Type::HTML) || (m_content_type == "application/xhtml+xml")) {
        namespace_ = infra::HTML_NAMESPACE;
    }

    // 6. Return the result of creating an element given this, localName,
    // namespace, null, is, and with the synchronous custom elements flag set.
    return create_element(local_name, namespace_, {}, is, true);
}

std::shared_ptr<Document> Document::_create(
    std::string debug_name, Type type, Mode mode, std::string content_type)
{
    std::shared_ptr<Document> node = std::make_shared<Document>(
        Constructor_Badge {}, debug_name, type, mode, content_type);
    return node;
}

} // namespace yw::dom
