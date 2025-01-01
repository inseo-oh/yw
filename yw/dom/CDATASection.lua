--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local object = require "yw.common.object"
local Text = require "yw.dom.Text"


---https://dom.spec.whatwg.org/#cdatasection
---@class DOM_CDATASection : DOM_Text
local CDATASection = object.create(Text)

---@param document DOM_Document
---@param data string
---@return DOM_Text
function CDATASection:new(document, data)
    local o = Text.new(self, document, data) --[[@as DOM_Text]]
    o.isCDATASection = true
    return o
end

return CDATASection
