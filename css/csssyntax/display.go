package csssyntax

import (
	"errors"

	"github.com/inseo-oh/yw/css/display"
)

// https://www.w3.org/TR/css-display-3/#typedef-display-outside
func (ts *tokenStream) parseDisplayOutside() (display.OuterMode, error) {
	if err := ts.consumeIdentTokenWith("block"); err == nil {
		return display.Block, nil
	} else if err := ts.consumeIdentTokenWith("inline"); err == nil {
		return display.Inline, nil
	} else if err := ts.consumeIdentTokenWith("run-in"); err == nil {
		return display.RunIn, nil
	}
	return 0, errors.New("invalid display-outside value")
}

// https://www.w3.org/TR/css-display-3/#typedef-display-inside
func (ts *tokenStream) parseDisplayInside() (display.InnerMode, error) {
	if err := ts.consumeIdentTokenWith("flow"); err == nil {
		return display.Flow, nil
	} else if err := ts.consumeIdentTokenWith("flow-root"); err == nil {
		return display.FlowRoot, nil
	} else if err := ts.consumeIdentTokenWith("table"); err == nil {
		return display.Table, nil
	} else if err := ts.consumeIdentTokenWith("flex"); err == nil {
		return display.Flex, nil
	} else if err := ts.consumeIdentTokenWith("grid"); err == nil {
		return display.Grid, nil
	} else if err := ts.consumeIdentTokenWith("ruby"); err == nil {
		return display.Ruby, nil
	}
	return 0, errors.New("invalid display-inside value")
}

// https://www.w3.org/TR/css-display-3/#propdef-display
func (ts *tokenStream) parseDisplay() (display.Display, error) {
	// Try legacy keyword first --------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-legacy
	if err := ts.consumeIdentTokenWith("inline-block"); err == nil {
		return display.Display{Mode: display.OuterInnerMode, OuterMode: display.Inline, InnerMode: display.FlowRoot}, nil
	} else if err := ts.consumeIdentTokenWith("inline-table"); err == nil {
		return display.Display{Mode: display.OuterInnerMode, OuterMode: display.Inline, InnerMode: display.Table}, nil
	} else if err := ts.consumeIdentTokenWith("inline-flex"); err == nil {
		return display.Display{Mode: display.OuterInnerMode, OuterMode: display.Inline, InnerMode: display.Flex}, nil
	} else if err := ts.consumeIdentTokenWith("inline-grid"); err == nil {
		return display.Display{Mode: display.OuterInnerMode, OuterMode: display.Inline, InnerMode: display.Grid}, nil
	}
	// Try <display-outside> <display-inside> ------------------------------
	gotOuterMode, gotInnerMode := false, false
	var outerMode display.OuterMode
	var innerMode display.InnerMode
	for !gotOuterMode || !gotInnerMode {
		gotSomething := false
		var err error
		if !gotOuterMode {
			outerMode, err = ts.parseDisplayOutside()
			if err == nil {
				gotSomething = true
				gotOuterMode = true
			}
		}
		if !gotInnerMode {
			innerMode, err = ts.parseDisplayInside()
			if err == nil {
				gotSomething = true
				gotInnerMode = true
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
		return display.Display{Mode: display.OuterInnerMode, OuterMode: outerMode, InnerMode: innerMode}, nil
	}
	// Try display-listitem ------------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-listitem
	// TODO

	// Try display-internal ------------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-internal

	if err := ts.consumeIdentTokenWith("table-row-group"); err == nil {
		return display.Display{Mode: display.TableRowGroup}, nil
	} else if err := ts.consumeIdentTokenWith("table-header-group"); err == nil {
		return display.Display{Mode: display.TableHeaderGroup}, nil
	} else if err := ts.consumeIdentTokenWith("table-footer-group"); err == nil {
		return display.Display{Mode: display.TableFooterGroup}, nil
	} else if err := ts.consumeIdentTokenWith("table-row"); err == nil {
		return display.Display{Mode: display.TableRow}, nil
	} else if err := ts.consumeIdentTokenWith("table-cell"); err == nil {
		return display.Display{Mode: display.TableCell}, nil
	} else if err := ts.consumeIdentTokenWith("table-column-group"); err == nil {
		return display.Display{Mode: display.TableColumnGroup}, nil
	} else if err := ts.consumeIdentTokenWith("table-column"); err == nil {
		return display.Display{Mode: display.TableColumn}, nil
	} else if err := ts.consumeIdentTokenWith("table-caption"); err == nil {
		return display.Display{Mode: display.TableCaption}, nil
	} else if err := ts.consumeIdentTokenWith("ruby-base"); err == nil {
		return display.Display{Mode: display.RubyBase}, nil
	} else if err := ts.consumeIdentTokenWith("ruby-text"); err == nil {
		return display.Display{Mode: display.RubyText}, nil
	} else if err := ts.consumeIdentTokenWith("ruby-base-container"); err == nil {
		return display.Display{Mode: display.RubyBaseContainer}, nil
	} else if err := ts.consumeIdentTokenWith("ruby-text-container"); err == nil {
		return display.Display{Mode: display.RubyTextContainer}, nil
	}

	// Try display-box -----------------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-box

	if err := ts.consumeIdentTokenWith("contents"); err == nil {
		return display.Display{Mode: display.Contents}, nil
	} else if err := ts.consumeIdentTokenWith("none"); err == nil {
		return display.Display{Mode: display.DisplayNone}, nil
	}

	return display.Display{}, errors.New("invalid display value")
}

func (ts *tokenStream) parseVisibility() (display.Visibility, error) {
	// Try legacy keyword first --------------------------------------------
	// https://www.w3.org/TR/css-display-3/#typedef-display-legacy
	if err := ts.consumeIdentTokenWith("visible"); err == nil {
		return display.Visible, nil
	} else if err := ts.consumeIdentTokenWith("hidden"); err == nil {
		return display.Hidden, nil
	} else if err := ts.consumeIdentTokenWith("collapse"); err == nil {
		return display.Collapse, nil
	}
	return 0, errors.New("invalid visibility value")
}
