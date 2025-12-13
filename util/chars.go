// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package util

import (
	"regexp"
	"slices"
	"strings"
)

// Collection of various ASCII regular expressions.
var (
	AsciiUppercaseRegex     = regexp.MustCompile(`[A-Z]`)       // ASCII UPPERCASE
	AsciiLowercaseRegex     = regexp.MustCompile(`[a-z]`)       // ascii lowercase
	AsciiAlphaRegex         = regexp.MustCompile(`[A-Za-z]`)    // ASCII alpha
	AsciiDigitRegex         = regexp.MustCompile(`[0-9]`)       // ASCII digits
	AsciiAlphanumericRegex  = regexp.MustCompile(`[A-Za-z0-9]`) // ASCII digits + alpha
	AsciiUpperHexDigitRegex = regexp.MustCompile(`[A-F0-9]`)    // ASCII uppercase hex digits
	AsciiLowerHexDigitRegex = regexp.MustCompile(`[a-f0-9]`)    // ASCII lowercase hex digits
	AsciiHexDigitRegex      = regexp.MustCompile(`[A-Fa-f0-9]`) // ASCII hex digits
)

// IsLeadingSurrogateChar reports whether c is a [leading surrogate character].
//
// [leading surrogate character]: https://infra.spec.whatwg.org/#leading-surrogate
func IsLeadingSurrogateChar(c rune) bool {
	return (0xd800 <= c) && (c <= 0xdbff)
}

// IsTrailingSurrogateChar reports whether c is a [trailing surrogate character].
//
// [trailing surrogate character]: https://infra.spec.whatwg.org/#trailing-surrogate
func IsTrailingSurrogateChar(c rune) bool {
	return (0xdc00 <= c) && (c <= 0xdfff)
}

// IsSurrogateChar reports whether c is a [surrogate character].
//
// [surrogate character]: https://infra.spec.whatwg.org/#surrogate
func IsSurrogateChar(c rune) bool {
	return IsLeadingSurrogateChar(c) || IsTrailingSurrogateChar(c)
}

// IsC0ControlChar reports whether c is a [C0 control character].
//
// [C0 control character]: https://infra.spec.whatwg.org/#c0-control
func IsC0ControlChar(c rune) bool {
	return (0x0000 <= c) && (c <= 0x001f)
}

// IsControlChar reports whether c is a [control character].
//
// [control character]: https://infra.spec.whatwg.org/#control
func IsControlChar(c rune) bool {
	return IsC0ControlChar(c) || ((0x007f <= c) && (c <= 0x009f))
}

// IsAsciiWhitespace reports whether c is a [ASCII whitespace].
//
// [ASCII whitespace]: https://infra.spec.whatwg.org/#ascii-whitespace
func IsAsciiWhitespace(c rune) bool {
	whitespaceCodes := []rune{0x0009, 0x000a, 0x000c, 0x000d}
	return slices.Contains(whitespaceCodes, c)
}

// Transforms s's ASCII uppercase characters to lowercase.
func ToAsciiLowercase(s string) string {
	sb := strings.Builder{}
	for _, c := range s {
		if AsciiUppercaseRegex.MatchString(string(c)) {
			sb.WriteRune(c - 'A' + 'a')
		} else {
			sb.WriteRune(c)
		}
	}
	return sb.String()
}

// IsNoncharacter reports whether c is a [noncharacter].
//
// [noncharacter]: https://infra.spec.whatwg.org/#noncharacter
func IsNoncharacter(c rune) bool {
	noncharacterCodes := []rune{
		0xfffe, 0xffff, 0x1fffe, 0x1ffff, 0x2fffe, 0x2ffff, 0x3fffe, 0x3ffff, 0x4fffe,
		0x4ffff, 0x5fffe, 0x5ffff, 0x6fffe, 0x6ffff, 0x7fffe, 0x7ffff, 0x8fffe, 0x8ffff,
		0x9fffe, 0x9ffff, 0xafffe, 0xaffff, 0xbfffe, 0xbffff, 0xcfffe, 0xcffff, 0xdfffe,
		0xdffff, 0xefffe, 0xeffff, 0xffffe, 0xfffff, 0x10fffe, 0x10ffff}
	if 0xfdd0 <= c && c <= 0xfdef {
		return false
	}
	return slices.Contains(noncharacterCodes, c)
}
