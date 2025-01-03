package javadatastream

import (
	"fmt"
	"io"
	"math"
)

type IDataOutputStream interface {
	Write(p []byte) (n int, err error)
	WriteByte(byte) error
	WriteBoolean(bool) error
	WriteUShort(uint16) error
	WriteShort(int16) error
	WriteChar(rune) error
	WriteInt(int32) error
	WriteLong(int64) error
	WriteFloat(float32) error
	WriteDouble(float64) error
	WriteUTF(string) error
}

// DataOutputStream implements IDataOutputStream
type DataOutputStream struct {
	w       io.Writer
	wClosed bool
	wbuf    []byte
}

func NewWriter(w io.Writer) DataOutputStream {
	return DataOutputStream{
		w:       w,
		wClosed: false,
		wbuf:    make([]byte, 8),
	}
}

func (w DataOutputStream) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

func (w DataOutputStream) WriteByte(b byte) error {
	_, err := w.w.Write([]byte{b})
	return err
}

func (w DataOutputStream) WriteBoolean(v bool) error {
	if v {
		_, err := w.Write([]byte{1})
		return err
	}
	_, err := w.Write([]byte{0})
	return err
}

func (w DataOutputStream) WriteUShort(v uint16) error {
	w.wbuf[0] = byte(v >> 8)
	w.wbuf[1] = byte(v >> 0)
	_, err := w.Write(w.wbuf[0:2])
	return err
}

func (w DataOutputStream) WriteShort(v int16) error {
	w.wbuf[0] = byte(v >> 8)
	w.wbuf[1] = byte(v >> 0)
	_, err := w.Write(w.wbuf[0:2])
	return err
}

func (w DataOutputStream) WriteChar(v rune) error {
	chr := uint16(v)
	w.wbuf[1] = byte(chr >> 8)
	w.wbuf[0] = byte(chr >> 0)
	_, err := w.Write(w.wbuf[0:2])
	return err
}

func (w DataOutputStream) WriteInt(v int32) error {
	w.wbuf[0] = byte(v >> 24)
	w.wbuf[1] = byte(v >> 16)
	w.wbuf[2] = byte(v >> 8)
	w.wbuf[3] = byte(v >> 0)
	_, err := w.Write(w.wbuf[0:4])
	return err
}

func (w DataOutputStream) WriteLong(v int64) error {
	w.wbuf[0] = byte(v >> 56)
	w.wbuf[1] = byte(v >> 48)
	w.wbuf[2] = byte(v >> 40)
	w.wbuf[3] = byte(v >> 32)
	w.wbuf[4] = byte(v >> 24)
	w.wbuf[5] = byte(v >> 16)
	w.wbuf[6] = byte(v >> 8)
	w.wbuf[7] = byte(v >> 0)
	_, err := w.Write(w.wbuf[0:8])
	return err
}

func (w DataOutputStream) WriteFloat(v float32) error {
	bits := math.Float32bits(v)
	w.wbuf[0] = byte(bits >> 24)
	w.wbuf[1] = byte(bits >> 16)
	w.wbuf[2] = byte(bits >> 8)
	w.wbuf[3] = byte(bits >> 0)
	_, err := w.Write(w.wbuf[0:4])
	return err
}

func (w DataOutputStream) WriteDouble(v float64) error {
	bits := math.Float64bits(v)
	w.wbuf[0] = byte(bits >> 56)
	w.wbuf[1] = byte(bits >> 48)
	w.wbuf[2] = byte(bits >> 40)
	w.wbuf[3] = byte(bits >> 32)
	w.wbuf[4] = byte(bits >> 24)
	w.wbuf[5] = byte(bits >> 16)
	w.wbuf[6] = byte(bits >> 8)
	w.wbuf[7] = byte(bits >> 0)
	_, err := w.Write(w.wbuf[0:8])
	return err
}

func (w DataOutputStream) WriteUTF(s string) error {
	if len(s) > 65535 {
		return fmt.Errorf("unable to write utf: string exceeds length")
	}
	buf := make([]byte, len(s))
	count := 0
	chars := []rune(s)
	var i uint64
	for i = 0; i < uint64(len(s)); i++ {
		c := chars[i]
		if c >= 0x80 || c == 0 {
			break
		}
		buf[i] = byte(c)
	}

	for ; i < uint64(len(s)); i++ {
		c := chars[i]
		if c < 0x80 && c != 0 {
			buf[count] = uint8(c)
			count++
		} else if c >= 0x800 {
			buf[count] = byte(0xE0 | (c>>12)&0x0F)
			buf[count+1] = byte(0x80 | (c>>6)&0x3F)
			buf[count+2] = byte(0x80 | c&0x3F)
			count += 3
		} else {
			buf[count] = byte(0xC0 | (c>>6)&0x1F)
			buf[count+1] = byte(0x80 | c&0x3F)
			count += 2
		}
	}

	length := len(buf)
	if length > 65535 {
		return fmt.Errorf("unable to write utf: string exceeds length")
	}
	if err := w.WriteUShort(uint16(length)); err != nil {
		return err
	}
	_, err := w.w.Write(buf)
	return err
}
