// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause

#pragma once
#include "sourcelocation.hh"
#include <ostream>

namespace yw {

enum class Log_Tag { DEBUG, TODO };

std::ostream& log_stream_impl(Source_Location location, Log_Tag tag);

#define LOG_TODO                                                               \
    ::yw::log_stream_impl(CURRENT_SOURCE_LOCATION, ::yw::Log_Tag::TODO)
#define LOG_DEBUG                                                              \
    ::yw::log_stream_impl(CURRENT_SOURCE_LOCATION, ::yw::Log_Tag::DEBUG)

} // namespace yw