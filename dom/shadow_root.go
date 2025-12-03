// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

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
