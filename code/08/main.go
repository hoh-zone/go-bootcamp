package main

import (
	"fmt"
	"strconv"
)

type Printer interface {
	Print() string
}

type User struct {
	Name string
	Age  int
}

func (u User) Print() string {
	return fmt.Sprintf("User(name=%s, age=%d)", u.Name, u.Age)
}

type Product struct {
	ID         int
	Name       string
	PriceCents int
}

func (p Product) Print() string {
	return fmt.Sprintf("Product[%d]: %s ($%.2f)", p.ID, p.Name, float64(p.PriceCents)/100)
}

func LogAll(ps []Printer) {
	for _, p := range ps {
		fmt.Println(p.Print())
	}
}

type Box struct {
	Val any //any = interface{}
}

func (b Box) AsString() string {
	switch v := b.Val.(type) {
	case nil:
		return "<nil>"
	case fmt.Stringer:
		return v.String()
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func main() {
	users := []User{{Name: "Alice", Age: 20}, {Name: "Bob", Age: 30}}
	products := []Product{
		{ID: 1, Name: "Book", PriceCents: 1500},
		{ID: 2, Name: "Pen", PriceCents: 299},
	}

	var items []Printer
	for _, u := range users {
		items = append(items, u)
	}
	for _, p := range products {
		items = append(items, p)
	}

	LogAll(items)

	boxes := []Box{
		{Val: "hello"},
		{Val: 42},
		{Val: 3.14},
		{Val: users[0]},
		{Val: nil},
	}
	for i, b := range boxes {
		fmt.Printf("box %d: %s\n", i, b.AsString())
	}
}
