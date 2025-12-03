package dom

import "github.com/inseo-oh/yw/namespaces"

// CustomElementRegistry represents a [HTML CustomElementRegistry].
//
// [HTML CustomElementRegistry]: https://html.spec.whatwg.org/multipage/custom-elements.html#customelementregistry
//
// TODO(ois): CustomElementRegistry is currently STUB.
type CustomElementRegistry struct {
	IsScoped          bool       // https://html.spec.whatwg.org/multipage/custom-elements.html#is-scoped
	ScopedDocumentSet []Document // https://html.spec.whatwg.org/multipage/custom-elements.html#scoped-document-set
}

// CustomElementDefinition represents a [HTML custom element definition].
//
// [HTML custom element definition]: https://html.spec.whatwg.org/multipage/custom-elements.html#custom-element-definition
type CustomElementDefinition struct {
	Name      string
	LocalName string
	// STUB
}

// LookupCustomElementDefinition looks up custom element definition.
//
// namespace, is may be nil if absent.
//
// TODO(ois): LookupCustomElementDefinition is currently just STUB, as there's no custom elements support.
//
// Spec: https://html.spec.whatwg.org/multipage/custom-elements.html#look-up-a-custom-element-definition
func (reg *CustomElementRegistry) LookupCustomElementDefinition(namespace *namespaces.Namespace, localname string, is *string) *CustomElementDefinition {
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
