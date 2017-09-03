package bits

import (
	"bytes"
	"io"
	"testing"
)

func TestNewReader(t *testing.T) {
	var buff bytes.Buffer
	reader := *NewReader(&buff)

	if reader.reader != &buff || len(reader.buff) != 1 || reader.bits_left != 0 {
		t.Errorf("Wrong NewReader() values: %v", reader)
	}
}

func TestReaderReadBit(t *testing.T) {
	var buff bytes.Buffer
	reader := NewReader(&buff)

	_, err := reader.ReadBit()
	if err != io.EOF {
		t.Errorf("Expected 'EOF', got %v", err)
	}

	// 1010 0001
	buff.WriteByte(0xA1)

	expectedOutput := []bool{true, false, true, false, false, false, false, true}

	for index, expected := range expectedOutput {
		got, _ := reader.ReadBit()

		if got != expected {
			t.Errorf("At index %d: expected %t, got %t", index, got, expected)
		}
	}

	_, err = reader.ReadBit()
	if err != io.EOF {
		t.Errorf("Expected 'EOF', got %v", err)
	}
}

func TestReaderCountUnflushedBits(t *testing.T) {
	var buff bytes.Buffer
	reader := NewReader(&buff)

	// 1010 0001
	buff.WriteByte(0xA1)

	if reader.CountUnflushedBits() != 0 {
		t.Errorf("Expected 0, got %d", reader.CountUnflushedBits())
	}

	expectedOutput := 7
	for {
		if _, err := reader.ReadBit(); err != nil {
			if err != io.EOF {
				t.Errorf("Got unexpected error %s", err)
			}
			break
		}
		output := reader.CountUnflushedBits()
		if output != expectedOutput {
			t.Errorf("Got %d, expected %d", output, expectedOutput)
		}
		expectedOutput--
	}
}

func TestReaderFlushBits(t *testing.T) {
	var buff bytes.Buffer
	reader := NewReader(&buff)

	// 1010 0001
	buff.WriteByte(0xA1)

	reader.ReadBit()

	reader.FlushBits()

	if reader.bits_left != 0 {
		t.Errorf("Expected 0, got %d", reader.bits_left)
	}
}
