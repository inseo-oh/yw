// This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
// SPDX-License-Identifier: BSD-3-Clause
// See LICENSE for details, and LICENSE.thirdparty.md for third-party license information.

package encoding

// https://encoding.spec.whatwg.org/#utf-8-decoder
type utf8Decoder struct {
	codepoint     rune
	bytesSeen     int
	bytesNeeded   int
	lowerBoundary byte
	upperBoundary byte
}

var utf8Encoding = encoding{
	makeDecoder: func() decoder {
		return &utf8Decoder{
			lowerBoundary: 0x80,
			upperBoundary: 0xbf,
		}
	},
}

func (dec *utf8Decoder) handler(queue *IoQueue, byteItem IoQueueItem) handlerResult {
	if byteItem.IsEndOfQueue() {
		if dec.bytesNeeded != 0 {
			dec.bytesNeeded = 0
			return handlerResultError
		} else {
			return handlerResultFinished
		}
	}
	byt := byteItem.V.(uint8)
	if dec.bytesNeeded == 0 {
		if byt <= 0x7f {
			return handlerResult(byt)
		} else if 0xc2 <= byt && byt <= 0xdf {
			dec.bytesNeeded = 1
			dec.codepoint = rune(byt & 0x1f)
		} else if 0xe0 <= byt && byt <= 0xef {
			switch byt {
			case 0xe0:
				dec.lowerBoundary = 0xa0
			case 0xed:
				dec.upperBoundary = 0x9f
			}
			dec.bytesNeeded = 2
			dec.codepoint = rune(byt & 0xf)
		} else if 0xf0 <= byt && byt <= 0xf4 {
			switch byt {
			case 0xf0:
				dec.lowerBoundary = 0x90
			case 0xf4:
				dec.upperBoundary = 0x8f
			}
			dec.bytesNeeded = 3
			dec.codepoint = rune(byt & 0x7)
		} else {
			return handlerResultError
		}
		return handlerResultContinue
	}
	if byt < dec.lowerBoundary || dec.upperBoundary < byt {
		dec.codepoint = 0
		dec.bytesNeeded = 0
		dec.bytesSeen = 0
		dec.lowerBoundary = 0x80
		dec.upperBoundary = 0xbf
		queue.RestoreOne(byteItem)
		return handlerResultError
	}
	dec.lowerBoundary = 0x80
	dec.upperBoundary = 0xbf
	dec.codepoint = (dec.codepoint << 6) | rune(byt&0x3f)
	dec.bytesSeen++
	if dec.bytesSeen != dec.bytesNeeded {
		return handlerResultContinue
	}
	cp := dec.codepoint
	dec.codepoint = 0
	dec.bytesNeeded = 0
	dec.bytesSeen = 0
	return handlerResult(cp)
}
