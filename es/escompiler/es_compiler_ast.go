// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license information.

package escompiler

import (
	"github.com/inseo-oh/yw/es"
	"github.com/inseo-oh/yw/es/vm"
)

func makeCodeForAstNodes(nodes []astNode) []vm.Instr {
	instrs := []vm.Instr{}
	for _, node := range nodes {
		instrs = append(instrs, node.makeCode()...)
	}
	return instrs
}

type astNode interface {
	nodeCursorFrom() int
	nodeCursorTo() int
	makeCode() []vm.Instr
}

type astIdentifierReferenceNode struct {
	cursorFrom, cursorTo int

	name string
}

func (n astIdentifierReferenceNode) nodeCursorFrom() int  { return n.cursorFrom }
func (n astIdentifierReferenceNode) nodeCursorTo() int    { return n.cursorTo }
func (n astIdentifierReferenceNode) makeCode() []vm.Instr { panic("TODO") }

type astBindingIdentifierNode struct {
	cursorFrom, cursorTo int

	name string
}

func (n astBindingIdentifierNode) nodeCursorFrom() int  { return n.cursorFrom }
func (n astBindingIdentifierNode) nodeCursorTo() int    { return n.cursorTo }
func (n astBindingIdentifierNode) makeCode() []vm.Instr { panic("TODO") }

type astParenExprNode struct {
	cursorFrom, cursorTo int

	node astNode
}

func (n astParenExprNode) nodeCursorFrom() int { return n.cursorFrom }
func (n astParenExprNode) nodeCursorTo() int   { return n.cursorTo }
func (n astParenExprNode) makeCode() []vm.Instr {
	return n.node.makeCode()
}

type astLiteralNode struct {
	cursorFrom, cursorTo int

	value es.Value
}

func (n astLiteralNode) nodeCursorFrom() int { return n.cursorFrom }
func (n astLiteralNode) nodeCursorTo() int   { return n.cursorTo }
func (n astLiteralNode) makeCode() []vm.Instr {
	return []vm.Instr{{Op: vm.OpcodePush, Value: n.value}}
}

type astCallExprNode struct {
	cursorFrom, cursorTo int

	callee astNode
	args   astArguments
}
type astArguments struct {
	cursorFrom, cursorTo int

	args     []astNode
	restArgs astNode // May be nil
}

func (n astCallExprNode) nodeCursorFrom() int  { return n.cursorFrom }
func (n astCallExprNode) nodeCursorTo() int    { return n.cursorTo }
func (n astCallExprNode) makeCode() []vm.Instr { panic("TODO") }

type astUnaryOpNode struct {
	cursorFrom, cursorTo int

	tp   astUnaryOpType
	node astNode
}
type astUnaryOpType uint8

const (
	astOpTypePlus astUnaryOpType = iota
	astOpTypeNeg
	astOpTypeBNot
	astOpTypeLNot
	astOpTypeAwait
)

func (n astUnaryOpNode) nodeCursorFrom() int { return n.cursorFrom }
func (n astUnaryOpNode) nodeCursorTo() int   { return n.cursorTo }

var astUnaryOpToVmMap = map[astUnaryOpType]vm.Opcode{
	astOpTypePlus:  vm.OpcodePlus,
	astOpTypeNeg:   vm.OpcodeNeg,
	astOpTypeBNot:  vm.OpcodeBNot,
	astOpTypeLNot:  vm.OpcodeLNot,
	astOpTypeAwait: vm.OpcodeAwait,
}

func (n astUnaryOpNode) makeCode() []vm.Instr {
	instrs := []vm.Instr{}
	instrs = append(instrs, n.node.makeCode()...)
	opcode := astUnaryOpToVmMap[n.tp]
	instrs = append(instrs, vm.Instr{opcode, nil})
	return instrs
}

type astBinaryOpNode struct {
	cursorFrom, cursorTo int

	tp      astBinaryOpType
	lhsNode astNode
	rhsNode astNode
}
type astBinaryOpType uint8

const (
	astOpTypeExponent astBinaryOpType = iota
	astOpTypeMultiply
	astOpTypeDivide
	astOpTypeModulo
	astOpTypeAdd
	astOpTypeSubtract
	astOpTypeLeftShift
	astOpTypeRightLShift
	astOpTypeRightAShift
	astOpTypeLessThan
	astOpTypeGreaterThan
	astOpTypeLessThanOrEqual
	astOpTypeGreaterThanOrEqual
	astOpTypeInstanceof
	astOpTypeEqual
	astOpTypeNotEqual
	astOpTypeStrictEqual
	astOpTypeStrictNotEqual
	astOpTypeBAnd
	astOpTypeBXor
	astOPTypeBOr
	astOpTypeLAnd
	astOpTypeLOr
	astOpTypeCoalesce
)

func (n astBinaryOpNode) nodeCursorFrom() int { return n.cursorFrom }
func (n astBinaryOpNode) nodeCursorTo() int   { return n.cursorTo }

