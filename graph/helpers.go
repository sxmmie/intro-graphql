package graph

import (
	"github.com/sxmmie/intro-graphql/graph/model"
	"github.com/sxmmie/intro-graphql/models"
)

func convertTodo(todo *models.Todo) *model.Todo {
	if todo == nil {
		return nil
	}

	return &model.Todo{
		ID:        todo.ID,
		Text:      todo.Text,
		Done:      todo.Done,
		CreatedAt: todo.CreatedAt,
	}
}

func convertTodos(todos []*models.Todo) []*model.Todo {
	result := make([]*model.Todo, len(todos))
	for i, todo := range todos {
		result[i] = convertTodo(todo)
	}

	return result
}
