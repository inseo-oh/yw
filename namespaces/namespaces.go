// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

// Package namespaces provides various namespace values.
package namespaces

import "fmt"

// Namespace is simply a string containing URL. And while they are URL, it's not
// used to pull data from it. URL is simply used as unique identifier value.
type Namespace string

var (
	Html   Namespace = "http://www.w3.org/1999/xhtml"
	Mathml Namespace = "http://www.w3.org/1998/Math/MathML"
	Svg    Namespace = "http://www.w3.org/2000/svg"
	Xlink  Namespace = "http://www.w3.org/1999/xlink"
	Xml    Namespace = "http://www.w3.org/XML/1998/namespace"
	Xmlns  Namespace = "http://www.w3.org/2000/xmlns/"
)

// String description of the namespace. For known namespaces, their names will
// be returned instead of the URL.
func (n Namespace) String() string {
	switch n {
	case Html:
		return "html"
	case Mathml:
		return "mathml"
	case Svg:
		return "svg"
	case Xlink:
		return "xlink"
	case Xml:
		return "xml"
	case Xmlns:
		return "xmlns"
	default:
		return fmt.Sprintf("<namespace %s>", string(n))
	}
}
