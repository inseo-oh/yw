// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package dom

import (
	"fmt"
	"strings"
)

// DocumentType represents a [DOM doctype](a.k.a DocumentFragment).
//
// [DOM doctype]: https://dom.spec.whatwg.org/#concept-doctype
type DocumentType interface {
	Node

	// Name returns [name] of the doctype.
	//
	// [name]: https://dom.spec.whatwg.org/#concept-doctype-name
	Name() string

	// PublicId returns [public ID] of the doctype.
	//
	// [public ID]: https://dom.spec.whatwg.org/#concept-doctype-publicid
	PublicId() string

	// SystemId returns [system ID] of the doctype.
	//
	// [system ID]: https://dom.spec.whatwg.org/#concept-doctype-systemid
	SystemId() string
}
type documentTypeImpl struct {
	Node
	name     string
	publicId string
	systemId string
}

// NewDocumentType constructs a new [DocumentType] node.
func NewDocumentType(doc Document, name, publicId, systemId string) DocumentType {
	return documentTypeImpl{
		NewNode(doc),
		name, publicId, systemId,
	}
}
func (dt documentTypeImpl) String() string {
	sb := strings.Builder{}
	sb.WriteString("<!DOCTYPE")
	if dt.name != "" {
		sb.WriteString(fmt.Sprintf(" %s", dt.name))
	}
	if dt.publicId != "" {
		sb.WriteString(fmt.Sprintf(" PUBLIC %s", dt.publicId))
	}
	if dt.systemId != "" {
		sb.WriteString(fmt.Sprintf(" SYSTEM %s", dt.systemId))
	}
	sb.WriteString(">")
	return sb.String()
}
func (dt documentTypeImpl) Name() string     { return dt.name }
func (dt documentTypeImpl) PublicId() string { return dt.publicId }
func (dt documentTypeImpl) SystemId() string { return dt.systemId }
