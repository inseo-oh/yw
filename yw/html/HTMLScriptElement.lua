--[[
    Copyright (c) 2025, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local object = require "yw.common.object"
local HTMLElement = require "yw.html.HTMLElement"

---https://html.spec.whatwg.org/multipage/scripting.html#the-script-element
---@class HTML_HTMLScriptElement : HTML_HTMLElement
---@field parserDocument HTML_Document?   https://html.spec.whatwg.org/multipage/scripting.html#parser-document
---@field forceAsync     boolean          https://html.spec.whatwg.org/multipage/scripting.html#script-force-async
---@field alreadyStarted boolean          https://html.spec.whatwg.org/multipage/scripting.html#already-started
local HTMLScriptElement = object.create(HTMLElement)

---@param attributes DOM_Attr[]
---@param namespace string?
---@param namespacePrefix string?
---@param localName string
---@param customElementState DOM_Element_CustomElementState
---@param customElementDefinition HTML_CustomElementDefinition?
---@param is string?
---@param document DOM_Document
---@return HTML_HTMLScriptElement
function HTMLScriptElement:new(attributes, namespace, namespacePrefix, localName, customElementState, customElementDefinition,
    is, document)
    local o = HTMLElement.new(self, attributes, namespace, namespacePrefix, localName, customElementState,
        customElementDefinition, is, document) --[[@as HTML_HTMLScriptElement]]
    o.parserDocument = nil

    return o
end


return HTMLScriptElement
