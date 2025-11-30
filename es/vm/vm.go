package vm

import (
	"log"
	"math"

	"github.com/inseo-oh/yw/es"
)

type vm struct {
	stack      []any
	lastResult es.Value
}

func MakeVm() vm {
	return vm{}
}

func (vm vm) peekStackTop() any {
	v := vm.stack[len(vm.stack)-1]
	return v
}
func (vm vm) peekStackTopValue() es.Value {
	return es.GetValue(vm.peekStackTop())
}
func (vm *vm) stackPop() any {
	v := vm.peekStackTop()
	vm.stack = vm.stack[:len(vm.stack)-1]
	return v
}
func (vm *vm) stackPush(v any) {
	vm.stack = append(vm.stack, v)
}

// Returns result of the last expression statement.
func (vm *vm) Exec(instrs []Instr) es.Value {
	vm.lastResult = es.NewUndefinedValue()
	for _, instr := range instrs {
		valueMustBePresent := func() {
			if instr.Value == nil {
				panic(".Value must not be nil")
			}
		}
		switch instr.Op {
		case OpcodePush:
			valueMustBePresent()
			vm.stackPush(instr.Value)
		case OpcodeGetValue:
			v := es.GetValue(vm.stackPop())
			vm.stackPush(v)
		case OpcodeMovToRes:
			v := es.GetValue(vm.stackPop())
			vm.lastResult = v
		case OpcodeDiscard:
			vm.stackPop()
		case OpcodeCond:
			cond := es.GetValue(vm.stackPop()).ExpectBoolean()
			trueV := es.GetValue(vm.stackPop())
			falseV := es.GetValue(vm.stackPop())
			if cond {
				vm.stackPush(trueV)
			} else {
				vm.stackPush(falseV)
			}
		case OpcodeExponent:
			lhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			rhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			res := math.Pow(lhs, rhs)
			vm.stackPush(es.NewNumberValueF(res))
		case OpcodeMultiply:
			lhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			rhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			res := lhs * rhs
			vm.stackPush(es.NewNumberValueF(res))
		case OpcodeDivide:
			lhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			rhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			res := lhs / rhs
			vm.stackPush(es.NewNumberValueF(res))
		case OpcodeModulo:
			lhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			rhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			res := math.Mod(lhs, rhs)
			vm.stackPush(es.NewNumberValueF(res))
		case OpcodeAdd:
			lhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			rhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			res := lhs + rhs
			vm.stackPush(es.NewNumberValueF(res))
		case OpcodeSubtract:
			lhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			rhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			res := lhs - rhs
			vm.stackPush(es.NewNumberValueF(res))
		case OpcodeLeftShift:
			lhs := int32(es.GetValue(vm.stackPop()).ExpectNumberI())
			rhs := int32(es.GetValue(vm.stackPop()).ExpectNumberI())
			res := lhs << rhs
			vm.stackPush(es.NewNumberValueI(int64(res)))
		case OpcodeRightAShift:
			lhs := int32(es.GetValue(vm.stackPop()).ExpectNumberI())
			rhs := int32(es.GetValue(vm.stackPop()).ExpectNumberI())
			res := lhs >> rhs
			vm.stackPush(es.NewNumberValueI(int64(res)))
		case OpcodeRightLShift:
			lhs := uint32(es.GetValue(vm.stackPop()).ExpectNumberI())
			rhs := int32(es.GetValue(vm.stackPop()).ExpectNumberI())
			res := lhs >> rhs
			vm.stackPush(es.NewNumberValueI(int64(int32(res))))
		case OpcodeBAnd:
			lhs := int32(es.GetValue(vm.stackPop()).ExpectNumberI())
			rhs := int32(es.GetValue(vm.stackPop()).ExpectNumberI())
			res := lhs & rhs
			vm.stackPush(es.NewNumberValueI(int64(res)))
		case OpcodeBXor:
			lhs := int32(es.GetValue(vm.stackPop()).ExpectNumberI())
			rhs := int32(es.GetValue(vm.stackPop()).ExpectNumberI())
			res := lhs ^ rhs
			vm.stackPush(es.NewNumberValueI(int64(res)))
		case OpcodeBOr:
			lhs := int32(es.GetValue(vm.stackPop()).ExpectNumberI())
			rhs := int32(es.GetValue(vm.stackPop()).ExpectNumberI())
			res := lhs | rhs
			vm.stackPush(es.NewNumberValueI(int64(res)))
		case OpcodeLAnd:
			lhs := es.GetValue(vm.stackPop()).ExpectBoolean()
			rhs := es.GetValue(vm.stackPop()).ExpectBoolean()
			res := lhs && rhs
			vm.stackPush(es.NewBooleanValue(res))
		case OpcodeLOr:
			lhs := es.GetValue(vm.stackPop()).ExpectBoolean()
			rhs := es.GetValue(vm.stackPop()).ExpectBoolean()
			res := lhs || rhs
			vm.stackPush(es.NewBooleanValue(res))
		case OpcodeCoalesce:
			lhs := es.GetValue(vm.stackPop())
			rhs := es.GetValue(vm.stackPop())
			if lhs.Type == es.ValueTypeNull || lhs.Type == es.ValueTypeUndefined {
				vm.stackPush(rhs)
			} else {
				vm.stackPush(lhs)
			}
		case OpcodeLessThan:
			lhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			rhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			res := lhs < rhs
			vm.stackPush(es.NewBooleanValue(res))
		case OpcodeLessThanOrEqual:
			lhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			rhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			res := lhs <= rhs
			vm.stackPush(es.NewBooleanValue(res))
		case OpcodeGreaterThan:
			lhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			rhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			res := lhs > rhs
			vm.stackPush(es.NewBooleanValue(res))
		case OpcodeGreaterThanOrEqual:
			lhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			rhs := es.GetValue(vm.stackPop()).ExpectNumberF()
			res := lhs >= rhs
			vm.stackPush(es.NewBooleanValue(res))
		case OpcodeEqual:
			// TODO: Compare between different types
			lhsV := es.GetValue(vm.stackPop())
			rhsV := es.GetValue(vm.stackPop())
			var res bool
			if lhsV.Type != rhsV.Type {
				panic("TODO")
			} else if lhsV.Type == es.ValueTypeNumber {
				res = lhsV.ExpectNumberF() == rhsV.ExpectNumberF()
			} else {
				res = lhsV.Value == rhsV.Value
			}
			vm.stackPush(es.NewBooleanValue(res))
		case OpcodeStrictEqual:
			panic("TODO")
		case OpcodeNotEqual:
			// TODO: Compare between different types
			lhsV := es.GetValue(vm.stackPop())
			rhsV := es.GetValue(vm.stackPop())
			var res bool
			if lhsV.Type != rhsV.Type {
				panic("TODO")
			} else if lhsV.Type == es.ValueTypeNumber {
				res = lhsV.ExpectNumberF() != rhsV.ExpectNumberF()
			} else {
				res = lhsV.Value != rhsV.Value
			}
			vm.stackPush(es.NewBooleanValue(res))
		case OpcodeStrictNotEqual:
			panic("TODO")
		case OpcodePlus:
			v := es.GetValue(vm.stackPop()).ExpectNumberF()
			vm.stackPush(es.NewNumberValueF(v))
		case OpcodeNeg:
			v := es.GetValue(vm.stackPop()).ExpectNumberF()
			vm.stackPush(es.NewNumberValueF(-v))
		case OpcodeBNot:
			v := uint32(int32(es.GetValue(vm.stackPop()).ExpectNumberI()))
			vm.stackPush(es.NewNumberValueI(int64(int32(^v))))
		case OpcodeLNot:
			v := es.GetValue(vm.stackPop()).ExpectBoolean()
			vm.stackPush(es.NewBooleanValue(!v))
		case OpcodeAwait:
			panic("TODO")
		default:
			log.Panicf("unExpected opcode %d", instr.Op)
		}
	}
	return vm.lastResult
}

