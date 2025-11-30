// Implementation of the Encoding Standard (https://encoding.spec.whatwg.org)
package encoding

import (
	"errors"
	"slices"
	"strings"

	"github.com/inseo-oh/yw/util"
)

type Type uint8

const (
	Utf8         = Type(iota) // https://encoding.spec.whatwg.org/#utf-8
	Ibm866                    // https://encoding.spec.whatwg.org/#ibm866
	Iso8859_2                 // https://encoding.spec.whatwg.org/#iso-8859-2
	Iso8859_3                 // https://encoding.spec.whatwg.org/#iso-8859-3
	Iso8859_4                 // https://encoding.spec.whatwg.org/#iso-8859-4
	Iso8859_5                 // https://encoding.spec.whatwg.org/#iso-8859-5
	Iso8859_6                 // https://encoding.spec.whatwg.org/#iso-8859-6
	Iso8859_7                 // https://encoding.spec.whatwg.org/#iso-8859-7
	Iso8859_8                 // https://encoding.spec.whatwg.org/#iso-8859-8
	Iso8859_8I                // https://encoding.spec.whatwg.org/#iso-8859-8-i
	Iso8859_10                // https://encoding.spec.whatwg.org/#iso-8859-10
	Iso8859_13                // https://encoding.spec.whatwg.org/#iso-8859-13
	Iso8859_14                // https://encoding.spec.whatwg.org/#iso-8859-14
	Iso8859_15                // https://encoding.spec.whatwg.org/#iso-8859-15
	Iso8859_16                // https://encoding.spec.whatwg.org/#iso-8859-16
	Koi8R                     // https://encoding.spec.whatwg.org/#koi8-r
	Koi8U                     // https://encoding.spec.whatwg.org/#koi8-u
	Macintosh                 // https://encoding.spec.whatwg.org/#macintosh
	Windows874                // https://encoding.spec.whatwg.org/#windows-874
	Windows1250               // https://encoding.spec.whatwg.org/#windows-1250
	Windows1251               // https://encoding.spec.whatwg.org/#windows-1251
	Windows1252               // https://encoding.spec.whatwg.org/#windows-1252
	Windows1253               // https://encoding.spec.whatwg.org/#windows-1253
	Windows1254               // https://encoding.spec.whatwg.org/#windows-1254
	Windows1255               // https://encoding.spec.whatwg.org/#windows-1255
	Windows1256               // https://encoding.spec.whatwg.org/#windows-1256
	Windows1257               // https://encoding.spec.whatwg.org/#windows-1257
	Windows1258               // https://encoding.spec.whatwg.org/#windows-1258
	XMacCyrillic              // https://encoding.spec.whatwg.org/#x-mac-cyrillic
	Gbk                       // https://encoding.spec.whatwg.org/#gbk
	Gb18030                   // https://encoding.spec.whatwg.org/#gb18030
	Big5                      // https://encoding.spec.whatwg.org/#big5
	EucJp                     // https://encoding.spec.whatwg.org/#euc-jp
	Iso2022Jp                 // https://encoding.spec.whatwg.org/#iso-2022-jp
	ShiftJis                  // https://encoding.spec.whatwg.org/#shift_jis
	EucKr                     // https://encoding.spec.whatwg.org/#euc-kr
	Replacement               // https://encoding.spec.whatwg.org/#replacement
	Utf16Be                   // https://encoding.spec.whatwg.org/#utf-16be
	Utf16Le                   // https://encoding.spec.whatwg.org/#utf-16le
	XUserDefined              // https://encoding.spec.whatwg.org/#x-user-defined
)

