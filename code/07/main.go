package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	PriceCents int    `json:"price_cents"`
}

// Discount 按百分比折扣，最小值 0。
func (p *Product) Discount(percent int) {
	if percent <= 0 {
		return
	}
	if percent > 100 {
		percent = 100
	}
	p.PriceCents = p.PriceCents * (100 - percent) / 100
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Tags []string
}

func (u *User) AddTag(tag string) {
	u.Tags = append(u.Tags, tag)
}

type Audit struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Employee struct {
	User
	Audit
	Title string `json:"title"`
}

type Post struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content,omitempty"`
	Author  string    `json:"author,omitempty"`
	Created time.Time `json:"created_at"`
}

func NewPost(title string, content string) (*Post, error) {
	if title == "" {
		return nil, errors.New("title required")
	}
	return &Post{
		Title:   title,
		Content: content,
		Created: time.Now(),
	}, nil
}

func demoProduct() {
	p := Product{ID: 1, Name: "Book", PriceCents: 1000}
	p.Discount(25)
	fmt.Println("discounted product:", p)
}

func demoEmployee() {
	e := Employee{
		User:  User{Name: "Alice", Age: 30},
		Audit: Audit{CreatedAt: time.Now()},
		Title: "Engineer",
	}
	e.AddTag("gopher") // 提升的方法
	e.UpdatedAt = time.Now()
	fmt.Println("employee name:", e.Name, "title:", e.Title)
}

func demoPost() {
	post, err := NewPost("hello", "content here")
	if err != nil {
		fmt.Println("post error:", err)
		return
	}
	b, _ := json.Marshal(post)
	fmt.Println("post json:", string(b))
}

func main() {
	demoProduct()
	demoEmployee()
	demoPost()
}
