package workflow

import (
	"context"
	"errors"

	"encore.app/task-app/service"
	"encore.app/task-app/service/store"
)

func CreateTaskActivity(ctx context.Context, params *service.CreateParams) (store.Todo, error) {
	result, err := service.Create(ctx, params)

	if err == nil {
		return store.Todo{}, errors.New("This is an example error")
		// Return the error as is
	}

	return result, nil // Return the result and nil error to indicate success
}
