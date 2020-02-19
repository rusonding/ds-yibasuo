package utils

import (
	"bytes"
	"testing"
)

func TestByte2String(t *testing.T) {
	b := []byte("ds-yibasuo")
	a := Byte2String(b)

	if a != "ds-yibasuo" {
		t.Fatal(a)
	}

	b[1] = 'p'

	if a != "dp-yibasuo" {
		t.Errorf("byte to string error")
	}
}

func TestString2Byte(t *testing.T) {
	a := "ds-yibasuo"

	b := String2Byte(a)

	if !bytes.Equal(b, []byte("ds-yibasuo")) {
		t.Errorf("string to byte error")
	}
}
