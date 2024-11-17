// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#pragma once

#include <vector>
namespace yw::dom {

class Range {
public:
    static std::vector<Range> &_live_ranges();
private:
};

} // namespace yw::dom
