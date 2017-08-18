package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
)

type HuffmanTreeNode struct {
	item   string
	weight int
	left   *HuffmanTreeNode
	right  *HuffmanTreeNode
}

func GetFrequencies(str string, max_len int) (freqs []HuffmanTreeNode) {
	freq_map := make(map[string]int, 10)

	for start := range str {
		for end := start + 1; end < len(str) && end-start <= max_len; end++ {
			slice := str[start:end]
			if _, ok := freq_map[slice]; ok {
				freq_map[slice] += 1
			} else {
				freq_map[slice] = 1
			}
		}
	}

	freqs = make([]HuffmanTreeNode, len(freq_map))
	i := 0
	for item, freq := range freq_map {
		freqs[i] = HuffmanTreeNode{
			item:   item,
			weight: freq,
			left:   nil,
			right:  nil}
		i += 1
	}

	sort.Slice(freqs, func(i, j int) bool { return freqs[i].weight > freqs[j].weight })
	return
}

func GenHuffmanTree(nodes []HuffmanTreeNode) (tree *HuffmanTreeNode) {
	for len(nodes) > 1 {
		sort.Slice(nodes, func(i, j int) bool { return nodes[i].weight > nodes[j].weight })

		left := new(HuffmanTreeNode)
		right := new(HuffmanTreeNode)

		*left = nodes[len(nodes)-1]
		*right = nodes[len(nodes)-2]

		nodes = nodes[:len(nodes)-1]
		nodes[len(nodes)-1] = HuffmanTreeNode{
			item:   "",
			weight: left.weight + right.weight,
			left:   left,
			right:  right}

	}
	tree = new(HuffmanTreeNode)
	*tree = nodes[0]
	return
}

func (node *HuffmanTreeNode) print(code string) {
	if node.item == "" {
		node.left.print(code + "0")
		node.right.print(code + "1")
	} else {
		fmt.Printf("%d\t'%s'\t%s\n", node.weight, node.item, code)
	}
}

func (node *HuffmanTreeNode) Print() {
	node.print("")
}

func main() {

	file_name := flag.String("i", "", "Input file")
	max_freq_group_len := flag.Int("mfgl", 2, "Maximum length in bytes for grouping in algorithm")
	flag.Parse()

	var file *os.File

	if *file_name == "" {
		file = os.Stdin
	} else {
		var err interface{}
		file, err = os.Open(*file_name)
		if err != nil {
			panic(err)
		}
	}

	if *max_freq_group_len < 1 || *max_freq_group_len >= 8 {
		panic("Specified max_freq_group_len is not allowed")
	}

	var input_buff bytes.Buffer
	io.Copy(&input_buff, file)

	input := input_buff.String()

	nodes := GetFrequencies(input, *max_freq_group_len)

	tree := GenHuffmanTree(nodes)
	tree.Print()

}
