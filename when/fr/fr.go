package fr

import (
	"sort"
	"strings"

	"github.com/olebedev/when/rules"
)

var All = []rules.Rule{
	// Deadline(rules.Override),
	Hour(rules.Override),
	HourMinute(rules.Override),
}

var INTEGER_WORDS = func() map[string]int {
	twoToNine := map[int]string{
		2: "deux",
		3: "trois",
		4: "quatre",
		5: "cinq",
		6: "six",
		7: "sept",
		8: "huit",
		9: "neuf",
	}

	words := map[string]int{
		"zéro":     0,
		"un":       1,
		"une":      1,
		"dix":      10,
		"onze":     11,
		"douze":    12,
		"treize":   13,
		"quatorze": 14,
		"quinze":   15,
		"seize":    16,
		"dix-sept": 17,
		"dix sept": 17,
		"dix-huit": 18,
		"dix huit": 18,
		"dix-neuf": 19,
		"dix neuf": 19,
	}

	tenth := map[string]int{
		"vingt":     20,
		"trente":    30,
		"quarante":  40,
		"cinquante": 50,
	}

	for k, x := range twoToNine {
		words[x] = k
	}

	for k, x := range tenth {
		words[k] = x
		words[(k + "-et-un")] = x + 1
		words[(k + " et un")] = x + 1
		for i := 2; i < 10; i++ {
			words[(k + "-" + twoToNine[i])] = x + i
			words[(k + " " + twoToNine[i])] = x + i
		}
	}

	return words
}()

var INTEGER_WORDS_PATTERN = func() string {
	var words []string
	for k := range INTEGER_WORDS {
		words = append(words, k)
	}
	sort.SliceStable(words, func(i, j int) bool {
		return len(words[i]) > len(words[j])
	})
	return `(?:` + strings.Join(words, "|") + `)`
}()
