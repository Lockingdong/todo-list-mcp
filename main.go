package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Todo 代表一個待辦事項
// Each todo item contains an ID, title, completion status, and creation timestamp
type Todo struct {
	ID        string    `json:"id"`         // Unique identifier for the todo item
	Title     string    `json:"title"`      // The text content of the todo item
	Completed bool      `json:"completed"`  // Whether the todo item is completed
	CreatedAt time.Time `json:"created_at"` // When the todo item was created
}

// TodoStore 是一個用於儲存待辦事項的記憶體儲存器
// It uses a mutex to ensure thread-safe operations on the todos map
type TodoStore struct {
	sync.RWMutex                  // Mutex for thread-safe operations
	todos        map[string]*Todo // Map of todo items indexed by their IDs
}

// NewTodoStore creates a new TodoStore
// Returns an initialized TodoStore with an empty map
func NewTodoStore() *TodoStore {
	return &TodoStore{
		todos: make(map[string]*Todo),
	}
}

// Add 新增一個待辦事項
// The todo is stored in memory and indexed by its ID
// Returns nil on success
func (s *TodoStore) Add(todo *Todo) error {
	s.Lock()
	defer s.Unlock()
	s.todos[todo.ID] = todo
	return nil
}

// Delete 刪除一個待辦事項
// Returns an error if the todo doesn't exist
func (s *TodoStore) Delete(id string) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.todos[id]; !exists {
		return errors.New("找不到待辦事項")
	}

	delete(s.todos, id)
	return nil
}

// Get 取得所有待辦事項
// Returns a slice of all todos in no particular order
func (s *TodoStore) Get() []*Todo {
	s.RLock()
	defer s.RUnlock()

	todos := make([]*Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}
	return todos
}

// Update 更新待辦事項的完成狀態
// Returns an error if the todo with the given ID doesn't exist
func (s *TodoStore) Update(id string, completed bool) error {
	s.Lock()
	defer s.Unlock()

	todo, exists := s.todos[id]
	if !exists {
		return errors.New("找不到待辦事項")
	}

	todo.Completed = completed
	return nil
}

func main() {
	// 建立一個新的 MCP 伺服器
	// The server handles the communication between the client and the todo store
	s := server.NewMCPServer(
		"Todo List Demo", // Application name
		"1.0.0",          // Version number
		server.WithResourceCapabilities(true, true), // Enable resource capabilities
		server.WithLogging(),                        // Enable logging
		server.WithRecovery(),                       // Enable panic recovery
	)

	// 建立待辦事項儲存器
	store := NewTodoStore()

	// 新增待辦事項工具
	// This tool creates a new todo item with the given title
	addTodoTool := mcp.NewTool("add_todo",
		mcp.WithDescription("新增一個待辦事項"),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("待辦事項的標題"),
		),
	)

	// Add the add_todo tool to the server
	// This tool creates a new todo item with the given title
	s.AddTool(addTodoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title := request.Params.Arguments["title"].(string)
		if title == "" {
			return nil, errors.New("標題不能為空")
		}

		// Create a new todo with a timestamp-based ID
		todo := &Todo{
			ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
			Title:     title,
			Completed: false,
			CreatedAt: time.Now(),
		}
		if err := store.Add(todo); err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(fmt.Sprintf("已新增待辦事項：%s，ID：%s", todo.Title, todo.ID)), nil
	})

	// 刪除待辦事項工具
	// This tool removes a todo item by its ID
	deleteTodoTool := mcp.NewTool("delete_todo",
		mcp.WithDescription("刪除一個待辦事項"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("要刪除的待辦事項 ID"),
		),
	)

	s.AddTool(deleteTodoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := request.Params.Arguments["id"].(string)
		if err := store.Delete(id); err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(fmt.Sprintf("已刪除 ID 為 %s 的待辦事項", id)), nil
	})

	// 取得待辦事項工具
	// This tool retrieves and displays all todo items
	getTodosTool := mcp.NewTool("get_todos",
		mcp.WithDescription("取得所有待辦事項"),
	)

	s.AddTool(getTodosTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		todos := store.Get()
		if len(todos) == 0 {
			return mcp.NewToolResultText("目前沒有待辦事項"), nil
		}

		// Format the todo list with checkboxes
		var result string
		for _, todo := range todos {
			status := "[ ]" // Unchecked box
			if todo.Completed {
				status = "[x]" // Checked box
			}
			result += fmt.Sprintf("%s %s (ID: %s)\n", status, todo.Title, todo.ID)
		}
		return mcp.NewToolResultText(result), nil
	})

	// 更新待辦事項工具
	// This tool updates the completion status of a todo item
	updateTodoTool := mcp.NewTool("update_todo",
		mcp.WithDescription("更新待辦事項的完成狀態"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("要更新的待辦事項 ID"),
		),
		mcp.WithBoolean("completed",
			mcp.Required(),
			mcp.Description("新的完成狀態"),
		),
	)

	s.AddTool(updateTodoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := request.Params.Arguments["id"].(string)
		completed := request.Params.Arguments["completed"].(bool)

		if err := store.Update(id, completed); err != nil {
			return nil, err
		}

		status := "未完成"
		if completed {
			status = "已完成"
		}
		return mcp.NewToolResultText(fmt.Sprintf("已將待辦事項 %s 更新為%s", id, status)), nil
	})

	// Start the server and handle any errors
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("伺服器錯誤：%v\n", err)
	}
}
