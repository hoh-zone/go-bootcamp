package main

import (
	"errors"
	"fmt"
)

// Average 计算可变参数的平均值；若为空返回错误。
func Average(nums ...float64) (float64, error) {
	if len(nums) == 0 {
		return 0, errors.New("no numbers provided")
	}
	var sum float64
	for _, n := range nums {
		sum += n
	}
	return sum / float64(len(nums)), nil
}

// Pred 谓词函数类型，用于过滤切片。
type Pred[T any] func(T) bool

// Filter 返回满足谓词的元素新切片。
func Filter[T any](xs []T, pred Pred[T]) []T {
	var out []T
	for _, v := range xs {
		if pred(v) {
			out = append(out, v)
		}
	}
	return out
}

// NewCounter 返回一个闭包，每次调用返回递增值（首次返回 start）。
func NewCounter(start int) func() int {
	current := start
	return func() int {
		val := current
		current++
		return val
	}
}

func main() {
	avg, err := Average(10, 20, 30)
	if err != nil {
		fmt.Println("average error:", err)
		return
	}
	fmt.Println("average:", avg) // 20

	evens := Filter([]int{1, 2, 3, 4, 5, 6}, func(n int) bool { return n%2 == 0 })
	fmt.Println("evens:", evens) // [2 4 6]

	counter := NewCounter(3)
	fmt.Println("counter:", counter(), counter(), counter()) // 3 4 5
}
