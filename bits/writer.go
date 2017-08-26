package bits

import (
	"io"
)

type Writer struct {
	writer io.Writer
	slice  Slice
}

// Creates a new bitsWriter
func NewWriter(writer io.Writer) (bitsWriter *Writer) {
	return &Writer{
		writer: writer,
		slice:  *NewSlice(0, 0x0)}
}

// Writes a bit
func (bitsWriter *Writer) WriteBit(bit bool) (bits_written int64, err error) {
	bitsWriter.slice.AppendBit(bit)
	bits_written, err = bitsWriter.doWrite()
	return
}

// Writes a slice of bits
func (bitsWriter *Writer) WriteSlice(slice *Slice) (bits_written int64, err error) {
	bitsWriter.slice.AppendSlice(slice)
	bits_written, err = bitsWriter.doWrite()
	return
}

// Pad and write last bits
func (bitsWriter *Writer) FlushRemainingBits() (bits_written int64, err error) {
	bits_left := int64(bitsWriter.slice.length)

	bitsWriter.slice.AppendPadding()

	var bytes_written int64
	bytes_written, err = bitsWriter.doWrite()
	if bytes_written == 0 {
		return 0, err
	}
	return bits_left, err
}

func (bitsWriter *Writer) doWrite() (bits_written int64, err error) {
	var bytes_written int
	for {
		b, ok := bitsWriter.slice.PopLeadingByte()
		if !ok {
			return
		}
		bytes_written, err = bitsWriter.writer.Write([]byte{b})
		bits_written += int64(8 * bytes_written)
		if err != nil {
			return
		}
	}
}
