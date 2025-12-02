package dom

// DocumentFragment represents a [DOM DocumentFragment].
//
// [DOM DocumentFragment]: https://dom.spec.whatwg.org/#documentfragment
type DocumentFragment interface {
	Node

	// Host returns [host] of the DocumentFragment.
	//
	// [host]: https://dom.spec.whatwg.org/#concept-documentfragment-host
	Host() Node
}
