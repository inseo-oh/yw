package elements

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"yw/dom"
)

// STUB
type html_source_set struct{}

type html_HTMLLinkElement interface {
	HTMLElement
	process_link()
}
type html_HTMLLinkElement_s struct {
	HTMLElement
	source_set []html_source_set
}

func NewHTMLLinkElement(options dom.ElementCreationCommonOptions) html_HTMLLinkElement {
	elem := &html_HTMLLinkElement_s{
		HTMLElement: NewHTMLElement(options),
	}

	cbs := elem.Callbacks()

	// HTML Spec defines precisely when link element should be processed, but this will do the job for now.
	// (Example: https://html.spec.whatwg.org/multipage/links.html#link-type-stylesheet)
	cbs.PoppedFromStackOfOpenElements = func() {
		elem.process_link()
	}
	return elem
}

func (elem html_HTMLLinkElement_s) process_link() {
	rel, ok := elem.AttrWithoutNamespace("rel")
	if !ok {
		return
	}
	var (
		fetch_and_process_linked_resource func()
		linked_resource_fetch_setup_steps func() bool
		process_linked_resource           func(success bool, response *http.Response, response_bytes []byte)
	)
	switch rel {
	case rel:
		fetch_and_process_linked_resource,
			linked_resource_fetch_setup_steps,
			process_linked_resource = elem.process_link_type_stylesheet()
	}

	// https://html.spec.whatwg.org/multipage/semantics.html#link-processing-options
	type link_processing_options struct {
		href              string                          // https://html.spec.whatwg.org/multipage/semantics.html#link-options-href
		initiator         string                          // https://html.spec.whatwg.org/multipage/semantics.html#link-options-initiator
		integrity         string                          // https://html.spec.whatwg.org/multipage/semantics.html#link-options-integrity
		tp                string                          // https://html.spec.whatwg.org/multipage/semantics.html#link-options-type
		nonce             string                          // https://html.spec.whatwg.org/multipage/semantics.html#link-options-nonce
		destination       string                          // https://html.spec.whatwg.org/multipage/semantics.html#link-options-destination
		crossorigin       html_cors_settings              // https://html.spec.whatwg.org/multipage/semantics.html#link-options-crossorigin
		referrer_policy   any                             // [STUB] https://html.spec.whatwg.org/multipage/semantics.html#link-options-referrer-policy
		source_set        []html_source_set               // https://html.spec.whatwg.org/multipage/semantics.html#link-options-source-set
		base_url          url.URL                         // https://html.spec.whatwg.org/multipage/semantics.html#link-options-base-url
		origin            dom.DocumentOrigin              // https://html.spec.whatwg.org/multipage/semantics.html#link-options-origin
		environment       dom.DocumentEnvironmentSettings // https://html.spec.whatwg.org/multipage/semantics.html#link-options-environment
		policy_container  dom.DocumentPolicyContainer     // https://html.spec.whatwg.org/multipage/semantics.html#link-options-policy-container
		document          dom.Document                    // https://html.spec.whatwg.org/multipage/semantics.html#link-options-document
		on_document_ready func(doc dom.Document)          // https://html.spec.whatwg.org/multipage/semantics.html#link-options-on-document-ready
		fetch_priority    html_fetch_priority             // https://html.spec.whatwg.org/multipage/semantics.html#link-options-fetch-priority
	}
	default_link_processing_options := func() link_processing_options {
		return link_processing_options{
			href:              "nil",
			initiator:         "link",
			integrity:         "",
			tp:                "",
			nonce:             "",
			destination:       "",
			crossorigin:       html_cors_settings_no_cors,
			referrer_policy:   nil,
			source_set:        nil,
			document:          nil,
			on_document_ready: nil,
			fetch_priority:    html_fetch_priority_auto,
		}
	}

	// https://html.spec.whatwg.org/multipage/semantics.html#create-link-options-from-element
	create_link_options := func() link_processing_options {
		document := elem.NodeDocument()
		options := default_link_processing_options()
		options.crossorigin = html_cors_settings_from_attr(elem, "crossorigin")
		options.referrer_policy = nil // TODO
		options.source_set = elem.source_set
		options.base_url = document.BaseURL()
		options.origin = document.Origin()
		options.environment = document.RelevantSettings()
		options.policy_container = document.PolicyContainer()
		options.document = document
		options.nonce = "" // TODO
		options.fetch_priority = html_fetch_priority_from_attr(elem, "fetchpriority")
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
	create_link_request := func(options link_processing_options) (*http.Request, error) {
		// STUB
		// NOTE: We don't use JoinPath() because the "path" part of URL may not be a real filesystem path.
		return http.NewRequest("GET", options.base_url.String()+"/"+options.href, nil)
	}

	if process_linked_resource == nil {
		process_linked_resource = func(success bool, response *http.Response, response_bytes []byte) {}
	}
	if fetch_and_process_linked_resource == nil {
		// https://html.spec.whatwg.org/multipage/semantics.html#default-fetch-and-process-the-linked-resource
		fetch_and_process_linked_resource = func() {
			// STUB
			options := create_link_options()
			request, err := create_link_request(options)
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
			process_linked_resource(success, resp, bytes)
		}
	}
	if linked_resource_fetch_setup_steps == nil {
		linked_resource_fetch_setup_steps = func() bool { return true }
	}
	if !linked_resource_fetch_setup_steps() {
		log.Printf("<link>: linked_resource_fetch_setup_steps() failed")
		return
	}
	fetch_and_process_linked_resource()
}

// https://html.spec.whatwg.org/multipage/links.html#link-type-stylesheet
func (elem html_HTMLLinkElement_s) process_link_type_stylesheet() (
	fetch_and_process_linked_resource func(),
	linked_resource_fetch_setup_steps func() bool,
	process_linked_resource func(success bool, response *http.Response, response_bytes []byte),
) {
	process_linked_resource = func(success bool, response *http.Response, response_bytes []byte) {
		// NOTE: All the step numbers(S#.) are based on spec from when this was initially written(2025.11.25)

		// S1.
		// TODO: If the resource's Content-Type metadata is not text/css, then set success to false.

		// S2.
		// TODO: If el no longer creates an external resource link that contributes to the styling processing model, or if, since the resource in question was fetched, it has become appropriate to fetch it again, then:

		// S3.
		if sheet := css_associated_stylesheet(elem); sheet != nil {
			css_remove_stylesheet(sheet)
		}

		// S4.
		if success {
			url_str := response.Request.URL.String()
			// S4-1.
			text := css_decode_bytes(response_bytes)
			tokens, err := css_tokenize(text)
			if err != nil {
				log.Printf("<link %s>: failed to tokenize stylesheet: %v", url_str, err)
				return
			}
			stylesheet := css_parse_stylesheet(tokens, &url_str)
			stylesheet.tp = "text/css"
			stylesheet.owner_node = elem
			stylesheet.location = &url_str
			// TODO: Set stylesheet.media once we implement that
			if dom_node_is_in_document_tree(elem) {
				if attr, ok := elem.AttrWithoutNamespace("title"); ok {
					stylesheet.title = attr
				}
			}
			stylesheet.alternate_flag = false   // TODO: Set if the link is an alternative style sheet and el's explicitly enabled is false; unset otherwise.
			stylesheet.origin_clean_flag = true // TODO: Set if the resource is CORS-same-origin; unset otherwise.
			stylesheet.parent_stylesheet = nil
			stylesheet.owner_rule = nil
			css_add_stylesheet(&stylesheet)
			log.Printf("<link %s>: stylesheet loaded", url_str)
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
	return nil, nil, process_linked_resource
}
