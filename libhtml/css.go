package libhtml

import (
	"fmt"
	"log"
)

type css_number struct {
	tp    css_number_type
	value any
}

type css_number_type uint8

const (
	css_number_type_int = css_number_type(iota)
	css_number_type_float
)

func (n css_number) to_int() int64 {
	if n.tp == css_number_type_float {
		return int64(n.to_float())
	} else {
		return n.value.(int64)
	}
}
func (n css_number) to_float() float64 {
	if n.tp == css_number_type_int {
		return float64(n.to_int())
	} else {
		return n.value.(float64)
	}
}

func (n css_number) equals(other css_number) bool {
	is_float := (n.tp == css_number_type_float) || (other.tp == css_number_type_float)
	if is_float {
		return n.to_float() == other.to_float()
	} else {
		return n.to_int() == other.to_int()
	}
}

func (n css_number) clamp(min, max css_number) css_number {
	// Using float all the time is probably fine, but let's avoid it if we can.
	is_float := (n.tp == css_number_type_float) || (min.tp == css_number_type_float) || (max.tp == css_number_type_float)
	if is_float {
		if n.to_float() < min.to_float() {
			return min
		} else if max.to_float() < n.to_float() {
			return max
		}
	} else {
		if n.to_int() < min.to_int() {
			return min
		} else if max.to_int() < n.to_int() {
			return max
		}
	}
	return n
}

func (n css_number) String() string {
	if n.tp == css_number_type_float {
		return fmt.Sprintf("%f", n.to_float())
	}
	return fmt.Sprintf("%d", n.to_int())
}

func css_number_from_int(v int64) css_number {
	return css_number{css_number_type_int, v}
}
func css_number_from_float(v float64) css_number {
	return css_number{css_number_type_float, v}
}

type css_style_rule struct {
	selector_list []css_complex_selector
	declarations  []css_declaration
	at_rules      []css_at_rule
}

func (r css_style_rule) apply_style_rules(roots []dom_Node) {
	// First we figure out where this style rule should be applied to
	selected_elements := css_match_selector_against_tree(r.selector_list, roots)

	for _, node := range selected_elements {
		elem := node.(dom_Element)
		for _, decl := range r.declarations {
			decl.apply_style_rules(elem)
		}

	}
}

type css_declaration struct {
	name  string
	value css_property_value
}

func (d css_declaration) apply_style_rules(elem dom_Element) {
	desc := css_property_descriptors_map[d.name]
	if desc.apply_func == nil {
		log.Printf("TODO: CSS Property %s is recognized but not supported yet. (Missing apply_func() function)", d.name)
		return
	}
	desc.apply_func(elem.get_computed_style_set(), d.value)
}

type css_at_rule struct {
	name    string
	prelude []css_token
	value   []css_token
}

type css_computed_style_set struct {
	display *css_display
	width   *css_size_value
	height  *css_size_value
}

func (css *css_computed_style_set) init_with_initial_value(name string) css_property_value {
	desc, ok := css_property_descriptors_map[name]
	if !ok {
		log.Panicf("attempted to initialize property '%s', but there's no such property", name)
	}
	return desc.initial
}
func (css *css_computed_style_set) get_display() css_display {
	if css.display == nil {
		v := css.init_with_initial_value("display").(css_display)
		css.display = &v
	}
	return *css.display
}
func (css *css_computed_style_set) get_width() css_size_value {
	if css.width == nil {
		v := css.init_with_initial_value("width").(css_size_value)
		css.width = &v
	}
	return *css.width
}
func (css *css_computed_style_set) get_height() css_size_value {
	if css.height == nil {
		v := css.init_with_initial_value("height").(css_size_value)
		css.height = &v
	}
	return *css.height
}
