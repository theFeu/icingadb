package utils

import (
	"io"
	"math"
	"sort"
)

// ObjectPacker serializes its #Object consistently to an io.Writer (similar to BSON),
// suitable for making checksums of objects.
type ObjectPacker struct {
	Object interface{}
}

var _ io.WriterTo = ObjectPacker{}

func (o ObjectPacker) WriteTo(w io.Writer) (n int64, err error) {
	err = packAny(o.Object, w, &n)
	return
}

// packAny serializes o to w and increments *n by the bytes written.
func packAny(o interface{}, w io.Writer, n *int64) (err error) {
	switch v := o.(type) {
	case string:
		if err = writeHelper([]byte{4}, w, n); err != nil {
			return
		}

		{
			length := packUInt64BE(uint64(len(v)))
			if err = writeHelper(length[:], w, n); err != nil {
				return
			}
		}

		return writeHelper([]byte(v), w, n)

	case float64:
		return packFloat(v, w, n)

	case float32:
		return packFloat(float64(v), w, n)

	case int:
		return packFloat(float64(v), w, n)

	case int64:
		return packFloat(float64(v), w, n)

	case int32:
		return packFloat(float64(v), w, n)

	case int16:
		return packFloat(float64(v), w, n)

	case int8:
		return packFloat(float64(v), w, n)

	case uint:
		return packFloat(float64(v), w, n)

	case uint64:
		return packFloat(float64(v), w, n)

	case uint32:
		return packFloat(float64(v), w, n)

	case uint16:
		return packFloat(float64(v), w, n)

	case uint8:
		return packFloat(float64(v), w, n)

	case bool:
		if v {
			return writeHelper([]byte{2}, w, n)
		}

		return writeHelper([]byte{1}, w, n)

	case nil:
		return writeHelper([]byte{0}, w, n)

	case map[string]interface{}:
		if err = writeHelper([]byte{6}, w, n); err != nil {
			return
		}

		{
			length := packUInt64BE(uint64(len(v)))
			if err = writeHelper(length[:], w, n); err != nil {
				return
			}
		}

		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			{
				length := packUInt64BE(uint64(len(k)))
				if err = writeHelper(length[:], w, n); err != nil {
					return
				}
			}

			if err = writeHelper([]byte(k), w, n); err != nil {
				return
			}

			if err = packAny(v[k], w, n); err != nil {
				return
			}
		}

		return

	case []interface{}:
		if err = writeHelper([]byte{5}, w, n); err != nil {
			return
		}

		{
			length := packUInt64BE(uint64(len(v)))
			if err = writeHelper(length[:], w, n); err != nil {
				return
			}
		}

		for _, v := range v {
			if err = packAny(v, w, n); err != nil {
				return
			}
		}

		return

	default:
		return writeHelper([]byte{0}, w, n)
	}
}

// packFloat serializes f to w and increments *n by the bytes written.
func packFloat(f float64, w io.Writer, n *int64) (err error) {
	if err = writeHelper([]byte{3}, w, n); err != nil {
		return
	}

	// This requires the float edianness to be the same as the integer one
	//
	// LE machine:
	// math.Float64bits(float64[0x0102030405060708]) = uint64[0x0102030405060708]
	// packUInt64BE(uint64[0x0102030405060708]) = [8]byte[0x0807060504030201]
	//
	// BE machine:
	// math.Float64bits(float64[0x0807060504030201]) = uint64[0x0807060504030201]
	// packUInt64BE(uint64[0x0807060504030201]) = [8]byte[0x0807060504030201]
	bits := packUInt64BE(math.Float64bits(f))
	return writeHelper(bits[:], w, n)
}

// writeHelper writes p to w and increments *n by the bytes written.
func writeHelper(p []byte, w io.Writer, n *int64) (err error) {
	var m int

	m, err = w.Write(p)
	*n += int64(m)

	return
}

// packUInt64BE converts i to [8]byte (big endian).
func packUInt64BE(i uint64) [8]byte {
	return [8]byte{
		byte(i >> 56),
		byte((i >> 48) & 255),
		byte((i >> 40) & 255),
		byte((i >> 32) & 255),
		byte((i >> 24) & 255),
		byte((i >> 16) & 255),
		byte((i >> 8) & 255),
		byte(i & 255),
	}
}
