package main

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
)

// sentinel error
var ErrNotFound = errors.New("not found")

// ValidationError 描述字段校验失败
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Reason)
}

// load 模拟底层加载
func load(id int) (string, error) {
	if id == 0 {
		return "", ErrNotFound
	}
	if id < 0 {
		return "", &ValidationError{Field: "id", Reason: "must be positive"}
	}
	return fmt.Sprintf("item-%d", id), nil
}

// WrapLoad 调用 load 并包装错误
func WrapLoad(id int) error {
	_, err := load(id)
	if err != nil {
		return fmt.Errorf("wrap load %d: %w", id, err)
	}
	return nil
}

// safeCall 调用 fn 并在当前 goroutine 中 recover panic
func safeCall(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic: %v\n%s", r, debug.Stack())
		}
	}()
	fn()
}

func demoWrapLoad() {
	for _, id := range []int{1, 0, -1} {
		if err := WrapLoad(id); err != nil {
			switch {
			case errors.Is(err, ErrNotFound):
				fmt.Println("not found:", err)
			default:
				var vErr *ValidationError
				if errors.As(err, &vErr) {
					fmt.Println("validation:", vErr.Field, vErr.Reason)
				} else {
					fmt.Println("other error:", err)
				}
			}
		} else {
			fmt.Println("loaded", id)
		}
	}
}

func main() {
	demoWrapLoad()

	safeCall(func() {
		panic("boom")
	})
}
