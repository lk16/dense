package huffman

import (
	"bytes"
	"dense/bits"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sort"
)

const (
	BLOCK_ID_SHAPE  = 0
	BLOCK_ID_LEAVES = 1
	BLOCK_ID_DATA   = 2
)

func Encode(reader io.Reader, writer io.Writer) (err error) {

	var freq_tree_reader, encode_reader bytes.Buffer

	multi_writer := io.MultiWriter(&freq_tree_reader, &encode_reader)
	io.Copy(multi_writer, reader)

	tree, err := generateTree(&freq_tree_reader)

	if err != nil {
		return
	}

	if err = tree.encodeTreeShape(writer); err != nil {
		return
	}

	if err = tree.encodeTreeLeaves(writer); err != nil {
		return
	}

	table := tree.getEncodingTable()
	err = tree.encodeBody(&encode_reader, writer, table)
	return
}

type decodeState struct {
	tree   *HuffmanTree
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
				break
			}
			return
		}
	}

	return

	/*tree, err := decodeTreeShape(reader)
	if err != nil {
		return err
	}

	err = tree.decodeTreeLeaves(reader)
	if err != nil {
		return err
	}

	err = tree.decodeBody(reader, writer)
	return*/
}

type HuffmanTree struct {
	data   byte
	weight int64
	left   *HuffmanTree
	right  *HuffmanTree
}

func generateTree(reader io.Reader) (tree *HuffmanTree, err error) {

	buff := make([]byte, 4096)
	table := make(map[byte]int64, 256)

	for {
		var read_bytes int
		read_bytes, err = reader.Read(buff)
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}

		for _, b := range buff[:read_bytes] {
			if _, ok := table[b]; !ok {
				table[b] = 0
			}
			table[b]++
		}
	}

	slice := make([]HuffmanTree, 0)

	for key, value := range table {

		if value == 0 {
			continue
		}

		slice = append(slice, HuffmanTree{
			data:   key,
			weight: value,
			left:   nil,
			right:  nil})
	}

	if len(slice) == 0 {
		// No input data. Put a useless root in the tree
		tree = &HuffmanTree{}
		return
	}
	if len(slice) == 1 {
		// Only one byte value in input. Put dummy node as sibling.
		tree = &HuffmanTree{
			left:   &slice[0],
			right:  &HuffmanTree{},
			weight: slice[0].weight}
		return
	}

	for len(slice) > 1 {
		sort.Slice(slice, func(i, j int) bool { return slice[i].weight > slice[j].weight })

		left := new(HuffmanTree)
		right := new(HuffmanTree)

		*left = slice[len(slice)-1]
		*right = slice[len(slice)-2]

		slice = slice[:len(slice)-1]
		slice[len(slice)-1] = HuffmanTree{
			data:   byte(0),
			weight: left.weight + right.weight,
			left:   left,
			right:  right}

	}
	tree = new(HuffmanTree)

	*tree = slice[0]
	return
}

func (state *decodeState) decodeTreeShape(buff bytes.Buffer) (err error) {

	bits_reader := bits.NewReader(&buff)
	var decodeShape func() (*HuffmanTree, error)

	decodeShape = func() (node *HuffmanTree, err error) {

		bit, err := bits_reader.ReadBit()

		if err != nil {
			return nil, err
		}

		if !bit {
			return nil, nil
		}

		node = &HuffmanTree{}
		if node.left, err = decodeShape(); err != nil {
			return nil, err
		}

		node.right, err = decodeShape()
		return
	}

	state.tree, err = decodeShape()
	fmt.Printf("%v\n", state.tree.String())
	return
}

func (node *HuffmanTree) encodeTreeShape(writer io.Writer) (err error) {

	var shape_buff bytes.Buffer
	shape_buff_writer := bits.NewWriter(&shape_buff)

	node.encodeTreeShapeRecursive(shape_buff_writer)
	shape_buff_writer.WritePaddingBits()

	_, err = writeBlock(writer, BLOCK_ID_SHAPE, shape_buff)
	return
}

func (state *decodeState) decodeTreeLeaves(buff bytes.Buffer) (err error) {

	var assign_leaves func(*HuffmanTree, *int)
	bytes := buff.Bytes()

	assign_leaves = func(tree *HuffmanTree, index *int) {
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

func (node *HuffmanTree) encodeTreeShapeRecursive(bits_writer *bits.Writer) {

	if node.left == nil {
		bits_writer.WriteBit(false)
		return
	}

	bits_writer.WriteBit(true)
	node.left.encodeTreeShapeRecursive(bits_writer)
	node.right.encodeTreeShapeRecursive(bits_writer)
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

func (node *HuffmanTree) encodeTreeLeaves(writer io.Writer) (err error) {

	var leaves_buff bytes.Buffer
	node.encodeTreeLeavesRecursive(&leaves_buff)

	_, err = writeBlock(writer, BLOCK_ID_LEAVES, leaves_buff)
	return
}

func (node *HuffmanTree) encodeTreeLeavesRecursive(buff *bytes.Buffer) {
	if node.left != nil {
		node.left.encodeTreeLeavesRecursive(buff)
		node.right.encodeTreeLeavesRecursive(buff)
		return
	}

	buff.WriteByte(node.data)
}

func (node *HuffmanTree) getEncodingTableRecursive(table *map[byte]bits.Slice, slice bits.Slice) {
	if node.left != nil {
		left := slice
		left.AppendBit(false)
		node.left.getEncodingTableRecursive(table, left)

		right := slice
		right.AppendBit(true)
		node.right.getEncodingTableRecursive(table, right)
		return
	}

	(*table)[node.data] = slice
}

func (node *HuffmanTree) getEncodingTable() (table map[byte]bits.Slice) {
	table = make(map[byte]bits.Slice, 20)
	node.getEncodingTableRecursive(&table, *bits.NewSlice(0, 0x0))
	return
}

func (node *HuffmanTree) encodeBody(reader io.Reader, writer io.Writer, table map[byte]bits.Slice) (err error) {

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

func (state *decodeState) decodeBody(buff bytes.Buffer) (err error) {

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
			state.writer.Write([]byte{node.data})
			node = state.tree
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
	data_len := binary.LittleEndian.Uint64(block_len_buff)

	var block_data bytes.Buffer
	_, err = io.CopyN(&block_data, reader, int64(data_len))
	if err != nil {
		return n, err
	}
	n += int(data_len)

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

func (tree *HuffmanTree) String() (str string) {

	var toString func(*HuffmanTree, string) string

	toString = func(node *HuffmanTree, prefix string) (str string) {
		if node == nil {
			return ""
		}

		str = prefix + fmt.Sprintf(" %d (%d)\n", node.weight, node.data)
		str += toString(node.left, prefix+"  ")
		str += toString(node.right, prefix+"  ")
		return
	}

	str = toString(tree, "")
	return
}
