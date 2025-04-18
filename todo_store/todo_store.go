package todo_store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// TodoStore 是一個用於儲存待辦事項的JSON檔案儲存器
type TodoStore struct {
	filePath string // Path to the JSON storage file
}

// NewTodoStore creates a new TodoStore
func NewTodoStore(storeFilePath string) (*TodoStore, error) {
	storePath := filepath.Join(storeFilePath, "todos.json")

	// Initialize empty JSON file if it doesn't exist
	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		if err := os.WriteFile(storePath, []byte("{}"), 0644); err != nil {
			return nil, err
		}
	}

	return &TodoStore{
		filePath: storePath,
	}, nil
}

// Add 新增一個待辦事項
func (s *TodoStore) Add(todo *Todo) error {
	todos, err := s.readTodos()
	if err != nil {
		return err
	}

	todos[todo.ID] = todo
	return s.writeTodos(todos)
}

// Get 取得所有待辦事項
func (s *TodoStore) Get() ([]*Todo, error) {
	todos, err := s.readTodos()
	if err != nil {
		return nil, err
	}

	result := make([]*Todo, 0, len(todos))
	for _, todo := range todos {
		result = append(result, todo)
	}
	return result, nil
}

// Update 更新待辦事項的完成狀態
func (s *TodoStore) Update(id string, completed bool) error {
	todos, err := s.readTodos()
	if err != nil {
		return err
	}

	todo, exists := todos[id]
	if !exists {
		return errors.New("找不到待辦事項")
	}

	todo.Completed = completed
	return s.writeTodos(todos)
}

// Delete 刪除一個待辦事項
func (s *TodoStore) Delete(id string) error {
	todos, err := s.readTodos()
	if err != nil {
		return err
	}

	if _, exists := todos[id]; !exists {
		return errors.New("找不到待辦事項")
	}

	delete(todos, id)
	return s.writeTodos(todos)
}

// readTodos reads todos from the JSON file
func (s *TodoStore) readTodos() (map[string]*Todo, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	todos := make(map[string]*Todo)
	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

// writeTodos writes todos to the JSON file
func (s *TodoStore) writeTodos(todos map[string]*Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}
