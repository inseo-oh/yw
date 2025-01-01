--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local Node   = require "yw.dom.Node"
local object = require "yw.common.object"


---https://dom.spec.whatwg.org/#documentfragment
---@class DOM_DocumentFragment : DOM_Node
---@field host DOM_Node
local DocumentFragment = object.create(Node)

---@param document DOM_Document
---@return DOM_DocumentFragment
function DocumentFragment:new(document)
    local o = Node.new(self, document) --[[@as DOM_DocumentFragment]]
    o.isDocumentFragment = true
    o.host = nil
    return o
end

return DocumentFragment
