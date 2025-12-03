package csssyntax

import (
	"github.com/inseo-oh/yw/css/display"
)

// https://www.w3.org/TR/css-display-3/#typedef-display-outside
func (ts *tokenStream) parseDisplayOutside() (display.OuterMode, bool) {
	if err := ts.consumeIdentTokenWith("block"); err == nil {
		return display.Block, true
	} else if err := ts.consumeIdentTokenWith("inline"); err == nil {
		return display.Inline, true
	} else if err := ts.consumeIdentTokenWith("run-in"); err == nil {
		return display.RunIn, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-display-3/#typedef-display-inside
func (ts *tokenStream) parseDisplayInside() (display.InnerMode, bool) {
	if err := ts.consumeIdentTokenWith("flow"); err == nil {
		return display.Flow, true
	} else if err := ts.consumeIdentTokenWith("flow-root"); err == nil {
		return display.FlowRoot, true
	} else if err := ts.consumeIdentTokenWith("table"); err == nil {
		return display.Table, true
	} else if err := ts.consumeIdentTokenWith("flex"); err == nil {
		return display.Flex, true
	} else if err := ts.consumeIdentTokenWith("grid"); err == nil {
		return display.Grid, true
	} else if err := ts.consumeIdentTokenWith("ruby"); err == nil {
		return display.Ruby, true
	}
	return 0, false
}

// https://www.w3.org/TR/css-display-3/#propdef-display
func (ts *tokenStream) parseDisplay() (display.Display, bool) {
	// Try legacy keyword first --------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-legacy
	if err := ts.consumeIdentTokenWith("inline-block"); err == nil {
		return display.Display{Mode: display.OuterInnerMode, OuterMode: display.Inline, InnerMode: display.FlowRoot}, true
	} else if err := ts.consumeIdentTokenWith("inline-table"); err == nil {
		return display.Display{Mode: display.OuterInnerMode, OuterMode: display.Inline, InnerMode: display.Table}, true
	} else if err := ts.consumeIdentTokenWith("inline-flex"); err == nil {
		return display.Display{Mode: display.OuterInnerMode, OuterMode: display.Inline, InnerMode: display.Flex}, true
	} else if err := ts.consumeIdentTokenWith("inline-grid"); err == nil {
		return display.Display{Mode: display.OuterInnerMode, OuterMode: display.Inline, InnerMode: display.Grid}, true
	}
	// Try <display-outside> <display-inside> ------------------------------
	gotOuterMode, gotInnerMode := false, false
	var outerMode display.OuterMode
	var innerMode display.InnerMode
	for !gotOuterMode || !gotInnerMode {
		gotSomething := false
		if !gotOuterMode {
			outerMode, gotOuterMode = ts.parseDisplayOutside()
			if gotOuterMode {
				gotSomething = true
			}
		}
		if !gotInnerMode {
			innerMode, gotInnerMode = ts.parseDisplayInside()
			if gotInnerMode {
				gotSomething = true
			}
		}
		if !gotSomething {
			break
		}
	}
	if gotOuterMode || gotInnerMode {
		if !gotInnerMode {
			innerMode = display.Flow
		} else if !gotOuterMode {
			if innerMode == display.Ruby {
				outerMode = display.Inline
			} else {
				outerMode = display.Block
			}
		}
		return display.Display{Mode: display.OuterInnerMode, OuterMode: outerMode, InnerMode: innerMode}, true
	}
	// Try display-listitem ------------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-listitem
	// TODO

	// Try display-internal ------------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-internal

	if err := ts.consumeIdentTokenWith("table-row-group"); err == nil {
		return display.Display{Mode: display.TableRowGroup}, true
	} else if err := ts.consumeIdentTokenWith("table-header-group"); err == nil {
		return display.Display{Mode: display.TableHeaderGroup}, true
	} else if err := ts.consumeIdentTokenWith("table-footer-group"); err == nil {
		return display.Display{Mode: display.TableFooterGroup}, true
	} else if err := ts.consumeIdentTokenWith("table-row"); err == nil {
		return display.Display{Mode: display.TableRow}, true
	} else if err := ts.consumeIdentTokenWith("table-cell"); err == nil {
		return display.Display{Mode: display.TableCell}, true
	} else if err := ts.consumeIdentTokenWith("table-column-group"); err == nil {
		return display.Display{Mode: display.TableColumnGroup}, true
	} else if err := ts.consumeIdentTokenWith("table-column"); err == nil {
		return display.Display{Mode: display.TableColumn}, true
	} else if err := ts.consumeIdentTokenWith("table-caption"); err == nil {
		return display.Display{Mode: display.TableCaption}, true
	} else if err := ts.consumeIdentTokenWith("ruby-base"); err == nil {
		return display.Display{Mode: display.RubyBase}, true
	} else if err := ts.consumeIdentTokenWith("ruby-text"); err == nil {
		return display.Display{Mode: display.RubyText}, true
	} else if err := ts.consumeIdentTokenWith("ruby-base-container"); err == nil {
		return display.Display{Mode: display.RubyBaseContainer}, true
	} else if err := ts.consumeIdentTokenWith("ruby-text-container"); err == nil {
		return display.Display{Mode: display.RubyTextContainer}, true
	}

	// Try display-box -----------------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-box

	if err := ts.consumeIdentTokenWith("contents"); err == nil {
		return display.Display{Mode: display.Contents}, true
	} else if err := ts.consumeIdentTokenWith("none"); err == nil {
		return display.Display{Mode: display.DisplayNone}, true
	}

	return display.Display{}, false
}

func (ts *tokenStream) parseVisibility() (display.Visibility, bool) {
	// Try legacy keyword first --------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-legacy
	if err := ts.consumeIdentTokenWith("visible"); err == nil {
		return display.Visible, true
	} else if err := ts.consumeIdentTokenWith("hidden"); err == nil {
		return display.Hidden, true
	} else if err := ts.consumeIdentTokenWith("collapse"); err == nil {
		return display.Collapse, true
	}
	return 0, false
}
