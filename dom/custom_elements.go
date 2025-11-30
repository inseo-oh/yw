package dom

import "github.com/inseo-oh/yw/namespaces"

type CustomElementRegistry struct {
	// STUB
	IsScoped          bool       // https://html.spec.whatwg.org/multipage/custom-elements.html#is-scoped
	ScopedDocumentSet []Document // https://html.spec.whatwg.org/multipage/custom-elements.html#scoped-document-set
}
type CustomElementDefinition struct {
	Name      string
	LocalName string
	// STUB
}

// https://html.spec.whatwg.org/multipage/custom-elements.html#look-up-a-custom-element-definition
//
// namespace, is may be nil.
func (reg *CustomElementRegistry) LookupCustomElementDefinition(namespace *namespaces.Namespace, localname string, is *string) *CustomElementDefinition {
	// STUB
	return nil
}

// https://html.spec.whatwg.org/multipage/custom-elements.html#concept-try-upgrade
func tryUpgradeElement(element Element) {
	var ns *namespaces.Namespace
	if v, ok := element.Namespace(); ok {
		ns = &v
	}
	var is *string
	if v, ok := element.Is(); ok {
		is = &v
	}
	definition := element.CustomElementRegistry().LookupCustomElementDefinition(ns, element.LocalName(), is)
	if definition != nil {
		// TODO: enqueue a custom element upgrade reaction given element and definition.
		panic("TODO[https://html.spec.whatwg.org/multipage/custom-elements.html#concept-try-upgrade]")
	}
}
