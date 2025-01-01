--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local Node   = require "yw.dom.Node"
local object = require "yw.common.object"


---https://dom.spec.whatwg.org/#attr
---@class DOM_Attr : DOM_Node
---@field namespace       string?      https://dom.spec.whatwg.org/#concept-attribute-namespace
---@field namespacePrefix string?      https://dom.spec.whatwg.org/#concept-attribute-namespace-prefix
---@field localName       string       https://dom.spec.whatwg.org/#concept-attribute-local-name
---@field value           string       https://dom.spec.whatwg.org/#concept-attribute-value
---@field element         DOM_Element? https://dom.spec.whatwg.org/#concept-attribute-element
local Attr = object.create(Node)

---@param document DOM_Document
---@param namespace string?
---@param namespacePrefix string?
---@param localName string
---@param value string
---@param element DOM_Element?
---@return DOM_Attr
function Attr:new(document, namespace, namespacePrefix, localName, value, element)
    local o = Node.new(self, document) --[[@as DOM_Attr]]
    o.isAttr = true
    o.namespace = namespace
    o.namespacePrefix = namespacePrefix
    o.localName = localName
    o.value = value
    o.element = element
    return o
end

---https://dom.spec.whatwg.org/#concept-attribute-qualified-name
---@return string
function Attr:qualifiedName()
    -- An attribute’s qualified name is its local name if its namespace prefix is null,
    if self.namespacePrefix == nil then
        return self.localName;
    end
    -- and its namespace prefix, followed by ":", followed by its local name, otherwise.
    return self.namespacePrefix .. ":" .. self.localName;
end

return Attr
