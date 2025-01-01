--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local Node   = require "yw.dom.Node"
local object = require "yw.common.object"


---https://dom.spec.whatwg.org/#characterdata
---@class DOM_CharacterData : DOM_Node
---@field data       string
local CharacterData = object.create(Node)

---@param document DOM_Document
---@param data string
---@return DOM_CharacterData
function CharacterData:new(document, data)
    local o = Node.new(self, document) --[[@as DOM_CharacterData]]
    o.isCharacterData = true
    o.data = data
    return o
end

return CharacterData
