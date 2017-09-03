package bits

import (
	"io"
)

type Writer struct {
	writer io.Writer
	slice  Slice
}

// Creates a new writer
func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer: writer,
		slice:  *NewSlice(0, 0x0)}
}

// Writes a bit
func (writer *Writer) WriteBit(bit bool) error {
	writer.slice.AppendBit(bit)
	return writer.doWrite()
}

// Writes a slice of bits
func (writer *Writer) WriteSlice(slice *Slice) error {
	writer.slice.AppendSlice(slice)
	return writer.doWrite()
}

// Count number of unflushed bits since the last written byte
func (writer *Writer) CountUnflushedBits() (count int) {
	count = writer.slice.length
	return
}

// Pad and write last bits
func (writer *Writer) FlushBits() (err error) {
	writer.slice.AppendPadding()
	return writer.doWrite()
}

func (writer *Writer) doWrite() (err error) {
	bytes := writer.slice.PopLeadingBytes()
	_, err = writer.writer.Write(bytes)
	return
}
