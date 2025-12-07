package main

import (
	"errors"
	"testing"
)

type fakeFetcher struct {
	data map[string]string
	err  error
}

func (f *fakeFetcher) Fetch(id string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	val, ok := f.data[id]
	if !ok {
		return "", errors.New("not found")
	}
	return val, nil
}

func TestProcess(t *testing.T) {
	f := &fakeFetcher{
		data: map[string]string{"1": "hello", "2": "go"},
	}
	got, err := Process(f, "1")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got != "HELLO" {
		t.Fatalf("want HELLO, got %s", got)
	}
}

func TestProcessError(t *testing.T) {
	f := &fakeFetcher{err: errors.New("boom")}
	_, err := Process(f, "1")
	if err == nil {
		t.Fatalf("expected error")
	}
}
