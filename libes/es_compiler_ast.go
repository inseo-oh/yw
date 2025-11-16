package libes

func es_make_code_for_ast_nodes(nodes []es_ast_node) []es_vm_instr {
	instrs := []es_vm_instr{}
	for _, node := range nodes {
		instrs = append(instrs, node.make_code()...)
	}
	return instrs
}

type es_ast_node interface {
	get_cursor_from() int
	get_cursor_to() int
	make_code() []es_vm_instr
}

type es_ast_identifier_reference_node struct {
	cursor_from, cursor_to int

	name string
}

func (n es_ast_identifier_reference_node) get_cursor_from() int     { return n.cursor_from }
func (n es_ast_identifier_reference_node) get_cursor_to() int       { return n.cursor_to }
func (n es_ast_identifier_reference_node) make_code() []es_vm_instr { panic("TODO") }

type es_ast_binding_identifier_node struct {
	cursor_from, cursor_to int

	name string
}

func (n es_ast_binding_identifier_node) get_cursor_from() int     { return n.cursor_from }
func (n es_ast_binding_identifier_node) get_cursor_to() int       { return n.cursor_to }
func (n es_ast_binding_identifier_node) make_code() []es_vm_instr { panic("TODO") }

type es_ast_paren_expr_node struct {
	cursor_from, cursor_to int

	node es_ast_node
}

func (n es_ast_paren_expr_node) get_cursor_from() int { return n.cursor_from }
func (n es_ast_paren_expr_node) get_cursor_to() int   { return n.cursor_to }
func (n es_ast_paren_expr_node) make_code() []es_vm_instr {
	return n.node.make_code()
}

type es_ast_literal_node struct {
	cursor_from, cursor_to int

	value es_value
}

func (n es_ast_literal_node) get_cursor_from() int { return n.cursor_from }
func (n es_ast_literal_node) get_cursor_to() int   { return n.cursor_to }
func (n es_ast_literal_node) make_code() []es_vm_instr {
	return []es_vm_instr{
		{es_vm_opcode_push, n.value},
	}
}

type es_ast_call_expr_node struct {
	cursor_from, cursor_to int

	callee es_ast_node
	args   es_ast_arguments
}
type es_ast_arguments struct {
	cursor_from, cursor_to int

	args      []es_ast_node
	rest_args es_ast_node // May be nil
}

func (n es_ast_call_expr_node) get_cursor_from() int     { return n.cursor_from }
func (n es_ast_call_expr_node) get_cursor_to() int       { return n.cursor_to }
func (n es_ast_call_expr_node) make_code() []es_vm_instr { panic("TODO") }

type es_ast_unary_op_node struct {
	cursor_from, cursor_to int

	tp   es_ast_unary_op_type
	node es_ast_node
}
type es_ast_unary_op_type uint8

const (
	es_ast_unary_op_type_plus = es_ast_unary_op_type(iota)
	es_ast_unary_op_type_neg
	es_ast_unary_op_type_bnot
	es_ast_unary_op_type_lnot
	es_ast_unary_op_type_await
)

func (n es_ast_unary_op_node) get_cursor_from() int { return n.cursor_from }
func (n es_ast_unary_op_node) get_cursor_to() int   { return n.cursor_to }

var es_ast_unary_op_to_vm_map = map[es_ast_unary_op_type]es_vm_opcode{
	es_ast_unary_op_type_plus:  es_vm_opcode_plus,
	es_ast_unary_op_type_neg:   es_vm_opcode_neg,
	es_ast_unary_op_type_bnot:  es_vm_opcode_bnot,
	es_ast_unary_op_type_lnot:  es_vm_opcode_lnot,
	es_ast_unary_op_type_await: es_vm_opcode_await,
}

func (n es_ast_unary_op_node) make_code() []es_vm_instr {
	instrs := []es_vm_instr{}
	instrs = append(instrs, n.node.make_code()...)
	opcode := es_ast_unary_op_to_vm_map[n.tp]
	instrs = append(instrs, es_vm_instr{opcode, nil})
	return instrs
}

type es_ast_binary_op_node struct {
	cursor_from, cursor_to int

	tp       es_ast_binary_op_type
	lhs_node es_ast_node
	rhs_node es_ast_node
}
type es_ast_binary_op_type uint8

const (
	es_ast_binary_op_type_exponent = es_ast_binary_op_type(iota)
	es_ast_binary_op_type_multiply
	es_ast_binary_op_type_divide
	es_ast_binary_op_type_modulo
	es_ast_binary_op_type_add
	es_ast_binary_op_type_subtract
	es_ast_binary_op_type_left_shift
	es_ast_binary_op_type_right_lshift
	es_ast_binary_op_type_right_ashift
	es_ast_binary_op_type_less_than
	es_ast_binary_op_type_greater_than
	es_ast_binary_op_type_less_than_or_equal
	es_ast_binary_op_type_greater_than_or_equal
	es_ast_binary_op_type_instanceof
	es_ast_binary_op_type_equal
	es_ast_binary_op_type_not_equal
	es_ast_binary_op_type_strict_equal
	es_ast_binary_op_type_strict_not_equal
	es_ast_binary_op_type_band
	es_ast_binary_op_type_bxor
	es_ast_binary_op_type_bor
	es_ast_binary_op_type_land
	es_ast_binary_op_type_lor
	es_ast_binary_op_type_coalesce
)

func (n es_ast_binary_op_node) get_cursor_from() int { return n.cursor_from }
func (n es_ast_binary_op_node) get_cursor_to() int   { return n.cursor_to }

