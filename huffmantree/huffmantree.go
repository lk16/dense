package huffmantree

import (
	"bytes"
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
	node.toShapeRecursive(&shape_buff)
	len_buff := make([]byte, 2)
	binary.LittleEndian.PutUint16(len_buff, uint16(shape_buff.Len()))
	buff.Write(len_buff)
	buff.Write(shape_buff.Bytes())
	return
}

func (node *HuffmanTree) toShapeRecursive(buff *bytes.Buffer) {
	if node.left == nil {
		buff.WriteByte(byte('0'))
		return
	}

	buff.WriteByte(byte('1'))
	node.left.toShapeRecursive(buff)
	node.right.toShapeRecursive(buff)
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
