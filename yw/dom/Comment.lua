--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local object        = require "yw.common.object"
local CharacterData = require "yw.dom.CharacterData"


---https://dom.spec.whatwg.org/#comment
---@class DOM_Comment : DOM_CharacterData
local Comment = object.create(CharacterData)

---@param document DOM_Document
---@param data string
---@return DOM_Comment
function Comment:new(document, data)
    local o = CharacterData.new(self, document, data) --[[@as DOM_Comment]]
    o.isComment = true
    return o
end

return Comment
