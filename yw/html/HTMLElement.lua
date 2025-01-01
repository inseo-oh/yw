--[[
    Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
    SPDX-License-Identifier: BSD-3-Clause
    This software may contain third-party material. For more info, see README.
]]
local Element = require "yw.dom.Element"
local object = require "yw.common.object"

---https://html.spec.whatwg.org/multipage/dom.html#htmlelement
---@class HTML_HTMLElement : DOM_Element
local HTMLElement = object.create(Element)

---@param attributes DOM_Attr[]
---@param namespace string?
---@param namespacePrefix string?
---@param localName string
---@param customElementState DOM_Element_CustomElementState
---@param customElementDefinition HTML_CustomElementDefinition?
---@param is string?
---@param document DOM_Document
---@return HTML_HTMLElement
function HTMLElement:new(attributes, namespace, namespacePrefix, localName, customElementState, customElementDefinition,
                         is, document)
    local o = Element.new(self, attributes, namespace, namespacePrefix, localName, customElementState,
        customElementDefinition, is, document) --[[@as HTML_HTMLElement]]
    return o
end

---https://html.spec.whatwg.org/multipage/custom-elements.html#form-associated-custom-element
---@return boolean
function HTMLElement:isFormAssociatedCustomElement()
    -- TODO: This is STUB
    return false
end


---https://html.spec.whatwg.org/multipage/forms.html#form-associated-element
HTMLElement.FORM_ASSOCIATED_ELEMENTS = {
    "button",
    "fieldset",
    "input",
    "object",
    "output",
    "select",
    "textarea",
    "img",
}

---https://html.spec.whatwg.org/multipage/forms.html#category-listed
HTMLElement.FORM_LISTED_ELEMENTS = {
    "button",
    "fieldset",
    "input",
    "object",
    "output",
    "select",
    "textarea"
}
---https://html.spec.whatwg.org/multipage/forms.html#category-submit
HTMLElement.FORM_SUBMIT_ELEMENTS = {
    "button",
    "input",
    "select",
    "textarea"
}
---https://html.spec.whatwg.org/multipage/forms.html#category-reset
HTMLElement.FORM_RESETTABLE_ELEMENTS = {
    "input",
    "output",
    "select",
    "textarea"
}
---https://html.spec.whatwg.org/multipage/forms.html#category-autocapitalize
HTMLElement.FORM_AUTOCAPITALIZE_ELEMENTS =  {
    "button",
    "fieldset",
    "input",
    "output",
    "select",
    "textarea"
}

---https://html.spec.whatwg.org/multipage/forms.html#form-associated-element
---@param category string[]?  See `HTMLElement.FORM_~_ELEMENTS`
---@return boolean
function HTMLElement:isFormAssociatedElement(category)
    if self:isFormAssociatedCustomElement() then
        return true
    end
    if category == nil then
        category = HTMLElement.FORM_ASSOCIATED_ELEMENTS
    end
    for _, element in ipairs(category) do
        if self.localName == element then
            return true
        end
    end
    return true
end

---https://html.spec.whatwg.org/multipage/form-control-infrastructure.html#concept-form-reset-control
function HTMLElement:reset()
    error("Not implemented")
end

return HTMLElement
