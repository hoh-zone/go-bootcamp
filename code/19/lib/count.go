package lib

// Count 返回字符串中每个字符出现的次数。
func Count(text string) map[rune]int {
	stats := make(map[rune]int)
	for _, r := range text {
		stats[r]++
	}
	return stats
}
