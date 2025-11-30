package libhtml

type dom_element_creation_common_options struct {
	node_document        dom_Document
	local_name           string
	namespace            *namespace                    // May be nil
	prefix               *string                       // May be nil
	custom_elem_registry *html_custom_element_registry // May be nil
	custom_element_state dom_custom_element_state
	is                   *string
	tag_token            html_tag_token
}

// This is just an unique pointer value, not an actual registry!
var dom_default_custom_element_registry = &html_custom_element_registry{}

// https://dom.spec.whatwg.org/#concept-create-element
//
// registry can also be nil, or dom_default_custom_element_registry(= use current document's registry)
func dom_create_element(
	document dom_Document, local_name string, namespace *namespace, prefix *string, is *string,
	synchronous_custom_elements bool, registry *html_custom_element_registry,
	tag_token html_tag_token,
) dom_Element {
	var res dom_Element
	is_default_registry := (registry == dom_default_custom_element_registry)
	if is_default_registry {
		registry = document.get_custom_element_registry()
	}
	definition := registry.look_up_custom_element_definition(namespace, local_name, is)
	if definition != nil && definition.name != definition.local_name {
		panic("TODO[https://dom.spec.whatwg.org/#concept-create-element]")
	} else if definition != nil {
		panic("TODO[https://dom.spec.whatwg.org/#concept-create-element]")
	} else {
		factory_fn := func(opt dom_element_creation_common_options) dom_Element {
			return html_make_HTMLElement(opt, dom_element_callbacks{})
		}
		if namespace != nil && *namespace == html_namespace && local_name == "html" {
			factory_fn = func(opt dom_element_creation_common_options) dom_Element { return html_make_HTMLHtmlElement(opt) }
		} else if namespace != nil && *namespace == html_namespace && local_name == "body" {
			factory_fn = func(opt dom_element_creation_common_options) dom_Element { return html_make_HTMLBodyElement(opt) }
		} else if namespace != nil && *namespace == html_namespace && local_name == "link" {
			factory_fn = func(opt dom_element_creation_common_options) dom_Element { return html_make_HTMLLinkElement(opt) }
		} else if namespace != nil && *namespace == html_namespace && local_name == "style" {
			factory_fn = func(opt dom_element_creation_common_options) dom_Element { return html_make_HTMLStyleElement(opt) }
		}

		res = factory_fn(dom_element_creation_common_options{
			node_document:        document,
			namespace:            namespace,
			prefix:               prefix,
			local_name:           local_name,
			tag_token:            tag_token,
			custom_elem_registry: registry,
			custom_element_state: dom_custom_element_state_uncustomized,
			is:                   is,
		})
	}
	return res
}

// https://dom.spec.whatwg.org/#is-a-global-custom-element-registry
//
// registry may be nil
func dom_is_global_custom_element_registry(registry *html_custom_element_registry) bool {
	return registry != nil && !registry.is_scoped
}
