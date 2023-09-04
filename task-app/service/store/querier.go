// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package store

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateTodo(ctx context.Context, name string) (Todo, error)
	DeleteTodo(ctx context.Context, id uuid.UUID) error
	GetTodo(ctx context.Context, id uuid.UUID) (Todo, error)
	ListTodos(ctx context.Context) ([]Todo, error)
	UpdateTodo(ctx context.Context, arg UpdateTodoParams) (Todo, error)
}

var _ Querier = (*Queries)(nil)
