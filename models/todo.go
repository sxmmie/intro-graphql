package models

import (
	"fmt"
	"sync"
	"time"
)

type Todo struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at"`
}

type TodoStore struct {
	mu     sync.RWMutex
	todos  map[string]*Todo
	nextID int
}

func NewTodoStore() *TodoStore {
	return &TodoStore{
		todos:  make(map[string]*Todo),
		nextID: 1,
	}
}

// returns all todos
func (s *TodoStore) GetAll() []*Todo {
	s.mu.RLock() // read lock
	defer s.mu.RUnlock()

	todos := make([]*Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}

	return todos
}

func (s *TodoStore) GetByID(id string) *Todo {
	s.mu.Lock()
	defer s.mu.RUnlock()

	return s.todos[id]
}

// GetByStatus returns a todo based on the done status
func (s *TodoStore) GetByStatus(done bool) []*Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filtered []*Todo
	for _, todo := range s.todos {
		if todo.Done == done {
			filtered = append(filtered, todo)
		}
	}

	return filtered
}

// Create adds new todo
func (s *TodoStore) Create(text string) *Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo := &Todo{
		ID:        fmt.Sprintf("%d", s.nextID),
		Text:      text,
		Done:      false,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	s.nextID++
	s.todos[todo.ID] = todo
	return todo
}

func (s *TodoStore) Update(id string, text *string, done *bool) *Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	// grab the todo from the list(or db)
	todo := s.todos[id]
	if todo == nil {
		return nil
	}

	if text != nil {
		todo.Text = *text
	}

	if done != nil {
		todo.Done = *done
	}

	return todo
}

func (s *TodoStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.todos[id]; !exists {
		return false
	}
	delete(s.todos, id)

	return true
}

// Switches the done status between true/false
func (s *TodoStore) Toggle(id string) *Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo := s.todos[id]
	if todo == nil {
		return nil
	}

	todo.Done = !todo.Done
	return todo
}
