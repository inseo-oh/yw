--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local Logger = require "yw.common.Logger"
local SourceCode = require "yw.common.SourceCode"
local codepoints = require "yw.common.codepoints"
local strings    = require "yw.common.strings"

local json = {}

local L = Logger:new("yw.json")

---@param sourceText string
---@return table
function json.parse(sourceText)
    local sourceCode = SourceCode:new(sourceText)

    local function reportError(msg, startLine, startCol, endLine, endCol)
        L:e("JSON error at %d:%d: %s", sourceCode.line, sourceCode.column, msg)
        sourceCode:printRange(L, "error", startLine, startCol, endLine, endCol)
        error("attempted to parse invalid JSON")
    end

    local function skipWhitespaces()
        while true do
            local c = sourceCode:peek()
            if c ~= "\t" and c ~= "\n" and c ~= "\r" and c ~= " " then
                break
            end
            sourceCode:consume()
        end
    end

    local consumeValue, consumeObject, consumeArray, consumeNumber, consumeString

    -- NOTE: This function will throw an error if a value cannot be found
    consumeValue = function()
        local savedCur = sourceCode:saveCursor()
        local res
        res = consumeObject()
        if res ~= nil then
            return res
        end
        sourceCode:restoreCursor(savedCur)
        res = consumeArray()
        if res ~= nil then
            return res
        end
        sourceCode:restoreCursor(savedCur)
        res = consumeNumber()
        if res ~= nil then
            return res
        end
        sourceCode:restoreCursor(savedCur)
        res = consumeString()
        if res ~= nil then
            return res
        end
        sourceCode:restoreCursor(savedCur)
        if sourceCode:consume(#"true") == "true" then
            return true
        end
        sourceCode:restoreCursor(savedCur)
        if sourceCode:consume(#"false") == "false" then
            return false
        end
        sourceCode:restoreCursor(savedCur)
        if sourceCode:consume(#"null") == "null" then
            return nil
        end
        sourceCode:restoreCursor(savedCur)
        return reportError("Expected value", sourceCode.line, sourceCode.column, sourceCode.line, sourceCode.column)
    end

    consumeObject = function()
        if sourceCode:peek(1) ~= "{" then
            return nil
        end
        sourceCode:consume()
        local result = {}
        while true do
            skipWhitespaces()
            local startLine = sourceCode.line
            local startCol = sourceCode.column
            local key = consumeString()
            if key == nil then
                return reportError("Expected name string", startLine, startCol, startLine, startCol)
            end
            skipWhitespaces()
            startLine = sourceCode.line
            startCol = sourceCode.column
            if sourceCode:consume() ~= ":" then
                return reportError("Expected colon", startLine, startCol, startLine, startCol)
            end
            skipWhitespaces()
            local value = consumeValue()
            result[key] = value
            skipWhitespaces()
            if sourceCode:peek(1) ~= "," then
                break
            end
            sourceCode:consume()
        end
        local line = sourceCode.line
        local col = sourceCode.column
        if sourceCode:consume() ~= "}" then
            return reportError("Expected }", line, col, line, col)
        end
        return result
    end

    consumeArray = function()
        local result = {}
        if sourceCode:peek(1) ~= "[" then
            return nil
        end
        sourceCode:consume()
        while true do
            local value = consumeValue()
            table.insert(result, value)
            skipWhitespaces()
            if sourceCode:peek(1) ~= "," then
                break
            end
            sourceCode:consume()
            skipWhitespaces()
        end
        local line = sourceCode.line
        local col = sourceCode.column
        if sourceCode:consume() ~= "]" then
            return reportError("Expected }", line, col, line, col)
        end
        return result
    end

    consumeNumber = function()
        local savedCur = sourceCode:saveCursor()
        local numStr = ""
        ------------------------------------------------------------------------
        -- (Optional) Negative sign
        ------------------------------------------------------------------------
        if sourceCode:peek() == "-" then
            sourceCode:consume(1)
            numStr = numStr .. "-"
        end

        local function consumeDigits()
            local found = false
            while true do
                local digit = sourceCode:peek()
                if digit == nil or not codepoints.isAsciiDigit(utf8.codepoint(digit)) then
                    break
                end
                sourceCode:consume()
                numStr = numStr .. string.char(utf8.codepoint(digit))
                found = true
            end
            return found
        end

        local digit = sourceCode:consume()

        ------------------------------------------------------------------------
        -- 0 or digits followed by non-zero digit
        ------------------------------------------------------------------------
        if digit == "0" then
            sourceCode:consume(1)
            numStr = numStr .. "0"
        elseif codepoints.isAsciiDigit(utf8.codepoint(digit)) then
            numStr = numStr .. string.char(utf8.codepoint(digit))
            consumeDigits()
        else
            -- Not a number
            sourceCode:restoreCursor(savedCur)
            return nil
        end

        ------------------------------------------------------------------------
        -- (Optional) Fractional part followed by decimal point
        ------------------------------------------------------------------------
        local line = sourceCode.line
        local col = sourceCode.column
        if sourceCode:peek() == "." then
            numStr = numStr .. "."
            if not consumeDigits() then
                return reportError("Expected digits", line, col, line, col)
            end
        end

        ------------------------------------------------------------------------
        -- (Optional) Exponent
        ------------------------------------------------------------------------
        local exponentChar = sourceCode:peek()
        if exponentChar == "e" or exponentChar == "E" then
            sourceCode:consume()
            numStr = numStr .. "e"
            -- Exponent sign
            local signChar = sourceCode:peek()
            if signChar == "+" then
                sourceCode:consume()
                numStr = numStr .. "+"
            elseif signChar == "-" then
                sourceCode:consume()
                numStr = numStr .. "-"
            end
            -- Exponent digits
            if not consumeDigits() then
                return reportError("Expected e digits", line, col, line, col)
            end
        end

        local n = tonumber(numStr)
        if n == nil then
            error("internal error: attempted to parse " .. numStr .. " as numeric string")
        end
        return n
    end

    consumeString = function()
        if sourceCode:peek(1) ~= "\"" then
            return nil
        end
        sourceCode:consume()
        local result = {}
        while true do
            local startLine = sourceCode.line
            local startCol = sourceCode.column
            local c = sourceCode:consume()
            if c == nil then
                reportError("Unexpected EOF", startLine, startCol, startLine, startCol)
            elseif c == "\"" then
                break
            elseif c == "\\" then
                local endLine = sourceCode.line
                local endCol = sourceCode.column
                -- Escape codes
                local c2 = sourceCode:consume()
                if c2 == "\"" or c2 == "\\" or c2 == "/" then
                    -- \", \\, \/
                    table.insert(result, utf8.codepoint(c2))
                elseif c2 == "b" then
                    -- \b
                    table.insert(result, 0x0008)
                elseif c2 == "f" then
                    -- \f
                    table.insert(result, 0x000c)
                elseif c2 == "n" then
                    -- \n
                    table.insert(result, 0x000a)
                elseif c2 == "r" then
                    -- \r
                    table.insert(result, 0x000d)
                elseif c2 == "t" then
                    -- \t
                    table.insert(result, 0x0009)
                elseif c2 == "u" then
                    -- \u
                    local charcode = 0
                    for _ = 1, 4 do
                        endLine = sourceCode.line
                        endCol = sourceCode.column
                        local c3 = sourceCode:consume()
                        if c3 == nil then
                            return reportError("Unexpected EOF", startLine, startCol, endLine, endCol)
                        end
                        local codepoint = utf8.codepoint(c3)
                        if not codepoints.isAsciiHexDigit(codepoint) then
                            return reportError("Expected hex digit", startLine, startCol, endLine, endCol)
                        end
                        charcode = charcode * 16 + tonumber(string.char(codepoint), 16)
                    end
                    table.insert(result, charcode)
                else
                    return reportError(
                        string.format("Unexpected escape character \\%s", tostring(c2)),
                        startLine, startCol, endLine, endCol
                    )
                end
            else
                table.insert(result, utf8.codepoint(c))
            end
        end
        return strings.fromCharCodes(result)
    end

    skipWhitespaces()
    local line = sourceCode.line
    local col = sourceCode.column
    local o = consumeObject()
    if o == nil then
        reportError("Expected JSON object", line, col, line, col)
        return {} -- Unreachable
    end
    return o
end

return json
