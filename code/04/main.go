package main

import (
	"fmt"
	"sort"
)

func main() {
	words := []string{"go", "rust", "go", "python", "go", "rust"}
	fmt.Println("dedup input :", words)
	fmt.Println("dedup output:", dedup(words))

	fmt.Println("\ncountWords:")
	counts := countWords(words)
	printCounts(counts)
}

// dedup 去重但保持第一次出现的顺序。
func dedup(in []string) []string {
	seen := make(map[string]bool, len(in))
	out := make([]string, 0, len(in))
	for _, v := range in {
		if seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	return out
}

func countWords(words []string) map[string]int {
	m := make(map[string]int, len(words))
	for _, w := range words {
		m[w]++
	}
	return m
}

// 打印 map，键按字典序排序，方便观察结果。
func printCounts(m map[string]int) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s -> %d\n", k, m[k])
	}
}
