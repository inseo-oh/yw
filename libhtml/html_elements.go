package libhtml

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	cm "yw/libcommon"
)

// ------------------------------------------------------------------------------
// HTMLElement
// ------------------------------------------------------------------------------

type html_HTMLElement interface {
	dom_Element
	is_form_associated_custom_element() bool
	is_form_associated_element() bool
	is_form_listed_element() bool
	is_form_submittable_element() bool
	is_form_resettable_element() bool
	is_form_autocapitalize_and_autocorrect_inheriting_element() bool
	is_form_labelable_element() bool
	contributes_script_blocking_style_sheet() bool
}
type html_HTMLElement_s struct{ dom_Element }

func html_make_HTMLElement(options dom_element_creation_common_options, callbacks dom_element_callbacks) html_HTMLElement {
	return html_HTMLElement_s{dom_make_Element(options, callbacks, true)}
}

// https://html.spec.whatwg.org/multipage/custom-elements.html#form-associated-custom-element
func (elem html_HTMLElement_s) is_form_associated_custom_element() bool {
	// STUB
	return false
}

// https://html.spec.whatwg.org/multipage/forms.html#form-associated-element
func (elem html_HTMLElement_s) is_form_associated_element() bool {
	html_elems := []string{
		"button", "fieldset", "input", "object", "output", "select", "textarea", "img",
	}
	return elem.is_form_associated_custom_element() ||
		slices.ContainsFunc(html_elems, elem.is_html_element)
}

// https://html.spec.whatwg.org/multipage/forms.html#category-listed
func (elem html_HTMLElement_s) is_form_listed_element() bool {
	html_elems := []string{
		"button", "fieldset", "input", "object", "output", "select", "textarea",
	}
	return elem.is_form_associated_custom_element() ||
		slices.ContainsFunc(html_elems, elem.is_html_element)
}

// https://html.spec.whatwg.org/multipage/forms.html#category-submit
func (elem html_HTMLElement_s) is_form_submittable_element() bool {
	html_elems := []string{"button", "input", "select", "textarea"}
	return elem.is_form_associated_custom_element() ||
		slices.ContainsFunc(html_elems, elem.is_html_element)
}

// https://html.spec.whatwg.org/multipage/forms.html#category-reset
func (elem html_HTMLElement_s) is_form_resettable_element() bool {
	html_elems := []string{"input", "output", "select", "textarea"}
	return elem.is_form_associated_custom_element() ||
		slices.ContainsFunc(html_elems, elem.is_html_element)
}

// https://html.spec.whatwg.org/multipage/forms.html#category-autocapitalize
func (elem html_HTMLElement_s) is_form_autocapitalize_and_autocorrect_inheriting_element() bool {
	html_elems := []string{
		"button", "fieldset", "input", "output", "select", "textarea",
	}
	return elem.is_form_associated_custom_element() ||
		slices.ContainsFunc(html_elems, elem.is_html_element)
}

// https://html.spec.whatwg.org/multipage/forms.html#category-label
func (elem html_HTMLElement_s) is_form_labelable_element() bool {
	html_elems := []string{
		"button", "meter", "output", "progress", "select", "textarea",
	}
	if elem.is_form_associated_custom_element() {
		return true
	}
	if slices.ContainsFunc(html_elems, elem.is_html_element) {
		return true
	}
	if elem.is_html_element("input") {
		if attr, ok := elem.get_attribute_without_namespace("type"); ok && attr == "hidden" {
			return true
		}
	}
	return false
}

// https://html.spec.whatwg.org/multipage/semantics.html#contributes-a-script-blocking-style-sheet
func (elem html_HTMLElement_s) contributes_script_blocking_style_sheet() bool {
	// STUB
	return false
}

// ------------------------------------------------------------------------------
// HTMLHtmlElement
// ------------------------------------------------------------------------------

type html_HTMLHtmlElement interface{ html_HTMLElement }
type html_HTMLHtmlElement_s struct{ html_HTMLElement }

func html_make_HTMLHtmlElement(options dom_element_creation_common_options) html_HTMLHtmlElement {
	return html_HTMLHtmlElement_s{
		html_HTMLElement: html_make_HTMLElement(options, dom_element_callbacks{
			get_intrinsic_size: func() (width float64, height float64) {
				// XXX: Do we need intrinsic size for this?
				return 0, 0
			},
		}),
	}
}

