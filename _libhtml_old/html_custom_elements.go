package libhtml

type html_custom_element_registry struct {
	// STUB
	is_scoped           bool           // https://html.spec.whatwg.org/multipage/custom-elements.html#is-scoped
	scoped_document_set []dom_Document // https://html.spec.whatwg.org/multipage/custom-elements.html#scoped-document-set
}
type dom_custom_element_definition struct {
	name       string
	local_name string
	// STUB
}

// https://html.spec.whatwg.org/multipage/custom-elements.html#look-up-a-custom-element-definition
//
// namespace, is may be nil.
func (reg *html_custom_element_registry) look_up_custom_element_definition(namespace *namespace, local_name string, is *string) *dom_custom_element_definition {
	// STUB
	return nil
}

// https://html.spec.whatwg.org/multipage/custom-elements.html#concept-try-upgrade
func html_try_upgrade_element(element dom_Element) {
	var ns *namespace
	if v, ok := element.get_namespace(); ok {
		ns = &v
	}
	var is *string
	if v, ok := element.get_is(); ok {
		is = &v
	}
	definition := element.get_custom_element_registry().look_up_custom_element_definition(ns, element.get_local_name(), is)
	if definition != nil {
		// TODO: enqueue a custom element upgrade reaction given element and definition.
		panic("TODO[https://html.spec.whatwg.org/multipage/custom-elements.html#concept-try-upgrade]")
	}
}
