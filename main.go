package main

import (
	"bytes"
	"dense/huffmantree"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {

	file_name := flag.String("i", "", "Input file")
	max_group_len := flag.Int("mfgl", 2, "Maximum length in bytes for grouping in algorithm")
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

	if *max_group_len < 1 || *max_group_len >= 8 {
		panic("Specified max_group_len is not allowed")
	}

	var buf1, buf2 bytes.Buffer
	multi_writer := io.MultiWriter(&buf1, &buf2)
	io.Copy(multi_writer, file)

	tree := huffmantree.NewHuffmanTree(&buf1, *max_group_len)
	tree.Print()

	shape_buff := tree.ToShapeBuff()
	fmt.Printf("\n")
	fmt.Printf("%s\n", shape_buff.String())

	value_buff := tree.ToValueBuff()
	fmt.Printf("\n")
	fmt.Printf("%s\n", value_buff.String())

}
