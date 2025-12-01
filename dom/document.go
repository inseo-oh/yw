package dom

import (
	"fmt"
	"net/url"
)

// STUB types
type DocumentOrigin struct{}
type DocumentEnvironmentSettings struct{}
type DocumentPolicyContainer struct{}

type Document interface {
	Node

	Mode() DocumentMode
	SetMode(mode DocumentMode)
	IsParserCannotChangeMode() bool
	IsIframeSrcdocDocument() bool
	EffectiveGLobalCustomElementRegistry() *CustomElementRegistry
	CustomElementRegistry() *CustomElementRegistry

	Origin() DocumentOrigin
	SetOrigin(origin DocumentOrigin)
	RelevantSettings() DocumentEnvironmentSettings
	SetReleavntSettings(settings DocumentEnvironmentSettings)
	PolicyContainer() DocumentPolicyContainer
	SetPolicyContainer(policyContainer DocumentPolicyContainer)
	BaseURL() url.URL
	SetBaseURL(url url.URL)
}

type DocumentMode uint8

const (
	NoQuirks DocumentMode = iota
	Quirks
	LimitedQuirks
)

type documentImpl struct {
	Node
	mode                   DocumentMode
	iframeSrcdocDocument   bool
	parserCannotChangeMode bool
	customElementRegistry  CustomElementRegistry
	baseURL                url.URL

	// Below are STUB
	origin              DocumentOrigin
	environmentSettings DocumentEnvironmentSettings
	policyContainer     DocumentPolicyContainer
}

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
func (doc documentImpl) EffectiveGLobalCustomElementRegistry() *CustomElementRegistry {
	if IsGlobalCustomElementReigstry(&doc.customElementRegistry) {
		return &doc.customElementRegistry
	}
	return nil
}
