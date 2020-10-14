package binary

import (
	"encoding/binary"
	"io"
	"math"
)

type Writer struct {
	buf [binary.MaxVarintLen64]byte

	byteOrder binary.ByteOrder
	wr        io.Writer
}

func NewBigEndianWriter(w io.Writer) *Writer {
	return NewWriter(binary.BigEndian, w)
}

func NewLittleEndianWriter(w io.Writer) *Writer {
	return NewWriter(binary.LittleEndian, w)
}

func NewWriter(byteOrder binary.ByteOrder, w io.Writer) *Writer {
	return &Writer{
		byteOrder: byteOrder,
		wr:        w,
	}
}

func (w *Writer) Reset(wr io.Writer) {
	w.wr = wr
}

func (w *Writer) Write(p []byte) (int, error) {
	return w.wr.Write(p)
}

func (w *Writer) WriteByte(b byte) error {
	return w.WriteUint8(b)
}

// func (w *Writer) WriteAny(v interface{}) error {} // XXX

func (w *Writer) WriteInt8(v int8) error {
	return w.WriteUint8(uint8(v))
}

func (w *Writer) WriteInt16(v int16) error {
	return w.WriteUint16(uint16(v))
}

func (w *Writer) WriteInt32(v int32) error {
	return w.WriteUint32(uint32(v))
}

func (w *Writer) WriteInt64(v int64) error {
	return w.WriteUint64(uint64(v))
}

func (w *Writer) WriteUint8(v uint8) error {
	_, err := w.wr.Write([]byte{v})
	return err
}

func (w *Writer) WriteUint16(v uint16) error {
	buf := w.buf[:2]
	w.byteOrder.PutUint16(buf, v)
	_, err := w.Write(buf)
	return err
}

func (w *Writer) WriteUint32(v uint32) error {
	buf := w.buf[:4]
	w.byteOrder.PutUint32(buf, v)
	_, err := w.Write(buf)
	return err
}

func (w *Writer) WriteUint64(v uint64) error {
	buf := w.buf[:8]
	w.byteOrder.PutUint64(buf, v)
	_, err := w.Write(buf)
	return err
}

func (w *Writer) WriteFloat32(v float32) error {
	bits := math.Float32bits(v)
	if err := w.WriteUint32(bits); err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteFloat64(v float64) error {
	bits := math.Float64bits(v)
	if err := w.WriteUint64(bits); err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteVarint(v int64) error {
	buf := w.buf[:binary.MaxVarintLen64]
	offset := binary.PutVarint(buf, v)
	_, err := w.Write(buf[:offset])
	return err
}

func (w *Writer) WriteUvarint(v uint64) error {
	buf := w.buf[:binary.MaxVarintLen64]
	offset := binary.PutUvarint(buf, v)
	_, err := w.Write(buf[:offset])
	return err
}
