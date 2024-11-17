// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#pragma once

namespace yw {

// This emulates C++20's std::source_location, though constructing such value
// using constructor requires compiler-specific functions, so we use good-old
// __FILE__, __LINE__, __func__, combined with a macro.
// (Also for similar reason we don't have column information, but line number is
// enough for debugging)
struct Source_Location {
    int line;
    char const* file_name;
    char const* function_name;
};

#define CURRENT_SOURCE_LOCATION                                                \
    ::yw::Source_Location                                                      \
    {                                                                          \
        __LINE__, __FILE__, __func__                                           \
    }

} // namespace yw
