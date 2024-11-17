// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#include "range.hh"
#include <vector>

namespace yw::dom {

std::vector<Range>& Range::_live_ranges()
{
    static std::vector<Range> live_ranges {};

    return live_ranges;
}

} // namespace yw::dom
