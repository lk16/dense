package main

import (
	"dense/huffmantree"
	"flag"
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

	tree := huffmantree.New(file, *max_group_len)
	tree.Print()

}
