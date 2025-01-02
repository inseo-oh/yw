--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local Node                 = require "yw.dom.Node"
local namespaces           = require "yw.common.namespaces"
local HTMLElement          = require "yw.html.HTMLElement"
local HTMLAnchorElement    = require "yw.html.HTMLAnchorElement"
local HTMLBodyElement      = require "yw.html.HTMLBodyElement"
local HTMLBRElement        = require "yw.html.HTMLBRElement"
local HTMLButtonElement    = require "yw.html.HTMLButtonElement"
local HTMLDetailsElement   = require "yw.html.HTMLDetailsElement"
local HTMLDialogElement    = require "yw.html.HTMLDialogElement"
local HTMLDivElement       = require "yw.html.HTMLDivElement"
local HTMLDListElement     = require "yw.html.HTMLDListElement"
local HTMLFieldSetElement  = require "yw.html.HTMLFieldSetElement"
local HTMLFontElement      = require "yw.html.HTMLFontElement"
local HTMLFormElement      = require "yw.html.HTMLFormElement"
local HTMLHeadElement      = require "yw.html.HTMLHeadElement"
local HTMLHeadingElement   = require "yw.html.HTMLHeadingElement"
local HTMLHRElement        = require "yw.html.HTMLHRElement"
local HTMLHtmlElement      = require "yw.html.HTMLHtmlElement"
local HTMLImageElement     = require "yw.html.HTMLImageElement"
local HTMLLIElement        = require "yw.html.HTMLLIElement"
local HTMLLinkElement      = require "yw.html.HTMLLinkElement"
local HTMLMenuElement      = require "yw.html.HTMLMenuElement"
local HTMLMetaElement      = require "yw.html.HTMLMetaElement"
local HTMLObjectElement    = require "yw.html.HTMLObjectElement"
local HTMLOListElement     = require "yw.html.HTMLOListElement"
local HTMLParagraphElement = require "yw.html.HTMLParagraphElement"
local HTMLPreElement       = require "yw.html.HTMLPreElement"
local HTMLQuoteElement     = require "yw.html.HTMLQuoteElement"
local HTMLSpanElement      = require "yw.html.HTMLSpanElement"
local HTMLTableElement     = require "yw.html.HTMLTableElement"
local HTMLTitleElement     = require "yw.html.HTMLTitleElement"
local HTMLUListElement     = require "yw.html.HTMLUListElement"
local object               = require "yw.common.object"
local strings              = require "yw.common.strings"


---https://dom.spec.whatwg.org/#concept-element-custom-element-state
---@alias DOM_Element_CustomElementState "undefined"|"failed"|"uncustomized"|"precustomized"|"custom"

---https://dom.spec.whatwg.org/#concept-element
---@class DOM_Element: DOM_Node
---@field namespace               string?                        https://dom.spec.whatwg.org/#concept-element-namespace
---@field namespacePrefix         string?                        https://dom.spec.whatwg.org/#concept-element-namespace-prefix
---@field localName               string                         https://dom.spec.whatwg.org/#concept-element-local-name
---@field customElementState      DOM_Element_CustomElementState https://dom.spec.whatwg.org/#concept-element-custom-element-state
---@field customElementDefinition HTML_CustomElementDefinition?  https://dom.spec.whatwg.org/#concept-element-custom-element-definition
---@field is                      string?                        https://dom.spec.whatwg.org/#concept-element-is-value
---@field attributes              DOM_Attr                       https://dom.spec.whatwg.org/#concept-element-attribute
---@field shadowRoot              DOM_Node?                      https://dom.spec.whatwg.org/#concept-element-shadow-root
local Element = object.create(Node)

---@param attributes DOM_Attr[]
---@param namespace string?
---@param namespacePrefix string?
---@param localName string
---@param customElementState DOM_Element_CustomElementState
---@param customElementDefinition HTML_CustomElementDefinition?
---@param is string?
---@param document DOM_Document
---@return DOM_Element
function Element:new(attributes, namespace, namespacePrefix, localName, customElementState, customElementDefinition, is,
                     document)
    local o = Node.new(self, document) --[[@as DOM_Element]]
    o.isElement = true
    o.shadowRoot = nil
    o.attributes = attributes
    o.namespace = namespace
    o.namespacePrefix = namespacePrefix
    o.localName = localName
    o.customElementState = customElementState
    o.customElementDefinition = customElementDefinition
    o.is = is
    return o
end

