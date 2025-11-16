// Implementation of the UTF-8 encoding (https://encoding.spec.whatwg.org/#utf-8)
package libencoding

// https://encoding.spec.whatwg.org/#utf-8-decoder
type utf8_decoder struct {
	code_point     rune
	bytes_seen     int
	bytes_needed   int
	lower_boundary byte
	upper_boundary byte
}

var utf8_encoding = encoding{
	make_decoder: func() decoder {
		return &utf8_decoder{
			lower_boundary: 0x80,
			upper_boundary: 0xbf,
		}
	},
}

func (dec *utf8_decoder) handler(queue *IoQueue, byte_item IoQueueItem) handler_result {
	if byte_item.IsEndOfQueue() {
		if dec.bytes_needed != 0 {
			dec.bytes_needed = 0
			return handler_result_error
		} else {
			return handler_result_finished
		}
	}
	byt := byte_item.V.(uint8)
	if dec.bytes_needed == 0 {
		if byt <= 0x7f {
			return handler_result(byt)
		} else if 0xc2 <= byt && byt <= 0xdf {
			dec.bytes_needed = 1
			dec.code_point = rune(byt & 0x1f)
		} else if 0xe0 <= byt && byt <= 0xef {
			switch byt {
			case 0xe0:
				dec.lower_boundary = 0xa0
			case 0xed:
				dec.upper_boundary = 0x9f
			}
			dec.bytes_needed = 2
			dec.code_point = rune(byt & 0xf)
		} else if 0xf0 <= byt && byt <= 0xf4 {
			switch byt {
			case 0xf0:
				dec.lower_boundary = 0x90
			case 0xf4:
				dec.upper_boundary = 0x8f
			}
			dec.bytes_needed = 3
			dec.code_point = rune(byt & 0x7)
		} else {
			return handler_result_error
		}
		return handler_result_continue
	}
	if byt < dec.lower_boundary || dec.upper_boundary < byt {
		dec.code_point = 0
		dec.bytes_needed = 0
		dec.bytes_seen = 0
		dec.lower_boundary = 0x80
		dec.upper_boundary = 0xbf
		queue.RestoreOne(byte_item)
		return handler_result_error
	}
	dec.lower_boundary = 0x80
	dec.upper_boundary = 0xbf
	dec.code_point = (dec.code_point << 6) | rune(byt&0x3f)
	dec.bytes_seen++
	if dec.bytes_seen != dec.bytes_needed {
		return handler_result_continue
	}
	cp := dec.code_point
	dec.code_point = 0
	dec.bytes_needed = 0
	dec.bytes_seen = 0
	return handler_result(cp)
}
