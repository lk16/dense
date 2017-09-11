package huffman

import (
	"bytes"
	"dense/bits"
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"
)

func TestGenerateTree(t *testing.T) {

	var buff bytes.Buffer

	// simple test
	buff.Write(make([]byte, 100))
	tree, err := generateTree(&buff)

	expected_node := HuffmanTree{
		left: &HuffmanTree{
			weight: 100,
			data:   byte(0x0)},
		right:  &HuffmanTree{},
		weight: 100,
		data:   byte(0)}

	if !reflect.DeepEqual(*tree, expected_node) {
		t.Errorf("generateTree failed. Got tree %v, expected %v",
			*tree, expected_node)
	}
	if err != nil {
		t.Errorf("generateTree failed. Got error %s", err)
	}

	// buffer bigger than 4k, not rounded to 4k boundary
	buff.Write(make([]byte, 20000))
	tree, err = generateTree(&buff)

	expected_node = HuffmanTree{
		left: &HuffmanTree{
			weight: 20000,
			data:   byte(0x0)},
		right:  &HuffmanTree{},
		weight: 20000,
		data:   byte(0)}

	if !reflect.DeepEqual(*tree, expected_node) {

		t.Errorf("generateTree failed. Got tree %v, expected %v",
			*tree, expected_node)
	}
	if err != nil {
		t.Errorf("generateTree failed. Got error %s", err)
	}

	// test huffman tree properties
	byte_slice := []byte{0x1, 0x1, 0x2, 0x3}

	buff.Write(byte_slice)
	tree, err = generateTree(&buff)

	if err != nil {
		t.Errorf("generateTree failed. Got error %s", err)
	}

	if tree.weight != 4 || tree.left.weight != 2 || tree.right.weight != 2 {
		t.Errorf("generateTree failed. Got weights %d %d %d",
			tree.weight, tree.left.weight, tree.right.weight)
	}
}

func TestHufmannTreeEncodeTreeShape(t *testing.T) {

	// root with no children
	tree := HuffmanTree{} // 0

	len_buff := make([]byte, 8)
	var expected_output bytes.Buffer

	expected_output.WriteByte(BLOCK_ID_SHAPE)

	binary.LittleEndian.PutUint64(len_buff, 1)
	expected_output.Write(len_buff)
	// 0 for tree
	// 0000000 for padding
	expected_output.Write([]byte{0x00})

	var output bytes.Buffer
	err := tree.encodeTreeShape(&output)

	if err != nil {
		t.Errorf("encodeTreeShape failed. Got error %s", err)
	}

	if !bytes.Equal(expected_output.Bytes(), output.Bytes()) {
		t.Errorf("Expected %v, got %v", expected_output.Bytes(), output.Bytes())
	}

	output.Reset()
	expected_output.Reset()

	// root with 2 children, no grandchildren
	tree = HuffmanTree{ // 1
		left:  &HuffmanTree{}, // 0
		right: &HuffmanTree{}} // 0

	expected_output.WriteByte(BLOCK_ID_SHAPE)

	binary.LittleEndian.PutUint64(len_buff, 1)
	expected_output.Write(len_buff)
	// 100 for tree
	// 00000 for padding
	expected_output.Write([]byte{0x80})

	err = tree.encodeTreeShape(&output)

	if err != nil {
		t.Errorf("encodeTreeShape failed. Got error %s", err)
	}

	if !bytes.Equal(expected_output.Bytes(), output.Bytes()) {
		t.Errorf("Expected %v, got %v", expected_output.Bytes(), output.Bytes())
	}

	output.Reset()
	expected_output.Reset()

	// complicated tree
	tree = HuffmanTree{ // 1
		left: &HuffmanTree{ // 1
			left: &HuffmanTree{ // 1
				left:  &HuffmanTree{},  // 0
				right: &HuffmanTree{}}, // 0
			right: &HuffmanTree{ // 1
				left:  &HuffmanTree{},   // 0
				right: &HuffmanTree{}}}, // 0
		right: &HuffmanTree{}} // 0

	expected_output.WriteByte(BLOCK_ID_SHAPE)

	binary.LittleEndian.PutUint64(len_buff, 2)
	expected_output.Write(len_buff)
	// 111001000 for tree
	// 0000000 for padding
	expected_output.Write([]byte{0xE4, 0x00})

	err = tree.encodeTreeShape(&output)

	if err != nil {
		t.Errorf("encodeTreeShape failed. Got error %s", err)
	}

	if !bytes.Equal(expected_output.Bytes(), output.Bytes()) {
		t.Errorf("Expected %v, got %v", expected_output.Bytes(), output.Bytes())
	}

}

