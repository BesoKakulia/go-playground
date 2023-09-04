package todos_api

import (
	"context"
	"fmt"

	"encore.app/task-app/api/workflow"
	"encore.app/task-app/service"
	"encore.app/task-app/service/store"
	"encore.dev/rlog"

	"encore.dev"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var (
	envName       = encore.Meta().Environment.Name
	todoTaskQueue = envName + "-todos"
)

//encore:service
type Service struct {
	client client.Client
	worker worker.Worker
}

//lint:ignore U1000 initService is executed by Encore.
func initService() (svc *Service, err error) {

	c, err := client.Dial(client.Options{})
	if err != nil {
		return nil, fmt.Errorf("create temporal client: %v", err)
	}

	w := worker.New(c, todoTaskQueue, worker.Options{})

	w.RegisterWorkflow(workflow.CreateTask)
	w.RegisterActivity(workflow.CreateTaskActivity)

	err = w.Start()
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("start temporal worker: %v", err)
	}

	svc = &Service{client: c, worker: w}

	return svc, nil
}

func (s *Service) Shutdown(force context.Context) {
	s.client.Close()
	s.worker.Stop()
}

// ==================================================================

type Response struct {
	Message string
}

//encore:api public method=POST path=/api/todos
func (s *Service) createTodo(ctx context.Context, params *service.CreateParams) (*Response, error) {

	options := client.StartWorkflowOptions{
		TaskQueue: todoTaskQueue,
	}

	we, err := s.client.ExecuteWorkflow(ctx, options, workflow.CreateTask, params)
	if err != nil {
		return nil, err
	}
	rlog.Info("started workflow", "id", we.GetID(), "run_id", we.GetRunID())

	// Get the results
	var task store.Todo
	err = we.Get(ctx, &task)
	if err != nil {
		return nil, err
	}

	return &Response{Message: "created"}, err
}

//encore:api public method=GET path=/api/todos
func getTodos(ctx context.Context) (*service.TodosResponse, error) {
	result, err := service.List(ctx)

	return result, err
}

//encore:api public method=GET path=/api/todos/:id
func getTodo(ctx context.Context, id string) (*service.Todo, error) {
	result, err := service.Get(ctx, &service.SingleParams{ID: id})

	return result, err
}

//encore:api public method=DELETE path=/api/todos/:id
func deleteTodo(ctx context.Context, id string) (*service.SingleParams, error) {

	result, err := service.Delete(ctx, &service.SingleParams{ID: id})

	return result, err
}

//encore:api public method=PUT path=/api/todos/:id
func updateTodo(ctx context.Context, id string, params *service.CreateParams) (*service.UpdateParams, error) {

	result, err := service.Update(ctx, &service.UpdateParams{Id: id, Name: params.Name})

	return result, err
}

// ==================================================================

// Encore comes with a built-in development dashboard for
// exploring your API, viewing documentation, debugging with
// distributed tracing, and more. Visit your API URL in the browser:
//
//     http://localhost:4000
//

// ==================================================================

// Next steps
//
// 1. Deploy your application to the cloud with a single command:
//
//     git push encore
//
// 2. To continue exploring Encore, check out one of these topics:
//
//    Building a Slack bot:  https://encore.dev/docs/tutorials/slack-bot
//    Building a REST API:   https://encore.dev/docs/tutorials/rest-api
//    Using SQL databases:   https://encore.dev/docs/develop/databases
//    Authenticating users:  https://encore.dev/docs/develop/auth
