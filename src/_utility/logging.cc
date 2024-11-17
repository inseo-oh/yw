// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause

#include "logging.hh"
#include "sourcelocation.hh"
#include <iostream>
#include <ostream>

namespace yw {

std::ostream& log_stream_impl(Source_Location location, Log_Tag tag)
{
    std::ostream& stream = std::cerr << "[" << location.file_name << ":"
                                     << location.line << "("
                                     << location.function_name << ")" << "] ";
    switch(tag) {
        case Log_Tag::DEBUG:
            return stream << "\x1b[34;1mDEBUG\x1b[0m: ";
        case Log_Tag::TODO:
            return stream << "\x1b[33;1mTODO\x1b[0m: ";
    }
}

} // namespace yw