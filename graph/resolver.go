package graph

import "github.com/sxmmie/intro-graphql/models"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	TodoStore *models.TodoStore
}

// type Mutation struct {
// 	TodoStore *models.TodoStore
// }
