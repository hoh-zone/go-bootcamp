package main

import (
	"fmt"
	"strings"
)

//go:generate mockgen -destination mock_fetcher_test.go -package main . Fetcher

// Fetcher 抽象数据源，用于演示接口驱动设计与测试。
type Fetcher interface {
	Fetch(id string) (string, error)
}

// Process 通过 Fetcher 获取数据并转换为大写。
func Process(f Fetcher, id string) (string, error) {
	data, err := f.Fetch(id)
	if err != nil {
		return "", fmt.Errorf("fetch %s: %w", id, err)
	}
	return strings.ToUpper(data), nil
}
