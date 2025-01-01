--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local object           = require "yw.common.object"
local DocumentFragment = require "yw.dom.DocumentFragment"


---https://dom.spec.whatwg.org/#concept-shadow-root
---@class DOM_ShadowRoot : DOM_DocumentFragment
---@field declarative                 boolean  https://dom.spec.whatwg.org/#shadowroot-declarative
---@field availableToElementInternals boolean  https://dom.spec.whatwg.org/#shadowroot-available-to-element-internals
local ShadowRoot = object.create(DocumentFragment)

---@param document DOM_Document
---@return DOM_ShadowRoot
function ShadowRoot:new(document)
    local o = DocumentFragment.new(self, document) --[[@as DOM_ShadowRoot]]
    o.isShadowRoot = true
    -- Shadow roots have an associated declarative (a boolean). It is initially set to false.
    o.declarative = false
    -- Shadow roots have an associated available to element internals. It is initially set to false.
    o.availableToElementInternals = false
    return o
end

return ShadowRoot
