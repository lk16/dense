package bits

import (
	"bytes"
	"testing"
)

func TestNewBitsWriter(t *testing.T) {
	var buff bytes.Buffer
	bw := NewBitsWriter(&buff)

	if bw.writer != &buff || bw.slice != *NewSlice(0, 0x0) {
		t.Errorf("NewBitsWriter failed")
	}
}

func TestBitsWriterWriteBit(t *testing.T) {
	var buff bytes.Buffer
	bw := NewBitsWriter(&buff)

	for i := 0; i < 7; i++ {
		bits_written, err := bw.WriteBit(true)
		if bits_written != 0 {
			t.Errorf("Expected 0, got %d", bits_written)
		}
		if err != nil {
			t.Errorf("Got error '%s'", err)
		}
	}

	bits_written, err := bw.WriteBit(true)
	if bits_written != 8 {
		t.Errorf("Expected 8, got %d", bits_written)
	}
	if err != nil {
		t.Errorf("Got error '%s'", err)
	}

	bytes := buff.Bytes()
	if len(bytes) != 1 || bytes[0] != 0xFF {
		t.Errorf("Expected [0xFF], got %v", bytes)
	}

}

func TestBitsWriterWriteSlice(t *testing.T) {
	var buff bytes.Buffer
	bw := NewBitsWriter(&buff)

	slice := *NewSlice(4, 0xF)

	bits_written, err := bw.WriteSlice(&slice)
	if bits_written != 0 {
		t.Errorf("Expected 0, got %d", bits_written)
	}
	if err != nil {
		t.Errorf("Got error '%s'", err)
	}

	bits_written, err = bw.WriteSlice(&slice)
	if bits_written != 8 {
		t.Errorf("Expected 0, got %d", bits_written)
	}
	if err != nil {
		t.Errorf("Got error '%s'", err)
	}

	bytes := buff.Bytes()
	if len(bytes) != 1 || bytes[0] != 0xFF {
		t.Errorf("Expected [0xFF], got %v", bytes)
	}

}

func TestBitsWriterFlushRemainingBits(t *testing.T) {
	var buff bytes.Buffer
	bw := NewBitsWriter(&buff)

	slice := *NewSlice(4, 0xF)

	bits_written, err := bw.WriteSlice(&slice)
	if bits_written != 0 {
		t.Errorf("Expected 0, got %d", bits_written)
	}
	if err != nil {
		t.Errorf("Got error '%s'", err)
	}

	bw.FlushRemainingBits()

	bytes := buff.Bytes()
	if len(bytes) != 1 || bytes[0] != 0xF0 {
		t.Errorf("Expected [0xFF], got %v", bytes)
	}

}
