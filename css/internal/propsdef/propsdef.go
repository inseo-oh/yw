package propsdef

import (
	"fmt"
	"strings"
	"unicode"
)

func camelCaseName(name string, leadingCharUpper bool) string {
	sb := strings.Builder{}
	nextUpper := leadingCharUpper
	for _, c := range name {
		if c == '_' || c == '-' {
			nextUpper = true
		} else if nextUpper {
			sb.WriteRune(unicode.ToUpper(c))
			nextUpper = false
		} else {
			sb.WriteRune(c)
		}
	}
	return sb.String()
}

type CssType struct {
	TypeName        string // Go type name
	ParseMethodName string // Parser method in csssyntax.tokenStream
}

type CssProp interface {
	PropName() string
	// outsidePropsPkg should be set if code that's generated is used outside of props package.
	// This appends "props." to every shorthand types, which are defined inside the props package.
	//
	// outsidePropsPkg has no effect for non-shorthand types.
	PropType(outsidePropsPkg bool) CssType
	PropInitialValue(outsidePropsPkg bool) string
	IsInheritable() bool
	IsShorthand() bool
}

// A CSS property
type SimpleProp struct {
	Name         string  // Name of the property
	Type         CssType // Type of the property value
	InitialValue string  // Initial value (Go expression)
	Inheritable  bool    // Can inherit?
}

func (p SimpleProp) PropName() string                             { return p.Name }
func (p SimpleProp) PropType(outsidePropsPkg bool) CssType        { return p.Type }
func (p SimpleProp) PropInitialValue(outsidePropsPkg bool) string { return p.InitialValue }
func (p SimpleProp) IsInheritable() bool                          { return p.Inheritable }
func (p SimpleProp) IsShorthand() bool                            { return false }

// Shorthand properties for top, right, bottom, left values.
// These take 1~4 values of the same type, and assigned to top, right, bottom, left accordingly.
//
// EXCEPTION: When this is used inside shorthandAnyProp, it accepts single value that sets all four sides.
//
// Examples: border-color, padding, margin
type ShorthandSidesProp struct {
	Name        string  // Name of the property
	PropTop     CssProp // Top property
	PropRight   CssProp // Right property
	PropBottom  CssProp // Bottom property
	PropLeft    CssProp // Left property
	Inheritable bool    // Can inherit?
}

func (p ShorthandSidesProp) TypeName(outsidePropsPkg bool) string {
	prefix := ""
	if outsidePropsPkg {
		prefix = "props."
	}
	return fmt.Sprintf("%s%sShorthand", prefix, camelCaseName(p.Name, true))
}
func (p ShorthandSidesProp) ParseMethodName() string {
	return fmt.Sprintf("parse%s", camelCaseName(p.TypeName(false), true))
}
func (p ShorthandSidesProp) PropName() string { return p.Name }
func (p ShorthandSidesProp) PropType(outsidePropsPkg bool) CssType {
	return CssType{
		TypeName:        p.TypeName(outsidePropsPkg),
		ParseMethodName: p.ParseMethodName(),
	}
}
func (p ShorthandSidesProp) PropInitialValue(outsidePropsPkg bool) string {
	return fmt.Sprintf(
		"%s{Left: %s, Top: %s, Right: %s, Bottom: %s}",
		p.TypeName(outsidePropsPkg),
		p.PropTop.PropInitialValue(outsidePropsPkg),
		p.PropRight.PropInitialValue(outsidePropsPkg),
		p.PropBottom.PropInitialValue(outsidePropsPkg),
		p.PropLeft.PropInitialValue(outsidePropsPkg),
	)
}
func (p ShorthandSidesProp) IsInheritable() bool { return p.Inheritable }
func (p ShorthandSidesProp) IsShorthand() bool   { return true }

// Shorthand properties for multiple types of properties
// These simply take any of accepted properties in any other, and fills with default values for absent ones.
// Examples: border, font
type ShorthandAnyProp struct {
	Name        string // Name of the property
	Props       []CssProp
	Inheritable bool // Can inherit?
}

func (p ShorthandAnyProp) TypeName(outsidePropsPkg bool) string {
	prefix := ""
	if outsidePropsPkg {
		prefix = "props."
	}
	return fmt.Sprintf("%s%sShorthand", prefix, camelCaseName(p.Name, true))
}
func (p ShorthandAnyProp) ParseMethodName() string {

	return fmt.Sprintf("parse%s", camelCaseName(p.TypeName(false), true))
}
func (p ShorthandAnyProp) PropName() string { return p.Name }
func (p ShorthandAnyProp) PropType(outsidePropsPkg bool) CssType {
	return CssType{
		TypeName:        p.TypeName(outsidePropsPkg),
		ParseMethodName: p.ParseMethodName(),
	}
}
func (p ShorthandAnyProp) PropInitialValue(outsidePropsPkg bool) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("%s{", p.TypeName(outsidePropsPkg)))
	for i, prop := range p.Props {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%s: %s", GoIdentNameOfProp(prop), prop.PropInitialValue(outsidePropsPkg)))
	}
	sb.WriteString("}")
	return sb.String()
}
func (p ShorthandAnyProp) IsInheritable() bool { return p.Inheritable }
func (p ShorthandAnyProp) IsShorthand() bool   { return true }

func GoIdentNameOfProp(prop CssProp) string {
	s := prop.PropName()
	if prop.IsShorthand() {
		// We explicitly name shorthand properties ~Shorthand to avoid confusion.
		return camelCaseName(s, true) + "Shorthand"
	} else {
		return camelCaseName(s, true)
	}
}
