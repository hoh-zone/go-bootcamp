package main

import (
	"testing"

	"go.uber.org/mock/gomock"
)

func TestProcessWithMockFetcher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockFetcher(ctrl)
	m.EXPECT().
		Fetch("1").
		Return("hello", nil)

	out, err := Process(m, "1")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if out != "HELLO" {
		t.Fatalf("want HELLO, got %s", out)
	}
}
