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

func TestGetTodos(t *testing.T) {
	store := models.NewTodoStore()
	resolver := &Resolver{TodoStore: store}

	store.Create("Todo 1")
	store.Create("Todo 2")

	todos, err := resolver.Query().Todos(context.Background())
	if err != nil {
		t.Fatalf("failed to get all todos", err)
	}

	assert.NoError(t, err)
	assert.Len(t, todos, 2)
}

func generateString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = 'a' + byte(rand.IntN(26))
	}
	return string(b)
}
