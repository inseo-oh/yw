// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#include "testlib.hh"
#include "../_utility/sourcelocation.hh"

#include <cstddef>
#include <functional>
#include <iostream>

namespace yw {

void test_assert_impl(bool x, char const* test, Source_Location location)
{
    if (!x) {
        std::cerr << "Test assertion \"" << test
                  << "\"\x1b[31;1mFAILED\x1b[0m at " << location.file_name
                  << ":" << location.line << "(" << location.function_name
                  << ")\n";
        throw TestFailedException();
    }
}

void TestManager::register_test(
    char const* name, std::function<void()> const& test)
{
    m_tests.push_back({ name, test });
}

void TestManager::run_tests()
{
    size_t failed_count = 0;
    for (Test const& test : m_tests) {
        try {
            test.entry();
        } catch (TestFailedException const&) {
            std::cerr << "Test " << test.name << " failed\n";
            failed_count++;
        }
    }
    size_t const total_count = m_tests.size();
    if (failed_count != 0) {
        std::cerr << "Ran " << total_count << " tests, \x1b[31;1m"
                  << failed_count << " FAILED\x1b[0m\n";
    } else {
        std::cerr << "Ran " << total_count
                  << " tests, \x1b[32;1mall OK\x1b[0m\n";
    }
}

}; // namespace yw
