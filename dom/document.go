// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package dom

import (
	"fmt"
	"net/url"
)

// TODO(ois): DocumentOrigin is currently a STUB type.
type DocumentOrigin struct{}

// TODO(ois): DocumentEnvironmentSettings is currently a STUB type.
type DocumentEnvironmentSettings struct{}

// TODO(ois): DocumentPolicyContainer is currently a STUB type.
type DocumentPolicyContainer struct{}

// Document represents a [DOM document].
//
// TODO(ois): Should we move some of Document's HTML related extensions to the html package instead?
//
// [DOM document]: https://dom.spec.whatwg.org/#concept-document
type Document interface {
	Node

	// Origin returns [origin] of the document.
	//
	// [origin]: https://dom.spec.whatwg.org/#concept-document-origin
	Origin() DocumentOrigin

	// SetOrigin sets [origin] of the document to origin.
	//
	// [origin]: https://dom.spec.whatwg.org/#concept-document-origin
	SetOrigin(origin DocumentOrigin)

	// Mode returns [mode] of the document.
	//
	// [mode]: https://dom.spec.whatwg.org/#concept-document-mode
	Mode() DocumentMode

	// SetMode sets [mode] of the document to mode.
	//
	// [mode]: https://dom.spec.whatwg.org/#concept-document-mode
	SetMode(mode DocumentMode)

	// CustomElementRegistry returns [custom element registry] of the document.
	//
	// [custom element registry]: https://dom.spec.whatwg.org/#document-custom-element-registry
	CustomElementRegistry() *CustomElementRegistry

	// EffectiveGlobalCustomElementRegistry returns [effective global custom element registry] of the document.
	//
	// [effective global custom element registry]: https://dom.spec.whatwg.org/#effective-global-custom-element-registry
	EffectiveGlobalCustomElementRegistry() *CustomElementRegistry

	//==========================================================================
	// HTML related extensions
	//==========================================================================

	// PolicyContainer returns [policy container] of the document.
	//
	// TODO(ois): Document.PolicyContainer is currently a STUB.
	//
	// [policy container]: https://html.spec.whatwg.org/multipage/dom.html#concept-document-policy-container
	PolicyContainer() DocumentPolicyContainer

	// SetPolicyContainer sets [policy container] of the document to policyContainer.
	//
	// TODO(ois): Document.SetPolicyContainer is currently a STUB.
	//
	// [policy container]: https://html.spec.whatwg.org/multipage/dom.html#concept-document-policy-container
	SetPolicyContainer(policyContainer DocumentPolicyContainer)

	// BaseURL returns [base URL] of the document.
	//
	// [base URL]: https://html.spec.whatwg.org/multipage/urls-and-fetching.html#document-base-url
	BaseURL() url.URL

	// SetBaseURL sets [base URL] of the document.
	//
	// BUG(ois): Document's Base URL should be read-only field, and Document.SetBaseURL should not exist. This only exists because BaseURL is currently implemented as read-write field, instead of algorithm that returns an URL.
	//
	// [base URL]: https://html.spec.whatwg.org/multipage/urls-and-fetching.html#document-base-url
	SetBaseURL(url url.URL)

	// IsParserCannotChangeMode reports whether [HTML parser cannot change the mode].
	//
	// [HTML parser cannot change the mode]: https://html.spec.whatwg.org/multipage/parsing.html#parser-cannot-change-the-mode-flag
	IsParserCannotChangeMode() bool

	// IsIframeSrcdocDocument reports whether document is an [iframe srcdoc document].
	//
	// [iframe srcdoc document]: https://html.spec.whatwg.org/multipage/iframe-embed-object.html#an-iframe-srcdoc-document
	IsIframeSrcdocDocument() bool

	// RelevantSettings returns [relevant settings object] of the document.
	//
	// TODO(ois): Document.RelevantSettings is currently a STUB.
	//
	// [relevant settings object]: https://html.spec.whatwg.org/multipage/webappapis.html#relevant-settings-object
	RelevantSettings() DocumentEnvironmentSettings

	// RelevantSettings sets [relevant settings object] of the document to settings.
	//
	// TODO(ois): Document.SetReleavntSettings is currently a STUB.
	//
	// [relevant settings object]: https://html.spec.whatwg.org/multipage/webappapis.html#relevant-settings-object
	SetReleavntSettings(settings DocumentEnvironmentSettings)
}

