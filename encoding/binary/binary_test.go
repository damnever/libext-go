package binary

import (
	"bytes"
	"crypto/rand"
	"io"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterface(t *testing.T) {
	var _ io.Reader = &Reader{}
	var _ io.Writer = &Writer{}
	var _ io.ByteReader = &Reader{}
	var _ io.ByteWriter = &Writer{}
}

func TestBigEndian(t *testing.T) {
	buf := &bytes.Buffer{}
	testReaderWriter(t, NewBigEndianReader(buf), NewBigEndianWriter(buf))
}

func TestLittleEndian(t *testing.T) {
	buf := &bytes.Buffer{}
	testReaderWriter(t, NewLittleEndianReader(buf), NewLittleEndianWriter(buf))
}

func testReaderWriter(t *testing.T, r *Reader, w *Writer) {
	{ // byte
		for _, b := range []byte("abcdefghijklmnopqrstuvwxyz") {
			require.Nil(t, w.WriteByte(b))
			actual, err := r.ReadByte()
			require.Nil(t, err)
			require.Equal(t, b, actual)
		}
	}
	{ // bytes
		for i := 8; i <= 128; i++ {
			buf := make([]byte, i, i)
			_, err := rand.Read(buf)
			require.Nil(t, err)
			_, err = w.Write(buf)
			require.Nil(t, err)
			actual := make([]byte, i, i)
			_, err = r.Read(actual)
			require.Nil(t, err)
			require.Equal(t, buf, actual)
		}
	}
	{ // int8
		for _, n := range []int8{math.MinInt8, -32, -64, 0, 1, 3, 8, 32, 65, math.MaxInt8} {
			require.Nil(t, w.WriteInt8(n))
			var actual int8
			require.Nil(t, r.ReadInt8(&actual))
			require.Equal(t, n, actual)
		}
	}
	{ // uint8
		for _, n := range []uint8{0, 1, 3, 8, 32, 65, math.MaxInt8, math.MaxUint8} {
			require.Nil(t, w.WriteUint8(n))
			var actual uint8
			require.Nil(t, r.ReadUint8(&actual))
			require.Equal(t, n, actual)
		}
	}
	{ // int16
		for _, n := range []int16{
			math.MinInt16, math.MinInt8, -32, -64, 0, 1, 3, 32, 65,
			math.MaxInt8, math.MaxUint8, math.MaxInt16,
		} {
			require.Nil(t, w.WriteInt16(n))
			var actual int16
			require.Nil(t, r.ReadInt16(&actual))
			require.Equal(t, n, actual)
		}
	}
	{ // uint16
		for _, n := range []uint16{
			0, 1, 3, 8, 32, 65, math.MaxInt8,
			math.MaxUint8, math.MaxInt16, math.MaxUint16,
		} {
			require.Nil(t, w.WriteUint16(n))
			var actual uint16
			require.Nil(t, r.ReadUint16(&actual))
			require.Equal(t, n, actual)
		}
	}
	{ // int32
		for _, n := range []int32{
			math.MinInt32, math.MinInt16, math.MinInt8, 0, 1,
			math.MaxInt8, math.MaxUint8, math.MaxInt16, math.MaxUint16, math.MaxInt32,
		} {
			require.Nil(t, w.WriteInt32(n))
			var actual int32
			require.Nil(t, r.ReadInt32(&actual))
			require.Equal(t, n, actual)
		}
	}
	{ // uint32
		for _, n := range []uint32{
			0, 1, 3, 8, 32, 65, math.MaxInt8, math.MaxUint8,
			math.MaxInt16, math.MaxUint16, math.MaxInt32, math.MaxUint32,
		} {
			require.Nil(t, w.WriteUint32(n))
			var actual uint32
			require.Nil(t, r.ReadUint32(&actual))
			require.Equal(t, n, actual)
		}
	}
	{ // int64
		for _, n := range []int64{
			math.MinInt64, math.MinInt32, math.MinInt16, math.MinInt8, 0, 1,
			math.MaxInt8, math.MaxUint8, math.MaxInt16, math.MaxUint16, math.MaxInt32,
			math.MaxUint32, math.MaxInt64,
		} {
			require.Nil(t, w.WriteInt64(n))
			var actual int64
			require.Nil(t, r.ReadInt64(&actual))
			require.Equal(t, n, actual)
		}
	}
	{ // uint64
		for _, n := range []uint64{
			0, 1, 3, 8, 32, 65, math.MaxInt8, math.MaxUint8,
			math.MaxInt16, math.MaxUint16, math.MaxInt32, math.MaxUint32,
			math.MaxInt64, math.MaxUint64,
		} {
			require.Nil(t, w.WriteUint64(n))
			var actual uint64
			require.Nil(t, r.ReadUint64(&actual))
			require.Equal(t, n, actual)
		}
	}
	{ // float32
		for _, n := range []float32{
			math.MinInt32, math.MinInt16, math.MinInt8, 0, 1.0 / 3.0, 1,
			math.MaxInt8, math.MaxUint8, math.MaxInt16, math.MaxUint16, math.MaxInt32,
		} {
			require.Nil(t, w.WriteFloat32(n))
			var actual float32
			require.Nil(t, r.ReadFloat32(&actual))
			require.Equal(t, n, actual)
		}
	}
	{ // float64
		for _, n := range []float64{
			math.MinInt64, math.MinInt32, math.MinInt16, math.MinInt8, 0, 1.0 / 3.0, 1,
			math.MaxInt8, math.MaxUint8, math.MaxInt16, math.MaxUint16, math.MaxInt32,
			math.MaxUint32, math.MaxInt64,
		} {
			require.Nil(t, w.WriteFloat64(n))
			var actual float64
			require.Nil(t, r.ReadFloat64(&actual))
			require.Equal(t, n, actual)
		}
	}
	{ // varint
		for _, n := range []int64{
			math.MinInt64, math.MinInt32, math.MinInt16, math.MinInt8, 0, 1,
			math.MaxInt8, math.MaxUint8, math.MaxInt16, math.MaxUint16, math.MaxInt32,
			math.MaxUint32, math.MaxInt64,
		} {
			require.Nil(t, w.WriteVarint(n))
			var actual int64
			require.Nil(t, r.ReadVarint(&actual))
			require.Equal(t, n, actual)
		}
	}
	{ // uvarint
		for _, n := range []uint64{
			0, 1, 3, 8, 32, 65, math.MaxInt8, math.MaxUint8,
			math.MaxInt16, math.MaxUint16, math.MaxInt32, math.MaxUint32,
			math.MaxInt64, math.MaxUint64,
		} {
			require.Nil(t, w.WriteUvarint(n))
			var actual uint64
			require.Nil(t, r.ReadUvarint(&actual))
			require.Equal(t, n, actual)
		}
	}
}
