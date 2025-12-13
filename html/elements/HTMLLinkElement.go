// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package elements

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/inseo-oh/yw/css/cssom"
	"github.com/inseo-oh/yw/css/csssyntax"
	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/html/fetch"
)

// STUB
type sourceSet struct{}

// HTMLLinkElement represents a [link] element.
//
// [link]: https://html.spec.whatwg.org/multipage/semantics.html#the-link-element
type HTMLLinkElement interface {
	HTMLElement
	processLink()
}
type htmlLinkElementImpl struct {
	HTMLElement
	sourceSet []sourceSet
}

// NewHTMLLinkElement constructs a new [HTMLLinkElement] node.
//
// [html]: https://html.spec.whatwg.org/multipage/semantics.html#the-html-element
func NewHTMLLinkElement(options dom.ElementCreationCommonOptions) HTMLLinkElement {
	elem := &htmlLinkElementImpl{
		HTMLElement: NewHTMLElement(options),
	}

	cbs := elem.Callbacks()

	// HTML Spec defines precisely when link element should be processed, but this will do the job for now.
	// (Example: https://html.spec.whatwg.org/multipage/links.html#link-type-stylesheet)
	cbs.PoppedFromStackOfOpenElements = func() {
		elem.processLink()
	}
	return elem
}

func (elem htmlLinkElementImpl) processLink() {
	rel, ok := elem.AttrWithoutNamespace("rel")
	if !ok {
		return
	}
	var (
		fetchAndProcessLinkedResource func()
		linkedResourceFetchSetupSteps func() bool
		processLinkedResource         func(success bool, response *http.Response, responseBytes []byte)
	)
	switch rel {
	case rel:
		fetchAndProcessLinkedResource, linkedResourceFetchSetupSteps, processLinkedResource = elem.processLinkTypeStylesheet()
	}

	// https://html.spec.whatwg.org/multipage/semantics.html#link-processing-options
	type linkProcessingOptions struct {
		href            string                          // https://html.spec.whatwg.org/multipage/semantics.html#link-options-href
		initiator       string                          // https://html.spec.whatwg.org/multipage/semantics.html#link-options-initiator
		integrity       string                          // https://html.spec.whatwg.org/multipage/semantics.html#link-options-integrity
		tp              string                          // https://html.spec.whatwg.org/multipage/semantics.html#link-options-type
		nonce           string                          // https://html.spec.whatwg.org/multipage/semantics.html#link-options-nonce
		destination     string                          // https://html.spec.whatwg.org/multipage/semantics.html#link-options-destination
		crossorigin     fetch.CorsSettings              // https://html.spec.whatwg.org/multipage/semantics.html#link-options-crossorigin
		referrerPolicy  any                             // [STUB] https://html.spec.whatwg.org/multipage/semantics.html#link-options-referrer-policy
		sourceSet       []sourceSet                     // https://html.spec.whatwg.org/multipage/semantics.html#link-options-source-set
		baseURL         url.URL                         // https://html.spec.whatwg.org/multipage/semantics.html#link-options-base-url
		origin          dom.DocumentOrigin              // https://html.spec.whatwg.org/multipage/semantics.html#link-options-origin
		environment     dom.DocumentEnvironmentSettings // https://html.spec.whatwg.org/multipage/semantics.html#link-options-environment
		policyContainer dom.DocumentPolicyContainer     // https://html.spec.whatwg.org/multipage/semantics.html#link-options-policy-container
		document        dom.Document                    // https://html.spec.whatwg.org/multipage/semantics.html#link-options-document
		onDocumentReady func(doc dom.Document)          // https://html.spec.whatwg.org/multipage/semantics.html#link-options-on-document-ready
		fetchPriority   fetch.FetchPriority             // https://html.spec.whatwg.org/multipage/semantics.html#link-options-fetch-priority
	}
	defaultLinkProcessingOptions := func() linkProcessingOptions {
		return linkProcessingOptions{
			href:            "nil",
			initiator:       "link",
			integrity:       "",
			tp:              "",
			nonce:           "",
			destination:     "",
			crossorigin:     fetch.CorsNone,
			referrerPolicy:  nil,
			sourceSet:       nil,
			document:        nil,
			onDocumentReady: nil,
			fetchPriority:   fetch.FetchPriorityAuto,
		}
	}

	// https://html.spec.whatwg.org/multipage/semantics.html#create-link-options-from-element
	createLinkOptions := func() linkProcessingOptions {
		document := elem.NodeDocument()
		options := defaultLinkProcessingOptions()
		options.crossorigin = fetch.CorsSettingsFromAttr(elem, "crossorigin")
		options.referrerPolicy = nil // TODO
		options.sourceSet = elem.sourceSet
		options.baseURL = document.BaseURL()
		options.origin = document.Origin()
		options.environment = document.RelevantSettings()
		options.policyContainer = document.PolicyContainer()
		options.document = document
		options.nonce = "" // TODO
		options.fetchPriority = fetch.FetchPriorityFromAttr(elem, "fetchpriority")
		if attr, ok := elem.AttrWithoutNamespace("href"); ok {
			options.href = attr
		}
		if attr, ok := elem.AttrWithoutNamespace("integrity"); ok {
			options.integrity = attr
		}
		if attr, ok := elem.AttrWithoutNamespace("type"); ok {
			options.tp = attr
		}
		return options
	}
	// https://html.spec.whatwg.org/multipage/semantics.html#create-a-link-request
	createLinkRequest := func(options linkProcessingOptions) (res *http.Request, err error) {
		// STUB
		// NOTE: We don't use JoinPath() because the "path" part of URL may not be a real filesystem path.
		req, err := http.NewRequest("GET", options.baseURL.String()+"/"+options.href, nil)
		if err != nil {
			return nil, err
		}
		// TODO: Set a real user agent
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
		return req, nil
	}

	if processLinkedResource == nil {
		processLinkedResource = func(success bool, response *http.Response, responseBytes []byte) {}
	}
	if fetchAndProcessLinkedResource == nil {
		// https://html.spec.whatwg.org/multipage/semantics.html#default-fetch-and-process-the-linked-resource
		fetchAndProcessLinkedResource = func() {
			// STUB
			options := createLinkOptions()
			request, err := createLinkRequest(options)
			var resp *http.Response
			var bytes []byte

			success := true
			if err != nil {
				log.Printf("<link>: %v", err)
				success = false
				goto process
			}
			resp, err = http.DefaultClient.Do(request)
			if err != nil {
				log.Printf("<link>: %v", err)
				success = false
				goto process
			}
			bytes, err = io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("<link>: %v", err)
				success = false
				goto process
			}
		process:
			processLinkedResource(success, resp, bytes)
		}
	}
	if linkedResourceFetchSetupSteps == nil {
		linkedResourceFetchSetupSteps = func() bool { return true }
	}
	if !linkedResourceFetchSetupSteps() {
		log.Printf("<link>: linkedResourceFetchSetupSteps() failed")
		return
	}
	fetchAndProcessLinkedResource()
}

