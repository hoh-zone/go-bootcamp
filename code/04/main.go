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

	fmt.Println("\ndemo shared slice vs full slice:")
	demoSharedSlice()
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

func demoSharedSlice() {
	nums := []int{1, 2, 3, 4}
	sub := nums[1:3] // 共享底层数组
	sub[0] = 20
	fmt.Printf("after modifying sub: nums=%v sub=%v\n", nums, sub)

	// 使用 full slice 表达式限制容量，后续 append 会分配新底层数组
	safe := nums[1:3:3] // 第三位是 cap 上界 c，cap= c - a，这里 a=1 c=3 得到 cap=2
	safe = append(safe, 99) // 触发扩容，安全独立
	safe[0] = 30
	fmt.Printf("after append on safe: nums=%v safe=%v\n", nums, safe)
}
