--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local object = require "yw.common.object"
local strings = require "yw.common.strings"

---@class SourceCode
---@field lines           string[]
---@field file            string
---@field line            integer
---@field column          integer
---@field charCount       integer
local SourceCode = {}

function SourceCode:resetCursor()
    self.line = 1
    self.column = 1
end

function SourceCode:saveCursor()
    return { line = self.line, column = self.column }
end

function SourceCode:restoreCursor(cursor)
    self.line = cursor.line
    self.column = cursor.column
end

---@param sourceText string
---@return SourceCode
function SourceCode:new(sourceText)
    local o = object.create(self)
    o.lines = {}
    o.charCount = 0
    for line in string.gmatch(strings.normalizeNewlines(sourceText), "([^\n]+)") do
        table.insert(o.lines, line)
        o.charCount = o.charCount + utf8.len(line)
    end
    o.line = 1
    o.column = 1
    return o
end

---@return integer line
---@return integer column
function SourceCode:nextCursorPos()
    return self.line, self.column
end

---@return string result
---@return integer line
---@return integer column
local function peekInternal(self, count)
    if count == nil then count = 1 end
    local result = ""
    local lineNum = self.line
    local currentLine = self.lines[lineNum]
    local startCol = self.column
    local endCol = self.column
    if currentLine == nil then
        return "", lineNum, endCol
    end
    local lineCharCount = utf8.len(currentLine)
    local emitEol = false
    for n = 1, count do
        if emitEol then
            result = result .. "\n"
            emitEol = false
        else
            if currentLine == nil then
                break
            end
            local eol = lineCharCount <= endCol
            endCol = endCol + 1
            if eol or n == count then
                local finalEndCol = endCol
                if not eol and n == count then
                    finalEndCol = finalEndCol - 1
                end
                result = result ..
                    string.sub(currentLine, utf8.offset(currentLine, startCol), utf8.offset(currentLine, finalEndCol))
            end
            if eol then
                if lineNum ~= #self.lines then
                    emitEol = true
                end
                startCol = 1
                endCol = 0
                lineNum = lineNum + 1
                currentLine = self.lines[lineNum]
                if currentLine ~= nil then
                    lineCharCount = utf8.len(currentLine)
                end
            end
        end
    end
    return result, lineNum, endCol
end

---@param count number?
---@return string
function SourceCode:consume(count)
    if count == nil then count = 1 end
    local ret, l, c = peekInternal(self, count)
    self.line = l
    self.column = c
    return ret

end

---@param count number?
---@return string
function SourceCode:peek(count)
    if count == nil then count = 1 end
    local ret = peekInternal(self, count)
    return ret
end

---@return string
function SourceCode:peekAll()
    local ret = peekInternal(self, self.charCount)
    return ret
end

---comment
---@param logger Logger
---@param loglevel LogLevel
---@param line integer
function SourceCode:printLine(logger, loglevel, line)
    logger:log(loglevel, string.format("% 4d | %s", line, self.lines[line]))
end

---@param logger Logger
---@param loglevel LogLevel
---@param startLine integer
---@param startCol integer
---@param endLine integer
---@param endCol integer
function SourceCode:printRange(logger, loglevel, startLine, startCol, endLine, endCol)
    for l = startLine, endLine do
        self:printLine(logger, loglevel, l)
        local firstCol, lastCol = 1, utf8.len(self.lines[l])
        if l == startLine then
            firstCol = startCol
        end
        if l == endLine then
            lastCol = endCol
        end
        local colCount = lastCol - firstCol + 1
        -- FIXME: Curently this assumes line number is printed with "% 4d" format,
        -- which is enough for right now, but that won't be the case if we start handling longer source files.
        local underlineStr = string.rep(" ", 7 + (firstCol - 1), "")
        for _ = 1, colCount do
            underlineStr = underlineStr .. "~"
        end
        logger:log(loglevel, underlineStr)
    end
end

return SourceCode
