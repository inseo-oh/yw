// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#include "_test/testrun.hh"
#include "dom/_debug/node.hh"
#include "dom/document.hh"
#include "dom/node.hh"
#include <cassert>
#include <iostream>
#include <memory>
#include <string>

namespace yw {

void test()
{
    std::shared_ptr<dom::Document> document
        = dom::Document::_create("Document", dom::Document::Type::HTML,
            dom::Document::Mode::NO_QUIRKS, "application/xhtml+xml");
    std::shared_ptr<dom::Node> html = document->create_element("html");
    document->append_child(html).should_not_fail();

    dom::dump_node(std::static_pointer_cast<dom::Node>(document));

    std::cout << "Hello, world!" << "\n";
    test_run();
}

} // namespace yw

int main()
{
    yw::test();
    return 0;
}