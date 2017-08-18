package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
)

type ItemFrequency struct {
	item string
	freq int
}

func GetFrequencies(str string, max_len int) (freqs []ItemFrequency) {
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

	freqs = make([]ItemFrequency, len(freq_map))
	i := 0
	for item, freq := range freq_map {
		freqs[i] = ItemFrequency{
			item: item,
			freq: freq}
		i += 1
	}

	sort.Slice(freqs, func(i, j int) bool { return freqs[i].freq > freqs[j].freq })
	return
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

	freqs := GetFrequencies(input, *max_freq_group_len)

	for _, pair := range freqs {
		fmt.Printf("%d\t'%s'\n", pair.freq, pair.item)
	}

}
