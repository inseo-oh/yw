package urlfetch

import (
	"yw/dom"
	"yw/util"
)

// https://html.spec.whatwg.org/multipage/urls-and-fetching.html#cors-settings-attribute
type CorsSettings uint8

const (
	CorsNone           = CorsSettings(iota) // https://html.spec.whatwg.org/multipage/urls-and-fetching.html#attr-crossorigin-none
	CorsAnonymous                           // https://html.spec.whatwg.org/multipage/urls-and-fetching.html#attr-crossorigin-anonymous
	CorsUseCredentials                      // https://html.spec.whatwg.org/multipage/urls-and-fetching.html#attr-crossorigin-use-credentials
)

func CorsSettingsFromAttr(elem dom.Element, attrName string) CorsSettings {
	if attr, ok := elem.AttrWithoutNamespace(attrName); ok {
		if attr == "" {
			return CorsAnonymous
		}
		switch util.ToAsciiLowercase(attr) {
		case "anonymous":
			return CorsAnonymous
		case "use-credentials":
			return CorsUseCredentials
		}
		return CorsAnonymous
	} else {
		return CorsNone
	}
}

// https://html.spec.whatwg.org/multipage/urls-and-fetching.html#fetch-priority-attribute
type FetchPriority uint8

const (
	FetchPriorityHigh = FetchPriority(iota) // https://html.spec.whatwg.org/multipage/urls-and-fetching.html#attr-fetchpriority-high-state
	FetchPriorityLow                        // https://html.spec.whatwg.org/multipage/urls-and-fetching.html#attr-fetchpriority-low-state
	FetchPriorityAuto                       // https://html.spec.whatwg.org/multipage/urls-and-fetching.html#attr-fetchpriority-auto-state
)

func FetchPriorityFromAttr(elem dom.Element, attrName string) FetchPriority {
	if attr, ok := elem.AttrWithoutNamespace(attrName); ok {
		switch util.ToAsciiLowercase(attr) {
		case "high":
			return FetchPriorityHigh
		case "low":
			return FetchPriorityLow
		case "auto":
			return FetchPriorityAuto
		}
		return FetchPriorityAuto
	} else {
		return FetchPriorityAuto
	}
}
