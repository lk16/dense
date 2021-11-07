package huffman

import (
    "io"
    "bytes"
    "dense/bits"
    "encoding/binary"
)


func Encode(reader io.Reader, writer io.Writer) (err error) {

	var freq_tree_reader, encode_reader bytes.Buffer

	multi_writer := io.MultiWriter(&freq_tree_reader, &encode_reader)
	io.Copy(multi_writer, reader)

	tree, err := generateTree(&freq_tree_reader)

	if err != nil {
		return
	}

	if err = tree.encodeShape(writer); err != nil {
		return
	}

	if err = tree.encodeShape(writer); err != nil {
		return
	}

	table := tree.getEncodingTable()
	err = tree.encodeBody(&encode_reader, writer, table)
	return
}

func writeBlock(writer io.Writer, block_id byte, block_data bytes.Buffer) (n int, err error) {

	block_len := make([]byte, 8)
	binary.LittleEndian.PutUint64(block_len, uint64(block_data.Len()))

	block := [][]byte{
		[]byte{block_id},
		block_len,
		block_data.Bytes()}

	for _, bytes := range block {
		nn, err := writer.Write(bytes)
		n += nn
		if err != nil {
			return n, err
		}
	}

	return
}

func (node *Tree) encodeTreeLeaves(writer io.Writer) (err error) {

	var leaves_buff bytes.Buffer
	node.encodeTreeLeavesRecursive(&leaves_buff)

	_, err = writeBlock(writer, BLOCK_ID_LEAVES, leaves_buff)
	return
}

func (node *Tree) encodeTreeLeavesRecursive(buff *bytes.Buffer) {
	if node.left != nil {
		node.left.encodeTreeLeavesRecursive(buff)
		node.right.encodeTreeLeavesRecursive(buff)
		return
	}

    data, ok := node.data.(byte)
    if ! ok {
        panic("node.data casting failed")
    }

	buff.WriteByte(data)
}

func (node *Tree) encodeBody(reader io.Reader, writer io.Writer, table map[byte]bits.Slice) (err error) {

	var body_buff, body_data_buff bytes.Buffer
	bits_writer := bits.NewWriter(&body_data_buff)

	key_buff := make([]byte, 1)

	for {

		if _, err = reader.Read(key_buff); err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		slice := table[key_buff[0]]
		bits_writer.WriteSlice(&slice)
	}

	padding_bits_count := byte(bits_writer.CountPaddingBits())
	bits_writer.WritePaddingBits()

	body_buff.WriteByte(padding_bits_count)
	body_buff.Write(body_data_buff.Bytes())

	_, err = writeBlock(writer, BLOCK_ID_DATA, body_buff)
	return
}
