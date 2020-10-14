package binary

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

var ErrVarintOverflow = fmt.Errorf("libext-go/encoding/binary: varint overflows a 64-bit integer")

type Reader struct {
	buf [8]byte

	byteOrder binary.ByteOrder
	rd        io.Reader
}

func NewBigEndianReader(r io.Reader) *Reader {
	return NewReader(binary.BigEndian, r)
}

func NewLittleEndianReader(r io.Reader) *Reader {
	return NewReader(binary.LittleEndian, r)
}

func NewReader(byteOrder binary.ByteOrder, r io.Reader) *Reader {
	return &Reader{
		byteOrder: byteOrder,
		rd:        r,
	}
}

func (r *Reader) Reset(rd io.Reader) {
	r.rd = rd
}

func (r *Reader) Read(p []byte) (int, error) {
	return r.rd.Read(p)
}

func (r *Reader) ReadByte() (b byte, err error) {
	err = r.ReadUint8(&b)
	return
}

// func (r *Reader) ReadAny(v interface{}) error {} // XXX

func (r *Reader) ReadInt8(v *int8) error {
	var uv uint8
	if err := r.ReadUint8(&uv); err != nil {
		return err
	}
	*v = int8(uv)
	return nil
}

func (r *Reader) ReadInt16(v *int16) error {
	var uv uint16
	if err := r.ReadUint16(&uv); err != nil {
		return err
	}
	*v = int16(uv)
	return nil
}

func (r *Reader) ReadInt32(v *int32) error {
	var uv uint32
	if err := r.ReadUint32(&uv); err != nil {
		return err
	}
	*v = int32(uv)
	return nil
}

func (r *Reader) ReadInt64(v *int64) error {
	var uv uint64
	if err := r.ReadUint64(&uv); err != nil {
		return err
	}
	*v = int64(uv)
	return nil
}

func (r *Reader) ReadUint8(v *uint8) error {
	buf := r.buf[:1]
	if _, err := r.rd.Read(buf); err != nil {
		return err
	}
	*v = buf[0]
	return nil
}

func (r *Reader) ReadUint16(v *uint16) error {
	buf := r.buf[:2]
	if _, err := r.rd.Read(buf); err != nil {
		return err
	}
	*v = r.byteOrder.Uint16(buf)
	return nil
}

func (r *Reader) ReadUint32(v *uint32) error {
	buf := r.buf[:4]
	if _, err := r.rd.Read(buf); err != nil {
		return err
	}
	*v = r.byteOrder.Uint32(buf)
	return nil
}

func (r *Reader) ReadUint64(v *uint64) error {
	buf := r.buf[:8]
	if _, err := r.rd.Read(buf); err != nil {
		return err
	}
	*v = r.byteOrder.Uint64(buf)
	return nil
}

func (r *Reader) ReadFloat32(v *float32) error {
	var bits uint32
	if err := r.ReadUint32(&bits); err != nil {
		return err
	}
	*v = math.Float32frombits(bits)
	return nil
}

func (r *Reader) ReadFloat64(v *float64) error {
	var bits uint64
	if err := r.ReadUint64(&bits); err != nil {
		return err
	}
	*v = math.Float64frombits(bits)
	return nil
}

func (r *Reader) ReadVarint(v *int64) error {
	var uv uint64
	if err := r.ReadUvarint(&uv); err != nil {
		return err
	}
	x := int64(uv >> 1)
	if uv&1 != 0 {
		x = ^x
	}
	*v = x
	return nil
}

func (r *Reader) ReadUvarint(v *uint64) error {
	// Copy and modified from golang source code.
	var x uint64
	var s uint
	var b uint8
	for i := 0; i < binary.MaxVarintLen64; i++ {
		if err := r.ReadUint8(&b); err != nil {
			return err
		}
		if b < 0x80 {
			if i == 9 && b > 1 {
				return ErrVarintOverflow
			}
			*v = x | uint64(b)<<s
			return nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return ErrVarintOverflow
}
