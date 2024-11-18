// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#pragma once
#include <memory>

namespace yw::dom {

class Slot;

// https://dom.spec.whatwg.org/#light-tree-slotables
class Slottable {
public:
    [[nodiscard]] bool is_assigned() const;

private:
    std::shared_ptr<Slot> m_assigned_slot;
};

} // namespace yw::dom