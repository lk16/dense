package huffman

import (
	"bytes"
	"dense/bits"
	"encoding/binary"
	"errors"
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

func Decode(reader io.Reader, writer io.Writer) (err error) {

	tree, err := decodeTreeShape(reader)
	if err != nil {
		return err
	}

	err = tree.decodeTreeLeaves(reader)
	if err != nil {
		return err
	}

	err = tree.decodeBody(reader, writer)
	return
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

func decodeTreeShape(reader io.Reader) (tree *HuffmanTree, err error) {
	block_id_buff := make([]byte, 1)

	if _, err = reader.Read(block_id_buff); err != nil {
		return
	}

	if block_id_buff[0] != BLOCK_ID_SHAPE {
		err = errors.New("Unexpected block ID")
		return
	}

	len_buff := make([]byte, 8)
	if _, err = reader.Read(len_buff); err != nil {
		return
	}
	shape_buff_len := binary.LittleEndian.Uint64(len_buff)

	var shape_buff bytes.Buffer

	_, err = io.CopyN(&shape_buff, reader, int64(shape_buff_len))
	if err != nil {
		return
	}

	tree = &HuffmanTree{}

	stack := []**HuffmanTree{
		&tree}

	bits_reader := bits.NewReader(&shape_buff)

	for len(stack) > 0 {

		visiting_node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		bit, err := bits_reader.ReadBit()

		if err != nil {
			return tree, err
		}

		if bit {
			(*visiting_node).left = &HuffmanTree{}
			(*visiting_node).right = &HuffmanTree{}
			stack = append(stack, &((*visiting_node).right))
			stack = append(stack, &((*visiting_node).left))
		}
	}
	return
}

func (node *HuffmanTree) encodeTreeShape(writer io.Writer) (err error) {

	var shape_buff bytes.Buffer
	shape_buff_writer := bits.NewWriter(&shape_buff)

	node.encodeTreeShapeRecursive(shape_buff_writer)
	shape_buff_writer.FlushBits()

	writer.Write([]byte{BLOCK_ID_SHAPE})

	len_buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(len_buff, uint64(shape_buff.Len()))
	if _, err = writer.Write(len_buff); err != nil {
		return
	}

	_, err = writer.Write(shape_buff.Bytes())
	return
}

func (tree *HuffmanTree) decodeTreeLeaves(reader io.Reader) (err error) {
	block_id_buff := make([]byte, 1)

	if _, err = reader.Read(block_id_buff); err != nil {
		return
	}

	if block_id_buff[0] != BLOCK_ID_LEAVES {
		err = errors.New("Unexpected block ID")
		return
	}

	len_buff := make([]byte, 8)
	if _, err = reader.Read(len_buff); err != nil {
		return
	}
	leaves_buff_len := binary.LittleEndian.Uint64(len_buff)

	leaves_buff := make([]byte, leaves_buff_len)

	if _, err = reader.Read(leaves_buff); err != nil {
		return
	}

	var assign_leaves func(*HuffmanTree, *int)

	assign_leaves = func(tree *HuffmanTree, index *int) {
		if tree.left == nil {
			tree.data = leaves_buff[*index]
			*index++
			return
		}

		assign_leaves(tree.left, index)
		assign_leaves(tree.right, index)
	}

	index := 0
	assign_leaves(tree, &index)

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

func (node *HuffmanTree) encodeTreeLeaves(writer io.Writer) (err error) {

	var leaves_buff bytes.Buffer
	node.encodeTreeLeavesRecursive(&leaves_buff)

	writer.Write([]byte{BLOCK_ID_LEAVES})

	len_buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(len_buff, uint64(leaves_buff.Len()))
	if _, err = writer.Write(len_buff); err != nil {
		return
	}

	_, err = writer.Write(leaves_buff.Bytes())
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

	var body_buff bytes.Buffer
	bits_writer := bits.NewWriter(&body_buff)

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

	len_buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(len_buff, uint64(body_buff.Len()))

	trailing_bit_count := byte(bits_writer.CountUnflushedBits())
	bits_writer.FlushBits()

	buffers := [][]byte{
		[]byte{BLOCK_ID_DATA},
		len_buff,
		[]byte{trailing_bit_count},
		body_buff.Bytes()}

	for _, buffer := range buffers {
		if _, err = writer.Write(buffer); err != nil {
			return
		}
	}

	return
}

func (tree *HuffmanTree) decodeBody(reader io.Reader, writer io.Writer) (err error) {
	block_id_buff := make([]byte, 1)

	if _, err = reader.Read(block_id_buff); err != nil {
		return
	}

	if block_id_buff[0] != BLOCK_ID_DATA {
		err = errors.New("Unexpected block ID")
		return
	}

	len_buff := make([]byte, 8)
	if _, err = reader.Read(len_buff); err != nil {
		return
	}
	data_len := binary.LittleEndian.Uint64(len_buff)

	trailing_bit_count_buff := make([]byte, 1)
	if _, err = reader.Read(trailing_bit_count_buff); err != nil {
		return
	}

	trailing_bit_count := trailing_bit_count_buff[0]

	bits_left := (8 * data_len) + uint64(trailing_bit_count)

	if trailing_bit_count != 0 {
		data_len++
	}

	var data_buff bytes.Buffer

	io.CopyN(&data_buff, reader, int64(data_len))

	bit_reader := bits.NewReader(&data_buff)

	node := tree

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
			writer.Write([]byte{node.data})
			node = tree
		}
	}

	return
}
