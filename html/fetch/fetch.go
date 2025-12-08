// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_THIRDPARTY for third-party license information.

// Package urlfetch provides types and values related to [HTML Spec 2.5 Fetching resources].
//
// [HTML Spec 2.5 Fetching resources]: https://html.spec.whatwg.org/multipage/urls-and-fetching.html#fetching-resources
package fetch

import (
	"github.com/inseo-oh/yw/dom"
	"github.com/inseo-oh/yw/util"
)

// CorsSettings represents a [CORS settings attribute].
//
// [CORS settings attribute]: https://html.spec.whatwg.org/multipage/urls-and-fetching.html#cors-settings-attribute
type CorsSettings uint8

const (
	CorsNone           CorsSettings = iota // <missing value default>
	CorsAnonymous                          // anonymous <invalid value default, empty value default>
	CorsUseCredentials                     // use-credentials
)

// Returns [CorsSettings] from element elem's attribute named attrName.
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

// FetchPriority represents a [fetch priority attribute].
//
// [fetch priority attribute]: https://html.spec.whatwg.org/multipage/urls-and-fetching.html#fetch-priority-attribute
type FetchPriority uint8

const (
	FetchPriorityHigh FetchPriority = iota // high
	FetchPriorityLow                       // low
	FetchPriorityAuto                      // auto <missing value default, invalid value default>
)

// Returns [FetchPriority] from element elem's attribute named attrName.
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
