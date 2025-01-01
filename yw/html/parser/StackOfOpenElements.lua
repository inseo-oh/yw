--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local object = require "yw.common.object"

---@class HTML_Parser_StackOfOpenElements
---@field elements DOM_Element[]
local StackOfOpenElements = {}

function StackOfOpenElements:new()
    local o = object.create(StackOfOpenElements)
    o.elements = {}
    return o
end

---@param e DOM_Element
function StackOfOpenElements:push(e)
    table.insert(self.elements, e)
end

---@return DOM_Element
function StackOfOpenElements:pop()
    return table.remove(self.elements, #self.elements)
end

---@param has DOM_Element
---@return boolean
function StackOfOpenElements:has(has)
    for _, e in ipairs(self.elements) do
        if e == has then
            return true
        end
    end
    return false
end

---https://html.spec.whatwg.org/multipage/parsing.html#current-node
function StackOfOpenElements:currentNode()
    -- The current node is the bottommost node in this stack of open elements.
    return self.elements[#self.elements]
end

---@param localName string
---@return boolean
function StackOfOpenElements:hasHTMLElement(localName)
    for _, e in ipairs(self.elements) do
        if e:isHTMLElement(localName) then
            return true
        end
    end
    return false
end

