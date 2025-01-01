--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
    ]]
local ioutil = {}

---@param L Logger
---@param path string
---@return string
function ioutil.readFile(L, path)
    L:d("Reading file " .. path)
    local file = io.open(path, "rb")
    if file == nil then
        error("failed to open file" .. path)
    end
    assert(file ~= nil)
    local text = file:read("*a")
    file:close()
    L:d("Reading file " .. path .. " complete(" .. #text .. " bytes)")
    return text
end

return ioutil