var encodingLabelMap = map[string]Type{
	"unicode-1-1-utf-8": Utf8,
	"unicode11utf8":     Utf8,
	"unicode20utf8":     Utf8,
	"utf-8":             Utf8,
	"utf8":              Utf8,
	"x-unicode20utf8":   Utf8,

	"866":      Ibm866,
	"cp866":    Ibm866,
	"csibm866": Ibm866,
	"ibm866":   Ibm866,

	"csisolatin2":     Iso8859_2,
	"iso-8859-2":      Iso8859_2,
	"iso-ir-101":      Iso8859_2,
	"iso8859-2":       Iso8859_2,
	"iso88592":        Iso8859_2,
	"iso_8859-2":      Iso8859_2,
	"iso_8859-2:1987": Iso8859_2,
	"l2":              Iso8859_2,
	"latin2":          Iso8859_2,

	"csisolatin3":     Iso8859_3,
	"iso-8859-3":      Iso8859_3,
	"iso-ir-109":      Iso8859_3,
	"iso8859-3":       Iso8859_3,
	"iso88593":        Iso8859_3,
	"iso_8859-3":      Iso8859_3,
	"iso_8859-3:1988": Iso8859_3,
	"l3":              Iso8859_3,
	"latin3":          Iso8859_3,

	"csisolatin4":     Iso8859_4,
	"iso-8859-4":      Iso8859_4,
	"iso-ir-110":      Iso8859_4,
	"iso8859-4":       Iso8859_4,
	"iso88594":        Iso8859_4,
	"iso_8859-4":      Iso8859_4,
	"iso_8859-4:1988": Iso8859_4,
	"l4":              Iso8859_4,
	"latin4":          Iso8859_4,

	"csisolatincyrillic": Iso8859_5,
	"cyrillic":           Iso8859_5,
	"iso-8859-5":         Iso8859_5,
	"iso-ir-144":         Iso8859_5,
	"iso8859-5":          Iso8859_5,
	"iso88595":           Iso8859_5,
	"iso_8859-5":         Iso8859_5,
	"iso_8859-5:1988":    Iso8859_5,

	"arabic":           Iso8859_6,
	"asmo-708":         Iso8859_6,
	"csiso88596e":      Iso8859_6,
	"csiso88596i":      Iso8859_6,
	"csisolatinarabic": Iso8859_6,
	"ecma-114":         Iso8859_6,
	"iso-8859-6":       Iso8859_6,
	"iso-8859-6-e":     Iso8859_6,
	"iso-8859-6-i":     Iso8859_6,
	"iso-ir-127":       Iso8859_6,
	"iso8859-6":        Iso8859_6,
	"iso88596":         Iso8859_6,
	"iso_8859-6":       Iso8859_6,
	"iso_8859-6:1987":  Iso8859_6,

	"csisolatingreek": Iso8859_7,
	"ecma-118":        Iso8859_7,
	"elot_928":        Iso8859_7,
	"greek":           Iso8859_7,
	"greek8":          Iso8859_7,
	"iso-8859-7":      Iso8859_7,
	"iso-ir-126":      Iso8859_7,
	"iso8859-7":       Iso8859_7,
	"iso88597":        Iso8859_7,
	"iso_8859-7":      Iso8859_7,
	"iso_8859-7:1987": Iso8859_7,
	"sun_eu_greek":    Iso8859_7,

	"csiso88598e":      Iso8859_8,
	"csisolatinhebrew": Iso8859_8,
	"hebrew":           Iso8859_8,
	"iso-8859-8":       Iso8859_8,
	"iso-8859-8-e":     Iso8859_8,
	"iso-ir-138":       Iso8859_8,
	"iso8859-8":        Iso8859_8,
	"iso88598":         Iso8859_8,
	"iso_8859-8":       Iso8859_8,
	"iso_8859-8:1988":  Iso8859_8,
	"visual":           Iso8859_8,

	"csiso88598i":  Iso8859_8I,
	"iso-8859-8-i": Iso8859_8I,
	"logical":      Iso8859_8I,

	"csisolatin6": Iso8859_10,
	"iso-8859-10": Iso8859_10,
	"iso-ir-157":  Iso8859_10,
	"iso8859-10":  Iso8859_10,
	"iso885910":   Iso8859_10,
	"l6":          Iso8859_10,
	"latin6":      Iso8859_10,

	"iso-8859-13": Iso8859_13,
	"iso8859-13":  Iso8859_13,
	"iso885913":   Iso8859_13,

	"iso-8859-14": Iso8859_14,
	"iso8859-14":  Iso8859_14,
	"iso885914":   Iso8859_14,

	"csisolatin9": Iso8859_15,
	"iso-8859-15": Iso8859_15,
	"iso8859-15":  Iso8859_15,
	"iso885915":   Iso8859_15,
	"iso_8859-15": Iso8859_15,
	"l9":          Iso8859_15,

	"iso-8859-16": Iso8859_16,

	"cskoi8r": Koi8R,
	"koi":     Koi8R,
	"koi8":    Koi8R,
	"koi8-r":  Koi8R,
	"koi8_r":  Koi8R,

	"koi8-ru": Koi8U,
	"koi8-u":  Koi8U,

	"csmacintosh": Macintosh,
	"mac":         Macintosh,
	"macintosh":   Macintosh,
	"x-mac-roman": Macintosh,

	"dos-874":     Windows874,
	"iso-8859-11": Windows874,
	"iso8859-11":  Windows874,
	"iso885911":   Windows874,
	"tis-620":     Windows874,
	"windows-874": Windows874,

	"cp1250":       Windows1250,
	"windows-1250": Windows1250,
	"x-cp1250":     Windows1250,

	"cp1251":       Windows1251,
	"windows-1251": Windows1251,
	"x-cp1251":     Windows1251,

	"ansi_x3.4-1968":  Windows1252,
	"ascii":           Windows1252,
	"cp1252":          Windows1252,
	"cp819":           Windows1252,
	"csisolatin1":     Windows1252,
	"ibm819":          Windows1252,
	"iso-8859-1":      Windows1252,
	"iso-ir-100":      Windows1252,
	"iso8859-1":       Windows1252,
	"iso88591":        Windows1252,
	"iso_8859-1":      Windows1252,
	"iso_8859-1:1987": Windows1252,
	"l1":              Windows1252,
	"latin1":          Windows1252,
	"us-ascii":        Windows1252,
	"windows-1252":    Windows1252,
	"x-cp1252":        Windows1252,

	"cp1253":       Windows1253,
	"windows-1253": Windows1253,
	"x-cp1253":     Windows1253,

	"cp1254":          Windows1254,
	"csisolatin5":     Windows1254,
	"iso-8859-9":      Windows1254,
	"iso-ir-148":      Windows1254,
	"iso8859-9":       Windows1254,
	"iso88599":        Windows1254,
	"iso_8859-9":      Windows1254,
	"iso_8859-9:1989": Windows1254,
	"l5":              Windows1254,
	"latin5":          Windows1254,
	"windows-1254":    Windows1254,
	"x-cp1254":        Windows1254,

	"cp1255":       Windows1255,
	"windows-1255": Windows1255,
	"x-cp1255":     Windows1255,

	"cp1256":       Windows1256,
	"windows-1256": Windows1256,
	"x-cp1256":     Windows1256,

	"cp1257":       Windows1257,
	"windows-1257": Windows1257,
	"x-cp1257":     Windows1257,

	"cp1258":       Windows1258,
	"windows-1258": Windows1258,
	"x-cp1258":     Windows1258,

	"x-mac-cyrillic":  XMacCyrillic,
	"x-mac-ukrainian": XMacCyrillic,

	"chinese":         Gbk,
	"csgb2312":        Gbk,
	"csiso58gb231280": Gbk,
	"gb2312":          Gbk,
	"gb_2312":         Gbk,
	"gb_2312-80":      Gbk,
	"gbk":             Gbk,
	"iso-ir-58":       Gbk,
	"x-gbk":           Gbk,

	"gb18030": Gb18030,

	"big5":       Big5,
	"big5-hkscs": Big5,
	"cn-big5":    Big5,
	"csbig5":     Big5,
	"x-x-big5":   Big5,

	"cseucpkdfmtjapanese": EucJp,
	"euc-jp":              EucJp,
	"x-euc-jp":            EucJp,

	"csiso2022jp": Iso2022Jp,
	"iso-2022-jp": Iso2022Jp,

	"csshiftjis":  ShiftJis,
	"ms932":       ShiftJis,
	"ms_kanji":    ShiftJis,
	"shift-jis":   ShiftJis,
	"shift_jis":   ShiftJis,
	"sjis":        ShiftJis,
	"windows-31j": ShiftJis,
	"x-sjis":      ShiftJis,

	"cseuckr":        EucKr,
	"csksc56011987":  EucKr,
	"euc-kr":         EucKr,
	"iso-ir-149":     EucKr,
	"korean":         EucKr,
	"ks_c_5601-1987": EucKr,
	"ks_c_5601-1989": EucKr,
	"ksc5601":        EucKr,
	"ksc_5601":       EucKr,
	"windows-949":    EucKr,

	"csiso2022kr":     Replacement,
	"hz-gb-2312":      Replacement,
	"iso-2022-cn":     Replacement,
	"iso-2022-cn-ext": Replacement,
	"iso-2022-kr":     Replacement,
	"replacement":     Replacement,

	"unicodefffe": Utf16Be,
	"utf-16be":    Utf16Be,

	"csunicode":       Utf16Le,
	"iso-10646-ucs-2": Utf16Le,
	"ucs-2":           Utf16Le,
	"unicode":         Utf16Le,
	"unicodefeff":     Utf16Le,
	"utf-16":          Utf16Le,
	"utf-16le":        Utf16Le,

	"x-user-defined": XUserDefined,
}

