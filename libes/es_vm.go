package libes

import (
	"log"
	"math"
)

type es_vm struct {
	stack       []any
	last_result es_value
}

func MakeVm() es_vm {
	return es_vm{}
}

func (vm es_vm) peek_stack_top() any {
	v := vm.stack[len(vm.stack)-1]
	return v
}
func (vm es_vm) peek_stack_top_value() es_value {
	return es_get_value(vm.peek_stack_top())
}
func (vm *es_vm) stack_pop() any {
	v := vm.peek_stack_top()
	vm.stack = vm.stack[:len(vm.stack)-1]
	return v
}
func (vm *es_vm) stack_push(v any) {
	vm.stack = append(vm.stack, v)
}

// Returns result of the last expression statement.
func (vm *es_vm) Exec(instrs []es_vm_instr) es_value {
	vm.last_result = es_make_undefined_value()
	for _, instr := range instrs {
		value_must_be_present := func() {
			if instr.value == nil {
				panic(".value must not be nil")
			}
		}
		switch instr.op {
		case es_vm_opcode_push:
			value_must_be_present()
			vm.stack_push(instr.value)
		case es_vm_opcode_get_value:
			v := es_get_value(vm.stack_pop())
			vm.stack_push(v)
		case es_vm_opcode_mov_to_res:
			v := es_get_value(vm.stack_pop())
			vm.last_result = v
		case es_vm_opcode_discard:
			vm.stack_pop()
		case es_vm_opcode_cond:
			cond := es_get_value(vm.stack_pop()).expect_boolean()
			true_v := es_get_value(vm.stack_pop())
			false_v := es_get_value(vm.stack_pop())
			if cond {
				vm.stack_push(true_v)
			} else {
				vm.stack_push(false_v)
			}
		case es_vm_opcode_exponent:
			lhs := es_get_value(vm.stack_pop()).expect_number_f()
			rhs := es_get_value(vm.stack_pop()).expect_number_f()
			res := math.Pow(lhs, rhs)
			vm.stack_push(es_make_number_value_f(res))
		case es_vm_opcode_multiply:
			lhs := es_get_value(vm.stack_pop()).expect_number_f()
			rhs := es_get_value(vm.stack_pop()).expect_number_f()
			res := lhs * rhs
			vm.stack_push(es_make_number_value_f(res))
		case es_vm_opcode_divide:
			lhs := es_get_value(vm.stack_pop()).expect_number_f()
			rhs := es_get_value(vm.stack_pop()).expect_number_f()
			res := lhs / rhs
			vm.stack_push(es_make_number_value_f(res))
		case es_vm_opcode_modulo:
			lhs := es_get_value(vm.stack_pop()).expect_number_f()
			rhs := es_get_value(vm.stack_pop()).expect_number_f()
			res := math.Mod(lhs, rhs)
			vm.stack_push(es_make_number_value_f(res))
		case es_vm_opcode_add:
			lhs := es_get_value(vm.stack_pop()).expect_number_f()
			rhs := es_get_value(vm.stack_pop()).expect_number_f()
			res := lhs + rhs
			vm.stack_push(es_make_number_value_f(res))
		case es_vm_opcode_subtract:
			lhs := es_get_value(vm.stack_pop()).expect_number_f()
			rhs := es_get_value(vm.stack_pop()).expect_number_f()
			res := lhs - rhs
			vm.stack_push(es_make_number_value_f(res))
		case es_vm_opcode_left_shift:
			lhs := int32(es_get_value(vm.stack_pop()).expect_number_i())
			rhs := int32(es_get_value(vm.stack_pop()).expect_number_i())
			res := lhs << rhs
			vm.stack_push(es_make_number_value_i(int64(res)))
		case es_vm_opcode_right_ashift:
			lhs := int32(es_get_value(vm.stack_pop()).expect_number_i())
			rhs := int32(es_get_value(vm.stack_pop()).expect_number_i())
			res := lhs >> rhs
			vm.stack_push(es_make_number_value_i(int64(res)))
		case es_vm_opcode_right_lshift:
			lhs := uint32(es_get_value(vm.stack_pop()).expect_number_i())
			rhs := int32(es_get_value(vm.stack_pop()).expect_number_i())
			res := lhs >> rhs
			vm.stack_push(es_make_number_value_i(int64(int32(res))))
		case es_vm_opcode_band:
			lhs := int32(es_get_value(vm.stack_pop()).expect_number_i())
			rhs := int32(es_get_value(vm.stack_pop()).expect_number_i())
			res := lhs & rhs
			vm.stack_push(es_make_number_value_i(int64(res)))
		case es_vm_opcode_bxor:
			lhs := int32(es_get_value(vm.stack_pop()).expect_number_i())
			rhs := int32(es_get_value(vm.stack_pop()).expect_number_i())
			res := lhs ^ rhs
			vm.stack_push(es_make_number_value_i(int64(res)))
		case es_vm_opcode_bor:
			lhs := int32(es_get_value(vm.stack_pop()).expect_number_i())
			rhs := int32(es_get_value(vm.stack_pop()).expect_number_i())
			res := lhs | rhs
			vm.stack_push(es_make_number_value_i(int64(res)))
		case es_vm_opcode_land:
			lhs := es_get_value(vm.stack_pop()).expect_boolean()
			rhs := es_get_value(vm.stack_pop()).expect_boolean()
			res := lhs && rhs
			vm.stack_push(es_make_boolean_value(res))
		case es_vm_opcode_lor:
			lhs := es_get_value(vm.stack_pop()).expect_boolean()
			rhs := es_get_value(vm.stack_pop()).expect_boolean()
			res := lhs || rhs
			vm.stack_push(es_make_boolean_value(res))
		case es_vm_opcode_coalesce:
			lhs := es_get_value(vm.stack_pop())
			rhs := es_get_value(vm.stack_pop())
			if lhs.tp == es_value_type_null || lhs.tp == es_value_type_undefined {
				vm.stack_push(rhs)
			} else {
				vm.stack_push(lhs)
			}
		case es_vm_opcode_less_than:
			lhs := es_get_value(vm.stack_pop()).expect_number_f()
			rhs := es_get_value(vm.stack_pop()).expect_number_f()
			res := lhs < rhs
			vm.stack_push(es_make_boolean_value(res))
		case es_vm_opcode_less_than_or_equal:
			lhs := es_get_value(vm.stack_pop()).expect_number_f()
			rhs := es_get_value(vm.stack_pop()).expect_number_f()
			res := lhs <= rhs
			vm.stack_push(es_make_boolean_value(res))
		case es_vm_opcode_greater_than:
			lhs := es_get_value(vm.stack_pop()).expect_number_f()
			rhs := es_get_value(vm.stack_pop()).expect_number_f()
			res := lhs > rhs
			vm.stack_push(es_make_boolean_value(res))
		case es_vm_opcode_greater_than_or_equal:
			lhs := es_get_value(vm.stack_pop()).expect_number_f()
			rhs := es_get_value(vm.stack_pop()).expect_number_f()
			res := lhs >= rhs
			vm.stack_push(es_make_boolean_value(res))
		case es_vm_opcode_equal:
			// TODO: Compare between different types
			lhs_v := es_get_value(vm.stack_pop())
			rhs_v := es_get_value(vm.stack_pop())
			var res bool
			if lhs_v.tp != rhs_v.tp {
				panic("TODO")
			} else if lhs_v.tp == es_value_type_number {
				res = lhs_v.expect_number_f() == rhs_v.expect_number_f()
			} else {
				res = lhs_v.value == rhs_v.value
			}
			vm.stack_push(es_make_boolean_value(res))
		case es_vm_opcode_strict_equal:
			panic("TODO")
		case es_vm_opcode_not_equal:
			// TODO: Compare between different types
			lhs_v := es_get_value(vm.stack_pop())
			rhs_v := es_get_value(vm.stack_pop())
			var res bool
			if lhs_v.tp != rhs_v.tp {
				panic("TODO")
			} else if lhs_v.tp == es_value_type_number {
				res = lhs_v.expect_number_f() != rhs_v.expect_number_f()
			} else {
				res = lhs_v.value != rhs_v.value
			}
			vm.stack_push(es_make_boolean_value(res))
		case es_vm_opcode_strict_not_equal:
			panic("TODO")
		case es_vm_opcode_plus:
			v := es_get_value(vm.stack_pop()).expect_number_f()
			vm.stack_push(es_make_number_value_f(v))
		case es_vm_opcode_neg:
			v := es_get_value(vm.stack_pop()).expect_number_f()
			vm.stack_push(es_make_number_value_f(-v))
		case es_vm_opcode_bnot:
			v := uint32(int32(es_get_value(vm.stack_pop()).expect_number_i()))
			vm.stack_push(es_make_number_value_i(int64(int32(^v))))
		case es_vm_opcode_lnot:
			v := es_get_value(vm.stack_pop()).expect_boolean()
			vm.stack_push(es_make_boolean_value(!v))
		case es_vm_opcode_await:
			panic("TODO")
		default:
			log.Panicf("unexpected opcode %d", instr.op)
		}
	}
	return vm.last_result
}