local HTML_ELEMENT_INTERFACES = {
    ["a"]          = HTMLAnchorElement,
    ["blockquote"] = HTMLQuoteElement,
    ["body"]       = HTMLBodyElement,
    ["br"]         = HTMLBRElement,
    ["button"]     = HTMLButtonElement,
    ["details"]    = HTMLDetailsElement,
    ["dialog"]     = HTMLDialogElement,
    ["div"]        = HTMLDivElement,
    ["dl"]         = HTMLDListElement,
    ["fieldset"]   = HTMLFieldSetElement,
    ["font"]       = HTMLFontElement,
    ["form"]       = HTMLFormElement,
    ["h1"]         = HTMLHeadingElement,
    ["h2"]         = HTMLHeadingElement,
    ["h3"]         = HTMLHeadingElement,
    ["h4"]         = HTMLHeadingElement,
    ["h5"]         = HTMLHeadingElement,
    ["h6"]         = HTMLHeadingElement,
    ["head"]       = HTMLHeadElement,
    ["hr"]         = HTMLHRElement,
    ["html"]       = HTMLHtmlElement,
    ["img"]        = HTMLImageElement,
    ["li"]         = HTMLLIElement,
    ["link"]       = HTMLLinkElement,
    ["menu"]       = HTMLMenuElement,
    ["meta"]       = HTMLMetaElement,
    ["object"]     = HTMLObjectElement,
    ["ol"]         = HTMLOListElement,
    ["p"]          = HTMLParagraphElement,
    ["pre"]        = HTMLPreElement,
    ["span"]       = HTMLSpanElement,
    ["table"]      = HTMLTableElement,
    ["title"]      = HTMLTitleElement,
    ["ul"]         = HTMLUListElement,
    ----------------------------------------------------------------------------
    ["address"]    = HTMLElement,
    ["article"]    = HTMLElement,
    ["aside"]      = HTMLElement,
    ["b"]          = HTMLElement,
    ["cite"]       = HTMLElement,
    ["code"]       = HTMLElement,
    ["dd"]         = HTMLElement,
    ["dfn"]        = HTMLElement,
    ["dt"]         = HTMLElement,
    ["em"]         = HTMLElement,
    ["figcaption"] = HTMLElement,
    ["figure"]     = HTMLElement,
    ["footer"]     = HTMLElement,
    ["header"]     = HTMLElement,
    ["hgroup"]     = HTMLElement,
    ["i"]          = HTMLElement,
    ["main"]       = HTMLElement,
    ["nav"]        = HTMLElement,
    ["s"]          = HTMLElement,
    ["search"]     = HTMLElement,
    ["section"]    = HTMLElement,
    ["small"]      = HTMLElement,
    ["strong"]     = HTMLElement,
    ["summary"]    = HTMLElement,
    ["u"]          = HTMLElement,
    ["var"]        = HTMLElement,
}

-- https://dom.spec.whatwg.org/#concept-element-interface
---@param localName string
---@param namespace string?
---@return DOM_Element
local function elementInterface(localName, namespace)
    local interfaces = nil
    if namespace == namespaces.HTML_NAMESPACE then
        interfaces = HTML_ELEMENT_INTERFACES
    elseif namespace == namespaces.SVG_NAMESPACE then
        error("TODO")
    elseif namespace == namespaces.MATHML_NAMESPACE then
        error("TODO")
    end
    if interfaces ~= nil then
        for name, interface in pairs(interfaces) do
            if name == localName then
                return interface
            end
        end
    end
    -- The element interface for any name and namespace is Element, unless stated otherwise.
    return Element
end


