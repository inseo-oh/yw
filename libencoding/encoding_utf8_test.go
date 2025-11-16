package libencoding

import (
	"slices"
	"testing"
)

func TestUtf8Decoder(t *testing.T) {
	cases := []struct {
		desc     string
		input    []uint8
		expected []rune
	}{
		{"Simple ASCII", []uint8{0x30, 0x31, 0x32, 0x33, 0x7e}, []rune("0123~")},
		{"Two byte characters", []uint8{0xc2, 0xa0, 0xde, 0xb1}, []rune{0x00a0, 0x07b1}},
		{"Three byte characters", []uint8{0xe0, 0xa4, 0x80, 0xed, 0x9f, 0xbb, 0xef, 0xad, 0x8f}, []rune{0x0900, 0xd7fb, 0xfb4f}},
		{"Four byte characters", []uint8{0xf0, 0x90, 0x91, 0x90, 0xf0, 0x9f, 0x83, 0xb5, 0xf4, 0x81, 0x8a, 0x8f}, []rune{0x10450, 0x1f0f5, 0x10128f}},
	}
	for _, cs := range cases {
		t.Run(cs.desc, func(t *testing.T) {
			input := IoQueueFromSlice(cs.input)
			output := IoQueueFromSlice[rune](nil)
			Decode(&input, Utf8, &output)
			got_runes := IoQueueToSlice[rune](output)
			if !slices.Equal(cs.expected, got_runes) {
				t.Errorf("expected %v, got %v", cs.expected, got_runes)
			}
		})
	}
}
