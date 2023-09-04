package workflow

import (
	"time"

	"encore.app/task-app/service"
	"encore.app/task-app/service/store"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func CreateTask(ctx workflow.Context, params *service.CreateParams) (store.Todo, error) {

	retrypolicy := &temporal.RetryPolicy{
		MaximumAttempts: 2,
	}

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
		RetryPolicy:         retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result store.Todo
	err := workflow.ExecuteActivity(ctx, CreateTaskActivity, params).Get(ctx, &result)

	return result, err
}