func TestHuffmanTreeEncodeLeaves(t *testing.T) {

	tree := HuffmanTree{
		data: 'A'}

	var output bytes.Buffer
	err := tree.encodeTreeLeaves(&output)

	var expected_output bytes.Buffer

	expected_output.WriteByte(BLOCK_ID_LEAVES)

	len_buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(len_buff, uint64(1))
	expected_output.Write(len_buff)

	expected_output.Write([]byte{'A'})

	if err != nil {
		t.Errorf("encodeTreeLeaves failed. Got error %s", err)
	}

	if !bytes.Equal(expected_output.Bytes(), output.Bytes()) {
		t.Errorf("Expected %v, got %v", expected_output.Bytes(), output.Bytes())
	}

	expected_output.Reset()
	output.Reset()

	tree = HuffmanTree{
		left: &HuffmanTree{
			left: &HuffmanTree{
				left: &HuffmanTree{
					data: 0x00},
				right: &HuffmanTree{
					data: 0xC0}},
			right: &HuffmanTree{
				left: &HuffmanTree{
					data: 0xFF},
				right: &HuffmanTree{
					data: 0xEE}}},
		right: &HuffmanTree{
			data: 0x01}}

	err = tree.encodeTreeLeaves(&output)

	expected_output.WriteByte(BLOCK_ID_LEAVES)

	binary.LittleEndian.PutUint64(len_buff, uint64(5))
	expected_output.Write(len_buff)

	expected_output.Write([]byte{0x00, 0xC0, 0xFF, 0xEE, 0x01})

	if err != nil {
		t.Errorf("encodeTreeLeaves failed. Got error %s", err)
	}

	if !bytes.Equal(expected_output.Bytes(), output.Bytes()) {
		t.Errorf("Expected %v, got %v", expected_output.Bytes(), output.Bytes())
	}
}

func TestHuffmanTreeGetEncodingTable(t *testing.T) {

	tree := HuffmanTree{
		left: &HuffmanTree{
			left: &HuffmanTree{
				left: &HuffmanTree{
					data: 0x00},
				right: &HuffmanTree{
					data: 0xC0}},
			right: &HuffmanTree{
				left: &HuffmanTree{
					data: 0xFF},
				right: &HuffmanTree{
					data: 0xEE}}},
		right: &HuffmanTree{
			data: 0x01}}

	table := tree.getEncodingTable()

	if len(table) != 5 {
		t.Errorf("getEncodingTable failed. Got table size %d", len(table))
	}

	expected_table := map[byte]bits.Slice{
		0x00: *bits.NewSlice(3, 0x0), // 000
		0xC0: *bits.NewSlice(3, 0x1), // 001
		0xFF: *bits.NewSlice(3, 0x2), // 010
		0xEE: *bits.NewSlice(3, 0x3), // 011
		0x01: *bits.NewSlice(1, 0x1)} // 1

	if !reflect.DeepEqual(table, expected_table) {
		t.Errorf("Expected %v, got %v", expected_table, table)
	}
}

