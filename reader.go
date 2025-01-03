package javadatastream

import (
	"errors"
	"fmt"
	"io"
	"math"
)

type IDataInputStream interface {
	Read(p []byte) (n int, err error)
	ReadByte() (byte, error)
	ReadBoolean() (bool, error)
	ReadUShort() (uint16, error)
	ReadShort() (int16, error)
	ReadUInt() (uint32, error)
	ReadInt() (int32, error)
	ReadULong() (uint64, error)
	ReadLong() (int64, error)
	ReadFloat() (float32, error)
	ReadDouble() (float64, error)
	ReadChar() (rune, error)
	ReadUTF() (string, error)
}

// DataInputStream implements java.io.DataInputStream
type DataInputStream struct {
	r    io.Reader
	buf1 []byte
	buf2 []byte
	buf4 []byte
	buf8 []byte
}

func NewReader(r io.Reader) *DataInputStream {
	return &DataInputStream{
		r:    r,
		buf1: make([]byte, 1),
		buf2: make([]byte, 2),
		buf4: make([]byte, 4),
		buf8: make([]byte, 8),
	}
}

func (r DataInputStream) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

func (r DataInputStream) ReadByte() (byte, error) {
	if _, err := io.ReadFull(r.r, r.buf1); err != nil {
		return 0, err
	}
	return r.buf1[0], nil
}

func (r DataInputStream) ReadBoolean() (bool, error) {
	if _, err := io.ReadFull(r.r, r.buf1); err != nil {
		return false, err
	}
	return r.buf1[0] != 0, nil
}

func (r DataInputStream) ReadUShort() (uint16, error) {
	if _, err := io.ReadFull(r.r, r.buf2); err != nil {
		return 0, fmt.Errorf("r.ReadUShort failed with: %w", err)
	}
	return uint16(r.buf2[1]) | uint16(r.buf2[0])<<8, nil
}

func (r DataInputStream) ReadShort() (int16, error) {
	v, err := r.ReadUShort()
	if err != nil {
		return 0, err
	}
	return int16(v), nil
}

func (r DataInputStream) ReadUInt() (uint32, error) {
	if _, err := io.ReadFull(r.r, r.buf4); err != nil {
		return 0, err
	}
	return uint32(r.buf4[3]) | uint32(r.buf4[2])<<8 | uint32(r.buf4[1])<<16 | uint32(r.buf4[0])<<24, nil
}

func (r DataInputStream) ReadInt() (int32, error) {
	v, err := r.ReadUInt()
	if err != nil {
		return 0, err
	}
	return int32(v), err
}

func (r DataInputStream) ReadULong() (uint64, error) {
	if _, err := io.ReadFull(r.r, r.buf8); err != nil {
		return 0, err
	}
	return uint64(r.buf8[7]) | uint64(r.buf8[6])<<8 | uint64(r.buf8[5])<<16 |
		uint64(r.buf8[4])<<24 | uint64(r.buf8[3])<<32 | uint64(r.buf8[2])<<40 |
		uint64(r.buf8[1])<<48 | uint64(r.buf8[0])<<56, nil
}

func (r DataInputStream) ReadLong() (int64, error) {
	v, err := r.ReadULong()
	if err != nil {
		return 0, err
	}
	return int64(v), err
}

func (r DataInputStream) ReadFloat() (float32, error) {
	v, err := r.ReadUInt()
	if err != nil {
		return 0.0, err
	}
	return math.Float32frombits(v), nil
}

func (r DataInputStream) ReadDouble() (float64, error) {
	v, err := r.ReadULong()
	if err != nil {
		return 0.0, err
	}
	return math.Float64frombits(v), nil
}

func (r DataInputStream) ReadChar() (rune, error) {
	v, err := r.ReadUShort()
	if err != nil {
		return '0', err
	}
	return rune(v), nil
}

func (r DataInputStream) ReadUTF() (string, error) {
	utflen, err := r.ReadUShort()
	if err != nil {
		return "", err
	}
	bArr := make([]byte, utflen)

	_, err = io.ReadFull(r.r, bArr)
	if err != nil {
		return "", err
	}

	cArr := make([]rune, utflen*3)
	var count uint16 = 0
	var cArrCount uint32 = 0
	var c2 byte
	var c3 byte
	for count < utflen {
		c := uint32(bArr[count] & 0xFF)
		switch c >> 4 {
		case 0, 1, 2, 3, 4, 5, 6, 7:
			count++
			cArr[cArrCount] = rune(c)
			cArrCount++

		case 12, 13:
			count += 2
			if count > utflen {
				return "", errors.New("malformed input: partial character at end")
			}
			c2 = bArr[count-1]
			if (c2 & 0xC0) != 0x80 {
				return "", fmt.Errorf("malformed input around byte %d", count)
			}
			cArr[cArrCount] = rune(int32((c&0x1F)<<6) | int32((c2 & 0x3F)))
			cArrCount++

		case 14:
			count += 3
			if count > utflen {
				return "", errors.New("malformed input: partial character at end")
			}
			c2 = bArr[count-2]
			c3 = bArr[count-1]
			if (c2&0xC0) != 0x80 || (c3&0xC0) != 0x80 {
				return "", fmt.Errorf("malformed input around byte ")
			}
			cArr[cArrCount] = rune(int32((c&0x0F)<<12) | int32((c2&0x3F)<<6) | int32(c3&0x3F))
			cArrCount++
		default:
			return "", fmt.Errorf("malformed input around byte ")
		}
	}

	return string(cArr[:cArrCount]), nil
}
