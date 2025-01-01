--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local Node   = require "yw.dom.Node"
local object = require "yw.common.object"


---https://dom.spec.whatwg.org/#concept-doctype
---@class DOM_DocumentType : DOM_Node
---@field name string      https://dom.spec.whatwg.org/#concept-doctype-name
---@field publicID string  https://dom.spec.whatwg.org/#concept-doctype-publicid
---@field systemID string  https://dom.spec.whatwg.org/#concept-doctype-systemid
local DocumentType = object.create(Node)

---@param document DOM_Document
---@param name string
---@param publicID string?
---@param systemID string?
---@return DOM_DocumentType
function DocumentType:new(document, name, publicID, systemID)
    local o = Node.new(self, document) --[[@as DOM_DocumentType]]
    o.isDocumentType = true
    -- When a doctype is created, its name is always given.
    o.name = name
    -- Unless explicitly given when a doctype is created, its public ID and system ID are the empty string.
    o.publicID = publicID or ""
    o.systemID = systemID or ""
    return o
end

return DocumentType
