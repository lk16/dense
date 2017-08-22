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

func (node *freqTree) toHuffmanTreeSliceRecursive(slice *[]HuffmanTree, prefix *[]byte) {
	*prefix = append(*prefix, byte(0))
	for b, child := range node.children {
		(*prefix)[len(*prefix)-1] = b

		bytes := make([]byte, len(*prefix))
		copy(bytes, *prefix)

		child.toHuffmanTreeSliceRecursive(slice, prefix)

		*slice = append(*slice, HuffmanTree{
			bytes:  bytes,
			weight: child.count,
			left:   nil,
			right:  nil})

	}
	*prefix = (*prefix)[:len(*prefix)-1]
}

func (node *freqTree) ToHuffmanTreeSlice() (slice []HuffmanTree) {
	slice = []HuffmanTree{}
	prefix := []byte{}
	node.toHuffmanTreeSliceRecursive(&slice, &prefix)
	return
}

func (node *freqTree) Prune() {
	node.pruneRecursive(node, 0)
}

func (node *freqTree) pruneRecursive(root *freqTree, depth int) {

	for child_byte, child := range node.children {
		// go depth first
		child.pruneRecursive(root, depth+1)

		if depth >= 1 && len(child.children) == 0 {

			unpruned_ratio := float32(child.count) / float32(node.count)
			pruned_ratio := float32(root.children[child_byte].count) / float32(root.count)

			if pruned_ratio >= unpruned_ratio {
				delete(node.children, child_byte)
			}

		}

	}
}
