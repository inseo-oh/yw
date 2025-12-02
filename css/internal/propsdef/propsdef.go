// Package propsdef contains definitions used by code generators that generates
// CSS property code.
//
// # Two kinds of properties
//
// There are two kinds of properties:
//
//   - Simple: Accepts single value of given type. See [SimpleProp].
//   - Shorthand: Shorthand for set of Simple properties. See [ShorthandSidesProp]
//     and [ShorthandAnyProp].
//     Note that these properties also generate new Go types and parser
//     function for the shorthand type.
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

// CssType represents type information for corresponding CSS value type.
type CssType struct {
	TypeName        string // Go type name
	ParseMethodName string // Parser method in csssyntax.tokenStream
}

// CssProp represents a CSS property. For most property, you probably want
// [SimpleProp], which implements this interface.
type CssProp interface {
	// Returns name of the property
	PropName() string

	// Returns type of the property value
	//
	// outsidePropsPkg should be set if code that's generated is used outside of props package.
	// This appends "props." to every shorthand types, which are defined inside the props package.
	//
	// outsidePropsPkg has no effect for non-shorthand types.
	PropType(outsidePropsPkg bool) CssType
	// Returns initial value (Go expression)
	PropInitialValue(outsidePropsPkg bool) string
	// Returns whether property is inherited when it's missing value.
	IsInheritable() bool
	// Returns whteher property is shorthand property.
	IsShorthand() bool
}

// Simple CSS property that accepts a single value.
type SimpleProp struct {
	Name         string  // Name of the property
	Type         CssType // Type of the property value
	InitialValue string  // Initial value (Go expression)
	Inheritable  bool    // Can inherit?
}

func (p SimpleProp) PropName() string { return p.Name }

func (p SimpleProp) PropType(outsidePropsPkg bool) CssType { return p.Type }

func (p SimpleProp) PropInitialValue(outsidePropsPkg bool) string { return p.InitialValue }

func (p SimpleProp) IsInheritable() bool { return p.Inheritable }

func (p SimpleProp) IsShorthand() bool { return false }

// ShorthandSidesProp represents shorthand property, accepting top, right, bottom, left values.
// These properties can accept 1~4 values:
//
//   - 1 value: <top-right-bottom-left> (Sets all four sides at once)
//   - 2 value: <top-bottom> <right-left>
//   - 3 value: <top> <right-left> <bottom>
//   - 4 value: <top> <right> <bottom> <left>
//
// EXCEPTION: When this is used inside [ShorthandAnyProp], it only accepts single value.
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

// TypeName returns name that will be used to generate Go types for the shorthand type.
func (p ShorthandSidesProp) TypeName(outsidePropsPkg bool) string {
	prefix := ""
	if outsidePropsPkg {
		prefix = "props."
	}
	return fmt.Sprintf("%s%sShorthand", prefix, camelCaseName(p.Name, true))
}

// ParseMethodName returns name that will be used to generate parser function for the shorthand type.
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

// ShorthandAnyProp represents shorthand property, accepting any of accepted
// types, in any order. If given value is missing values for certain properties,
// they are filled with default value for the property.
//
// Note that [ShorthandSidesProp] properties will only accept single value when
// nested within this type of property. (border is an example of this)
//
// Examples: border, font
type ShorthandAnyProp struct {
	Name        string // Name of the property
	Props       []CssProp
	Inheritable bool // Can inherit?
}

// TypeName returns name that will be used to generate Go types for the shorthand type.
func (p ShorthandAnyProp) TypeName(outsidePropsPkg bool) string {
	prefix := ""
	if outsidePropsPkg {
		prefix = "props."
	}
	return fmt.Sprintf("%s%sShorthand", prefix, camelCaseName(p.Name, true))
}

// ParseMethodName returns name that will be used to generate parser function for the shorthand type.
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

// Returns Go CamelCase identifier name for given CSS property's name.
// For example, this turns "text-decoration-color" into TextDecorationColor.
func GoIdentNameOfProp(prop CssProp) string {
	s := prop.PropName()
	if prop.IsShorthand() {
		// We explicitly name shorthand properties ~Shorthand to avoid confusion.
		return camelCaseName(s, true) + "Shorthand"
	} else {
		return camelCaseName(s, true)
	}
}
