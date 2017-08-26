package huffmantree

import (
	"bytes"
	"dense/bits"
	"encoding/binary"
	"fmt"
	"io"
	"sort"
)

type HuffmanTree struct {
	bytes  []byte
	weight int
	left   *HuffmanTree
	right  *HuffmanTree
}

func NewHuffmanTree(reader io.Reader, max_group_len int) (huffman_tree *HuffmanTree) {

	freq_tree := NewFrequencyTree(reader, max_group_len)
	freq_tree.Prune()
	huffman_tree = fromSlice(freq_tree.ToHuffmanTreeSlice())
	return
}

func fromSlice(slice []HuffmanTree) (tree *HuffmanTree) {
	for len(slice) > 1 {
		sort.Slice(slice, func(i, j int) bool { return slice[i].weight > slice[j].weight })

		left := new(HuffmanTree)
		right := new(HuffmanTree)

		*left = slice[len(slice)-1]
		*right = slice[len(slice)-2]

		slice = slice[:len(slice)-1]
		slice[len(slice)-1] = HuffmanTree{
			bytes:  []byte{},
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

func (node *HuffmanTree) ToShapeBuff() (buff bytes.Buffer) {
	var shape_buff bytes.Buffer
	shape_buff_writer := bits.NewWriter(&shape_buff)
	bits_written, err := node.toShapeRecursive(shape_buff_writer)
	if err != nil {
		panic(err)
	}

	len_buff := make([]byte, 2)
	binary.LittleEndian.PutUint16(len_buff, uint16(bits_written))
	buff.Write(len_buff)
	buff.Write(shape_buff.Bytes())
	return
}

func (node *HuffmanTree) toShapeRecursive(bits_writer *bits.Writer) (bits_written int64, err error) {
	if node.left == nil {
		bits_written, err = bits_writer.WriteBit(false)
		return
	}

	bits_written, err = bits_writer.WriteBit(true)
	if err != nil {
		return
	}

	var written int64

	written, err = node.left.toShapeRecursive(bits_writer)
	bits_written += written
	if err != nil {
		return
	}

	written, err = node.right.toShapeRecursive(bits_writer)
	bits_written += written
	return

}

func (node *HuffmanTree) ToValueBuff() (buff bytes.Buffer) {
	node.toValueRecursive(&buff)
	return
}

func (node *HuffmanTree) toValueRecursive(buff *bytes.Buffer) {
	if len(node.bytes) == 0 {
		node.left.toValueRecursive(buff)
		node.right.toValueRecursive(buff)
		return
	}

	if len(node.bytes) == 1 && node.bytes[0] != 'x' {
		buff.WriteByte(node.bytes[0])
	} else {
		buff.WriteByte('x')
		buff.WriteByte('0' + byte(len(node.bytes)))
		buff.Write(node.bytes)
		return
	}

	return
}

func (node *HuffmanTree) getEncodingTableRecursive(table *map[string]bits.Slice, slice bits.Slice) {
	if len(node.bytes) == 0 {
		left := slice
		left.AppendBit(false)
		node.left.getEncodingTableRecursive(table, left)

		right := slice
		right.AppendBit(true)
		node.right.getEncodingTableRecursive(table, right)

		return
	}

	(*table)[string(node.bytes)] = slice

}

func (node *HuffmanTree) GetEncodingTable() (table map[string]bits.Slice) {
	table = make(map[string]bits.Slice, 20)
	node.getEncodingTableRecursive(&table, *bits.NewSlice(0, 0x0))
	return
}

func (node *HuffmanTree) Encode(reader io.Reader, writer io.Writer, max_group_len int) (written_bits int64, write_err error) {
	bits_writer := bits.NewWriter(writer)
	table := node.GetEncodingTable()
	buff := make([]byte, 1)

	var read_err error
	var written int64

	for read_err == nil {
		key := []byte{}
		for {
			_, read_err = reader.Read(buff)
			if read_err != nil {
				break
			}
			key = append(key, buff[0])
			if _, ok := table[string(key)]; !ok {
				key = key[:len(key)-1]
				break
			}
		}
		slice := table[string(key)]
		written, write_err = bits_writer.WriteSlice(&slice)
		written_bits += written
		if write_err != nil {
			return
		}
	}

	written, write_err = bits_writer.FlushRemainingBits()
	written_bits += written
	return
}
