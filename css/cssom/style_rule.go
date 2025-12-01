package cssom

import "github.com/inseo-oh/yw/css/selector"

// StyleRule represents a CSS style rule (e.g. div { font-size: 40px; }).
type StyleRule struct {
	SelectorList []selector.Selector // List of selectors, used to select elements
	Declarations []Declaration       // List of declarations
	AtRules      []AtRule            // List of at-rules
}