var es_ast_binary_op_to_vm_map = map[es_ast_binary_op_type]es_vm_opcode{
	es_ast_binary_op_type_exponent:              es_vm_opcode_exponent,
	es_ast_binary_op_type_multiply:              es_vm_opcode_multiply,
	es_ast_binary_op_type_divide:                es_vm_opcode_divide,
	es_ast_binary_op_type_modulo:                es_vm_opcode_modulo,
	es_ast_binary_op_type_add:                   es_vm_opcode_add,
	es_ast_binary_op_type_subtract:              es_vm_opcode_subtract,
	es_ast_binary_op_type_left_shift:            es_vm_opcode_left_shift,
	es_ast_binary_op_type_right_ashift:          es_vm_opcode_right_ashift,
	es_ast_binary_op_type_right_lshift:          es_vm_opcode_right_lshift,
	es_ast_binary_op_type_band:                  es_vm_opcode_band,
	es_ast_binary_op_type_bxor:                  es_vm_opcode_bxor,
	es_ast_binary_op_type_bor:                   es_vm_opcode_bor,
	es_ast_binary_op_type_land:                  es_vm_opcode_land,
	es_ast_binary_op_type_lor:                   es_vm_opcode_lor,
	es_ast_binary_op_type_less_than:             es_vm_opcode_less_than,
	es_ast_binary_op_type_less_than_or_equal:    es_vm_opcode_less_than_or_equal,
	es_ast_binary_op_type_greater_than:          es_vm_opcode_greater_than,
	es_ast_binary_op_type_greater_than_or_equal: es_vm_opcode_greater_than_or_equal,
	es_ast_binary_op_type_equal:                 es_vm_opcode_equal,
	es_ast_binary_op_type_strict_equal:          es_vm_opcode_strict_equal,
	es_ast_binary_op_type_not_equal:             es_vm_opcode_not_equal,
	es_ast_binary_op_type_strict_not_equal:      es_vm_opcode_strict_not_equal,
	es_ast_binary_op_type_coalesce:              es_vm_opcode_coalesce,
}

func (n es_ast_binary_op_node) make_code() []es_vm_instr {
	instrs := []es_vm_instr{}
	instrs = append(instrs, n.rhs_node.make_code()...)
	instrs = append(instrs, n.lhs_node.make_code()...)
	opcode := es_ast_binary_op_to_vm_map[n.tp]
	instrs = append(instrs, es_vm_instr{opcode, nil})
	return instrs

}

type es_ast_cond_expr_node struct {
	cursor_from, cursor_to int

	cond_node  es_ast_node
	true_node  es_ast_node
	false_node es_ast_node
}

func (n es_ast_cond_expr_node) get_cursor_from() int { return n.cursor_from }
func (n es_ast_cond_expr_node) get_cursor_to() int   { return n.cursor_to }
func (n es_ast_cond_expr_node) make_code() []es_vm_instr {
	instrs := []es_vm_instr{}
	instrs = append(instrs, n.false_node.make_code()...)
	instrs = append(instrs, n.true_node.make_code()...)
	instrs = append(instrs, n.cond_node.make_code()...)
	instrs = append(instrs, es_vm_instr{es_vm_opcode_cond, nil})
	return instrs

}

type es_ast_comma_op_node struct {
	cursor_from, cursor_to int

	lhs_node es_ast_node
	rhs_node es_ast_node
}

func (n es_ast_comma_op_node) get_cursor_from() int { return n.cursor_from }
func (n es_ast_comma_op_node) get_cursor_to() int   { return n.cursor_to }
func (n es_ast_comma_op_node) make_code() []es_vm_instr {
	instrs := []es_vm_instr{}

	instrs = append(instrs, n.lhs_node.make_code()...)
	instrs = append(instrs, es_vm_instr{es_vm_opcode_get_value, nil})
	instrs = append(instrs, es_vm_instr{es_vm_opcode_discard, nil})
	instrs = append(instrs, n.rhs_node.make_code()...)
	instrs = append(instrs, es_vm_instr{es_vm_opcode_get_value, nil})

	return instrs
}

type es_ast_expr_statement_node struct {
	cursor_from, cursor_to int

	node es_ast_node // May be nil
}

func (n es_ast_expr_statement_node) get_cursor_from() int { return n.cursor_from }
func (n es_ast_expr_statement_node) get_cursor_to() int   { return n.cursor_to }
func (n es_ast_expr_statement_node) make_code() []es_vm_instr {
	instrs := []es_vm_instr{}

	instrs = append(instrs, n.node.make_code()...)
	instrs = append(instrs, es_vm_instr{es_vm_opcode_get_value, nil})
	instrs = append(instrs, es_vm_instr{es_vm_opcode_mov_to_res, nil})

	return instrs
}

type es_ast_return_statement_node struct {
	cursor_from, cursor_to int

	node es_ast_node // May be nil
}

func (n es_ast_return_statement_node) get_cursor_from() int     { return n.cursor_from }
func (n es_ast_return_statement_node) get_cursor_to() int       { return n.cursor_to }
func (n es_ast_return_statement_node) make_code() []es_vm_instr { panic("TODO") }

type es_ast_formal_parameters struct {
	params     []es_ast_node
	rest_param es_ast_node // May be nil
}

type es_ast_function_decl_node struct {
	cursor_from, cursor_to int

	ident  es_ast_binding_identifier_node
	params es_ast_formal_parameters
	body   []es_ast_node
}

func (n es_ast_function_decl_node) get_cursor_from() int     { return n.cursor_from }
func (n es_ast_function_decl_node) get_cursor_to() int       { return n.cursor_to }
func (n es_ast_function_decl_node) make_code() []es_vm_instr { panic("TODO") }
