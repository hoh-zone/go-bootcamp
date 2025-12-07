package main

import (
	"fmt"
	"testing"
)

func TestCount(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected map[rune]int
	}{
		{"empty", "", map[rune]int{}},
		{"single", "a", map[rune]int{'a': 1}},
		{"repeat", "aba", map[rune]int{'a': 2, 'b': 1}},
		{"unicode", "你好你", map[rune]int{'你': 2, '好': 1}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := Count(tc.input)
			if len(got) != len(tc.expected) {
				t.Fatalf("len mismatch, expected %d got %d", len(tc.expected), len(got))
			}
			for r, want := range tc.expected {
				if got[r] != want {
					t.Fatalf("rune %q: want %d got %d", r, want, got[r])
				}
			}
		})
	}
}

func BenchmarkCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Count("the quick brown fox jumps over the lazy dog")
	}
}

func ExampleCount() {
	fmt.Println(Count("aba"))
	// Output: map[97:2 98:1]
}
