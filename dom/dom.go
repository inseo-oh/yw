// Package dom implements the [DOM standard].
//
// # Introduction to hierarchical type system
//
// DOM's type system is very hierarchical. For example, every DOM types([Attr],
// [Document], [Element], ...) inherit from DOM [Node], and HTMLElement is
// DOM [Element], and so on. This is very easy to implement in OO languages,
// but Go wasn't designed for this.
//
// Fortunately, Go's interfaces can work like this. You can embed other interfaces
// into a interface or struct, and interface values can be converted from one
// interface to another on the fly, as long as concrete type implements all of
// necessary functions (This does mean Go checks presence of functions at
// runtime, potentially slowing it down, but we won't worry about it for now).
//
// But for this to work flawlessly, we never return structs implementing such
// interfaces outside of this package. And this also means there's always two
// different types for a single DOM node:
//
//   - An interface exposed to outside of package
//   - And the struct that actually implements it - those are private and has
//     suffix "~Impl".
//
// And this also means in order to add a public function to such type, you have to
// add it to both struct and interface. And also means we have a lot of accessor
// functions, as Go interfaces doesn't have fields, only functions.
//
// So it's a bit awkard, but this is probably the most reasonable solution in my
// opinion. Also, technically web specs also use term "interface", and phrases like
// "Nodes are objects that implement Node", so they aren't too far off from how Go
// works.
//
// # Relationship with html package
//
// Generally speaking, things that are part of DOM spec goes into this package,
// and HTML spec things go to html package.
//
// But DOM and HTML standards work so close together, that it's not uncommon to
// see relying on each other's types or functionality. That is fine from people
// writing standards, but trying to replicate that in Go would cause cyclic
// dependency.
//
// So to avoid that, this package also implements serveral parts from
// [HTML standard] locally, such as [CustomElementRegistry], and this package
// never imports anything from the html package. [TagToken] is another example:
// that interface exists so that we can access some of HTML tag functionality
// without actually importing html package.
//
// # New~ functions
//
// As structs are not exposed publically, DOM types usually have New~ function,
// which constructs the value and returns corresponding interface value.
//
// [DOM standard]: https://dom.spec.whatwg.org
// [HTML standard]: https://html.spec.whatwg.org/multipage
package dom

import (
	"github.com/inseo-oh/yw/namespaces"
)

// TagToken is interface for the HTML token.
//
// html package's token implements this interface.
type TagToken interface {
	// Attr searches attribute from tag's attributes that matches the name, and
	// returns its value. ok is set to false if there's no such attribute.
	Attr(name string) (value string, ok bool)
}

// ElementCreationCommonOptions contains common options when creating an [Element].
type ElementCreationCommonOptions struct {
	NodeDocument          Document               // Node Document of the element.
	LocalName             string                 // Local name of the element.
	Namespace             *namespaces.Namespace  // Namespace of the element. Nil if absent.
	Prefix                *string                // Namespace prefix of the element. Nil if absent.
	CustomElementRegistry *CustomElementRegistry // Custom element registry of the element. Nil if not applicable.
	CustomElementState    CustomElementState     // Initial custom element state
	Is                    *string                // Is value
	TagToken              TagToken               // Associated HTML tag token
}

// DefaultCustomElementReigistry is an unique pointer value, that tells
// [CreateElement] to use document's registry.
//
// Note that DefaultCustomElementReigistry itself is not actually a
// custom element registry.
var DefaultCustomElementReigistry = &CustomElementRegistry{}

// Creates a DOM Element.
//
// This looks similar to [NewElement], but that only constructs the value.
// This one creates the right type based on return value of getFactoryFn.
//
// registry can be nil, or if [DefaultCustomElementReigistry] is passed, it will
// use document's registry.
//
// Spec: https://dom.spec.whatwg.org/#concept-create-element
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

// DefaultCustomElementReigistry reports whether registry is a global custom element registry.
//
// registry may be nil.
//
// Spec: https://dom.spec.whatwg.org/#is-a-global-custom-element-registry
func IsGlobalCustomElementReigstry(registry *CustomElementRegistry) bool {
	return registry != nil && !registry.IsScoped
}
