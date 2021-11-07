package huffman

import (
	"bytes"
	"dense/bits"
	"fmt"
	"io"
	"sort"
)

type Tree struct {
	weight int64
	left   *Tree
	right  *Tree
	data   interface{}
}

func NewTreeFromSlice(slice []*Tree) (tree *Tree) {
	if len(slice) == 0 {
		// No input data. Put a useless root in the tree
		tree = &Tree{}
		return
	}
	if len(slice) == 1 {
		// Only one byte value in input. Put dummy node as sibling.
		tree = &Tree{
			left:   slice[0],
			right:  &Tree{},
			weight: slice[0].weight}
		return
	}

	for len(slice) > 1 {

		sort.Slice(slice, func(i, j int) bool {
			return slice[i].weight > slice[j].weight
		})

		left := slice[len(slice)-1]
		right := slice[len(slice)-2]

		parent := &Tree{
			weight: left.weight + right.weight,
			left:   left,
			right:  right,
			data:   nil}

		slice = slice[:len(slice)-2]
		slice = append(slice, parent)
	}

	tree = slice[0]
	return
}

func (tree *Tree) encodeShape(writer io.Writer) (err error) {

	var shape_buff bytes.Buffer
	bits_writer := bits.NewWriter(&shape_buff)

	var encode func(*Tree)

	encode = func(tree *Tree) {

		if tree.left == nil {
			bits_writer.WriteBit(false)
			return
		}

		bits_writer.WriteBit(true)
		encode(tree.left)
		encode(tree.right)
	}

	encode(tree)
	bits_writer.WritePaddingBits()

	_, err = writeBlock(writer, BLOCK_ID_SHAPE, shape_buff)
	return
}

func (tree *Tree) decodeShape(reader io.Reader) (err error) {

	bits_reader := bits.NewReader(reader)
	var decode func() (*Tree, error)

	decode = func() (node *Tree, err error) {

		bit, err := bits_reader.ReadBit()

		if err != nil {
			return nil, err
		}

		if !bit {
			return nil, nil
		}

		node = &Tree{}
		if node.left, err = decode(); err != nil {
			return nil, err
		}

		node.right, err = decode()
		return
	}

	tree, err = decode()
	return
}

func (tree *Tree) String() (str string) {

	var toString func(*Tree, string) string

	toString = func(node *Tree, prefix string) (str string) {
		if node == nil {
			return ""
		}

		str = prefix + fmt.Sprintf(" %d (%v)\n", node.weight, node.data)
		str += toString(node.left, prefix+"  ")
		str += toString(node.right, prefix+"  ")
		return
	}

	str = toString(tree, "")
	return
}
