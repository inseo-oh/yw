// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#pragma once
#include "documentfragment.hh"
#include "node.hh"
#include <memory>
#include <string>

namespace yw::dom {

// https://dom.spec.whatwg.org/#interface-shadowroot
//
// https://developer.mozilla.org/en-US/docs/Web/API/ShadowRoot
class Shadow_Root : public Document_Fragment {
    struct Constructor_Badge { };

public:
    Shadow_Root([[maybe_unused]] Constructor_Badge badge,
        std::string debug_name, std::shared_ptr<Document> const& node_document);

    // https://dom.spec.whatwg.org/#dom-shadowroot-host
    //
    // https://developer.mozilla.org/en-US/docs/Web/API/ShadowRoot/host
    std::shared_ptr<Node> host() const;

    //--------------------------------------------------------------------------
    // YW Internal functions
    //--------------------------------------------------------------------------

    static std::shared_ptr<Shadow_Root> _create(
        std::string debug_name, std::shared_ptr<Document> const& node_document);
};

} // namespace yw::dom
