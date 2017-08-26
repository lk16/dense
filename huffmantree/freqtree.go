package huffmantree

import (
	"fmt"
	"io"
)

type FrequencyTree struct {
	children map[byte]*FrequencyTree
	count    int
}

func NewFrequencyTree(reader io.Reader, max_group_len int) (freq_root *FrequencyTree) {

	freq_root = &FrequencyTree{
		children: make(map[byte]*FrequencyTree, 10),
		count:    0}

	// visited nodes will only have nil values
	visited_nodes := make([]*FrequencyTree, max_group_len)

	buff := make([]byte, 4096)

	for {
		read_bytes, err := reader.Read(buff)
		if err != nil {
			break
		}

		for byte_index, the_byte := range buff[:read_bytes] {
			visited_nodes[byte_index%max_group_len] = freq_root
			freq_root.count++

			for visited_index, visited_node := range visited_nodes {
				if visited_node == nil {
					continue
				}

				var child *FrequencyTree
				var ok bool
				if child, ok = visited_node.children[the_byte]; !ok {
					visited_node.children[the_byte] = &FrequencyTree{
						children: make(map[byte]*FrequencyTree, 10),
						count:    0}
					child = visited_node.children[the_byte]
				}
				if visited_node != freq_root {
					visited_node.count--
				}
				child.count++
				visited_nodes[visited_index] = child

			}
		}
	}
	return
}

func (node *FrequencyTree) printRecursive(depth int) {
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

func (node *FrequencyTree) Print() {
	fmt.Printf("root (%d):\n", node.count)
	node.printRecursive(0)
}

func (node *FrequencyTree) appendHuffmanTreeSlice(slice *[]HuffmanTree, length int, prefix *[]byte) {
	*prefix = append(*prefix, byte(0))
	for b, child := range node.children {
		(*prefix)[length-1] = b

		bytes := make([]byte, len(*prefix))
		copy(bytes, *prefix)

		child.appendHuffmanTreeSlice(slice, length+1, prefix)

		if child.count != 0 {
			*slice = append(*slice, HuffmanTree{
				bytes:  bytes,
				weight: child.count * length,
				left:   nil,
				right:  nil})
		}

	}
	*prefix = (*prefix)[:length-1]
}

func (node *FrequencyTree) ToHuffmanTreeSlice() (slice []HuffmanTree) {
	slice = []HuffmanTree{}
	prefix := []byte{}
	node.appendHuffmanTreeSlice(&slice, 1, &prefix)
	return
}

func (node *FrequencyTree) Prune() {
	for _, child := range node.children {
		child.pruneRecursive()
	}
}

func (node *FrequencyTree) pruneRecursive() (weight int) {

	weight = node.count

	subtree_weights := map[byte]int{}

	for child_byte, child := range node.children {
		subtree_weight := child.pruneRecursive()

		subtree_weights[child_byte] = subtree_weight
		weight += subtree_weight
	}

	minimum_weight := 0.5 * float64(weight)

	for child_byte, subtree_weight := range subtree_weights {
		if float64(subtree_weight) < minimum_weight || subtree_weight < 20 {
			node.count += subtree_weight
			delete(node.children, child_byte)
		}
	}
	return
}
