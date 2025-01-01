--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local Node   = require "yw.dom.Node"
local object = require "yw.common.object"
local CharacterData = require "yw.dom.CharacterData"


---https://dom.spec.whatwg.org/#text
---@class DOM_Text : DOM_CharacterData
local Text = object.create(CharacterData)

---@param document DOM_Document
---@param data string
---@return DOM_Text
function Text:new(document, data)
    local o = CharacterData.new(self, document, data) --[[@as DOM_Text]]
    o.isText = true
    return o
end

return Text
