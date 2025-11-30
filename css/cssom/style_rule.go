package cssom

import "yw/css/selector"

type StyleRule struct {
	SelectorList []selector.Selector
	Declarations []Declaration
	AtRules      []AtRule
}
