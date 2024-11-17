// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#include "testlib.hh"
#include "tests.hh"

namespace yw {

void test_run()
{
    TestManager tm;
    tests_init_dom(tm);
    tm.run_tests();
}

}; // namespace yw
