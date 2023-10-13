package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

// по условию дополнительного задания:
// ищем либо слова со знаками препинания внутри, либо просто слова, либо знаки препинания, если их > 1 подряд.
var re = regexp.MustCompile(`(\p{L}+[\p{P}\p{S}]*\p{L}+|\p{L}|[\p{P}\p{S}]{2,})`)

func Top10(s string) []string {
	m := make(map[string]int)
	words := strings.Fields(s)

	for _, v := range words {
		match := re.FindStringSubmatch(v)
		if match == nil {
			continue
		}
		w := strings.ToLower(match[1])
		val, ok := m[w]
		if !ok {
			m[w] = 1
			continue
		}
		m[w] = val + 1
	}

	keys := make([]string, 0, len(words))

	for k := range m {
		keys = append(keys, k)
	}
	// сначала сортировка по частоте, потом лексикографическая
	sort.Slice(keys, func(i, j int) bool {
		if m[keys[i]] > m[keys[j]] {
			return true
		}
		if m[keys[i]] == m[keys[j]] {
			return keys[i] < keys[j]
		}
		return false
	})

	if len(keys) > 10 {
		keys = keys[:10]
	}

	return keys
}
