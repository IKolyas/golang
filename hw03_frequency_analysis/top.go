package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordCountPair struct {
	word string

	count int
}

func Top10(text string) []string {
	wordCount := make(map[string]int)

	words := strings.Fields(text)

	for _, word := range words {
		wordCount[word]++
	}

	pairs := make([]wordCountPair, 0, len(wordCount))

	for word, count := range wordCount {
		pairs = append(pairs, wordCountPair{word, count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count == pairs[j].count {
			return pairs[i].word < pairs[j].word
		}

		return pairs[i].count > pairs[j].count
	})

	if len(pairs) > 10 {
		pairs = pairs[:10]
	}

	topWords := make([]string, len(pairs))

	for i, pair := range pairs {
		topWords[i] = pair.word
	}

	return topWords
}
