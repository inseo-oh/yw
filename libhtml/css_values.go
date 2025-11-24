// Implementation of the CSS Values and Units Module Level 3 (https://www.w3.org/TR/css-values-3/)
package libhtml

import (
	"fmt"
	cm "yw/libcommon"
)

// Returns nil if not found
func (ts *css_token_stream) parse_number() *css_number {
	num_tk := ts.consume_token_with_type(css_token_type_number)
	if cm.IsNil(num_tk) {
		return nil
	}
	per := num_tk.(css_number_token)
	return &per.value
}

type css_length_resolvable interface {
	as_length(container_size css_number) css_length
	String() string
}

// https://www.w3.org/TR/css-values-3/#length-value
type css_length struct {
	value css_number
	unit  css_length_unit
}

type css_length_unit uint8

func css_length_from_px(px css_number) css_length {
	return css_length{px, css_length_unit_px}
}

func (l css_length) String() string                                 { return fmt.Sprintf("%v%v", l.value, l.unit) }
func (l css_length) as_length(container_size css_number) css_length { return l }
func (l css_length) length_to_px(container_size css_number) float64 {
	switch l.unit {
	case css_length_unit_px:
		return l.value.to_float()
	case css_length_unit_em:
		return container_size.to_float() * l.value.to_float()
	default:
		panic("TODO")
	}
}

// Returns nil if not found
//
// allow_zero_shorthand should not be set if the property(such as line-height) also accepts number token.
// (In that case, 0 should be parsed as <number 0>, not <length 0>)
func (ts *css_token_stream) parse_length(allow_zero_shorthand bool) (*css_length, error) {
	dim_tk := ts.consume_token_with_type(css_token_type_dimension)
	if cm.IsNil(dim_tk) {
		if allow_zero_shorthand {
			old_cursor := ts.cursor
			num_tk := ts.consume_token_with_type(css_token_type_number)
			if cm.IsNil(num_tk) || !num_tk.(css_number_token).value.equals(css_number_from_int(0)) {
				ts.cursor = old_cursor
			} else {
				return &css_length{css_number_from_int(0), css_length_unit_px}, nil
			}
		}

		return nil, nil
	}
	dim := dim_tk.(css_dimension_token)
	var unit css_length_unit
	switch dim.unit {
	case "em":
		unit = css_length_unit_em
	case "ex":
		unit = css_length_unit_ex
	case "ch":
		unit = css_length_unit_ch
	case "rem":
		unit = css_length_unit_rem
	case "vw":
		unit = css_length_unit_vw
	case "vh":
		unit = css_length_unit_vh
	case "vmin":
		unit = css_length_unit_vmin
	case "vmax":
		unit = css_length_unit_vmax
	case "cm":
		unit = css_length_unit_cm
	case "mm":
		unit = css_length_unit_mm
	case "q":
		unit = css_length_unit_q
	case "pc":
		unit = css_length_unit_pc
	case "pt":
		unit = css_length_unit_pt
	case "px":
		unit = css_length_unit_px
	default:
		return nil, fmt.Errorf("unrecognized unit %s", dim.unit)
	}
	return &css_length{dim.value, unit}, nil
}

const (
	// https://www.w3.org/TR/css-values-3/#relative-lengths
	css_length_unit_em = css_length_unit(iota)
	css_length_unit_ex
	css_length_unit_ch
	css_length_unit_rem
	css_length_unit_vw
	css_length_unit_vh
	css_length_unit_vmin
	css_length_unit_vmax
	// https://www.w3.org/TR/css-values-3/#absolute-lengths
	css_length_unit_cm
	css_length_unit_mm
	css_length_unit_q
	css_length_unit_pc
	css_length_unit_pt
	css_length_unit_px
)

func (u css_length_unit) String() string {
	switch u {
	case css_length_unit_em:
		return "em"
	case css_length_unit_ex:
		return "ex"
	case css_length_unit_ch:
		return "ch"
	case css_length_unit_rem:
		return "rem"
	case css_length_unit_vw:
		return "vw"
	case css_length_unit_vh:
		return "vh"
	case css_length_unit_vmin:
		return "vmin"
	case css_length_unit_vmax:
		return "vmax"
	case css_length_unit_cm:
		return "cm"
	case css_length_unit_mm:
		return "mm"
	case css_length_unit_q:
		return "q"
	case css_length_unit_pc:
		return "pc"
	case css_length_unit_pt:
		return "pt"
	case css_length_unit_px:
		return "px"
	}
	return fmt.Sprintf("<unknown css_length_unit %d>", u)
}

// https://www.w3.org/TR/css-values-3/#percentage-value
type css_percentage struct {
	value css_number
}

func (len css_percentage) String() string { return fmt.Sprintf("%v%%", len.value) }

func (len css_percentage) as_length(container_size css_number) css_length { panic("TODO") }

// Returns nil if not found
func (ts *css_token_stream) parse_percentage() *css_percentage {
	per_tk := ts.consume_token_with_type(css_token_type_percentage)
	if cm.IsNil(per_tk) {
		return nil
	}
	per := per_tk.(css_percentage_token)
	return &css_percentage{per.value}
}

// https://www.w3.org/TR/css-values-3/#typedef-length-percentage
func (ts *css_token_stream) parse_length_or_percentage(allow_zero_shorthand bool) (css_length_resolvable, error) {
	if len, err := ts.parse_length(allow_zero_shorthand); len != nil {
		return len, nil
	} else if err != nil {
		return nil, err
	}
	if per := ts.parse_percentage(); per != nil {
		return per, nil
	}
	return nil, nil
}