func TestHuffmanTreeEncodeBody(t *testing.T) {

	tree := HuffmanTree{
		left: &HuffmanTree{
			left: &HuffmanTree{
				left: &HuffmanTree{
					data: 0x00},
				right: &HuffmanTree{
					data: 0xC0}},
			right: &HuffmanTree{
				left: &HuffmanTree{
					data: 0xFF},
				right: &HuffmanTree{
					data: 0xEE}}},
		right: &HuffmanTree{
			data: 0x01}}

	table := tree.getEncodingTable()

	var input, output, expected_output bytes.Buffer

	// empty input
	err := tree.encodeBody(&input, &output, table)

	if err != nil {
		t.Errorf("encodeTreeLeaves failed. Got error %s", err)
	}

	len_buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(len_buff, uint64(0))

	expected_output.WriteByte(BLOCK_ID_DATA)
	expected_output.Write(len_buff)
	expected_output.WriteByte(0x0)

	if !bytes.Equal(expected_output.Bytes(), output.Bytes()) {
		t.Errorf("Expected %v, got %v", expected_output.Bytes(), output.Bytes())
	}

	output.Reset()
	expected_output.Reset()

	// short input, which should produce just over one byte of output
	// expected body 010 010 010 (+ 000 0000 padding)
	// or more readable: 0100 1001 0000 0000
	input.Write([]byte{0xFF, 0xFF, 0xFF})

	err = tree.encodeBody(&input, &output, table)

	if err != nil {
		t.Errorf("encodeTreeLeaves failed. Got error %s", err)
	}

	binary.LittleEndian.PutUint64(len_buff, uint64(1))

	expected_output.WriteByte(BLOCK_ID_DATA)
	expected_output.Write(len_buff)
	expected_output.WriteByte(0x1)
	expected_output.Write([]byte{0x49, 0x00})

	if !bytes.Equal(expected_output.Bytes(), output.Bytes()) {
		t.Errorf("Expected %v, got %v", expected_output.Bytes(), output.Bytes())
	}
}

func AssertPanic(t *testing.T) {
	if ok := recover(); ok == nil {
		t.Errorf("AssertPanic failed")
	}
}

func TestHuffmanDecodeTreeShape(t *testing.T) {

	var buff bytes.Buffer

	// invalid block id header
	buff.WriteByte(BLOCK_ID_DATA)
	_, err := decodeTreeShape(&buff)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	if err.Error() != "Unexpected block ID" {
		t.Errorf("Unexpected error %s", err)
	}

	buff.Reset()

	len_buff := make([]byte, 8)

	binary.LittleEndian.PutUint64(len_buff, uint64(1))

	buff.WriteByte(BLOCK_ID_SHAPE)
	buff.Write(len_buff)

	// 0 for tree
	// 0000000 for padding
	buff.WriteByte(0x00)

	tree, err := decodeTreeShape(&buff)

	if err != nil {
		t.Errorf("Got error %s", err)
	}

	if !(tree != nil && tree.left == nil && tree.right == nil) {
		t.Errorf("Unexpected tree: %+v", tree)
	}

	buff.Reset()

	buff.WriteByte(BLOCK_ID_SHAPE)

	binary.LittleEndian.PutUint64(len_buff, 1)
	buff.Write(len_buff)
	// 100 for tree
	// 00000 for padding
	buff.Write([]byte{0x80})

	tree, err = decodeTreeShape(&buff)

	if err != nil {
		t.Errorf("Got error %s", err)
	}

	if !(tree != nil &&
		tree.left != nil && tree.left.left == nil && tree.left.right == nil &&
		tree.right != nil && tree.right.left == nil && tree.right.right == nil) {
		t.Errorf("Unexpected tree: %v", tree)
	}

	buff.Reset()

	buff.WriteByte(BLOCK_ID_SHAPE)

	binary.LittleEndian.PutUint64(len_buff, 2)
	buff.Write(len_buff)
	// 111001000 for tree
	// 0000000 for padding
	buff.Write([]byte{0xE4, 0x00})

	tree, err = decodeTreeShape(&buff)

	if err != nil {
		t.Errorf("Got error %s", err)
	}

	expected_tree := &HuffmanTree{ // 1
		left: &HuffmanTree{ // 1
			left: &HuffmanTree{ // 1
				left:  &HuffmanTree{},  // 0
				right: &HuffmanTree{}}, // 0
			right: &HuffmanTree{ // 1
				left:  &HuffmanTree{},   // 0
				right: &HuffmanTree{}}}, // 0
		right: &HuffmanTree{}} // 0

	if !reflect.DeepEqual(tree, expected_tree) {
		t.Errorf("Unexpected tree: %v", tree)
	}
}

