package libhtml

import (
	cm "yw/libcommon"
)

// https://html.spec.whatwg.org/multipage/urls-and-fetching.html#cors-settings-attribute
type html_cors_settings uint8

const (
	html_cors_settings_no_cors         = html_cors_settings(iota) // https://html.spec.whatwg.org/multipage/urls-and-fetching.html#attr-crossorigin-none
	html_cors_settings_anonymous                                  // https://html.spec.whatwg.org/multipage/urls-and-fetching.html#attr-crossorigin-anonymous
	html_cors_settings_use_credentials                            // https://html.spec.whatwg.org/multipage/urls-and-fetching.html#attr-crossorigin-use-credentials
)

func html_cors_settings_from_attr(elem dom_Element, attr_name string) html_cors_settings {
	if attr, ok := elem.get_attribute_without_namespace(attr_name); ok {
		if attr == "" {
			return html_cors_settings_anonymous
		}
		switch cm.ToAsciiLowercase(attr) {
		case "anonymous":
			return html_cors_settings_anonymous
		case "use-credentials":
			return html_cors_settings_use_credentials
		}
		return html_cors_settings_anonymous
	} else {
		return html_cors_settings_no_cors
	}
}

// https://html.spec.whatwg.org/multipage/urls-and-fetching.html#fetch-priority-attribute
type html_fetch_priority uint8

const (
	html_fetch_priority_high = html_fetch_priority(iota) // https://html.spec.whatwg.org/multipage/urls-and-fetching.html#attr-fetchpriority-high-state
	html_fetch_priority_low                              // https://html.spec.whatwg.org/multipage/urls-and-fetching.html#attr-fetchpriority-low-state
	html_fetch_priority_auto                             // https://html.spec.whatwg.org/multipage/urls-and-fetching.html#attr-fetchpriority-auto-state
)

func html_fetch_priority_from_attr(elem dom_Element, attr_name string) html_fetch_priority {
	if attr, ok := elem.get_attribute_without_namespace(attr_name); ok {
		switch cm.ToAsciiLowercase(attr) {
		case "high":
			return html_fetch_priority_high
		case "low":
			return html_fetch_priority_low
		case "auto":
			return html_fetch_priority_auto
		}
		return html_fetch_priority_auto
	} else {
		return html_fetch_priority_auto
	}
}
