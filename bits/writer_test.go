package bits

import (
	"bytes"
	"testing"
)

func TestNewBitsWriter(t *testing.T) {
	var buff bytes.Buffer
	bw := NewWriter(&buff)

	if bw.writer != &buff || bw.slice != *NewSlice(0, 0x0) {
		t.Errorf("NewBitsWriter failed")
	}
}

func TestBitsWriterWriteBit(t *testing.T) {
	var buff bytes.Buffer
	bw := NewWriter(&buff)

	// writes 1111 0000 = 0xF0
	for i := 0; i < 8; i++ {
		bw.WriteBit(i < 4)
	}

	bytes := buff.Bytes()
	if len(bytes) != 1 || bytes[0] != 0xF0 {
		t.Errorf("Expected [0xF0], got %v", bytes)
	}

}

func TestBitsWriterWriteSlice(t *testing.T) {
	var buff bytes.Buffer
	bw := NewWriter(&buff)

	slice := *NewSlice(4, 0xF)

	bw.WriteSlice(&slice)
	bw.WriteSlice(&slice)

	bytes := buff.Bytes()
	if len(bytes) != 1 || bytes[0] != 0xFF {
		t.Errorf("Expected [0xFF], got %v", bytes)
	}

}

func TestBitsWriterFlushBits(t *testing.T) {
	var buff bytes.Buffer
	bw := NewWriter(&buff)

	slice := *NewSlice(4, 0xF)

	bw.WriteSlice(&slice)
	bw.FlushBits()

	bytes := buff.Bytes()
	if len(bytes) != 1 || bytes[0] != 0xF0 {
		t.Errorf("Expected [0xFF], got %v", bytes)
	}

}

func BitsWriterCountUnflushedBits(t *testing.T) {

	var buff bytes.Buffer
	bw := NewWriter(&buff)

	slice := *NewSlice(4, 0xF)

	if bw.CountUnflushedBits() != 0 {
		t.Errorf("Expected 0, got %d", bw.CountUnflushedBits())
	}

	bw.WriteBit(true)

	if bw.CountUnflushedBits() != 1 {
		t.Errorf("Expected 1, got %d", bw.CountUnflushedBits())
	}

	bw.WriteSlice(&slice)

	if bw.CountUnflushedBits() != 5 {
		t.Errorf("Expected 1, got %d", bw.CountUnflushedBits())
	}

	bw.WriteSlice(&slice)

	if bw.CountUnflushedBits() != 1 {
		t.Errorf("Expected 1, got %d", bw.CountUnflushedBits())
	}

	bw.FlushBits()

	if bw.CountUnflushedBits() != 0 {
		t.Errorf("Expected 0, got %d", bw.CountUnflushedBits())
	}
}
