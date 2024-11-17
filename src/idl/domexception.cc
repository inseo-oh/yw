// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#include "domexception.hh"

#include "../_utility/sourcelocation.hh"
#include <string>
#include <utility>

namespace yw::idl {

DOM_Exception::DOM_Exception(
    Source_Location location, std::string message, std::string name)
    : m_origin_location(location)
    , m_message(std::move(message))
    , m_name(std::move(name))
{
}

// https://webidl.spec.whatwg.org/#dom-domexception-name
//
// https://developer.mozilla.org/en-US/docs/Web/API/DOMException/name
[[nodiscard]] std::string DOM_Exception::name() const
{
    return m_name;
}

// https://webidl.spec.whatwg.org/#dom-domexception-message
//
// https://developer.mozilla.org/en-US/docs/Web/API/DOMException/message
[[nodiscard]] std::string DOM_Exception::message() const
{
    return m_message;
}

} // namespace yw::idl
