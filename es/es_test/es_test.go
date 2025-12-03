// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package es

import (
	"testing"

	"github.com/inseo-oh/yw/es"
	"github.com/inseo-oh/yw/es/escompiler"
	"github.com/inseo-oh/yw/es/vm"
)

func TestEs(t *testing.T) {
	testCases := []struct {
		input   string
		wantVal es.Value
	}{
		// Literal only --------------------------------------------------------
		{"123;", es.NewNumberValueF(123)},
		{"0x1234DEAD;", es.NewNumberValueI(0x1234dead)},
		{"0x1234dead;", es.NewNumberValueI(0x1234dead)},
		{"0b10101010;", es.NewNumberValueI(0xaa)},
		{"0o123;", es.NewNumberValueI(0o123)},
		{"0123;", es.NewNumberValueI(0o123)},
		{"0789;", es.NewNumberValueF(789)},
		{"true;", es.NewBooleanValue(true)},
		{"false;", es.NewBooleanValue(false)},
		{"null;", es.NewNullValue()},
		// TODO: String literal

		// Binary operators ----------------------------------------------------
		{"40 ** 2;", es.NewNumberValueF(1600)},
		{"40 * 20;", es.NewNumberValueF(800)},
		{"10 / 2;", es.NewNumberValueF(5)},
		{"10 % 3;", es.NewNumberValueF(1)},
		{"34 % 10;", es.NewNumberValueF(4)},
		{"34 + 35;", es.NewNumberValueF(69)},
		{"100 - 31;", es.NewNumberValueF(69)},
		{"100 << 2;", es.NewNumberValueI(400)},
		{"100 >> 2;", es.NewNumberValueI(25)},
		{"0xffffff00 >> 2;", es.NewNumberValueI(-64)},
		{"0xffffff00 >>> 2;", es.NewNumberValueI(0x3fffffc0)},
		{"0xab & 0xcd;", es.NewNumberValueI(0x89)},
		{"0xab ^ 0xcd;", es.NewNumberValueI(0x66)},
		{"0xab | 0xcd;", es.NewNumberValueI(0xef)},
		{"true && true;", es.NewBooleanValue(true)},
		{"true && false;", es.NewBooleanValue(false)},
		{"false && true;", es.NewBooleanValue(false)},
		{"false && false;", es.NewBooleanValue(false)},
		{"true || true;", es.NewBooleanValue(true)},
		{"true || false;", es.NewBooleanValue(true)},
		{"false || true;", es.NewBooleanValue(true)},
		{"false || false;", es.NewBooleanValue(false)},
		{"1 ?? 2;", es.NewNumberValueF(1)},
		{"null ?? 2;", es.NewNumberValueF(2)},
		{"40 < 30;", es.NewBooleanValue(false)},
		{"30 < 40;", es.NewBooleanValue(true)},
		{"30 <= 30;", es.NewBooleanValue(true)},
		{"40 > 30;", es.NewBooleanValue(true)},
		{"30 > 40;", es.NewBooleanValue(false)},
		{"30 >= 30;", es.NewBooleanValue(true)},
		{"30 == 30;", es.NewBooleanValue(true)},
		{"30 == 40;", es.NewBooleanValue(false)},
		// TODO: strict equal
		{"30 != 30;", es.NewBooleanValue(false)},
		{"30 != 40;", es.NewBooleanValue(true)},
		// TODO: strict not equal

		// Unary operators -----------------------------------------------------
		{"+123;", es.NewNumberValueF(123)},
		{"-123;", es.NewNumberValueF(-123)},
		{"~0;", es.NewNumberValueI(-1)},
		{"~0xffffffff;", es.NewNumberValueI(0)},
		{"!false;", es.NewBooleanValue(true)},
		{"!true;", es.NewBooleanValue(false)},

		// Conditional operator ------------------------------------------------
		{"true ? 34 : 35;", es.NewNumberValueF(34)},
		{"false ? 34 : 35;", es.NewNumberValueF(35)},

		// Comma operator ------------------------------------------------------
		{"34, 35, 69;", es.NewNumberValueF(69)},

		// Operator precedence -------------------------------------------------
		{"2 * 3 ** 2;", es.NewNumberValueF(18)},
		{"27 / 3 ** 2;", es.NewNumberValueF(3)},
		{"26 % 3 ** 2;", es.NewNumberValueF(8)},

		{"100 + 9 * 7;", es.NewNumberValueF(163)},
		{"100 - 9 * 7;", es.NewNumberValueF(37)},
		{"100 + 27 / 3;", es.NewNumberValueF(109)},
		{"100 - 27 / 3;", es.NewNumberValueF(91)},
		{"100 + 26 % 3;", es.NewNumberValueF(102)},
		{"100 - 26 % 3;", es.NewNumberValueF(98)},

		{"16 << 2 + 1;", es.NewNumberValueI(128)},
		{"16 << 2 - 1;", es.NewNumberValueI(32)},
		{"16 >>> 2 + 1;", es.NewNumberValueI(2)},
		{"16 >>> 2 - 1;", es.NewNumberValueI(8)},
		{"16 >> 2 + 1;", es.NewNumberValueI(2)},
		{"16 >> 2 - 1;", es.NewNumberValueI(8)},

		{"16 < 16 << 1;", es.NewBooleanValue(true)},
		{"16 < 16 >> 1;", es.NewBooleanValue(false)},
		{"16 < 16 >>> 1;", es.NewBooleanValue(false)},
		{"16 <= 16 << 1;", es.NewBooleanValue(true)},
		{"16 <= 16 >> 1;", es.NewBooleanValue(false)},
		{"16 <= 16 >>> 1;", es.NewBooleanValue(false)},
		{"16 > 16 << 1;", es.NewBooleanValue(false)},
		{"16 > 16 >> 1;", es.NewBooleanValue(true)},
		{"16 > 16 >>> 1;", es.NewBooleanValue(true)},
		{"16 >= 16 << 1;", es.NewBooleanValue(false)},
		{"16 >= 16 >> 1;", es.NewBooleanValue(true)},
		{"16 >= 16 >>> 1;", es.NewBooleanValue(true)},

		{"false == 16 < 16;", es.NewBooleanValue(true)},
		{"true == 16 <= 16;", es.NewBooleanValue(true)},
		{"false == 16 > 16;", es.NewBooleanValue(true)},
		{"true == 16 >= 16;", es.NewBooleanValue(true)},
		{"true != 16 < 16;", es.NewBooleanValue(true)},
		{"false != 16 <= 16;", es.NewBooleanValue(true)},
		{"true != 16 > 16;", es.NewBooleanValue(true)},
		{"false != 16 >= 16;", es.NewBooleanValue(true)},

		{"0xfc & 3 << 1;", es.NewNumberValueI(0x04)},
		{"0xff ^ 0xfc & 6;", es.NewNumberValueI(0xfb)},
		{"1 | 0xff ^ 0xfc;", es.NewNumberValueI(0x03)},

		{"true && true == false;", es.NewBooleanValue(false)},
		{"true && true == true;", es.NewBooleanValue(true)},
		{"true && true != false;", es.NewBooleanValue(true)},
		{"true && true != true;", es.NewBooleanValue(false)},

		{"false || true == false;", es.NewBooleanValue(false)},
		{"false || true == true;", es.NewBooleanValue(true)},
		{"false || true != false;", es.NewBooleanValue(true)},
		{"false || true != true;", es.NewBooleanValue(false)},
		{"true || true == false;", es.NewBooleanValue(true)},
		{"true || true == true;", es.NewBooleanValue(true)},
		{"true || true != false;", es.NewBooleanValue(true)},
		{"true || true != true;", es.NewBooleanValue(true)},

		// Parentheses ---------------------------------------------------------
		{"(69);", es.NewNumberValueF(69)},
		{"(6 + 9) * 7;", es.NewNumberValueF(105)},
		// TODO
	}
	for _, cs := range testCases {
		t.Run(cs.input, func(t *testing.T) {
			vm := vm.Vm{}
			code, err := escompiler.Compile(cs.input)
			if err != nil {
				t.Errorf("[%s] failed to compile: %v", cs.input, err)
				return
			}
			resVal := vm.Exec(code)
			if (resVal.Type != cs.wantVal.Type) || (resVal.Value != cs.wantVal.Value) {
				t.Errorf("[%s] got result %v, want %v", cs.input, resVal, cs.wantVal)
			}
		})

	}

}