// https://html.spec.whatwg.org/multipage/common-microsyntaxes.html#rules-for-parsing-a-legacy-colour-value
func html_parse_legacy_color_value(input string) (css_color, bool) {
	if input == "" {
		return css_color{}, false
	}
	input = strings.Trim(input, " ")
	if cm.ToAsciiLowercase(input) == "transparent" {
		// transparent
		return css_color{}, false
	}
	if color, ok := css_named_colors_map[cm.ToAsciiLowercase(input)]; ok {
		// CSS named colors
		return css_color_from_rgba(color.red, color.green, color.blue, color.alpha), true
	}
	input_cps := []rune(input)
	if len(input_cps) == 4 && input_cps[0] == '#' {
		// #rgb
		red, err1 := strconv.ParseInt(string(input_cps[1]), 16, 8)
		green, err2 := strconv.ParseInt(string(input_cps[2]), 16, 8)
		blue, err3 := strconv.ParseInt(string(input_cps[3]), 16, 8)
		if err1 == nil && err2 == nil && err3 == nil {
			return css_color_from_rgba(uint8(red), uint8(blue), uint8(green), 255), true
		}
	}
	// Now we assume the format is #rrggbb -------------------------------------
	new_input_cps := make([]rune, 0, len(input_cps))
	for i := range len(input_cps) {
		// Replace characters beyond BMP with "00"
		if input_cps[i] > 0xffff {
			new_input_cps = append(new_input_cps, '0')
			new_input_cps = append(new_input_cps, '0')
		} else {
			new_input_cps = append(new_input_cps, input_cps[i])
		}
	}
	input_cps = new_input_cps
	if 128 < len(input_cps) {
		input_cps = input_cps[:128]
	}
	if input_cps[0] == '#' {
		input_cps = input_cps[1:]
	}
	for i := range len(input_cps) {
		// Replace non-hex characters with '0'
		if _, err := strconv.ParseInt(string(input_cps[i]), 16, 8); err != nil {
			input_cps[i] = '0'
		}
	}
	// Length must be nonzero, and multiple of 3. If not, append '0's.
	for len(input_cps) == 0 || len(input_cps)%3 != 0 {
		input_cps = append(input_cps, '0')
	}
	comp_len := len(input_cps) / 3
	comps := [][]rune{
		input_cps[:comp_len*1],
		input_cps[comp_len*1 : comp_len*2],
		input_cps[comp_len*2 : comp_len*3],
	}
	if comp_len > 8 {
		for i := range 3 {
			comps[i] = comps[i][comp_len-8:]
		}
		comp_len = 8
	}
	for comp_len > 2 {
		for i := range 3 {
			if comps[i][0] == '0' {
				comps[i] = comps[i][1:]
			}
		}
		comp_len--
	}
	if comp_len > 2 {
		for i := range 3 {
			comps[i] = comps[i][:2]
		}
		comp_len = 2
	}
	red, err1 := strconv.ParseInt(string(comps[0]), 16, 16)
	green, err2 := strconv.ParseInt(string(comps[1]), 16, 16)
	blue, err3 := strconv.ParseInt(string(comps[2]), 16, 16)
	if err1 != nil || err2 != nil || err3 != nil {
		panic("unreachable")
	}
	return css_color_from_rgba(uint8(red), uint8(blue), uint8(green), 255), true
}

// ------------------------------------------------------------------------------
// HTMLBodyElement
// https://html.spec.whatwg.org/multipage/sections.html#the-body-element
// ------------------------------------------------------------------------------
func html_make_HTMLBodyElement(options dom_element_creation_common_options) html_HTMLElement {
	elem := html_make_HTMLElement(options, dom_element_callbacks{
		get_intrinsic_size: func() (width float64, height float64) {
			// XXX: Do we need intrinsic size for this?
			return 0, 0
		},
	})

	cbs := elem.get_callbacks()
	cbs.get_presentational_hints = func() []css_declaration {
		decls := []css_declaration{}

		// https://html.spec.whatwg.org/multipage/rendering.html#the-page
		if attr, ok := elem.get_attribute_without_namespace("bgcolor"); ok {
			color, ok := html_parse_legacy_color_value(attr)
			if ok {
				decls = append(decls, css_declaration{"background-color", color})
			}
		}
		if attr, ok := elem.get_attribute_without_namespace("text"); ok {
			color, ok := html_parse_legacy_color_value(attr)
			if ok {
				decls = append(decls, css_declaration{"color", color})
			}
		}
		return decls

	}
	return elem
}

