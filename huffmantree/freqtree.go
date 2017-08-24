package huffmantree

import (
	"fmt"
)

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

func (node *freqTree) Print() {
	fmt.Printf("root (%d):\n", node.count)
	node.printRecursive(0)
}

func (node *freqTree) toHuffmanTreeSliceRecursive(slice *[]HuffmanTree, length int, prefix *[]byte) {
	*prefix = append(*prefix, byte(0))
	for b, child := range node.children {
		(*prefix)[length-1] = b

		bytes := make([]byte, len(*prefix))
		copy(bytes, *prefix)

		child.toHuffmanTreeSliceRecursive(slice, length+1, prefix)

		if child.count == 0 {
			continue
		}

		*slice = append(*slice, HuffmanTree{
			bytes:  bytes,
			weight: child.count * (length + 1),
			left:   nil,
			right:  nil})

	}
	*prefix = (*prefix)[:length-1]
}

func (node *freqTree) ToHuffmanTreeSlice() (slice []HuffmanTree) {
	slice = []HuffmanTree{}
	prefix := []byte{}
	node.toHuffmanTreeSliceRecursive(&slice, 1, &prefix)
	return
}

func (node *freqTree) Prune() {
	for child_byte, child := range node.children {
		prefix := []byte{}
		prefix = append(prefix, child_byte)
		child.pruneRecursive(&prefix)
	}
}

func (node *freqTree) pruneRecursive(prefix *[]byte) (weight int) {

	weight = node.count

	subtree_weights := map[byte]int{}

	for child_byte, child := range node.children {
		*prefix = append(*prefix, child_byte)
		subtree_weight := child.pruneRecursive(prefix)
		*prefix = (*prefix)[:len(*prefix)-1]

		subtree_weights[child_byte] = subtree_weight
		weight += subtree_weight
	}

	minimum_ratio := 0.3

	for child_byte, subtree_weight := range subtree_weights {
		ratio := float64(subtree_weight) / float64(weight)
		if ratio < minimum_ratio || subtree_weight < 20 {
			node.count += subtree_weight
			delete(node.children, child_byte)
		}
	}
	return
}
