// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#pragma once
#include "../_utility/sourcelocation.hh"
#include <string>

namespace yw::idl {

// https://webidl.spec.whatwg.org/#idl-DOMException
//
// https://developer.mozilla.org/en-US/docs/Web/API/DOMException
//
// NOTE: Do not use the constructor directly. Use DOM_EXCEPTION() macro instead.
class DOM_Exception {
public:
    explicit DOM_Exception(Source_Location location, std::string message = "",
        std::string name = "Error");

    // https://webidl.spec.whatwg.org/#hierarchyrequesterror
    static constexpr char const* HIERARCHY_REQUEST_ERROR
        = "HierarchyRequestError";
    // https://webidl.spec.whatwg.org/#notfounderror
    static constexpr char const* NOT_FOUND_ERROR = "NotFoundError";

    //--------------------------------------------------------------------------
    // Functions with associated MDN docs
    //--------------------------------------------------------------------------

    // https://webidl.spec.whatwg.org/#dom-domexception-name
    //
    // https://developer.mozilla.org/en-US/docs/Web/API/DOMException/name
    [[nodiscard]] std::string name() const;

    // https://webidl.spec.whatwg.org/#dom-domexception-message
    //
    // https://developer.mozilla.org/en-US/docs/Web/API/DOMException/message
    [[nodiscard]] std::string message() const;

    //--------------------------------------------------------------------------
    // YW Internal functions
    //--------------------------------------------------------------------------

    [[nodiscard]] Source_Location _origin_location() const;

private:
    Source_Location m_origin_location;
    std::string m_message, m_name;
};

// NOLINTNEXTLINE(cppcoreguidelines-macro-usage)
#define DOM_EXCEPTION(...)                                                     \
    ::yw::idl::DOM_Exception(CURRENT_SOURCE_LOCATION, __VA_ARGS__)

} // namespace yw::idl