// https://html.spec.whatwg.org/multipage/links.html#link-type-stylesheet
func (elem htmlLinkElementImpl) processLinkTypeStylesheet() (
	fetchAndProcessLinkedResource func(),
	linkedResourceFetchSetupSteps func() bool,
	processLinkedResource func(success bool, response *http.Response, responseBytes []byte),
) {
	processLinkedResource = func(success bool, response *http.Response, responseBytes []byte) {
		// NOTE: All the step numbers(S#.) are based on spec from when this was initially written(2025.11.25)

		// S1.
		// TODO: If the resource's Content-Type metadata is not text/css, then set success to false.

		// S2.
		// TODO: If el no longer creates an external resource link that contributes to the styling processing model, or if, since the resource in question was fetched, it has become appropriate to fetch it again, then:

		// S3.
		if sheet := cssom.AssociatedStylesheet(elem); sheet != nil {
			cssom.RemoveStylesheet(sheet)
		}

		// S4.
		if success {
			urlStr := response.Request.URL.String()
			// S4-1.
			stylesheet, err := csssyntax.ParseStylesheet(responseBytes, &urlStr, urlStr)
			if err != nil {
				log.Printf("<link %s>: failed to tokenize stylesheet: %v", urlStr, err)
				return
			}
			stylesheet.Type = "text/css"
			stylesheet.OwnerNode = elem
			stylesheet.Location = &urlStr
			// TODO: Set stylesheet.media once we implement that
			if dom.IsInDocumentTree(elem) {
				if attr, ok := elem.AttrWithoutNamespace("title"); ok {
					stylesheet.Title = attr
				}
			}
			stylesheet.AlternateFlag = false  // TODO: Set if the link is an alternative style sheet and el's explicitly enabled is false; unset otherwise.
			stylesheet.OriginCleanFlag = true // TODO: Set if the resource is CORS-same-origin; unset otherwise.
			stylesheet.ParentStylesheet = nil
			stylesheet.OwnerRule = nil
			cssom.AddStylesheet(&stylesheet)
			log.Printf("<link %s>: stylesheet loaded", urlStr)
		} else {
			// S5.
			// TODO: Otherwise, fire an event named error at el.
		}

		// S6.
		if elem.ContributesScriptBlockingStylesheet() {
			// TODO: append element to its node document's script-blocking style sheet set.
			panic("TODO[https://html.spec.whatwg.org/multipage/links.html#link-type-stylesheet]")
		}

		// S7.
		// TODO: Unblock rendering on el.
	}
	return nil, nil, processLinkedResource
}
