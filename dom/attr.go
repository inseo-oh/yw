// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package dom

import (
	"fmt"

	"github.com/inseo-oh/yw/namespaces"
)

// Attr represents a [DOM attribute].
//
// [DOM attribute]: https://dom.spec.whatwg.org/#concept-attribute
type Attr interface {
	Node

	// Namespace returns [namespace] of the attribute. ok is set to false if it's absent.
	//
	// [namespace]: https://dom.spec.whatwg.org/#concept-attribute-namespace
	Namespace() (ns namespaces.Namespace, ok bool)

	// NamespacePrefix returns [namespace prefix] of the attribute. ok is set to false if it's absent.
	//
	// [namespace prefix]: https://dom.spec.whatwg.org/#concept-attribute-namespace-prefix
	NamespacePrefix() (pr string, ok bool)

	// LocalName returns the [local name] of the attribute.
	//
	// [local name]: https://dom.spec.whatwg.org/#concept-attribute-local-name
	LocalName() string

	// Value returns the [value] of the attribute.
	//
	// [value]: https://dom.spec.whatwg.org/#concept-attribute-value
	Value() string

	// Element returns [element] for the attribute.
	//
	// [element]: https://dom.spec.whatwg.org/#concept-attribute-element
	Element() Element
}

// AttrData is ligher weight version of [Attr]. This isn't a DOM node, but can
// hold data needed to create an [Attr].
type AttrData struct {
	LocalName       string                // local name of the attribute.
	Value           string                // value of the attribute.
	Namespace       *namespaces.Namespace // namespace of the attribute, nil if absent.
	NamespacePrefix *string               // namespace of the attribute, nil if absent.
}
type attrImpl struct {
	Node
	localName       string
	value           string
	namespace       *namespaces.Namespace // nil if absent.
	namespacePrefix *string               // nil if absent.
	element         Element
}

// NewAttr constructs a new [Attr] value.
//
// namespace, namespacePrefix may be nil if absent.
func NewAttr(localName string, value string, namespace *namespaces.Namespace, namespacePrefix *string, element Element) Attr {
	return &attrImpl{
		NewNode(element.NodeDocument()),
		localName, value, namespace, namespacePrefix, element,
	}
}

func (at attrImpl) String() string {
	if at.namespace != nil {
		return fmt.Sprintf("#attr(%v:%s = %s)", *at.namespace, at.localName, at.value)
	} else {
		return fmt.Sprintf("#attr(%s = %s)", at.localName, at.value)
	}
}

func (at attrImpl) LocalName() string { return at.localName }
func (at attrImpl) Value() string     { return at.value }
func (at attrImpl) Element() Element  { return at.element }
func (at attrImpl) Namespace() (namespaces.Namespace, bool) {
	if at.namespace == nil {
		return "", false
	}
	return *at.namespace, true
}
func (at attrImpl) NamespacePrefix() (string, bool) {
	if at.namespacePrefix == nil {
		return "", false
	}
	return *at.namespacePrefix, true
}