func GetEncodingFromLabel(label string) (Type, error) {
	label = util.ToAsciiLowercase(strings.TrimFunc(label, func(r rune) bool { return r == ' ' }))
	encoding, ok := encodingLabelMap[label]
	if !ok {
		return 0, errors.New("no such encoding")
	}
	return encoding, nil
}

type handlerResult int64 // Positive values are actual result, negative values are special results(see below)
const (
	handlerResultError    = handlerResult(-99) // https://encoding.spec.whatwg.org/#error
	handlerResultFinished = handlerResult(-98) // https://encoding.spec.whatwg.org/#finished
	handlerResultContinue = handlerResult(-97) // https://encoding.spec.whatwg.org/#continue
)

type encoding struct {
	makeDecoder func() decoder
}

var encodings = map[Type]encoding{
	Utf8: utf8Encoding,
}

type decoder interface {
	handler(queue *IoQueue, byteItem IoQueueItem) handlerResult
}

// https://encoding.spec.whatwg.org/#decode
func Decode(input *IoQueue, fallbackEncodingType Type, output *IoQueue) {
	encodingType := fallbackEncodingType
	bomEncoding, ok := bomSniff(*input)
	if ok {
		encodingType = bomEncoding
		if bomEncoding == Utf8 {
			input.Read(3)
		} else {
			input.Read(2)
		}
	}
	encoding, ok := encodings[encodingType]
	if !ok {
		panic("unsupported encoding")
	}
	decoder := encoding.makeDecoder()
	decode(decoder, input, output, errorModeReplacement)
}

