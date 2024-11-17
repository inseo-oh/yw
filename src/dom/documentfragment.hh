// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#pragma once
#include "node.hh"
#include <memory>
#include <string>

namespace yw::dom {

// https://dom.spec.whatwg.org/#interface-documentfragment
//
// https://developer.mozilla.org/en-US/docs/Web/API/DocumentFragment
class Document_Fragment : public Node {
    struct Constructor_Badge { };

public:
    Document_Fragment(Constructor_Badge, std::string debug_name, std::shared_ptr<Document> const &node_document);

    //--------------------------------------------------------------------------
    // YW Internal functions
    //--------------------------------------------------------------------------

    static std::shared_ptr<Document_Fragment> _create(std::string debug_name, std::shared_ptr<Document> const &node_document);
    std::shared_ptr<Node> _host() const;
    void _set_host(std::shared_ptr<Node>& host);

protected:
    static Constructor_Badge _constructor_badge();

private:
    std::weak_ptr<Node> m_host;
};

} // namespace yw::dom
