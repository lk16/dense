package huffman

import (
    "io"
    "bytes"
    "encoding/binary"
    "dense/bits"
    "errors"
    "fmt"
)

type decodeState struct {
	tree   *Tree
	writer io.Writer
}

func Decode(reader io.Reader, writer io.Writer) (err error) {

	state := decodeState{
		tree:   nil,
		writer: writer}

	for {
		_, err = state.readBlock(reader)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
            break
		}
	}

	return
}


func (state *decodeState) readBlock(reader io.Reader) (n int, err error) {

	block_id_buff := make([]byte, 1)
	block_len_buff := make([]byte, 8)
	var nn int

	if nn, err = reader.Read(block_id_buff); err != nil {
		return n, err
	}
	n += nn

	if nn, err = reader.Read(block_len_buff); err != nil {
		return
	}
	n += nn

	block_id := block_id_buff[0]
	block_len := binary.LittleEndian.Uint64(block_len_buff)
    fmt.Printf("Found block_id %d, length %d\n",block_id,block_len)

	var block_data bytes.Buffer
	_, err = io.CopyN(&block_data, reader, int64(block_len))
	if err != nil {
		return n, err
	}
	n += int(block_len)

	block_decoders := map[byte]func(bytes.Buffer) error{
		BLOCK_ID_DATA:   state.decodeBody,
		BLOCK_ID_LEAVES: state.decodeTreeLeaves,
		BLOCK_ID_SHAPE:  state.decodeTreeShape}

	if block_decoder, ok := block_decoders[block_id]; ok {
		err = block_decoder(block_data)
	} else {
		err = errors.New("Unexpected block ID")
	}
	return
}


func (state *decodeState) decodeTreeShape(buff bytes.Buffer) (err error) {
	err = state.tree.decodeShape(&buff)
	return
}

func (state *decodeState) decodeTreeLeaves(buff bytes.Buffer) (err error) {

    if state.tree == nil {
        panic("decodeTreeLeaves() called without having a tree")
    }

	var assign_leaves func(*Tree, *int)
	bytes := buff.Bytes()

	assign_leaves = func(tree *Tree, index *int) {
		if tree.left == nil {
			tree.data = bytes[*index]
			*index++
			return
		}

		assign_leaves(tree.left, index)
		assign_leaves(tree.right, index)
	}

	index := 0
	assign_leaves(state.tree, &index)

	return
}

func (state *decodeState) decodeBody(buff bytes.Buffer) (err error) {

    if state.tree == nil {
        panic("decodeBody() called without having a tree")
    }

	trailing_bit_count_buff := make([]byte, 1)
	if _, err = buff.Read(trailing_bit_count_buff); err != nil {
		return
	}

	padding_bits_count := trailing_bit_count_buff[0]

	bits_left := (8 * buff.Len()) - int(padding_bits_count)

	bit_reader := bits.NewReader(&buff)

	node := state.tree

	for bits_left != 0 {

		bits_left--
		bit, err := bit_reader.ReadBit()

		if err != nil {
			return err
		}

		if bit {
			node = node.right
		} else {
			node = node.left
		}

		if node.left == nil {
            data, ok := node.data.(byte)
            if !ok {
                panic("node.data is not a byte")
            }
			state.writer.Write([]byte{data})
			node = state.tree
		}
	}

	return
}
