package utils

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// 0x00 for nil or fallback
func TestPackNil(t *testing.T) {
	assertPackResult(t, nil, []byte{0})
	assertPackResult(t, struct{}{}, []byte{0})
}

// 0x01 for false
// 0x02 for true
func TestPackBool(t *testing.T) {
	assertPackResult(t, false, []byte{1})
	assertPackResult(t, true, []byte{2})
}

// 0x03 for type float, 64-bit IEEE 754 (big endian)
func TestPackNumber(t *testing.T) {
	for _, i := range []interface{}{uint(42), uint8(42), uint64(42), int(42), int8(42), int64(42)} {
		assertPackResult(t, i, []byte{
			// type
			3,
			// IEEE 754
			64, 69, 0, 0, 0, 0, 0, 0,
		})
	}

	for _, f := range []interface{}{float32(42.125), float64(42.125)} {
		assertPackResult(t, f, []byte{
			// type
			3,
			// IEEE 754
			64, 69, 16, 0, 0, 0, 0, 0,
		})
	}
}

// 0x04 for type string, length as 64-bit unsigned int (big endian), UTF-8 encoded (or possibly binary) payload
func TestPackString(t *testing.T) {
	in := "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f" + // ASCII (0 to 127)
		"\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f" +
		"\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2a\x2b\x2c\x2d\x2e\x2f" +
		"\x30\x31\x32\x33\x34\x35\x36\x37\x38\x39\x3a\x3b\x3c\x3d\x3e\x3f" +
		"\x40\x41\x42\x43\x44\x45\x46\x47\x48\x49\x4a\x4b\x4c\x4d\x4e\x4f" +
		"\x50\x51\x52\x53\x54\x55\x56\x57\x58\x59\x5a\x5b\x5c\x5d\x5e\x5f" +
		"\x60\x61\x62\x63\x64\x65\x66\x67\x68\x69\x6a\x6b\x6c\x6d\x6e\x6f" +
		"\x70\x71\x72\x73\x74\x75\x76\x77\x78\x79\x7a\x7b\x7c\x7d\x7e\x7f" +
		"áéíóú" // some keyboard-independent non-ASCII unicode characters

	out := []byte{
		// type
		4,
		// length
		0, 0, 0, 0, 0, 0, 0, 138,
		// ASCII
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
		80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95,
		96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111,
		112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127,
		// UTF-8
		195, 161, 195, 169, 195, 173, 195, 179, 195, 186,
	}

	assertPackResult(t, in, out)
}

// 0x05 for type array, length as 64-bit unsigned int (big endian), zero or more elements (any type, as given by length)
func TestPackArray(t *testing.T) {
	in := []interface{}{nil, false, true, 42.125, "foobar"}

	out := []byte{
		// type
		5,
		// length
		0, 0, 0, 0, 0, 0, 0, 5,
		// nil
		0,
		// false
		1,
		// true
		2,
		// 42.125
		3,
		64, 69, 16, 0, 0, 0, 0, 0,
		// "foobar"
		4,
		0, 0, 0, 0, 0, 0, 0, 6,
		102, 111, 111, 98, 97, 114,
	}

	assertPackResult(t, in, out)
}

// 0x06 for type dictionary, length as 64-bit unsigned int (big endian),
// zero or more key-value pairs (as given by length) ordered by key (collate: binary/latin1)
//
// key-value pair: key (string without type prefix), value (any)
func TestPackDict(t *testing.T) {
	in := map[string]interface{}{
		"nil":             nil,
		"false":           false,
		"true":            true,
		"42.125":          42.125,
		"foobar":          "foobar",
		"[]interface{}{}": []interface{}{},
	}

	out := []byte{
		// type
		6,
		// length
		0, 0, 0, 0, 0, 0, 0, 6,
		// "42.125"
		0, 0, 0, 0, 0, 0, 0, 6,
		52, 50, 46, 49, 50, 53,
		// 42.125
		3,
		64, 69, 16, 0, 0, 0, 0, 0,
		// "[]interface{}{}"
		0, 0, 0, 0, 0, 0, 0, 15,
		91, 93, 105, 110, 116, 101, 114, 102, 97, 99, 101, 123, 125, 123, 125,
		// []interface{}{}
		5,
		0, 0, 0, 0, 0, 0, 0, 0,
		// "false"
		0, 0, 0, 0, 0, 0, 0, 5,
		102, 97, 108, 115, 101,
		// false
		1,
		// "foobar"
		0, 0, 0, 0, 0, 0, 0, 6,
		102, 111, 111, 98, 97, 114,
		// "foobar"
		4,
		0, 0, 0, 0, 0, 0, 0, 6,
		102, 111, 111, 98, 97, 114,
		// "nil"
		0, 0, 0, 0, 0, 0, 0, 3,
		110, 105, 108,
		// nil
		0,
		// "true"
		0, 0, 0, 0, 0, 0, 0, 4,
		116, 114, 117, 101,
		// true
		2,
	}

	assertPackResult(t, in, out)
}

func assertPackResult(t *testing.T, in interface{}, out []byte) {
	t.Helper()

	var buf bytes.Buffer
	ObjectPacker{in}.WriteTo(&buf)

	res := buf.Bytes()
	if bytes.Compare(res, out) != 0 {
		t.Errorf(
			"got unexpected result while serialized ( %#v ):\n--- %s\n+++ %s\n",
			in,
			hex.EncodeToString(out),
			hex.EncodeToString(res),
		)
	}
}
