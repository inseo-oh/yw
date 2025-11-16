package libcommon

import (
	"fmt"
	"reflect"
	"strings"
)

func ConsumeLongestString(strs []string) string {
	longest := ""
	for _, str := range strs {
		if len(longest) < len(str) {
			longest = str
		}
	}
	return longest
}

func IsNil(t any) bool {
	// IsNil() will panic if the value is not supported by it (e.g. Struct).
	// So we recover() from the panic if that happens.
	defer func() { recover() }()

	return t == nil || reflect.ValueOf(t).IsNil()
}

func MakeStrPtr(s string) *string {
	return &s
}

func EncodeSingleQuoteString(s string) string {
	quote_replaced := strings.ReplaceAll(s, "'", "\\'")
	slash_replaced := strings.ReplaceAll(quote_replaced, "\\", "\\\\")
	newline_replaced := strings.ReplaceAll(slash_replaced, "\n", "\\\n'")
	return fmt.Sprintf("'%s'", newline_replaced)
}
