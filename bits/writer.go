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

// Count number of padding bits added if WritePaddingBits were called
func (writer *Writer) CountPaddingBits() (count int) {
	if writer.slice.length == 0 {
		count = 0
		return
	}

	count = 8 - writer.slice.length
	return
}

// Writes padding bits
func (writer *Writer) WritePaddingBits() (err error) {
	writer.slice.AppendPadding()
	return writer.doWrite()
}

func (writer *Writer) doWrite() (err error) {
	bytes := writer.slice.PopLeadingBytes()
	_, err = writer.writer.Write(bytes)
	return
}
