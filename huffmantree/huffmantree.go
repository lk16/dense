package huffmantree

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
)

type HuffmanTree struct {
	item   string
	weight int
	left   *HuffmanTree
	right  *HuffmanTree
}

func New(file *os.File, max_group_len int) (tree *HuffmanTree) {
	var input_buff bytes.Buffer
	io.Copy(&input_buff, file)
	input := input_buff.String()

	freqs := getFrequencies(input, max_group_len)
	freqs = trimFrequencies(freqs, max_group_len)
	tree = genHuffmanTree(freqs)
	return
}

func trimFrequencies(freqs []HuffmanTree, max_group_len int) (trimmed []HuffmanTree) {
	i := 0
	for i < len(freqs) {
		node := freqs[i]
		if len(node.item) > 1 && node.weight == 1 {
			last := len(freqs) - 1
			if i != last {
				freqs[i] = freqs[last]
			}
			freqs = freqs[:last-1]
			continue
		} else {
			i++
		}
	}
	trimmed = freqs
	return
}

func getFrequencies(str string, max_group_len int) (freqs []HuffmanTree) {
	freq_map := make(map[string]int, 10)

	for start := range str {
		for end := start + 1; end < len(str) && end-start <= max_group_len; end++ {
			slice := str[start:end]
			if _, ok := freq_map[slice]; ok {
				freq_map[slice] += 1
			} else {
				freq_map[slice] = 1
			}
		}
	}

	freqs = make([]HuffmanTree, len(freq_map))
	i := 0
	for item, freq := range freq_map {
		freqs[i] = HuffmanTree{
			item:   item,
			weight: freq,
			left:   nil,
			right:  nil}
		i += 1
	}

	sort.Slice(freqs, func(i, j int) bool { return freqs[i].weight > freqs[j].weight })
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
			item:   "",
			weight: left.weight + right.weight,
			left:   left,
			right:  right}

	}
	tree = new(HuffmanTree)
	*tree = nodes[0]
	return
}

func (node *HuffmanTree) print(code string) {
	if node.item == "" {
		node.left.print(code + "0")
		node.right.print(code + "1")
	} else {
		fmt.Printf("%d\t'%s'\t%s\n", node.weight, node.item, code)
	}
}

func (node *HuffmanTree) Print() {
	node.print("")
}
