// Package text provides types and values for [CSS Text Module Level 3].
//
// [CSS Text Module Level 3]: https://www.w3.org/TR/css-text-3
package text

import (
	"log"
	"strings"
)

// Transform represents value of [CSS text-transform] property.
//
// [CSS text-transform]: https://www.w3.org/TR/css-text-3/#propdef-text-transform
type Transform struct {
	Type         TransformType // Type of text transform
	FullWidth    bool          // Corresponds to "full-width" of text-transform
	FullSizeKana bool          // Corresponds to "full-size-kana" of text-transform
}

// Type of [Transform]
type TransformType uint8

const (
	NoTransform  TransformType = iota // No transform
	OriginalCaps                      // Apply transform, but don't change cases
	Capitalize                        // text-transform: Captialize
	Uppercase                         // text-transform: UPPERCASE
	Lowercase                         // text-transform: lowercase
)

func (t Transform) String() string {
	sb := strings.Builder{}
	switch t.Type {
	case NoTransform:
		return "none"
	case OriginalCaps:
		break
	case Capitalize:
		sb.WriteString("capitalize ")
	case Lowercase:
		sb.WriteString("lowercase ")
	case Uppercase:
		sb.WriteString("uppercase ")
	default:
		log.Panicf("<bad Transform type %v>", t.Type)
	}
	if t.FullWidth {
		sb.WriteString("full-width ")
	}
	if t.FullSizeKana {
		sb.WriteString("full-size-kana ")
	}
	return strings.TrimSpace(sb.String())
}

// Apply applies text transform to given text.
func (t Transform) Apply(text string) string {
	outText := text
	switch t.Type {
	case NoTransform:
		return text
	case OriginalCaps:
		break
	case Capitalize:
		runes := []rune(text)
		outText = strings.ToUpper(string(runes[0])) + string(runes[1:])
	case Lowercase:
		outText = strings.ToLower(text)
	case Uppercase:
		outText = strings.ToUpper(text)
	default:
		log.Panicf("<bad Transform type %v>", t.Type)
	}
	if t.FullWidth {
		log.Printf("TODO: %v: Support full-width", t)
	}
	if t.FullSizeKana {
		log.Printf("TODO: %v: Support full-size-kana", t)
	}
	return outText
}