type errorMode uint8

const (
	errorModeReplacement = errorMode(iota)
	errorModeHtml
	errorModeFatal
)

func decode(decoder decoder, input *IoQueue, output *IoQueue, mode errorMode) handlerResult {
	for {
		item := input.ReadOne()
		res := decodeItem(item, decoder, input, output, mode)
		if res != handlerResultContinue {
			return res
		}
	}
}
func decodeItem(item IoQueueItem, decoder decoder, input *IoQueue, output *IoQueue, mode errorMode) handlerResult {
	if mode == errorModeHtml {
		panic("invalid errorMode")
	}
	res := decoder.handler(input, item)
	if res == handlerResultFinished {
		output.PushOne(IoQueueItem{EndOfQueue{}})
		return res
	} else if 0 <= res {
		if util.IsSurrogateChar(rune(res)) {
			panic("result cannot contain surrogate char")
		}
		output.PushOne(IoQueueItem{rune(res)})
	} else if res == handlerResultError {
		switch mode {
		case errorModeReplacement:
			output.PushOne(IoQueueItem{rune(0xfffd)})
		case errorModeHtml:
			panic("unreachable")
		case errorModeFatal:
			return res
		}
	}
	return handlerResultContinue
}

// https://encoding.spec.whatwg.org/#bom-sniff
func bomSniff(queue IoQueue) (Type, bool) {
	bytes := IoQueueItemsToSlice[uint8](queue.Peek(3))
	if 3 <= len(bytes) && slices.Equal([]byte{0xef, 0xbb, 0xbf}, bytes[:3]) {
		return Utf8, true
	} else if 2 <= len(bytes) && slices.Equal([]byte{0xfe, 0xff}, bytes[:2]) {
		return Utf16Be, true
	} else if 2 <= len(bytes) && slices.Equal([]byte{0xff, 0xfe}, bytes[:2]) {
		return Utf16Le, true
	}
	return 0, false
}

