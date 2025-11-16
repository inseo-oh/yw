package libes

import (
	"testing"
)

func Test_Es(t *testing.T) {
	test_cases := []struct {
		input    string
		want_val es_value
	}{
		// Literal only --------------------------------------------------------
		{"123;", es_make_number_value_f(123)},
		{"0x1234DEAD;", es_make_number_value_i(0x1234dead)},
		{"0x1234dead;", es_make_number_value_i(0x1234dead)},
		{"0b10101010;", es_make_number_value_i(0xaa)},
		{"0o123;", es_make_number_value_i(0o123)},
		{"0123;", es_make_number_value_i(0o123)},
		{"0789;", es_make_number_value_f(789)},
		{"true;", es_make_boolean_value(true)},
		{"false;", es_make_boolean_value(false)},
		{"null;", es_make_null_value()},
		// TODO: String literal

		// Binary operators ----------------------------------------------------
		{"40 ** 2;", es_make_number_value_f(1600)},
		{"40 * 20;", es_make_number_value_f(800)},
		{"10 / 2;", es_make_number_value_f(5)},
		{"10 % 3;", es_make_number_value_f(1)},
		{"34 % 10;", es_make_number_value_f(4)},
		{"34 + 35;", es_make_number_value_f(69)},
		{"100 - 31;", es_make_number_value_f(69)},
		{"100 << 2;", es_make_number_value_i(400)},
		{"100 >> 2;", es_make_number_value_i(25)},
		{"0xffffff00 >> 2;", es_make_number_value_i(-64)},
		{"0xffffff00 >>> 2;", es_make_number_value_i(0x3fffffc0)},
		{"0xab & 0xcd;", es_make_number_value_i(0x89)},
		{"0xab ^ 0xcd;", es_make_number_value_i(0x66)},
		{"0xab | 0xcd;", es_make_number_value_i(0xef)},
		{"true && true;", es_make_boolean_value(true)},
		{"true && false;", es_make_boolean_value(false)},
		{"false && true;", es_make_boolean_value(false)},
		{"false && false;", es_make_boolean_value(false)},
		{"true || true;", es_make_boolean_value(true)},
		{"true || false;", es_make_boolean_value(true)},
		{"false || true;", es_make_boolean_value(true)},
		{"false || false;", es_make_boolean_value(false)},
		{"1 ?? 2;", es_make_number_value_f(1)},
		{"null ?? 2;", es_make_number_value_f(2)},
		{"40 < 30;", es_make_boolean_value(false)},
		{"30 < 40;", es_make_boolean_value(true)},
		{"30 <= 30;", es_make_boolean_value(true)},
		{"40 > 30;", es_make_boolean_value(true)},
		{"30 > 40;", es_make_boolean_value(false)},
		{"30 >= 30;", es_make_boolean_value(true)},
		{"30 == 30;", es_make_boolean_value(true)},
		{"30 == 40;", es_make_boolean_value(false)},
		// TODO: strict equal
		{"30 != 30;", es_make_boolean_value(false)},
		{"30 != 40;", es_make_boolean_value(true)},
		// TODO: strict not equal

		// Unary operators -----------------------------------------------------
		{"+123;", es_make_number_value_f(123)},
		{"-123;", es_make_number_value_f(-123)},
		{"~0;", es_make_number_value_i(-1)},
		{"~0xffffffff;", es_make_number_value_i(0)},
		{"!false;", es_make_boolean_value(true)},
		{"!true;", es_make_boolean_value(false)},

		// Conditional operator ------------------------------------------------
		{"true ? 34 : 35;", es_make_number_value_f(34)},
		{"false ? 34 : 35;", es_make_number_value_f(35)},

		// Comma operator ------------------------------------------------------
		{"34, 35, 69;", es_make_number_value_f(69)},

		// Operator precedence -------------------------------------------------
		{"2 * 3 ** 2;", es_make_number_value_f(18)},
		{"27 / 3 ** 2;", es_make_number_value_f(3)},
		{"26 % 3 ** 2;", es_make_number_value_f(8)},

		{"100 + 9 * 7;", es_make_number_value_f(163)},
		{"100 - 9 * 7;", es_make_number_value_f(37)},
		{"100 + 27 / 3;", es_make_number_value_f(109)},
		{"100 - 27 / 3;", es_make_number_value_f(91)},
		{"100 + 26 % 3;", es_make_number_value_f(102)},
		{"100 - 26 % 3;", es_make_number_value_f(98)},

		{"16 << 2 + 1;", es_make_number_value_i(128)},
		{"16 << 2 - 1;", es_make_number_value_i(32)},
		{"16 >>> 2 + 1;", es_make_number_value_i(2)},
		{"16 >>> 2 - 1;", es_make_number_value_i(8)},
		{"16 >> 2 + 1;", es_make_number_value_i(2)},
		{"16 >> 2 - 1;", es_make_number_value_i(8)},

		{"16 < 16 << 1;", es_make_boolean_value(true)},
		{"16 < 16 >> 1;", es_make_boolean_value(false)},
		{"16 < 16 >>> 1;", es_make_boolean_value(false)},
		{"16 <= 16 << 1;", es_make_boolean_value(true)},
		{"16 <= 16 >> 1;", es_make_boolean_value(false)},
		{"16 <= 16 >>> 1;", es_make_boolean_value(false)},
		{"16 > 16 << 1;", es_make_boolean_value(false)},
		{"16 > 16 >> 1;", es_make_boolean_value(true)},
		{"16 > 16 >>> 1;", es_make_boolean_value(true)},
		{"16 >= 16 << 1;", es_make_boolean_value(false)},
		{"16 >= 16 >> 1;", es_make_boolean_value(true)},
		{"16 >= 16 >>> 1;", es_make_boolean_value(true)},

		{"false == 16 < 16;", es_make_boolean_value(true)},
		{"true == 16 <= 16;", es_make_boolean_value(true)},
		{"false == 16 > 16;", es_make_boolean_value(true)},
		{"true == 16 >= 16;", es_make_boolean_value(true)},
		{"true != 16 < 16;", es_make_boolean_value(true)},
		{"false != 16 <= 16;", es_make_boolean_value(true)},
		{"true != 16 > 16;", es_make_boolean_value(true)},
		{"false != 16 >= 16;", es_make_boolean_value(true)},

		{"0xfc & 3 << 1;", es_make_number_value_i(0x04)},
		{"0xff ^ 0xfc & 6;", es_make_number_value_i(0xfb)},
		{"1 | 0xff ^ 0xfc;", es_make_number_value_i(0x03)},

		{"true && true == false;", es_make_boolean_value(false)},
		{"true && true == true;", es_make_boolean_value(true)},
		{"true && true != false;", es_make_boolean_value(true)},
		{"true && true != true;", es_make_boolean_value(false)},

		{"false || true == false;", es_make_boolean_value(false)},
		{"false || true == true;", es_make_boolean_value(true)},
		{"false || true != false;", es_make_boolean_value(true)},
		{"false || true != true;", es_make_boolean_value(false)},
		{"true || true == false;", es_make_boolean_value(true)},
		{"true || true == true;", es_make_boolean_value(true)},
		{"true || true != false;", es_make_boolean_value(true)},
		{"true || true != true;", es_make_boolean_value(true)},

		// Parentheses ---------------------------------------------------------
		{"(69);", es_make_number_value_f(69)},
		{"(6 + 9) * 7;", es_make_number_value_f(105)},
		// TODO
	}
	for _, cs := range test_cases {
		t.Run(cs.input, func(t *testing.T) {
			vm := MakeVm()
			code, err := Compile(cs.input)
			if err != nil {
				t.Errorf("[%s] failed to compile: %v", cs.input, err)
				return
			}
			res_val := vm.Exec(code)
			if (res_val.tp != cs.want_val.tp) || (res_val.value != cs.want_val.value) {
				t.Errorf("[%s] got result %v, want %v", cs.input, res_val, cs.want_val)
			}
		})

	}

}
