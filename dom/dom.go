package dom

import (
	"yw/namespaces"
)

type TagToken interface {
	Attr(name string) (string, bool)
}

type ElementCreationCommonOptions struct {
	NodeDocument          Document
	LocalName             string
	Namespace             *namespaces.Namespace  // May be nil
	Prefix                *string                // May be nil
	CustomElementRegistry *CustomElementRegistry // May be nil
	CustomElementState    CustomElementState
	Is                    *string
	TagToken              TagToken
}

// This is just an unique pointer value, not an actual registry!
var DefaultCustomElementReigistry = &CustomElementRegistry{}

// https://dom.spec.whatwg.org/#concept-create-element
//
// registry can also be nil, or DefaultCustomElementReigistry(= use current document's registry)
func CreateElement(
	document Document, localName string, namespace *namespaces.Namespace, prefix *string, is *string,
	synchronousCustomElements bool, registry *CustomElementRegistry,
	tagToken TagToken, getFactoryFn func(namespace *namespaces.Namespace, localName string) func(opt ElementCreationCommonOptions) Element,
) Element {
	var res Element
	isDefaultRegistry := (registry == DefaultCustomElementReigistry)
	if isDefaultRegistry {
		registry = document.CustomElementRegistry()
	}
	definition := registry.LookupCustomElementDefinition(namespace, localName, is)
	if definition != nil && definition.Name != definition.LocalName {
		panic("TODO[https://dom.spec.whatwg.org/#concept-create-element]")
	} else if definition != nil {
		panic("TODO[https://dom.spec.whatwg.org/#concept-create-element]")
	} else {
		factoryFn := getFactoryFn(namespace, localName)
		res = factoryFn(ElementCreationCommonOptions{
			NodeDocument:          document,
			Namespace:             namespace,
			Prefix:                prefix,
			LocalName:             localName,
			TagToken:              tagToken,
			CustomElementRegistry: registry,
			CustomElementState:    CustomElementUncustomized,
			Is:                    is,
		})
	}
	return res
}

// https://dom.spec.whatwg.org/#is-a-global-custom-element-registry
//
// registry may be nil
func IsGlobalCustomElementReigstry(registry *CustomElementRegistry) bool {
	return registry != nil && !registry.IsScoped
}
