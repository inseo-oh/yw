package util

import (
	"regexp"
	"strings"
)

type TokenizerHelper struct {
	Str    []rune
	Cursor int
}

func (t TokenizerHelper) GetRemainingChars() []rune {
	return t.Str[t.Cursor:]
}

func (t TokenizerHelper) IsEof() bool {
	return len(t.Str) <= int(t.Cursor)
}

func (t TokenizerHelper) PeekChar() rune {
	if t.IsEof() {
		return -1
	}
	return t.Str[t.Cursor]
}

// Returns -1 on eof
func (t *TokenizerHelper) ConsumeChar() rune {
	if t.IsEof() {
		return -1
	}
	char := t.PeekChar()
	t.Cursor++
	return char
}

// Returns -1 if not found.
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

// Returns -1 if not found.
func (t *TokenizerHelper) ConsumeCharIfMatches(char rune) rune {
	return t.ConsumeCharIfMatchesOneOf(string([]rune{char}))
}

type MatchFlags int

const (
	MatchFlagsAsciiCaseInsensitive = MatchFlags(1 << 0)
)

// Returns empty string if not found
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
	remainingChars := t.GetRemainingChars()
	if len(remainingChars) <= maxLen {
		maxLen = len(remainingChars)
	}
	remaining := string(remainingChars[:maxLen])
	if (flags & MatchFlagsAsciiCaseInsensitive) != 0 {
		remaining = ToAsciiLowercase(remaining)
	}

	matchMaxLen := 0
	resultStr := ""
	for _, s := range strs {
		if (flags & MatchFlagsAsciiCaseInsensitive) != 0 {
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

// Returns empty string if not found.
func (t *TokenizerHelper) ConsumeStrIfMatches(str string, flags MatchFlags) string {
	return t.ConsumeStrIfMatchesOneOf([]string{str}, flags)
}

// Returns empty string if not found
func (t *TokenizerHelper) ConsumeStrIfMatchesOneOfR(pats []regexp.Regexp) string {
	remaining := string(t.GetRemainingChars())

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

// Returns empty string if not found.
func (t *TokenizerHelper) ConsumeStrIfMatchesR(pat regexp.Regexp) string {
	return t.ConsumeStrIfMatchesOneOfR([]regexp.Regexp{pat})
}

func (t *TokenizerHelper) Lookahead(cb func()) {
	oldCursor := t.Cursor
	cb()
	t.Cursor = oldCursor
}
