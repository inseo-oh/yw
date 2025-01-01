--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local DOM_Document = require "yw.dom.Document"
local object       = require "yw.common.object"

-- https://html.spec.whatwg.org/multipage/dom.html#document
---@class HTML_Document : DOM_Document
---@field parserCannotChangeTheModeFlag         boolean   https://html.spec.whatwg.org/multipage/parsing.html#parser-cannot-change-the-mode-flag
---@field throwOnDynamicMarkupInsertionCounter  integer   https://html.spec.whatwg.org/multipage/dynamic-markup-insertion.html#throw-on-dynamic-markup-insertion-counter
---@field isIFrameSrcdocDocument                boolean
local Document     = object.create(DOM_Document)

---@return HTML_Document
function Document:new()
    local o = DOM_Document.new(self) --[[@as HTML_Document]]
    o.type = "html"
    o.contentType = "text/html"
    o.parserCannotChangeTheModeFlag = false
    o.isIFrameSrcdocDocument = false
    o.throwOnDynamicMarkupInsertionCounter = 0

    return o
end
