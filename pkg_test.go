package javadatastream

import (
	"bytes"
	"slices"
	"testing"
)

func TestWriteReadUTF(t *testing.T) {
	expected := "Hello World"
	buf := new(bytes.Buffer)
	w := NewWriter(buf)
	if err := w.WriteUTF(expected); err != nil {
		t.Fatal("expected WriteUTF() to be properly written but got: %w", err)
	}

	r := NewReader(buf)
	actual, err := r.ReadUTF()
	if err != nil {
		t.Fatalf("expected ReadUTF() to read out the string %#v but got error: %s", expected, err)
	}
	if actual != expected {
		t.Fatalf("expected ReadUTF() to read out the string %#v but got: %#v", expected, actual)
	}
}

func TestWriteEmptyString(t *testing.T) {
	expected := ""
	buf := new(bytes.Buffer)
	w := NewWriter(buf)
	if err := w.WriteUTF(expected); err != nil {
		t.Fatal("expected WriteUTF() to be properly written but got: %w", err)
	}
	r := NewReader(buf)
	actual, err := r.ReadUTF()
	if err != nil {
		t.Fatalf("expected ReadUTF() to read out the string %#v but got error: %s", expected, err)
	}
	if actual != expected {
		t.Fatalf("expected ReadUTF() to read out the string %#v but got: %#v", expected, actual)
	}
}

func TestWriteReadBoolean(t *testing.T) {
	expected := true
	buf := new(bytes.Buffer)
	w := NewWriter(buf)
	if err := w.WriteBoolean(expected); err != nil {
		t.Fatalf("expected WriteBoolean() to be properly written but got: %#v", err)
	}

	r := NewReader(buf)
	actual, err := r.ReadBoolean()
	if err != nil {
		t.Fatalf("expected ReadBoolean() to read out the boolean value %#v but got error: %s", expected, err)
	}
	if actual != expected {
		t.Fatalf("expected ReadBoolean() to read out the boolean value %#v but got: %#v", expected, actual)
	}
}

func TestWriteReadFloat(t *testing.T) {
	var expected float32 = 313.37
	buf := new(bytes.Buffer)
	w := NewWriter(buf)
	if err := w.WriteFloat(expected); err != nil {
		t.Fatalf("expected WriteFloat() to be properly written but got: %#v", err)
	}

	r := NewReader(buf)
	actual, err := r.ReadFloat()
	if err != nil {
		t.Fatalf("expected ReadFloat() to read float value %#v but got error: %s", expected, err)
	}
	if actual != expected {
		t.Fatalf("expected ReadFloat() to read float value %#v but got: %#v", expected, actual)
	}
}

func TestWriteNaNFloat(t *testing.T) {
	var NaN float32 = 0x7fc00000
	expected := []byte{0x4e, 0xff, 0x80, 0x00}
	buf := new(bytes.Buffer)
	w := NewWriter(buf)
	if err := w.WriteFloat(NaN); err != nil {
		t.Fatalf("expected WriteFloat() to be properly written but got: %#v", err)
	}
	if !slices.Equal(expected, buf.Bytes()) {
		t.Fatalf("float NaN was not written correctly, expected %#v but got: %#v", expected, buf.Bytes())
	}
}
