package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"example.com/todo-list-mcp/todo_store"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

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

	// 設定待辦事項儲存器路徑
	const storeFilePath = "/Users/lidongying/Documents/Projects/todo-list-mcp"

	// 建立待辦事項儲存器
	store, err := todo_store.NewTodoStore(storeFilePath)
	if err != nil {
		log.Fatalf("建立待辦事項儲存器失敗：%v", err)
	}

	// 新增待辦事項工具
	// This tool creates a new todo item with the given title
	addTodoTool := mcp.NewTool("add_todo",
		mcp.WithDescription("新增一個待辦事項"), // 新增待辦事項的描述
		mcp.WithString("title",
			mcp.Required(),             // 標題是必需的
			mcp.Description("待辦事項的標題"), // 標題的描述
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
		todo := todo_store.NewTodo(title)
		if err := store.Add(todo); err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(fmt.Sprintf("已新增待辦事項：%s，ID：%s", todo.Title, todo.ID)), nil
	})

	// 取得待辦事項工具
	// This tool retrieves and displays all todo items
	getTodosTool := mcp.NewTool("get_todos",
		mcp.WithDescription("取得所有待辦事項"),
	)

	s.AddTool(getTodosTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		todos, err := store.Get()
		if err != nil {
			return nil, err
		}
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

	// Start the server and handle any errors
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("伺服器錯誤：%v\n", err)
	}
}