//------------------------------------------------------------------------------
// HTMLLinkElement
// https://html.spec.whatwg.org/multipage/semantics.html#the-link-element
//------------------------------------------------------------------------------

// STUB
type html_source_set struct{}

type html_HTMLLinkElement interface {
	html_HTMLElement
	process_link()
}
type html_HTMLLinkElement_s struct {
	html_HTMLElement
	source_set []html_source_set
}

func html_make_HTMLLinkElement(options dom_element_creation_common_options) html_HTMLLinkElement {
	elem := &html_HTMLLinkElement_s{
		html_HTMLElement: html_make_HTMLElement(options, dom_element_callbacks{
			get_intrinsic_size: func() (width float64, height float64) {
				// XXX: Do we need intrinsic size for this?
				return 0, 0
			},
		}),
	}

	cbs := elem.get_callbacks()

	// HTML Spec defines precisely when link element should be processed, but this will do the job for now.
	// (Example: https://html.spec.whatwg.org/multipage/links.html#link-type-stylesheet)
	cbs.popped_from_stack_of_open_elements = func() {
		elem.process_link()
	}
	return elem
}

func (elem html_HTMLLinkElement_s) process_link() {
	rel, ok := elem.get_attribute_without_namespace("rel")
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
		href              string                        // https://html.spec.whatwg.org/multipage/semantics.html#link-options-href
		initiator         string                        // https://html.spec.whatwg.org/multipage/semantics.html#link-options-initiator
		integrity         string                        // https://html.spec.whatwg.org/multipage/semantics.html#link-options-integrity
		tp                string                        // https://html.spec.whatwg.org/multipage/semantics.html#link-options-type
		nonce             string                        // https://html.spec.whatwg.org/multipage/semantics.html#link-options-nonce
		destination       string                        // https://html.spec.whatwg.org/multipage/semantics.html#link-options-destination
		crossorigin       html_cors_settings            // https://html.spec.whatwg.org/multipage/semantics.html#link-options-crossorigin
		referrer_policy   any                           // [STUB] https://html.spec.whatwg.org/multipage/semantics.html#link-options-referrer-policy
		source_set        []html_source_set             // https://html.spec.whatwg.org/multipage/semantics.html#link-options-source-set
		base_url          url.URL                       // https://html.spec.whatwg.org/multipage/semantics.html#link-options-base-url
		origin            document_origin               // https://html.spec.whatwg.org/multipage/semantics.html#link-options-origin
		environment       document_environment_settings // https://html.spec.whatwg.org/multipage/semantics.html#link-options-environment
		policy_container  document_policy_container     // https://html.spec.whatwg.org/multipage/semantics.html#link-options-policy-container
		document          dom_Document                  // https://html.spec.whatwg.org/multipage/semantics.html#link-options-document
		on_document_ready func(doc dom_Document)        // https://html.spec.whatwg.org/multipage/semantics.html#link-options-on-document-ready
		fetch_priority    html_fetch_priority           // https://html.spec.whatwg.org/multipage/semantics.html#link-options-fetch-priority
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
		document := elem.get_node_document()
		options := default_link_processing_options()
		options.crossorigin = html_cors_settings_from_attr(elem, "crossorigin")
		options.referrer_policy = nil // TODO
		options.source_set = elem.source_set
		options.base_url = document.get_base_url()
		options.origin = document.get_origin()
		options.environment = document.get_releavant_settings()
		options.policy_container = document.get_policy_container()
		options.document = document
		options.nonce = "" // TODO
		options.fetch_priority = html_fetch_priority_from_attr(elem, "fetchpriority")
		if attr, ok := elem.get_attribute_without_namespace("href"); ok {
			options.href = attr
		}
		if attr, ok := elem.get_attribute_without_namespace("integrity"); ok {
			options.integrity = attr
		}
		if attr, ok := elem.get_attribute_without_namespace("type"); ok {
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
				if attr, ok := elem.get_attribute_without_namespace("title"); ok {
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
		if elem.contributes_script_blocking_style_sheet() {
			// TODO: append element to its node document's script-blocking style sheet set.
			panic("TODO[https://html.spec.whatwg.org/multipage/links.html#link-type-stylesheet]")
		}

		// S7.
		// TODO: Unblock rendering on el.
	}
	return nil, nil, process_linked_resource
}

//------------------------------------------------------------------------------
// HTMLStyleElement
// https://html.spec.whatwg.org/multipage/semantics.html#the-style-element
//------------------------------------------------------------------------------

type html_HTMLStyleElement interface{ html_HTMLElement }
type html_HTMLStyleElement_s struct {
	html_HTMLElement
}

func html_make_HTMLStyleElement(options dom_element_creation_common_options) html_HTMLStyleElement {
	elem := html_HTMLStyleElement_s{
		html_HTMLElement: html_make_HTMLElement(options, dom_element_callbacks{
			get_intrinsic_size: func() (width float64, height float64) {
				// XXX: Do we need intrinsic size for this?
				return 0, 0
			},
		}),
	}

	cbs := elem.get_callbacks()

	// From 4.2.6. The style element(https://html.spec.whatwg.org/multipage/semantics.html#the-style-element)
	// The user agent must run the update a style block algorithm whenever any of the following conditions occur:
	//  - The element is popped off the stack of open elements of an HTML parser or XML parser.
	cbs.popped_from_stack_of_open_elements = func() {
		elem.update_style_block()
	}
	//  - The element is not on the stack of open elements of an HTML parser or XML parser, and it becomes connected or disconnected.
	//  - The element's children changed steps run.
	cbs.run_children_changed_steps = func() {
		elem.update_style_block()
	}
	return elem
}

// https://html.spec.whatwg.org/multipage/semantics.html#update-a-style-block
func (elem *html_HTMLStyleElement_s) update_style_block() {
	// NOTE: All the step numbers(S#.) are based on spec from when this was initially written(2025.11.13)

	// S2.
	if sheet := css_associated_stylesheet(elem); sheet != nil {
		css_remove_stylesheet(sheet)
	}
	// S3.
	if !dom_node_is_connected(elem) {
		return
	}
	// S4.
	// TODO: If element's type attribute is present and its value is neither the empty string nor an ASCII case-insensitive match for "text/css", then return.
	// S5.
	// TODO: If the Should element's inline behavior be blocked by Content Security Policy? algorithm returns "Blocked" when executed upon the style element, "style", and the style element's child text content, then return. [CSP]
	// S6.
	text, ok := elem.get_child_text_node()
	if !ok {
		text = ""
	}
	tokens, err := css_tokenize(text)
	if err != nil {
		log.Printf("<style>: failed to tokenize stylesheet: %v", err)
		return
	}
	stylesheet := css_parse_stylesheet(tokens, nil)
	stylesheet.tp = "text/css"
	stylesheet.owner_node = elem
	// TODO: Set stylesheet.media once we implement that
	if dom_node_is_in_document_tree(elem) {
		if attr, ok := elem.get_attribute_without_namespace("title"); ok {
			stylesheet.title = attr
		}
	}
	stylesheet.alternate_flag = false
	stylesheet.origin_clean_flag = true
	stylesheet.location = nil
	stylesheet.parent_stylesheet = nil
	stylesheet.owner_rule = nil
	css_add_stylesheet(&stylesheet)
	log.Printf("<style>: stylesheet loaded")

	// S7.
	if elem.contributes_script_blocking_style_sheet() {
		// TODO: append element to its node document's script-blocking style sheet set.
		panic("TODO[https://html.spec.whatwg.org/multipage/semantics.html#update-a-style-block]")
	}
	// S8.
	// If element's media attribute's value matches the environment and element is potentially render-blocking, then block rendering on element.

	// TODO: Specs has extra steps after critical subresources has been loaded, but they don't seem *that* important right now
	// (Mostly related to render blocking)

}
