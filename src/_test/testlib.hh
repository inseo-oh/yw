// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#pragma once
#include "../_utility/sourcelocation.hh"

#include <exception>
#include <functional>
#include <vector>

namespace yw {

class TestFailedException : public std::exception { };

void test_assert_impl(bool x, char const* test, Source_Location location);

// NOLINTNEXTLINE(cppcoreguidelines-macro-usage)
#define TEST_ASSERT(_x) test_assert_impl((_x), #_x, CURRENT_SOURCE_LOCATION)

class TestManager {
public:
    void register_test(char const* name, std::function<void()> const& test);
    void run_tests();

private:
    struct Test {
        char const* name;
        std::function<void()> entry;
    };

    std::vector<Test> m_tests;
};

}; // namespace yw
