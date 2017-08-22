package huffmantree

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
)

type HuffmanTree struct {
	bytes  []byte
	weight int
	left   *HuffmanTree
	right  *HuffmanTree
}

type freqTree struct {
	children map[byte]*freqTree
	count    int
}

func (node *freqTree) printRecursive(depth int) {
	for key, child := range node.children {
		for i := 0; i < depth; i++ {
			fmt.Print("  ")
		}
		fmt.Printf("'%s' (%d):\n", string(key), child.count)
		child.printRecursive(depth + 1)
	}
	if len(node.children) == 0 {
		for i := 0; i < depth; i++ {
			fmt.Print("  ")
		}
		fmt.Printf("-\n")
	}
}

func (node *freqTree) print() {
	fmt.Printf("root (%d):\n", node.count)
	node.printRecursive(0)
}

func (node *freqTree) toHuffmanTreeSliceRecursive(slice *[]HuffmanTree, prefix *[]byte) {
	*prefix = append(*prefix, byte(0))
	for b, child := range node.children {
		(*prefix)[len(*prefix)-1] = b

		bytes := make([]byte, len(*prefix))
		copy(bytes, *prefix)

		*slice = append(*slice, HuffmanTree{
			bytes:  bytes,
			weight: child.count,
			left:   nil,
			right:  nil})
		child.toHuffmanTreeSliceRecursive(slice, prefix)
	}
	*prefix = (*prefix)[:len(*prefix)-1]
}

func (node *freqTree) toHuffmanTreeSlice() (slice []HuffmanTree) {
	slice = []HuffmanTree{}
	prefix := []byte{}
	node.toHuffmanTreeSliceRecursive(&slice, &prefix)
	return
}

func New(file *os.File, max_group_len int) (tree *HuffmanTree) {
	var input_buff bytes.Buffer
	io.Copy(&input_buff, file)
	input := input_buff.Bytes()

	freqs := getFrequencies(input, max_group_len)
	tree = genHuffmanTree(freqs)
	return
}

func getFrequencies(input []byte, max_group_len int) (freqs []HuffmanTree) {

	root := &freqTree{
		children: make(map[byte]*freqTree, 10),
		count:    0}

	// visited nodes will only have nil values
	visited_nodes := make([]*freqTree, max_group_len)

	for byte_index, the_byte := range input {
		visited_nodes[byte_index%max_group_len] = root
		root.count++

		for visited_index, visited_node := range visited_nodes {
			if visited_node == nil {
				continue
			}

			var child *freqTree
			var ok bool
			if child, ok = visited_node.children[the_byte]; !ok {
				visited_node.children[the_byte] = &freqTree{
					children: make(map[byte]*freqTree, 10),
					count:    0}
				child = visited_node.children[the_byte]
			}
			child.count++
			visited_nodes[visited_index] = child

		}

	}

	root.print()
	fmt.Printf("\n\n")

	// TODO
	freqs = root.toHuffmanTreeSlice()
	return
}

func genHuffmanTree(nodes []HuffmanTree) (tree *HuffmanTree) {
	for len(nodes) > 1 {
		sort.Slice(nodes, func(i, j int) bool { return nodes[i].weight > nodes[j].weight })

		left := new(HuffmanTree)
		right := new(HuffmanTree)

		*left = nodes[len(nodes)-1]
		*right = nodes[len(nodes)-2]

		nodes = nodes[:len(nodes)-1]
		nodes[len(nodes)-1] = HuffmanTree{
			bytes:  []byte{},
			weight: left.weight + right.weight,
			left:   left,
			right:  right}

	}
	tree = new(HuffmanTree)
	*tree = nodes[0]
	return
}

func (node *HuffmanTree) print(code string) {
	if node.left == nil {
		fmt.Printf("%d\t'%s'\t%s\n", node.weight, node.bytes, code)
		return
	}

	node.left.print(code + "0")
	node.right.print(code + "1")
}

func (node *HuffmanTree) Print() {
	node.print("")
}
