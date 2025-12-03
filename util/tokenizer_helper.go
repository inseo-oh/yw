// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package util

import (
	"regexp"
	"strings"
)

// TokenizerHelper is helper type that provides useful functions and state for
// tokenizers.
type TokenizerHelper struct {
	Str    []rune // List of characters
	Cursor int    // Next position of character that will be read from
}

// RemainingChars returns remaining characters.
func (t TokenizerHelper) RemainingChars() []rune {
	return t.Str[t.Cursor:]
}

// IsEof reports whether tokenizer has reached the end of string.
func (t TokenizerHelper) IsEof() bool {
	return len(t.Str) <= int(t.Cursor)
}

// PeekChar returns next character that will be read, without advancing
// t's Cursor.
func (t TokenizerHelper) PeekChar() rune {
	if t.IsEof() {
		return -1
	}
	return t.Str[t.Cursor]
}

// PeekChar consumes next character and advances cursor by one, and returns
// consumed character.
//
// When t has reached the end already, it returns -1.
func (t *TokenizerHelper) ConsumeChar() rune {
	if t.IsEof() {
		return -1
	}
	char := t.PeekChar()
	t.Cursor++
	return char
}

// ConsumeCharIfMatchesOneOf consumes next character if it matches one of
// characters in chars, and returns consumed character.
//
// When t has reached the end already, or the character was not found,
// it returns -1.
func (t *TokenizerHelper) ConsumeCharIfMatchesOneOf(chars string) rune {
	if t.IsEof() {
		return -1
	}
	char := t.PeekChar()
	for _, c := range chars {
		if c == char {
			t.ConsumeChar()
			return c
		}
	}
	return -1
}

// Same as [ConsumeCharIfMatchesOneOf], but only accepts single character.
//
// TODO(ois): Could we make ConsumeCharIfMatches return a boolean instead?
func (t *TokenizerHelper) ConsumeCharIfMatches(char rune) rune {
	return t.ConsumeCharIfMatchesOneOf(string([]rune{char}))
}

// MatchFlags is used to specify flags for [ConsumeStrIfMatchesOneOf] and
// [ConsumeStrIfMatches].
type MatchFlags int

const (
	AsciiCaseInsensitive MatchFlags = 1 << 0
)

// ConsumeStrIfMatchesOneOf consumes sequence of one or more next characters
// if the sequence matches one of strings in strs, and returns consumed string.
//
// When t has reached the end already, or none of strs match, it returns
// an empty string.
func (t *TokenizerHelper) ConsumeStrIfMatchesOneOf(strs []string, flags MatchFlags) string {
	if t.IsEof() {
		return ""
	}
	maxLen := 0
	for _, s := range strs {
		l := len([]rune(s))
		if maxLen < l {
			maxLen = l
		}
	}
	remainingChars := t.RemainingChars()
	if len(remainingChars) <= maxLen {
		maxLen = len(remainingChars)
	}
	remaining := string(remainingChars[:maxLen])
	if (flags & AsciiCaseInsensitive) != 0 {
		remaining = ToAsciiLowercase(remaining)
	}

	matchMaxLen := 0
	resultStr := ""
	for _, s := range strs {
		if (flags & AsciiCaseInsensitive) != 0 {
			s = ToAsciiLowercase(s)
		}
		if strings.HasPrefix(remaining, s) {
			myLen := len([]rune(s))
			if matchMaxLen < myLen {
				resultStr = s
				matchMaxLen = myLen
			}
		}
	}
	t.Cursor += matchMaxLen
	return resultStr
}

// Same as [ConsumeStrIfMatchesOneOf], but only accepts single string.
//
// TODO(ois): Could we make ConsumeStrIfMatches return a boolean instead?
func (t *TokenizerHelper) ConsumeStrIfMatches(str string, flags MatchFlags) string {
	return t.ConsumeStrIfMatchesOneOf([]string{str}, flags)
}

// ConsumeStrIfMatchesOneOfR consumes sequence of one or more next characters
// if the sequence matches one of regular expressions in pats, and returns
// consumed string.
//
// When t has reached the end already, or none of pats match, it returns
// an empty string.
func (t *TokenizerHelper) ConsumeStrIfMatchesOneOfR(pats []regexp.Regexp) string {
	remaining := string(t.RemainingChars())

	matchMaxLen := 0
	resultStr := ""
	for _, p := range pats {
		if loc := p.FindStringIndex(remaining); loc != nil && loc[0] == 0 {
			s := p.FindString(remaining)
			myLen := len([]rune(s))
			if matchMaxLen < myLen {
				resultStr = s
				matchMaxLen = myLen
			}
		}
	}
	t.Cursor += matchMaxLen
	return resultStr
}

// Same as [ConsumeStrIfMatchesOneOfR], but only accepts single character.
//
// TODO(ois): Could we make ConsumeStrIfMatchesR return a boolean instead?
func (t *TokenizerHelper) ConsumeStrIfMatchesR(pat regexp.Regexp) string {
	return t.ConsumeStrIfMatchesOneOfR([]regexp.Regexp{pat})
}

// Lookahead calls cb, which can freely use tokenizer, and Lookahead restores
// Cursor before returning, so t's state does not change from caller's standpoint.
//
//	t := TokenizerHelper{Str: "hello"}
//	hasHello := false
//	t.Lookahead(func() {
//		return t.ConsumeStrIfMatches("hello")
//	})
//	if !hasHello {
//		log.Error("Could not find hello")
//	}
func (t *TokenizerHelper) Lookahead(cb func()) {
	oldCursor := t.Cursor
	cb()
	t.Cursor = oldCursor
}
