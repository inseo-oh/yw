// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

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
