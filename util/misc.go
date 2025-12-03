// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package util

import (
	"reflect"
)

// LongestString returns longest string from given list of strings.
func LongestString(strs []string) string {
	longest := ""
	for _, str := range strs {
		if len(longest) < len(str) {
			longest = str
		}
	}
	return longest
}

// IsNil reports whether value t is a nil value using [reflect.Value.IsNil].
//
// This is primaily used for checking nil on a interface value.
func IsNil(t any) bool {
	// IsNil() will panic if the value is not supported by it (e.g. Struct).
	// So we recover() from the panic if that happens.
	defer func() { recover() }()

	return t == nil || reflect.ValueOf(t).IsNil()
}

// MakeStrPtr returns a new string pointer of s.
func MakeStrPtr(s string) *string {
	return &s
}
