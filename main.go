package main

import (
	"flag"
	"fmt"
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
	input_flag := flag.String("input", "", "Input string for the compression algorithm")

	flag.Parse()

	input := *input_flag

	freqs := GetFrequencies(input, 2)

	fmt.Printf("input = %s\n\n", input)
	fmt.Printf("freqs =\n")

	for _, pair := range freqs {
		fmt.Printf("%d\t'%s'\n", pair.freq, pair.item)
	}

}
