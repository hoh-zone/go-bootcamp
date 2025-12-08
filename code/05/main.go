package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func main() {
	log := func(msg string) {
		fmt.Printf("[%s] %s\n", time.Now().Format(time.RFC3339), msg)
	}

	name := "  Alice  "
	age := 18

	var err error
	if name, err = validateUser(name, age); err != nil {
		log(fmt.Sprintf("validate error: %v", err))
		return
	}

	switch {
	case age < 18:
		log(fmt.Sprintf("%s is underage", name))
	case age == 18:
		log(fmt.Sprintf("%s just became an adult", name))
	default:
		log(fmt.Sprintf("%s is an adult", name))
	}

	defer log("defer: goodbye") // 后进先出收尾

	msg, err := greet(name, age)
	if err != nil {
		log(fmt.Sprintf("greet error: %v", err))
		return
	}
	log(msg)

	DoWork()
}

func validateUser(name string, age int) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("name required")
	}
	if age < 0 || age > 130 {
		return "", fmt.Errorf("invalid age: %d", age)
	}
	return name, nil
}

func greet(name string, age int) (string, error) {
	if age < 18 {
		return "", fmt.Errorf("user %s is too young", name)
	}
	return fmt.Sprintf("hello, %s", name), nil
}

func DoWork() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from panic:", r)
		}
	}()

	for i := 0; i < 3; i++ {
		fmt.Println("loop", i)
		if i == 2 {
			panic("simulated panic")
		}
	}
}
