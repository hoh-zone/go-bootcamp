package main

// Count 返回字符串中每个字符出现的次数。
func Count(text string) map[rune]int {
	out := make(map[rune]int)
	for _, r := range text {
		out[r]++
	}
	return out
}
