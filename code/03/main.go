package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

func main() {
	log.SetPrefix("[demo] ")
	log.SetFlags(0)

	// 1) 打印当前时间与命令行参数
	fmt.Println("current time:", time.Now().Format(time.RFC3339))
	fmt.Println("args:", os.Args)

	// 2) 字符串的字节长度与 rune 数
	text := "Go语言"
	fmt.Printf("text=%q bytes=%d runes=%d\n", text, len(text), utf8.RuneCountInString(text))
	// idx 是 UTF-8 字节偏移，不是字符序号；中文占 3 字节，所以会出现 0,1,2,5 的偏移
	for idx, r := range text {
		fmt.Printf("idx=%d rune=%c\n", idx, r)
	}

	// 3) 调用返回 (result, error) 的函数并日志输出
	result, err := greetFromArgs(os.Args)
	if err != nil {
		log.Printf("greet error: %v", err)
		return
	}
	log.Printf("greet result: %s", result)
}

func greetFromArgs(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("need at least one name argument")
	}
	name := strings.TrimSpace(args[1])
	if name == "" {
		return "", errors.New("name cannot be empty")
	}
	return fmt.Sprintf("hello, %s", strings.ToUpper(name)), nil
}