type Instr struct {
	Op    Opcode
	Value any // Most of the time this is ignored, but some opcodes use it.
}
type Opcode uint8

const (
	// General, Misc
	OpcodePush     = Opcode(0x01) // Pushes given value. [.Value = Value to push]
	OpcodeGetValue = Opcode(0x02) // Pops a value, resolves binding if necessary, and pushes es.Value
	OpcodeMovToRes = Opcode(0x03) // Pops a value, saves result to internal "last result" register
	OpcodeDiscard  = Opcode(0x04) // Pops a value, and forgets about it
	OpcodeCond     = Opcode(0x05) // C ? T : F support. Pops <C>, <T>, <F>, and pushes <T> if <C> is true, <F> otherwise.

	// Binary operators - These pop LHS, RHS, and pushes calculation result.
	OpcodeExponent           = Opcode(0x10)
	OpcodeMultiply           = Opcode(0x11)
	OpcodeDivide             = Opcode(0x12)
	OpcodeModulo             = Opcode(0x13)
	OpcodeAdd                = Opcode(0x14)
	OpcodeSubtract           = Opcode(0x15)
	OpcodeLeftShift          = Opcode(0x16)
	OpcodeRightAShift        = Opcode(0x17) // Arithmetic shift (>>)
	OpcodeRightLShift        = Opcode(0x18) // Logical shift (>>>)
	OpcodeBAnd               = Opcode(0x19)
	OpcodeBXor               = Opcode(0x1a)
	OpcodeBOr                = Opcode(0x1b)
	OpcodeLAnd               = Opcode(0x1c)
	OpcodeLOr                = Opcode(0x1d)
	OpcodeCoalesce           = Opcode(0x1e)
	OpcodeLessThan           = Opcode(0x1f)
	OpcodeLessThanOrEqual    = Opcode(0x20)
	OpcodeGreaterThan        = Opcode(0x21)
	OpcodeGreaterThanOrEqual = Opcode(0x22)
	OpcodeEqual              = Opcode(0x23)
	OpcodeStrictEqual        = Opcode(0x24)
	OpcodeNotEqual           = Opcode(0x25)
	OpcodeStrictNotEqual     = Opcode(0x26)

	// Unary operators -  These pop a value, and pushes calculation result.
	OpcodePlus  = Opcode(0x30)
	OpcodeNeg   = Opcode(0x31)
	OpcodeBNot  = Opcode(0x32)
	OpcodeLNot  = Opcode(0x33)
	OpcodeAwait = Opcode(0x34)
)
