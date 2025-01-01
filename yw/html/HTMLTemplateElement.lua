--[[
    Copyright (c) 2025, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local object = require "yw.common.object"
local HTMLElement = require "yw.html.HTMLElement"

---https://html.spec.whatwg.org/multipage/scripting.html#the-template-element
---@class HTML_HTMLTemplateElement : HTML_HTMLElement
---@field templateContents  DOM_DocumentFragment  https://html.spec.whatwg.org/multipage/scripting.html#template-contents
local HTMLTemplateElement = object.create(HTMLElement)

---@param attributes DOM_Attr[]
---@param namespace string?
---@param namespacePrefix string?
---@param localName string
---@param customElementState DOM_Element_CustomElementState
---@param customElementDefinition HTML_CustomElementDefinition?
---@param is string?
---@param document DOM_Document
---@return HTML_HTMLTemplateElement
function HTMLTemplateElement:new(attributes, namespace, namespacePrefix, localName, customElementState, customElementDefinition,
    is, document)
    local o = HTMLElement.new(self, attributes, namespace, namespacePrefix, localName, customElementState,
        customElementDefinition, is, document) --[[@as HTML_HTMLTemplateElement]]
    
    -- When a template element is created, the user agent must run the following steps to establish the template contents:
    error("todo")
    -- 1. Let doc be the template element's node document's appropriate template contents owner document.
    -- 2. Create a DocumentFragment object whose node document is doc and host is the template element.
    -- 3. Set the template element's template contents to the newly created DocumentFragment object.

    return o
end


return HTMLTemplateElement