---https://dom.spec.whatwg.org/#concept-create-element
---@param document DOM_Document
---@param localName string
---@param namespace string?
---@param namespacePrefix string?
---@param is string?
---@param synchronousCustomElements? boolean
---@return DOM_Element
function Element.create(document, localName, namespace, namespacePrefix, is, synchronousCustomElements)
    if synchronousCustomElements == nil then synchronousCustomElements = false end

    --- 1. Let result be null.
    local result = nil

    -- 2. Let definition be the result of looking up a custom element definition given document, namespace, localName, and is.
    local definition = nil
    -- TODO

    -- 3. If definition is non-null, and definition’s name is not equal to its local name (i.e., definition represents a customized built-in element):
    if definition ~= nil then
        error("TODO")

        -- 1. Let interface be the element interface for localName and the HTML namespace.

        -- 2. Set result to a new element that implements interface, with no attributes, namespace set to the HTML namespace, namespace prefix set to prefix, local name set to localName, custom element state set to "undefined", custom element definition set to null, is value set to is, and node document set to document.

        -- 3. If synchronousCustomElements is true, then run this step while catching any exceptions:

        --        1. Upgrade result using definition.

        --    If this step threw an exception exception:

        --        1. Report exception for definition’s constructor’s corresponding JavaScript object’s associated realm’s global object.

        --        2. Set result’s custom element state to "failed".

        -- 4. Otherwise, enqueue a custom element upgrade reaction given result and definition.
    elseif definition ~= nil then --[[
        4. Otherwise, if definition is non-null:
    ]]
        -- 1. If synchronousCustomElements is true, then run these steps while catching any exceptions:

        --        1. Let C be definition’s constructor.

        --        2. Set result to the result of constructing C, with no arguments.

        --        3. Assert: result’s custom element state and custom element definition are initialized.

        --        4. Assert: result’s namespace is the HTML namespace.

        --        5. If result’s attribute list is not empty, then throw a "NotSupportedError" DOMException.

        --        6. If result has children, then throw a "NotSupportedError" DOMException.

        --        7. If result’s parent is not null, then throw a "NotSupportedError" DOMException.

        --        8. If result’s node document is not document, then throw a "NotSupportedError" DOMException.

        --        9. If result’s local name is not equal to localName, then throw a "NotSupportedError" DOMException.

        --        10. Set result’s namespace prefix to prefix.

        --        11. Set result’s is value to null.

        --    If any of these steps threw an exception exception:

        --        1. Report exception for definition’s constructor’s corresponding JavaScript object’s associated realm’s global object.

        --        2. Set result to a new element that implements the HTMLUnknownElement interface, with no attributes, namespace set to the HTML namespace, namespace prefix set to prefix, local name set to localName, custom element state set to "failed", custom element definition set to null, is value set to null, and node document set to document.

        -- 2. Otherwise:

        --        1. Set result to a new element that implements the HTMLElement interface, with no attributes, namespace set to the HTML namespace, namespace prefix set to prefix, local name set to localName, custom element state set to "undefined", custom element definition set to null, is value set to null, and node document set to document.

        --        2. Enqueue a custom element upgrade reaction given result and definition.
    else --[[
        5. Otherwise:
    ]]
        -- 1. Let interface be the element interface for localName and namespace.
        local interface = elementInterface(localName, namespace)

        -- 2. Set result to a new element that implements interface,
        result = interface:new(
            {},              -- with no attributes,
            namespace,       -- namespace set to namespace,
            namespacePrefix, -- namespace prefix set to prefix,
            localName,       -- local name set to localName,
            "uncustomized",  -- custom element state set to "uncustomized",
            nil,             -- custom element definition set to null,
            is,              -- is value set to is,
            document         -- and node document set to document.
        )

        -- 3. If namespace is the HTML namespace,
        if namespace == namespaces.HTML_NAMESPACE and
            -- and either localName is a valid custom element name
            (true --[[ TODO ]] or
                -- or is is non-null,
                is ~= nil
            )
        then
            -- then set result’s custom element state to "undefined".
            result.customElementState = "undefined"
        end
    end
    assert(result ~= nil)
    -- 6. Return result.
    return result
end

---https://dom.spec.whatwg.org/#element-shadow-host
---@return boolean
function Element:isShadowHost()
    -- An element is a shadow host if its shadow root is non-null.
    return self.shadowRoot ~= nil
end

---https://dom.spec.whatwg.org/#concept-element-qualified-name
---@return string
function Element:qualifiedName()
    -- An element’s qualified name is
    if self.namespacePrefix == nil then
        -- its local name if its namespace prefix is null;
        return self.localName
    end
    -- otherwise its namespace prefix, followed by ":", followed by its local name.
    return self.namespacePrefix .. ":" .. self.localName
end

---https://dom.spec.whatwg.org/#element-html-uppercased-qualified-name
---@return string
function Element:htmlUppercasedQualifiedName()
    -- 1. Let qualifiedName be this’s qualified name.
    local qualifiedName = self:qualifiedName()
    -- 2. If this is in the HTML namespace and its node document is an HTML document, then set qualifiedName to qualifiedName in ASCII uppercase.
    if (self.namespace == namespaces.HTML_NAMESPACE) and (self.nodeDocument.type == "html") then
        qualifiedName = strings.asciiLowercase(qualifiedName)
    end
    -- 3. Return qualifiedName.
    return qualifiedName
end

---https://dom.spec.whatwg.org/#concept-element-custom
---@return boolean
function Element:isCustom()
    -- An element whose custom element state is "custom" is said to be custom.
    return self.customElementState == "custom"
