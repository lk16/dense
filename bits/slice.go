package bits

type Slice struct {
	length int
	data   uint64
}

// Creates a new slice of Bits
func NewSlice(length int, data uint64) (slice *Slice) {
	if length > 64 {
		panic("bits.Slice too big")
	}
	return &Slice{
		length: length,
		data:   data}
}

// Appends a bit slice
func (slice *Slice) AppendSlice(rhs *Slice) {
	new_length := slice.length + rhs.length
	if new_length > 64 {
		panic("bits.Slice too big")
	}
	slice.data = (slice.data << uint(rhs.length)) | rhs.data
	slice.length = new_length
}

// Appends a bit
func (slice *Slice) AppendBit(bit bool) {
	new_length := slice.length + 1
	if new_length > 64 {
		panic("bits.Slice too big")
	}

	slice.data <<= 1
	if bit {
		slice.data |= 0x1
	}

	slice.length = new_length
}

// Appends padding to least significant bits.
// This ensures the slice length is a multiple of 8.
func (slice *Slice) AppendPadding() {
	if slice.length%8 == 0 {
		return
	}
	bit_padding := 8 - (slice.length % 8)

	slice.data <<= uint(bit_padding)
	slice.length += bit_padding
}

// Returns leading 8 bits as a byte
// Ok indicates whether there are 8 bytes to be returned.
func (slice *Slice) PopLeadingBytes() (bytes []byte) {

	bytes = []byte{}

	for slice.length >= 8 {
		b := byte(slice.data >> uint(slice.length-8))
		bytes = append(bytes, b)
		slice.length -= 8
	}
	return
}
