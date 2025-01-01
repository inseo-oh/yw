--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]

local object         = require "yw.common.object"
local SourceLocation = require "yw.common.SourceLocation"

---@alias HTML_Parser_TokenType  "eof"|"character"|"tag"|"comment"|"doctype"

---@alias HTML_Parser_Token_Union  HTML_Parser_EofToken|HTML_Parser_CharacterToken|HTML_Parser_TagToken|HTML_Parser_CommentToken|HTML_Parser_DoctypeToken

---@class HTML_Parser_Token
---@field type          HTML_Parser_TokenType
---@field startLocation SourceLocation
---@field endLocation   SourceLocation
local Token          = {}

---@param type HTML_Parser_TokenType
---@param startLocation SourceLocation
---@param endLocation   SourceLocation?
---@return HTML_Parser_Token
function Token:new(type, startLocation, endLocation)
    local o = object.create(self)
    o.type = type
    o.startLocation = startLocation
    o.endLocation = endLocation or SourceLocation:new("", 0, 0) -- Placeholder
    return o
end

--------------------------------------------------------------------------------

---@class HTML_Parser_EofToken : HTML_Parser_Token
local EofToken = object.create(Token)

---@param startLocation SourceLocation
---@param endLocation   SourceLocation?
---@return HTML_Parser_EofToken
function EofToken:new(startLocation, endLocation)
    return Token.new(self, "eof", startLocation, endLocation) --[[@as HTML_Parser_EofToken]]
end

function EofToken:toDebugHTML()
    return ""
end

--------------------------------------------------------------------------------

---@class HTML_Parser_CharacterToken : HTML_Parser_Token
---@field char integer
local CharacterToken = object.create(Token)

---@param char integer
---@param startLocation SourceLocation
---@param endLocation SourceLocation?
---@return HTML_Parser_CharacterToken
function CharacterToken:new(char, startLocation, endLocation)
    local o = Token.new(self, "character", startLocation, endLocation) --[[@as HTML_Parser_CharacterToken]]
    o.char = char
    return o
end

function CharacterToken:toDebugHTML()
    return tostring(utf8.char(self.char))
end

--------------------------------------------------------------------------------

---@alias HTML_Parser_TagKind  "start"|"end"
---@alias HTML_Parser_TagAttr  {startLocation: SourceLocation, endLocation: SourceLocation, name: string, value: string}

---@class HTML_Parser_TagToken : HTML_Parser_Token
---@field name                    string
---@field kind                    HTML_Parser_TagKind
---@field attributes              HTML_Parser_TagAttr[]
---@field selfClosing             boolean
---@field selfClosingAcknowledged boolean
local TagToken = object.create(Token)

---@param startLocation SourceLocation
---@param endLocation SourceLocation?
---@return HTML_Parser_TagAttr
function TagToken.makeAttr(startLocation, endLocation)
    return {
        name = "",
        value = "",
        startLocation = startLocation,
        endLocation = endLocation or SourceLocation:new("", 0, 0)
    }
end

---@param name string
---@param kind HTML_Parser_TagKind
---@return HTML_Parser_TagToken
---@param startLocation SourceLocation
---@param endLocation SourceLocation?
function TagToken:new(name, kind, startLocation, endLocation)
    local o = Token.new(self, "tag", startLocation, endLocation) --[[@as HTML_Parser_TagToken]]
    self.name = name
    self.kind = kind
    self.selfClosing = false
    return o
end

---https://html.spec.whatwg.org/multipage/parsing.html#acknowledge-self-closing-flag
function TagToken:acknowledgeSelfClosingTag()
    if not self.selfClosing then
        return
    end
    self.selfClosingAcknowledged = true
end

---@param name string
---@return string|nil
function TagToken:getAttribute(name)
    for _, attr in ipairs(self.attributes) do
        if attr.name == name then
            return attr.value
        end
    end
    return nil
end

function TagToken:toDebugHTML()
    local attrStr = ""

    for n, attr in ipairs(self.attributes) do
        attrStr = attrStr .. tostring(attr.name)

        if #attr.value ~= 0 then
            attrStr = attrStr .. "="
            if string.find(attr.value, "\"") then
                attrStr = attrStr .. "\'" .. tostring(attr.value) .. "\'"
            else
                attrStr = attrStr .. "\"" .. tostring(attr.value) .. "\""
            end
        end
        if n ~= #self.attributes then
            attrStr = attrStr .. " "
        end
    end
    local tagNameAndAttrs
    if #attrStr ~= 0 then
        tagNameAndAttrs = self.name .. " " .. attrStr
    else
        tagNameAndAttrs = self.name
    end
    if self.kind == "end" then
        return string.format("</%s>", tagNameAndAttrs)
    end
    return string.format("<%s>", tagNameAndAttrs)
end

--------------------------------------------------------------------------------

---@class HTML_Parser_CommentToken : HTML_Parser_Token
---@field data string
local CommentToken = object.create(Token)

---@param data string
---@param startLocation SourceLocation
---@param endLocation SourceLocation?
---@return HTML_Parser_CommentToken
function CommentToken:new(data, startLocation, endLocation)
    local o = Token.new(self, "comment", startLocation, endLocation) --[[@as HTML_Parser_CommentToken]]
    o.data = data
    return o
end

function CommentToken:toDebugHTML()
    return string.format("<!-- %s -->", tostring(self.data));
end

--------------------------------------------------------------------------------

---@class HTML_Parser_DoctypeToken : HTML_Parser_Token
---@field name             string?
---@field publicIdentifier string?
---@field systemIdentifier string?
---@field forceQuirks      boolean
local DoctypeToken = object.create(Token)

---@param startLocation SourceLocation
---@param endLocation SourceLocation?
---@return HTML_Parser_DoctypeToken
function DoctypeToken:new(startLocation, endLocation)
    local o = Token.new(self, "comment", startLocation, endLocation) --[[@as HTML_Parser_DoctypeToken]]
    o.name = nil
    o.publicIdentifier = nil
    o.systemIdentifier = nil
    o.forceQuirks = false
    return o
end

function DoctypeToken:toDebugHTML()
    -- TODO: Print Public and System ID if present
    return string.format("<!DOCTYPE %s>", tostring(self.name));
end

return {
    EofToken       = EofToken,
    CharacterToken = CharacterToken,
    TagToken       = TagToken,
    CommentToken   = CommentToken,
    DoctypeToken   = DoctypeToken,
}
