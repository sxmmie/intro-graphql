package graph

import (
	"context"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sxmmie/intro-graphql/graph/model"
	"github.com/sxmmie/intro-graphql/models"
)

func TestCreateTodo(t *testing.T) {
	store := models.NewTodoStore()
	resolver := &Resolver{TodoStore: store}

	input := model.NewTodo{Text: "Text todo"}

	todo, err := resolver.Mutation().CreateTodo(context.Background(), input)

	assert.NoError(t, err)
	assert.Equal(t, "Text todo", todo.Text)
	assert.False(t, todo.Done)
	assert.NotEmpty(t, todo.ID)
}

// When input is empty, it should return an error
func TestCreateTodo_Validation_Errors(t *testing.T) {
	store := models.NewTodoStore()
	resolver := &Resolver{TodoStore: store}

	input := model.NewTodo{Text: " "}

	todo, err := resolver.Mutation().CreateTodo(context.Background(), input)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	assert.Nil(t, todo)
	assert.Error(t, err)
	assert.Equal(t, "validation failed for text: text cannot be empty", err.Error())

	input = model.NewTodo{Text: generateString(256)}
	todo, err = resolver.Mutation().CreateTodo(context.Background(), input)

	assert.Nil(t, todo)
	assert.Error(t, err)
	assert.Equal(t, "validation failed for text: text cannot be too long", err.Error())
}

func TestDeleteTodo(t *testing.T) {
	store := models.NewTodoStore()
	resolver := &Resolver{TodoStore: store}

	// 1. Create a todo first
	input := model.NewTodo{Text: "Test todo"}
	created, err := resolver.Mutation().CreateTodo(context.Background(), input)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	// 2. Delete it using the returned ID
	ok, err := resolver.Mutation().DeleteTodo(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("failed to delete todo: %v", err)
	}

	// 3. Assert deletion was successful
	if !ok {
		t.Errorf("expected DeleteTodo to return true, got false")
	}
}

func TestUpdateTodo(t *testing.T) {
	store := models.NewTodoStore()
	resolver := &Resolver{TodoStore: store}

	// 1. Create a todo first
	input := model.NewTodo{Text: "Test todo"}
	created, err := resolver.Mutation().CreateTodo(context.Background(), input)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	// 2. Update it using the returned ID
	newText := "Updated todo"
	done := true
	updateInput := model.UpdateTodo{
		Text: &newText,
		Done: &done,
	}

	updated, err := resolver.Mutation().UpdateTodo(context.Background(), created.ID, updateInput)
	if err != nil {
		t.Fatalf("failed to update todo: %v", err)
	}

	// 3. Assert the fields were updated
	if updated.Text != newText {
		t.Errorf("expected text %q, got %q", newText, updated.Text)
	}

	if updated.Done != done {
		t.Errorf("expected done %v, got %v", done, updated.Done)
	}

	if updated.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, updated.ID)
	}
}

func TestGetTodos(t *testing.T) {
	store := models.NewTodoStore()
	resolver := &Resolver{TodoStore: store}

	store.Create("Todo 1")
	store.Create("Todo 2")

	todos, err := resolver.Query().Todos(context.Background())
	if err != nil {
		t.Fatalf("failed to get all todos: %v", err)
	}

	assert.NoError(t, err)
	assert.Len(t, todos, 2)
}

func TestGetTodo(t *testing.T) {
	store := models.NewTodoStore()
	resolver := &Resolver{TodoStore: store}

	created := store.Create("Todo 1")

	todo, err := resolver.Query().Todo(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("failed to get todo", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, "Todo 1", todo)
}

func TestTodoByStatus(t *testing.T) {
	store := models.NewTodoStore()
	resolver := &Resolver{TodoStore: store}

	// 1. Create a todo (starts as done=false by default)
	input := model.NewTodo{Text: "Test todo"}
	created, err := resolver.Mutation().CreateTodo(context.Background(), input)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	// 2. Update it to done=true
	done := true
	updateInput := model.UpdateTodo{Done: &done}
	_, err = resolver.Mutation().UpdateTodo(context.Background(), created.ID, updateInput)
	if err != nil {
		t.Fatalf("failed to update todo: %v", err)
	}

	// 3. Query for completed todos
	todos, err := resolver.Query().TodoByStatus(context.Background(), true)
	if err != nil {
		t.Fatalf("failed to get todos by status: %v", err)
	}

	// 4. Assert the updated todo appears in results
	if len(todos) == 0 {
		t.Fatal("expected at least one todo, got none")
	}

	if todos[0].ID != created.ID {
		t.Errorf("expected todo ID %s, got %s", created.ID, todos[0].ID)
	}

	if todos[0].Done != true {
		t.Errorf("expected done=true, got done=false")
	}
}

func generateString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = 'a' + byte(rand.IntN(26))
	}
	return string(b)
}