// https://encoding.spec.whatwg.org/#concept-stream
type IoQueue struct {
	items []IoQueueItem
}
type IoQueueItem struct {
	V any
}
type EndOfQueue struct{}

func (EndOfQueue) String() string { return "<end-of-queue>" }

func (item IoQueueItem) IsEndOfQueue() bool {
	if _, ok := item.V.(EndOfQueue); ok {
		return true
	}
	return false
}
func IoQueueItemsToSlice[T any](items []IoQueueItem) []T {
	res := []T{}
	for _, item := range items {
		if item.IsEndOfQueue() {
			break
		}
		res = append(res, item.V.(T))
	}
	return res
}
func IoQueueToSlice[T any](queue IoQueue) []T {
	return IoQueueItemsToSlice[T](queue.items)
}
func IoQueueFromSlice[T any](values []T) IoQueue {
	queue := IoQueue{items: []IoQueueItem{{EndOfQueue{}}}}
	for _, item := range values {
		queue.PushOne(IoQueueItem{item})
	}
	return queue
}

// https://encoding.spec.whatwg.org/#concept-stream-read
func (q *IoQueue) ReadOne() IoQueueItem {
	if len(q.items) == 0 {
		panic("The queue must not be empty")
	}
	item := q.items[0]
	if item.IsEndOfQueue() {
		return item
	}
	q.items = q.items[1:]
	return item
}

// https://encoding.spec.whatwg.org/#concept-stream-Read
func (q *IoQueue) Read(num int) []IoQueueItem {
	readItems := []IoQueueItem{}
	for range num {
		item := q.ReadOne()
		if !item.IsEndOfQueue() {
			readItems = append(readItems, item)
		}
	}
	return readItems
}

// https://encoding.spec.whatwg.org/#i-o-queue-Peek
func (q *IoQueue) Peek(num int) []IoQueueItem {
	prefix := []IoQueueItem{}
	for i := range num {
		item := q.items[i]
		if item.IsEndOfQueue() {
			break
		}
		prefix = append(prefix, item)
	}
	return prefix
}

// https://encoding.spec.whatwg.org/#concept-stream-push
func (q *IoQueue) PushOne(item IoQueueItem) {
	if len(q.items) == 0 || !q.items[len(q.items)-1].IsEndOfQueue() {
		panic("the last item must be end-of-queue")
	}
	if item.IsEndOfQueue() {
		return
	}
	if len(q.items)-1 < 0 {
		q.items = append([]IoQueueItem{}, item, q.items[len(q.items)-1])
	} else {
		q.items = append(q.items[:len(q.items)-1], item, q.items[len(q.items)-1])
	}
}

// https://encoding.spec.whatwg.org/#concept-stream-Push
func (q *IoQueue) Push(items []IoQueueItem) {
	for _, item := range items {
		q.PushOne(item)
	}
}

// https://encoding.spec.whatwg.org/#concept-stream-prepend
func (q *IoQueue) RestoreOne(item IoQueueItem) {
	if item.IsEndOfQueue() {
		panic("attempted to restore end-of-queue item")
	}
	q.items = append([]IoQueueItem{item}, q.items...)
}

// https://encoding.spec.whatwg.org/#concept-stream-prepend
func (q *IoQueue) Restore(items []IoQueueItem) {
	if slices.ContainsFunc(items, IoQueueItem.IsEndOfQueue) {
		panic("attempted to restore end-of-queue item")
	}
	q.items = append(items, q.items...)
}
