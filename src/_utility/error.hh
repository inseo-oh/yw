// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#pragma once
#include <optional>

namespace yw {

template <typename T>
class [[nodiscard]] Error: public std::optional<T> {
public:
    Error() = default;
    explicit constexpr Error(T const &val): std::optional<T>(val) {}
};

} // namespace yw
