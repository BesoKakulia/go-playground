package service

import (
	"context"
	"database/sql"

	"encore.app/task-app/service/store"
	"encore.dev/storage/sqldb"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var db = sqldb.NewDatabase("taks_app", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

//encore:service
type Service struct {
	db     *sql.DB
	store  *store.Queries
	client client.Client
	worker worker.Worker
}

//lint:ignore U1000 initService is executed by Encore.
func initService() (svc *Service, err error) {

	dbStdLib := db.Stdlib()

	store := store.New(dbStdLib)

	svc = &Service{db: dbStdLib, store: store}

	return svc, nil
}

func (s *Service) Shutdown(force context.Context) {
	s.client.Close()
	s.worker.Stop()
}

// ------------------------------------------------

type CreateParams struct {
	Name string `json:"name"`
}

//encore:api private
func (service *Service) Create(ctx context.Context, params *CreateParams) (store.Todo, error) {
	result, err := service.store.CreateTodo(ctx, params.Name)

	return result, err
}

type SingleParams struct {
	ID string
}

//encore:api private
func (service *Service) Get(ctx context.Context, params *SingleParams) (*Todo, error) {
	result, err := service.store.GetTodo(ctx, uuid.MustParse(params.ID))

	return &Todo{ID: result.ID, Name: result.Name}, err
}

//encore:api private
func (service *Service) Delete(ctx context.Context, params *SingleParams) (*SingleParams, error) {
	err := service.store.DeleteTodo(ctx, uuid.MustParse(params.ID))

	return params, err
}

type UpdateParams struct {
	Id   string
	Name string
}

//encore:api private
func (service *Service) Update(ctx context.Context, params *UpdateParams) (*UpdateParams, error) {
	_, err := service.store.UpdateTodo(ctx, store.UpdateTodoParams{
		ID:   uuid.MustParse(params.Id),
		Name: params.Name,
	})

	return params, err
}

type TodosResponse struct {
	Result []Todo
}

//encore:api private
func (service *Service) List(ctx context.Context) (*TodosResponse, error) {
	todos, err := service.store.ListTodos(ctx)

	var todoResponses []Todo

	for _, todo := range todos {
		todoResponse := Todo{
			ID:   todo.ID,
			Name: todo.Name,
		}
		todoResponses = append(todoResponses, todoResponse)
	}

	return &TodosResponse{Result: todoResponses}, err
}
