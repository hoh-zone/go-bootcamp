package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"unicode/utf8"
)

var (
	timeFormat = flag.String("timefmt", time.RFC3339, "layout for printing current time")
)

func init() {
	flag.Parse()
}

func main() {
	log.SetPrefix("[demo] ")
	log.SetFlags(0)

	// 1) 打印当前时间与命令行参数
	fmt.Println("current time:", time.Now().Format(*timeFormat))
	fmt.Println("args (raw):", os.Args)
	fmt.Println("args (after flag parsing):", flag.Args())

	// 2) 字符串的字节长度与 rune 数
	text := "Go语言"
	fmt.Printf("text=%q bytes=%d runes=%d\n", text, len(text), utf8.RuneCountInString(text))
	// idx 是 UTF-8 字节偏移，不是字符序号；中文占 3 字节，所以会出现 0,1,2,5 的偏移
	for idx, r := range text {
		fmt.Printf("idx=%d rune=%c\n", idx, r)
	}
}
