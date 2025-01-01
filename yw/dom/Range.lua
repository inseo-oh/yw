--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
if _G["yw.dom.Range.liveRanges"] == nil then
    _G["yw.dom.Range.liveRanges"] = {}
end

local liveRanges = _G["yw.dom.Range.liveRanges"]

---https://dom.spec.whatwg.org/#concept-live-range
---@class DOM_Range
local DOM_Range = {}

---@return table
function DOM_Range:liveRanges()
    return liveRanges
end

return DOM_Range