type es_vm_instr struct {
	op    es_vm_opcode
	value any // Most of the time this is ignored, but some opcodes use it.
}
type es_vm_opcode uint8

const (
	// General, Misc
	es_vm_opcode_push       = es_vm_opcode(0x01) // Pushes given value. [.value = Value to push]
	es_vm_opcode_get_value  = es_vm_opcode(0x02) // Pops a value, resolves binding if necessary, and pushes es_value
	es_vm_opcode_mov_to_res = es_vm_opcode(0x03) // Pops a value, saves result to internal "last result" register
	es_vm_opcode_discard    = es_vm_opcode(0x04) // Pops a value, and forgets about it
	es_vm_opcode_cond       = es_vm_opcode(0x05) // C ? T : F support. Pops <C>, <T>, <F>, and pushes <T> if <C> is true, <F> otherwise.

	// Binary operators - These pop LHS, RHS, and pushes calculation result.
	es_vm_opcode_exponent              = es_vm_opcode(0x10)
	es_vm_opcode_multiply              = es_vm_opcode(0x11)
	es_vm_opcode_divide                = es_vm_opcode(0x12)
	es_vm_opcode_modulo                = es_vm_opcode(0x13)
	es_vm_opcode_add                   = es_vm_opcode(0x14)
	es_vm_opcode_subtract              = es_vm_opcode(0x15)
	es_vm_opcode_left_shift            = es_vm_opcode(0x16)
	es_vm_opcode_right_ashift          = es_vm_opcode(0x17) // Arithmetic shift (>>)
	es_vm_opcode_right_lshift          = es_vm_opcode(0x18) // Logical shift (>>>)
	es_vm_opcode_band                  = es_vm_opcode(0x19)
	es_vm_opcode_bxor                  = es_vm_opcode(0x1a)
	es_vm_opcode_bor                   = es_vm_opcode(0x1b)
	es_vm_opcode_land                  = es_vm_opcode(0x1c)
	es_vm_opcode_lor                   = es_vm_opcode(0x1d)
	es_vm_opcode_coalesce              = es_vm_opcode(0x1e)
	es_vm_opcode_less_than             = es_vm_opcode(0x1f)
	es_vm_opcode_less_than_or_equal    = es_vm_opcode(0x20)
	es_vm_opcode_greater_than          = es_vm_opcode(0x21)
	es_vm_opcode_greater_than_or_equal = es_vm_opcode(0x22)
	es_vm_opcode_equal                 = es_vm_opcode(0x23)
	es_vm_opcode_strict_equal          = es_vm_opcode(0x24)
	es_vm_opcode_not_equal             = es_vm_opcode(0x25)
	es_vm_opcode_strict_not_equal      = es_vm_opcode(0x26)

	// Unary operators -  These pop a value, and pushes calculation result.
	es_vm_opcode_plus  = es_vm_opcode(0x30)
	es_vm_opcode_neg   = es_vm_opcode(0x31)
	es_vm_opcode_bnot  = es_vm_opcode(0x32)
	es_vm_opcode_lnot  = es_vm_opcode(0x33)
	es_vm_opcode_await = es_vm_opcode(0x34)
)
