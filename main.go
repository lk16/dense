package main

import (
	"dense/huffman"
	"flag"
	"os"
)

func main() {

	flag_input_file := flag.String("i", "", "Input file")
	flag_output_file := flag.String("o", "", "Output file")
	flag_decode := flag.Bool("d", false, "If used, specifies decompressing.")
	flag.Parse()

	input_file := os.Stdin
	output_file := os.Stdout
	var err error

	if *flag_input_file != "" {
		input_file, err = os.Open(*flag_input_file)
		if err != nil {
			panic(err)
		}
	}

	if *flag_output_file != "" {
		output_file, err = os.Open(*flag_output_file)
		if err != nil {
			panic(err)
		}
	}

	if *flag_decode {
		huffman.Decode(input_file, output_file)
	} else {
		huffman.Encode(input_file, output_file)
	}

}
