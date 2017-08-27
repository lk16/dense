package huffman

import (
	"bytes"
	"dense/bits"
	"encoding/binary"
	"fmt"
	"io"
	"sort"
)

func Decode(reader io.Reader, writer io.Writer) (written_bits int64, err error) {
	fmt.Printf("This is not implemented\n")
	return
}

func Encode(reader io.Reader, writer io.Writer) (written_bits int64, err error) {

	var freq_tree_reader, encode_reader bytes.Buffer

	multi_writer := io.MultiWriter(&freq_tree_reader, &encode_reader)
	io.Copy(multi_writer, reader)

	freq_table, err := getFrequencyTable(&freq_tree_reader)
	if err != nil {
		return
	}

	huffman_tree := fromFreqencyTable(freq_table)
	table := huffman_tree.GetEncodingTable()

	written_bits, err = huffman_tree.doEncode(&encode_reader, writer, table)
	return
}

type HuffmanTree struct {
	bytes  byte
	weight int64
	left   *HuffmanTree
	right  *HuffmanTree
}

func getFrequencyTable(reader io.Reader) (table map[byte]int64, err error) {

	buff := make([]byte, 4096)
	table = make(map[byte]int64, 256)

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

	for key, value := range table {
		if value == 0 {
			delete(table, key)
		}
	}

	return
}

func fromFreqencyTable(table map[byte]int64) (tree *HuffmanTree) {

	slice := make([]HuffmanTree, 0)

	for key, value := range table {
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
	if node.left == nil {
		node.left.toValueRecursive(buff)
		node.right.toValueRecursive(buff)
		return
	}

	buff.WriteByte(node.bytes)
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

func (node *HuffmanTree) doEncode(reader io.Reader, writer io.Writer, table map[byte]bits.Slice) (written_bits int64, write_err error) {
	bits_writer := bits.NewWriter(writer)
	key_buff := make([]byte, 1)

	var read_err error
	var written int64

	for {
		_, read_err = reader.Read(key_buff)
		if read_err != nil {
			if read_err != io.EOF {
				panic(read_err)
			}
			break
		}
		slice := table[key_buff[0]]
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
