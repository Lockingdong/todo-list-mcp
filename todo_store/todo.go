package todo_store

import (
	"fmt"
	"time"
)

// Todo 代表一個待辦事項
// Each todo item contains an ID, title, completion status, and creation timestamp
type Todo struct {
	ID        string    `json:"id"`         // Unique identifier for the todo item
	Title     string    `json:"title"`      // The text content of the todo item
	Completed bool      `json:"completed"`  // Whether the todo item is completed
	CreatedAt time.Time `json:"created_at"` // When the todo item was created
}

// NewTodo 創建一個新的待辦事項
func NewTodo(title string) *Todo {
	return &Todo{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
}