func TestHuffmanTreeEncodeDecodeRandom(t *testing.T) {

	random_tree := func(n int) (tree *HuffmanTree) {

		tree = &HuffmanTree{}

		for i := 0; i < n; i++ {
			node := tree

			for node.left != nil {
				if rand.Intn(2) == 0 {
					node = node.left
				} else {
					node = node.right
				}
			}

			node.left = &HuffmanTree{}
			node.right = &HuffmanTree{}
		}

		return tree
	}

	var buff bytes.Buffer

	for size := 0; size < 50; size++ {

		for n := 0; n < 10; n++ {

			input_tree := random_tree(size)
			input_tree.encodeTreeShape(&buff)

			output_tree, err := decodeTreeShape(&buff)

			if err != nil {
				t.Errorf("Got error %s", err)
			}

			if !reflect.DeepEqual(input_tree, output_tree) {
				t.Errorf("Unexpected tree: %v", output_tree)
			}
		}

	}

}

func TestHuffmanTreeDecodeLeaves(t *testing.T) {

	with_leaf_data := &HuffmanTree{
		data: 0x01}

	without_leaf_data := &HuffmanTree{}

	var buff bytes.Buffer

	with_leaf_data.encodeTreeLeaves(&buff)

	err := without_leaf_data.decodeTreeLeaves(&buff)

	if err != nil {
		t.Errorf("Got error %s", err)
	}

	if !reflect.DeepEqual(with_leaf_data, without_leaf_data) {
		t.Errorf("Unexpected tree: %v", without_leaf_data)
	}

	// zero byte as data
	with_leaf_data = &HuffmanTree{
		data: 0x00}

	without_leaf_data = &HuffmanTree{}

	with_leaf_data.encodeTreeLeaves(&buff)

	err = without_leaf_data.decodeTreeLeaves(&buff)

	if err != nil {
		t.Errorf("Got error %s", err)
	}

	if !reflect.DeepEqual(with_leaf_data, without_leaf_data) {
		t.Errorf("Unexpected tree: %v", without_leaf_data)
	}

	// bigger tree
	with_leaf_data = &HuffmanTree{
		left: &HuffmanTree{
			left: &HuffmanTree{
				left: &HuffmanTree{
					left: &HuffmanTree{
						data: 0xBB},
					right: &HuffmanTree{
						data: 0x99}},
				right: &HuffmanTree{
					data: 0xC0}},
			right: &HuffmanTree{
				left: &HuffmanTree{
					data: 0xFF},
				right: &HuffmanTree{
					data: 0x00}}},
		right: &HuffmanTree{
			data: 0x01}}

	without_leaf_data = &HuffmanTree{
		left: &HuffmanTree{
			left: &HuffmanTree{
				left: &HuffmanTree{
					left:  &HuffmanTree{},
					right: &HuffmanTree{}},
				right: &HuffmanTree{}},
			right: &HuffmanTree{
				left:  &HuffmanTree{},
				right: &HuffmanTree{}}},
		right: &HuffmanTree{}}

	with_leaf_data.encodeTreeLeaves(&buff)

	err = without_leaf_data.decodeTreeLeaves(&buff)

	if err != nil {
		t.Errorf("Got error %s", err)
	}

	if !reflect.DeepEqual(with_leaf_data, without_leaf_data) {
		t.Errorf("Unexpected tree: %v", without_leaf_data)
	}
}

func TestHuffmanEncodeDecode(t *testing.T) {

	encode_decode := func(input []byte) (output []byte, err error) {
		var reader, buff, writer bytes.Buffer
		reader.Write(input)

		err = Encode(&reader, &buff)
		if err != nil {
			return
		}

		err = Decode(&buff, &writer)
		output = writer.Bytes()
		return
	}

	for length := 0; length < 100; length++ {
		for n := 0; n < 100; n++ {

			input := make([]byte, length)

			for i, _ := range input {
				input[i] = byte(rand.Int())
			}

			output, err := encode_decode(input)
			if err != nil {
				t.Errorf("Got unexpected error '%s'", err.Error())
			}

			if !bytes.Equal(input, output) {
				t.Errorf("Expected '%v', got '%v'", input, output)
			}
		}
	}

}
