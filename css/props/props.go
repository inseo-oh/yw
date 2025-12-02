// Package props provides various CSS properties, and [ComputedStyleSet]
// containing computed values for properties.
package props

//go:generate go run ./gen

// Source of [ComputedStyleSet]. [ComputedStyleSet] is stored inside DOM element,
// and this type is used to avoid props package from depending on dom package.
type ComputedStyleSetSource interface {
	ComputedStyleSet() *ComputedStyleSet
	ParentSource() ComputedStyleSetSource
}

// PropertyValue represents an property value.
type PropertyValue interface {
	String() string
}

// Descriptor represents information about each property.
type Descriptor struct {
	Initial   PropertyValue
	ApplyFunc func(dest *ComputedStyleSet, value any)
}
