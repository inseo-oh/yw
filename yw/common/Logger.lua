local object = require "yw.common.object"
--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]

---@class Logger
---@field tag string
local Logger = {}

---@param tag string
---@return Logger
function Logger:new(tag)
    local o = object.create(self)
    o.tag = tag
    return o
end

---@param ... number
---@return string
local function sgr(...)
    local CSI = "\27["
    local SGR_SUFFIX = "m"

    local result = CSI
    local args = { ... }
    for n, arg in ipairs(args) do
        result = result .. arg
        if n ~= #args then
            result = result .. ";"
        end
    end
    result = result .. SGR_SUFFIX
    return result
end

local SGR_RESET      = 0
local SGR_BOLD       = 1
local SGR_FG_RED     = 31
local SGR_FG_GREEN   = 32
local SGR_FG_YELLOW  = 33
local SGR_FG_BLUE    = 34
local SGR_FG_MAGENTA = 35
local SGR_FG_CYAN    = 36
local SGR_FG_WHITE   = 37
local SGR_FG_DEFAULT = 39

---@param ... number SGR attributes to apply to log level
---@return string
local function logFormat(...)
    -- <date> [<tag>] <level>: message <LF>
    return "%s " .. sgr(SGR_FG_WHITE) .. "[%s]" .. sgr(SGR_BOLD, ...) .. " %s" ..
        sgr(SGR_RESET, SGR_FG_DEFAULT) .. ": %s\n"
end

local LOG_FORMATS = {
    info  = logFormat(SGR_FG_GREEN),
    warn  = logFormat(SGR_FG_YELLOW),
    error = logFormat(SGR_FG_RED),
    debug = logFormat(SGR_FG_CYAN),
    trace = logFormat(SGR_FG_MAGENTA),
}

---@alias LogLevel "info"|"warn"|"error"|"debug"|"trace"

---@param level LogLevel
function Logger:log(level, fmt, ...)
    local date = os.date()
    local logFmt = LOG_FORMATS[level]
    local msg = string.format(logFmt, date, self.tag, level, string.format(fmt, ...))
    io.write(msg)
end

function Logger:i(fmt, ...)
    return self:log("info", fmt, ...)
end

function Logger:w(fmt, ...)
    return self:log("warn", fmt, ...)
end

function Logger:e(fmt, ...)
    return self:log("error", fmt, ...)
end

function Logger:d(fmt, ...)
    return self:log("debug", fmt, ...)
end

function Logger:t(fmt, ...)
    return self:log("trace", fmt, ...)
end

return Logger
