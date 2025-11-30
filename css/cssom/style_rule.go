package cssom

import "github.com/inseo-oh/yw/css/selector"

type StyleRule struct {
	SelectorList []selector.Selector
	Declarations []Declaration
	AtRules      []AtRule
}
