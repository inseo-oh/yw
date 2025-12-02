package dom

// ShadowRoot represents a [DOM shadow root]
//
// [DOM shadow root]: https://dom.spec.whatwg.org/#concept-shadow-root
type ShadowRoot interface {
	DocumentFragment

	// CustomElementRegistry returns [custom element registry] of the shadow root.
	//
	// [custom element registry]: https://dom.spec.whatwg.org/#shadowroot-custom-element-registry
	CustomElementRegistry() *CustomElementRegistry

	// SetCustomElementRegistry sets [custom element registry] of the shadow root.
	//
	// [custom element registry]: https://dom.spec.whatwg.org/#shadowroot-custom-element-registry
	SetCustomElementRegistry(registry *CustomElementRegistry)
}
