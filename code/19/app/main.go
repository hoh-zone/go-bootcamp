package main

import (
	"fmt"
	"os"

	"example.com/go-class/19/lib"
)

func main() {
	text := "hello go"
	if len(os.Args) > 1 {
		text = os.Args[1]
	}

	stats := lib.Count(text)

	fmt.Printf("input: %q\n", text)
	fmt.Println("counts:")
	for r, n := range stats {
		fmt.Printf("  %q: %d\n", r, n)
	}
}