// DocumentMode represents Document's [mode]
//
// [mode]: https://dom.spec.whatwg.org/#concept-document-mode
type DocumentMode uint8

const (
	NoQuirks      DocumentMode = iota // "no-quirks"
	Quirks                            // "quirks"
	LimitedQuirks                     // "limited-quirks"
)

// BUG(ois): documentImpl.baseURL should not be a field. Instead it should follow steps at: https://html.spec.whatwg.org/multipage/urls-and-fetching.html#document-base-url
type documentImpl struct {
	Node

	// TODO: https://dom.spec.whatwg.org/#concept-document-encoding
	// TODO: https://dom.spec.whatwg.org/#concept-document-content-type
	// TODO: https://dom.spec.whatwg.org/#concept-document-url
	// TODO: https://dom.spec.whatwg.org/#concept-document-type
	// TODO: https://dom.spec.whatwg.org/#document-allow-declarative-shadow-roots
	origin                DocumentOrigin // STUB
	mode                  DocumentMode
	customElementRegistry CustomElementRegistry

	iframeSrcdocDocument   bool
	parserCannotChangeMode bool
	baseURL                url.URL

	// Below are STUB
	environmentSettings DocumentEnvironmentSettings
	policyContainer     DocumentPolicyContainer
}

// NewDocument constructs a new [Document].
func NewDocument() Document {
	doc := &documentImpl{}
	doc.Node = NewNode(doc)
	return doc
}
func (doc documentImpl) String() string {
	mode := "?"
	switch doc.mode {
	case NoQuirks:
		mode = "no-quirks"
	case Quirks:
		mode = "quirks"
	case LimitedQuirks:
		mode = "limited-quirks"
	}
	return fmt.Sprintf("#document(mode=%s)", mode)
}

func (doc documentImpl) Origin() DocumentOrigin {
	return doc.origin
}
func (doc *documentImpl) SetOrigin(origin DocumentOrigin) {
	doc.origin = origin
}
func (doc documentImpl) RelevantSettings() DocumentEnvironmentSettings {
	return doc.environmentSettings
}
func (doc *documentImpl) SetReleavntSettings(settings DocumentEnvironmentSettings) {
	doc.environmentSettings = settings
}
func (doc documentImpl) PolicyContainer() DocumentPolicyContainer {
	return doc.policyContainer
}
func (doc *documentImpl) SetPolicyContainer(policyContainer DocumentPolicyContainer) {
	doc.policyContainer = policyContainer
}
func (doc documentImpl) BaseURL() url.URL {
	return doc.baseURL
}
func (doc *documentImpl) SetBaseURL(url url.URL) {
	doc.baseURL = url
}

func (doc documentImpl) Mode() DocumentMode           { return doc.mode }
func (doc *documentImpl) SetMode(mode DocumentMode)   { doc.mode = mode }
func (d documentImpl) IsParserCannotChangeMode() bool { return d.parserCannotChangeMode }
func (d documentImpl) IsIframeSrcdocDocument() bool   { return d.iframeSrcdocDocument }
func (doc documentImpl) CustomElementRegistry() *CustomElementRegistry {
	return &doc.customElementRegistry
}

// https://dom.spec.whatwg.org/#effective-global-custom-element-registry
func (doc documentImpl) EffectiveGlobalCustomElementRegistry() *CustomElementRegistry {
	if IsGlobalCustomElementReigstry(&doc.customElementRegistry) {
		return &doc.customElementRegistry
	}
	return nil
}
