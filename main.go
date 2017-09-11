package main

import (
	"flag"
	"fmt"
	"github.com/lk16/dense/huffman"
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
		if _, err := os.Stat(*flag_output_file); !os.IsNotExist(err) {
			fmt.Printf("File '%s' exists already. Exiting.\n", *flag_output_file)
			return
		}

		output_file, err = os.Create(*flag_output_file)
		if err != nil {
			fmt.Printf("Could not create file '%s': %s\n", *flag_output_file, err)
			return
		}
	}

	if *flag_decode {
		err = huffman.Decode(input_file, output_file)
	} else {
		err = huffman.Encode(input_file, output_file)
	}

	if err != nil {
		fmt.Printf("An error occorred: '%s'\n", err)
	}

}
