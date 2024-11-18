// Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
//
// SPDX-License-Identifier: BSD-3-Clause
#include "slottable.hh"

namespace yw::dom {

// https://dom.spec.whatwg.org/#light-tree-slotables
bool Slottable::Slottable::is_assigned() const {
    return !!m_assigned_slot;
};

} // namespace yw::dom