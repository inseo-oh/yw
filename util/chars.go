package util

import (
	"regexp"
	"slices"
	"strings"
)

var AsciiUppercaseRegex = regexp.MustCompile(`[A-Z]`)
var AsciiLowercaseRegex = regexp.MustCompile(`[a-z]`)
var AsciiAlphaRegex = regexp.MustCompile(`[A-Za-z]`)
var AsciiDigitRegex = regexp.MustCompile(`[0-9]`)
var AsciiAlphanumericRegex = regexp.MustCompile(`[A-Za-z0-9]`)
var AsciiUpperHexDigitRegex = regexp.MustCompile(`[A-F0-9]`)
var AsciiLowerHexDigitRegex = regexp.MustCompile(`[a-f0-9]`)
var AsciiHexDigitRegex = regexp.MustCompile(`[A-Fa-f0-9]`)

func IsLeadingSurrogateChar(c rune) bool {
	return (0xd800 <= c) && (c <= 0xdbff)
}
func IsTrailingSurrogateChar(c rune) bool {
	return (0xdc00 <= c) && (c <= 0xdfff)
}
func IsSurrogateChar(c rune) bool {
	return IsLeadingSurrogateChar(c) || IsTrailingSurrogateChar(c)
}
func IsC0ControlChar(c rune) bool {
	return (0x0000 <= c) && (c <= 0x001f)
}
func IsControlChar(c rune) bool {
	return IsC0ControlChar(c) || ((0x007f <= c) && (c <= 0x009f))
}
func IsAsciiWhitespace(c rune) bool {
	whitespaceCodes := []rune{0x0009, 0x000a, 0x000c, 0x000d}
	return slices.Contains(whitespaceCodes, c)
}
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
