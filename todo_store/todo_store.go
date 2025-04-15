package todo_store

import (
	"errors"
)

// TodoStore 是一個用於儲存待辦事項的記憶體儲存器
type TodoStore struct {
	todos map[string]*Todo // Map of todo items indexed by their IDs
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
	s.todos[todo.ID] = todo
	return nil
}

// Delete 刪除一個待辦事項
// Returns an error if the todo doesn't exist
func (s *TodoStore) Delete(id string) error {
	if _, exists := s.todos[id]; !exists {
		return errors.New("找不到待辦事項")
	}

	delete(s.todos, id)
	return nil
}

// Get 取得所有待辦事項
// Returns a slice of all todos in no particular order
func (s *TodoStore) Get() []*Todo {
	todos := make([]*Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}
	return todos
}

// Update 更新待辦事項的完成狀態
// Returns an error if the todo with the given ID doesn't exist
func (s *TodoStore) Update(id string, completed bool) error {
	todo, exists := s.todos[id]
	if !exists {
		return errors.New("找不到待辦事項")
	}

	todo.Completed = completed
	return nil
}
