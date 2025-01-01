--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
if _G["yw.dom.NodeIterator.nodeIterators"] == nil then
    _G["yw.dom.NodeIterator.nodeIterators"] = {}
end

local nodeIterators = _G["yw.dom.NodeIterator.nodeIterators"]

---https://dom.spec.whatwg.org/#nodeiterator
---@class DOM_NodeIterator
local DOM_NodeIterator = {}

---@return table
function DOM_NodeIterator:nodeIterators()
    return nodeIterators
end

return DOM_NodeIterator