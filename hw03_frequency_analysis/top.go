package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func Top10(input string) []string {
	words := strings.Fields(input)
	wordsCount := map[string]int{}

	for _, w := range words {
		if w != "" {
			wordsCount[w]++
		}
	}

	keys := make([]string, 0, len(wordsCount))
	for key := range wordsCount {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	sort.SliceStable(keys, func(i, j int) bool {
		return wordsCount[keys[i]] > wordsCount[keys[j]]
	})

	return keys[:min(10, len(keys))]
}
