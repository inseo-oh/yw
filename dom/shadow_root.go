package dom

type ShadowRoot interface {
	DocumentFragment
	CustomElementRegistry() *CustomElementRegistry
	SetCustomElementRegistry(registry *CustomElementRegistry)
}
