// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#pragma once
#include <optional>

namespace yw {

template <typename T>
using Error = std::optional<T>;

} // namespace yw
