//go:generate go run ./gen
package props

type ComputedStyleSetSource interface {
	ComputedStyleSet() *ComputedStyleSet
	ParentSource() ComputedStyleSetSource
}

type PropertyValue interface {
	String() string
}

type Descriptor struct {
	Initial   PropertyValue
	ApplyFunc func(dest *ComputedStyleSet, value any)
}
