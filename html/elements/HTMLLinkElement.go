package elements

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"yw/css/cssom"
	"yw/css/csssyntax"
	"yw/dom"
	"yw/html/urlfetch"
)

// STUB
type SourceSet struct{}

type HTMLLinkElement interface {
	HTMLElement
	ProcessLink()
}
type htmlLinkElementImpl struct {
	HTMLElement
	sourceSet []SourceSet
}

func NewHTMLLinkElement(options dom.ElementCreationCommonOptions) HTMLLinkElement {
	elem := &htmlLinkElementImpl{
		HTMLElement: NewHTMLElement(options),
	}

	cbs := elem.Callbacks()

	// HTML Spec defines precisely when link element should be processed, but this will do the job for now.
	// (Example: https://html.spec.whatwg.org/multipage/links.html#link-type-stylesheet)
	cbs.PoppedFromStackOfOpenElements = func() {
		elem.ProcessLink()
	}
	return elem
}

func (elem htmlLinkElementImpl) ProcessLink() {
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
		crossorigin     urlfetch.CorsSettings           // https://html.spec.whatwg.org/multipage/semantics.html#link-options-crossorigin
		referrerPolicy  any                             // [STUB] https://html.spec.whatwg.org/multipage/semantics.html#link-options-referrer-policy
		sourceSet       []SourceSet                     // https://html.spec.whatwg.org/multipage/semantics.html#link-options-source-set
		baseURL         url.URL                         // https://html.spec.whatwg.org/multipage/semantics.html#link-options-base-url
		origin          dom.DocumentOrigin              // https://html.spec.whatwg.org/multipage/semantics.html#link-options-origin
		environment     dom.DocumentEnvironmentSettings // https://html.spec.whatwg.org/multipage/semantics.html#link-options-environment
		policyContainer dom.DocumentPolicyContainer     // https://html.spec.whatwg.org/multipage/semantics.html#link-options-policy-container
		document        dom.Document                    // https://html.spec.whatwg.org/multipage/semantics.html#link-options-document
		onDocumentReady func(doc dom.Document)          // https://html.spec.whatwg.org/multipage/semantics.html#link-options-on-document-ready
		fetchPriority   urlfetch.FetchPriority          // https://html.spec.whatwg.org/multipage/semantics.html#link-options-fetch-priority
	}
	defaultLinkProcessingOptions := func() linkProcessingOptions {
		return linkProcessingOptions{
			href:            "nil",
			initiator:       "link",
			integrity:       "",
			tp:              "",
			nonce:           "",
			destination:     "",
			crossorigin:     urlfetch.CorsNone,
			referrerPolicy:  nil,
			sourceSet:       nil,
			document:        nil,
			onDocumentReady: nil,
			fetchPriority:   urlfetch.FetchPriorityAuto,
		}
	}

	// https://html.spec.whatwg.org/multipage/semantics.html#create-link-options-from-element
	createLinkOptions := func() linkProcessingOptions {
		document := elem.NodeDocument()
		options := defaultLinkProcessingOptions()
		options.crossorigin = urlfetch.CorsSettingsFromAttr(elem, "crossorigin")
		options.referrerPolicy = nil // TODO
		options.sourceSet = elem.sourceSet
		options.baseURL = document.BaseURL()
		options.origin = document.Origin()
		options.environment = document.RelevantSettings()
		options.policyContainer = document.PolicyContainer()
		options.document = document
		options.nonce = "" // TODO
		options.fetchPriority = urlfetch.FetchPriorityFromAttr(elem, "fetchpriority")
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
	createLinkRequest := func(options linkProcessingOptions) (*http.Request, error) {
		// STUB
		// NOTE: We don't use JoinPath() because the "path" part of URL may not be a real filesystem path.
		return http.NewRequest("GET", options.baseURL.String()+"/"+options.href, nil)
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
			text := csssyntax.DecodeBytes(responseBytes)
			tokens, err := csssyntax.Tokenize(text)
			if err != nil {
				log.Printf("<link %s>: failed to tokenize stylesheet: %v", urlStr, err)
				return
			}
			stylesheet := csssyntax.ParseStylesheet(tokens, &urlStr)
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
