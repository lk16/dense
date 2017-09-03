package bits

import (
	"io"
)

type Reader struct {
	reader    io.Reader
	buff      []byte
	bits_left int
}

// Creates a new Reader
func NewReader(reader io.Reader) *Reader {
	return &Reader{
		reader:    reader,
		buff:      make([]byte, 1),
		bits_left: 0}
}

// Reads a bit
func (reader *Reader) ReadBit() (bit bool, err error) {
	if reader.bits_left == 0 {
		if err = reader.doRead(); err != nil {
			return
		}
	}
	reader.bits_left--
	bit = (reader.buff[0]&0x80 == 0x80)
	reader.buff[0] <<= 1
	return
}

// Count number of unflushed bits since the last read byte
func (reader *Reader) CountUnflushedBits() int {
	return reader.bits_left
}

// Read padding bits and discard them
func (reader *Reader) FlushBits() {
	reader.bits_left = 0
}

func (reader *Reader) doRead() (err error) {
	_, err = reader.reader.Read(reader.buff)
	if err != nil {
		return
	}
	reader.bits_left = 8
	return
}
