--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local Node       = require "yw.dom.Node"
local object     = require "yw.common.object"
local namespaces = require "yw.common.namespaces"
local Element    = require "yw.dom.Element"
local strings    = require "yw.common.strings"


---@alias DOM_Document_Type "xml"|"html"                           https://dom.spec.whatwg.org/#concept-document-type
---@alias DOM_Document_Mode "no-quirks"|"quirks"|"limited-quirks"  https://dom.spec.whatwg.org/#concept-document-mode

---https://dom.spec.whatwg.org/#concept-document
---@class DOM_Document : DOM_Node
---@field type                          DOM_Document_Type  https://dom.spec.whatwg.org/#concept-document-type
---@field mode                          DOM_Document_Mode  https://dom.spec.whatwg.org/#concept-document-mode
---@field contentType                   string             https://dom.spec.whatwg.org/#concept-document-content-type
---@field allowDeclarativeShadowRoots   boolean            https://dom.spec.whatwg.org/#document-allow-declarative-shadow-roots
---
local Document = object.create(Node)

---@return DOM_Document
function Document:new()
    local o = Node.new(self, self) --[[@as DOM_Document]]
    o.isDocument = true
    -- Unless stated otherwise, a document’s encoding is the utf-8 encoding,
    -- -> TODO
    -- content type is "application/xml",
    o.contentType = "application/xml"
    -- URL is "about:blank",
    -- -> TODO
    -- origin is an opaque origin,
    -- -> TODO
    -- type is "xml",
    o.type = "xml"
    -- mode is "no-quirks",
    o.mode = "no-quirks"
    -- and its allow declarative shadow roots is false.
    o.allowDeclarativeShadowRoots = false
    return o
end

---https://dom.spec.whatwg.org/#xml-document
---@return boolean
function Document:isXMLDocument()
    return self.type == "xml"
end

---https://dom.spec.whatwg.org/#html-document
---@return boolean
function Document:isHTMLDocument()
    return self.type ~= "xml"
end

---https://dom.spec.whatwg.org/#dom-document-createelement
---@param localName string
---@param options table<string, any>|nil
---@return DOM_Element
function Document:createElement(localName, options)
    -- 1. If localName does not match the Name production, then throw an "InvalidCharacterError" DOMException.
    -- TODO

    -- 2. If this is an HTML document, then set localName to localName in ASCII lowercase.
    if self:isHTMLDocument() then
        localName = strings.asciiLowercase(localName)
    end

    -- 3. Let is be null.
    local is = nil

    -- 4. If options is a dictionary and options["is"] exists, then set is to it.
    if options ~= nil then
        error("TODO")
    end

    -- 5. Let namespace be the HTML namespace, if this is an HTML document or this’s content type is "application/xhtml+xml"; otherwise null.
    local namespace = nil
    if (self.type == "html") or (self.contentType == "application/xhtml+xml") then
        namespace = namespaces.HTML_NAMESPACE
    end

    -- 6. Return the result of creating an element given
    return Element.create(
    -- this,
        self,
        -- localName,
        localName,
        -- namespace,
        namespace,
        -- null,
        nil,
        --is,
        is,
        --and with the synchronous custom elements flag set.
        true
    )
end

return Document