end

---https://dom.spec.whatwg.org/#concept-element-attributes-append
---@param attribute DOM_Attr
function Element:appendAttribute(attribute)
    -- 1. Append attribute to element’s attribute list.
    table.insert(self.attributes, attribute)
    -- 2. Set attribute’s element to element.
    attribute.element = self
    -- 3. Handle attribute changes for attribute with element, null, and attribute’s value.
    -- TODO
end

---@param name string
---@param namespace string?
---@param prefix string?
---@return DOM_Attr|nil
function Element:getAttribute(name, namespace, prefix)
    for _, attr in ipairs(self.attributes) do
        if attr.localName == name
            and (namespace == nil or attr.namespace == namespace)
            and (prefix == nil or attr.prefix == prefix)
        then
            return attr
        end
    end
    return nil
end

---@param name string
---@param namespace string?
---@param prefix string?
---@return string|nil
function Element:getAttributeValue(name, namespace, prefix)
    local attr = self:getAttribute(name, namespace, prefix)
    if attr == nil then
        return nil
    end
    return attr.value
end

---@param name string
---@param value string
---@param namespace string?
---@param prefix string?
function Element:setAttributeValue(name, value, namespace, prefix)
    local attr = self:getAttribute(name, namespace, prefix)
    if attr == nil then
        error("No such attribute: " .. tostring(name))
    end
    attr.value = value
end

---@param localName string?
---@return boolean
function Element:isHTMLElement(localName)
    if self.namespace ~= namespaces.HTML_NAMESPACE then
        return false
    end
    if localName ~= nil and self.localName ~= localName then
        return false
    end
    return true
end

---@param localName string?
---@return boolean
function Element:isSVGElement(localName)
    if self.namespace ~= namespaces.SVG_NAMESPACE then
        return false
    end
    if localName ~= nil and self.localName ~= localName then
        return false
    end
    return true
end

---@param localName string?
---@return boolean
function Element:isMathMLElement(localName)
    if self.namespace ~= namespaces.MATHML_NAMESPACE then
        return false
    end
    if localName ~= nil and self.localName ~= localName then
        return false
    end
    return true
end

---@param localNames string[]
---@return boolean
function Element:isOneOfHTMLElements(localNames)
    if self.namespace ~= namespaces.HTML_NAMESPACE then
        return false
    end
    for _, n in ipairs(localNames) do
        if n == self.localName then
            return true
        end
    end
    return false
end

local SPECIAL_CATEGORY_HTML_TAG_NAMES = {
    "address", "applet", "area", "article", "aside", "base", "basefont",
    "bgsound", "blockquote", "body", "br", "button", "caption", "center",
    "col", "colgroup", "dd", "details", "dir", "div", "dl", "dt", "embed",
    "fieldset", "figcaption", "figure", "footer", "form", "frame", "frameset",
    "h1", "h2", "h3", "h4", "h5", "h6", "head", "header", "hgroup", "hr",
    "html", "iframe", "img", "input", "keygen", "li", "link", "listing", "main",
    "marquee", "menu", "meta", "nav", "noembed", "noframes", "noscript",
    "object", "ol", "p", "param", "plaintext", "pre", "script", "search",
    "section", "select", "source", "style", "summary", "table", "tbody", "td",
    "template", "textarea", "tfoot", "th", "thead", "title", "tr", "track",
    "ul", "wbr", "xmp"
}
local SPECIAL_CATEGORY_MATHML_TAG_NAMES = { "mi", "mn", "ms", "mtext", "annotation-xml" }
local SPECIAL_CATEGORY_SVG_TAG_NAMES = { "foreignObject", "desc", "title" }

---https://html.spec.whatwg.org/multipage/parsing.html#special
---@return boolean
function Element:isInSpecialCategory()
    for _, t in ipairs(SPECIAL_CATEGORY_HTML_TAG_NAMES) do
        if self:isHTMLElement(t) then
            return true
        end
    end
    for _, t in ipairs(SPECIAL_CATEGORY_MATHML_TAG_NAMES) do
        if self:isMathMLElement(t) then
            return true
        end
    end
    for _, t in ipairs(SPECIAL_CATEGORY_SVG_TAG_NAMES) do
        if self:isSVGElement(t) then
            return true
        end
    end
    return false
end

---https://dom.spec.whatwg.org/#dom-element-tagname
---@return string
function Element:tagName()
    -- The tagName getter steps are to return this’s HTML-uppercased qualified name.
    return self:htmlUppercasedQualifiedName()
end

return Element
