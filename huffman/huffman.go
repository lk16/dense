package huffman

import (
	"bytes"
	"dense/bits"
	"encoding/binary"
	"fmt"
	"io"
	"sort"
)

func Encode(reader io.Reader, writer io.Writer) (err error) {

	var freq_tree_reader, encode_reader bytes.Buffer

	multi_writer := io.MultiWriter(&freq_tree_reader, &encode_reader)
	io.Copy(multi_writer, reader)

	tree, err := newTreeFromReader(&freq_tree_reader)
	if err != nil {
		return
	}

	if err = tree.encodeTreeShape(writer); err != nil {
		return
	}

	if err = tree.encodeTreeLeaves(writer); err != nil {
		return
	}

	table := tree.GetEncodingTable()
	err = tree.encodeBody(&encode_reader, writer, table)
	return
}

func Decode(reader io.Reader, writer io.Writer) (err error) {
	fmt.Printf("Decompressing is not yet implemented\n")
	return
}

type HuffmanTree struct {
	bytes  byte
	weight int64
	left   *HuffmanTree
	right  *HuffmanTree
}

func newTreeFromReader(reader io.Reader) (tree *HuffmanTree, err error) {

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
			bytes:  key,
			weight: value,
			left:   nil,
			right:  nil})
	}

	for len(slice) > 1 {
		sort.Slice(slice, func(i, j int) bool { return slice[i].weight > slice[j].weight })

		left := new(HuffmanTree)
		right := new(HuffmanTree)

		*left = slice[len(slice)-1]
		*right = slice[len(slice)-2]

		slice = slice[:len(slice)-1]
		slice[len(slice)-1] = HuffmanTree{
			bytes:  byte(0),
			weight: left.weight + right.weight,
			left:   left,
			right:  right}

	}
	tree = new(HuffmanTree)
	*tree = slice[0]
	return
}

func (node *HuffmanTree) print(code string) {
	if node.left == nil {
		fmt.Printf("%d\t'%s'\t%s\n", node.weight, string(node.bytes), code)
		return
	}

	node.left.print(code + "0")
	node.right.print(code + "1")
}

func (node *HuffmanTree) Print() {
	node.print("")
}

func (node *HuffmanTree) encodeTreeShape(writer io.Writer) (err error) {

	var shape_buff bytes.Buffer
	shape_buff_writer := bits.NewWriter(&shape_buff)

	if err = node.encodeTreeShapeRecursive(shape_buff_writer); err != nil {
		return
	}

	if _, err = shape_buff_writer.FlushRemainingBits(); err != nil {
		return
	}

	len_buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(len_buff, uint64(shape_buff.Len()))
	if _, err = writer.Write(len_buff); err != nil {
		return
	}

	_, err = writer.Write(shape_buff.Bytes())
	return
}

func (node *HuffmanTree) encodeTreeShapeRecursive(bits_writer *bits.Writer) (err error) {
	if node.left == nil {
		_, err = bits_writer.WriteBit(false)
		return
	}

	if _, err = bits_writer.WriteBit(true); err != nil {
		return
	}

	if err = node.left.encodeTreeShapeRecursive(bits_writer); err != nil {
		return
	}

	err = node.right.encodeTreeShapeRecursive(bits_writer)
	return

}

func (node *HuffmanTree) encodeTreeLeaves(writer io.Writer) (err error) {

	var leaves_buff bytes.Buffer
	if err = node.encodeTreeLeavesRecursive(&leaves_buff); err != nil {
		return
	}

	len_buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(len_buff, uint64(leaves_buff.Len()))
	if _, err = writer.Write(len_buff); err != nil {
		return
	}

	_, err = writer.Write(leaves_buff.Bytes())
	return
}

func (node *HuffmanTree) encodeTreeLeavesRecursive(buff *bytes.Buffer) (err error) {
	if node.left == nil {
		if err = node.left.encodeTreeLeavesRecursive(buff); err != nil {
			return
		}

		err = node.right.encodeTreeLeavesRecursive(buff)
		return
	}

	err = buff.WriteByte(node.bytes)
	return
}

func (node *HuffmanTree) getEncodingTableRecursive(table *map[byte]bits.Slice, slice bits.Slice) {
	if node.left == nil {
		left := slice
		left.AppendBit(false)
		node.left.getEncodingTableRecursive(table, left)

		right := slice
		right.AppendBit(true)
		node.right.getEncodingTableRecursive(table, right)

		return
	}

	(*table)[node.bytes] = slice

}

func (node *HuffmanTree) GetEncodingTable() (table map[byte]bits.Slice) {
	table = make(map[byte]bits.Slice, 20)
	node.getEncodingTableRecursive(&table, *bits.NewSlice(0, 0x0))
	return
}

func (node *HuffmanTree) encodeBody(reader io.Reader, writer io.Writer, table map[byte]bits.Slice) error {
	bits_writer := bits.NewWriter(writer)
	key_buff := make([]byte, 1)

	for {
		_, read_err := reader.Read(key_buff)
		if read_err != nil {
			if read_err != io.EOF {
				return read_err
			}
			break
		}

		slice := table[key_buff[0]]

		if _, write_err := bits_writer.WriteSlice(&slice); write_err != nil {
			return write_err
		}
	}

	_, write_err := bits_writer.FlushRemainingBits()
	return write_err
}