var astBinaryOpToVmMap = map[astBinaryOpType]vm.Opcode{
	astOpTypeExponent:           vm.OpcodeExponent,
	astOpTypeMultiply:           vm.OpcodeMultiply,
	astOpTypeDivide:             vm.OpcodeDivide,
	astOpTypeModulo:             vm.OpcodeModulo,
	astOpTypeAdd:                vm.OpcodeAdd,
	astOpTypeSubtract:           vm.OpcodeSubtract,
	astOpTypeLeftShift:          vm.OpcodeLeftShift,
	astOpTypeRightAShift:        vm.OpcodeRightAShift,
	astOpTypeRightLShift:        vm.OpcodeRightLShift,
	astOpTypeBAnd:               vm.OpcodeBAnd,
	astOpTypeBXor:               vm.OpcodeBXor,
	astOPTypeBOr:                vm.OpcodeBOr,
	astOpTypeLAnd:               vm.OpcodeLAnd,
	astOpTypeLOr:                vm.OpcodeLOr,
	astOpTypeLessThan:           vm.OpcodeLessThan,
	astOpTypeLessThanOrEqual:    vm.OpcodeLessThanOrEqual,
	astOpTypeGreaterThan:        vm.OpcodeGreaterThan,
	astOpTypeGreaterThanOrEqual: vm.OpcodeGreaterThanOrEqual,
	astOpTypeEqual:              vm.OpcodeEqual,
	astOpTypeStrictEqual:        vm.OpcodeStrictEqual,
	astOpTypeNotEqual:           vm.OpcodeNotEqual,
	astOpTypeStrictNotEqual:     vm.OpcodeStrictNotEqual,
	astOpTypeCoalesce:           vm.OpcodeCoalesce,
}

func (n astBinaryOpNode) makeCode() []vm.Instr {
	instrs := []vm.Instr{}
	instrs = append(instrs, n.rhsNode.makeCode()...)
	instrs = append(instrs, n.lhsNode.makeCode()...)
	opcode := astBinaryOpToVmMap[n.tp]
	instrs = append(instrs, vm.Instr{opcode, nil})
	return instrs

}

type astCondExprNode struct {
	cursorFrom, cursorTo int

	condNode  astNode
	trueNode  astNode
	falseNode astNode
}

func (n astCondExprNode) nodeCursorFrom() int { return n.cursorFrom }
func (n astCondExprNode) nodeCursorTo() int   { return n.cursorTo }
func (n astCondExprNode) makeCode() []vm.Instr {
	instrs := []vm.Instr{}
	instrs = append(instrs, n.falseNode.makeCode()...)
	instrs = append(instrs, n.trueNode.makeCode()...)
	instrs = append(instrs, n.condNode.makeCode()...)
	instrs = append(instrs, vm.Instr{vm.OpcodeCond, nil})
	return instrs

}

type astCommaOpNode struct {
	cursorFrom, cursorTo int

	lhsNode astNode
	rhsNode astNode
}

func (n astCommaOpNode) nodeCursorFrom() int { return n.cursorFrom }
func (n astCommaOpNode) nodeCursorTo() int   { return n.cursorTo }
func (n astCommaOpNode) makeCode() []vm.Instr {
	instrs := []vm.Instr{}

	instrs = append(instrs, n.lhsNode.makeCode()...)
	instrs = append(instrs, vm.Instr{vm.OpcodeGetValue, nil})
	instrs = append(instrs, vm.Instr{vm.OpcodeDiscard, nil})
	instrs = append(instrs, n.rhsNode.makeCode()...)
	instrs = append(instrs, vm.Instr{vm.OpcodeGetValue, nil})

	return instrs
}

type astExprStatementNode struct {
	cursorFrom, cursorTo int

	node astNode // May be nil
}

func (n astExprStatementNode) nodeCursorFrom() int { return n.cursorFrom }
func (n astExprStatementNode) nodeCursorTo() int   { return n.cursorTo }
func (n astExprStatementNode) makeCode() []vm.Instr {
	instrs := []vm.Instr{}

	instrs = append(instrs, n.node.makeCode()...)
	instrs = append(instrs, vm.Instr{vm.OpcodeGetValue, nil})
	instrs = append(instrs, vm.Instr{vm.OpcodeMovToRes, nil})

	return instrs
}

type astReturnStatementNode struct {
	cursorFrom, cursorTo int

	node astNode // May be nil
}

func (n astReturnStatementNode) nodeCursorFrom() int  { return n.cursorFrom }
func (n astReturnStatementNode) nodeCursorTo() int    { return n.cursorTo }
func (n astReturnStatementNode) makeCode() []vm.Instr { panic("TODO") }

type astFormalParameters struct {
	params    []astNode
	restParam astNode // May be nil
}

type astFunctionDeclNode struct {
	cursorFrom, cursorTo int

	ident  astBindingIdentifierNode
	params astFormalParameters
	body   []astNode
}

func (n astFunctionDeclNode) nodeCursorFrom() int  { return n.cursorFrom }
func (n astFunctionDeclNode) nodeCursorTo() int    { return n.cursorTo }
func (n astFunctionDeclNode) makeCode() []vm.Instr { panic("TODO") }
