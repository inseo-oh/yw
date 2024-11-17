// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#pragma once

#define YW_NON_COPYABLE(_n)                                                    \
    _n(_n const&) = delete;                                                    \
    _n& operator=(_n const&) = delete;

#define YW_NON_MOVEABLE(_n)                                                    \
    _n(_n&&) = delete;                                                         \
    _n& operator=(_n&&) = delete;
