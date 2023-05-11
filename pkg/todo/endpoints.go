package todo

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"

	"github.com/demeesterdev/todo-service/pkg/authorization"
)

type Endpoints struct {
	AddTodoEndpoint       endpoint.Endpoint
	GetTodoEndpoint       endpoint.Endpoint
	UpdateTodoEndpoint    endpoint.Endpoint
	DeleteTodoEndpoint    endpoint.Endpoint
	GetTodosEndpoint      endpoint.Endpoint
	GetTodosOwnedEndpoint endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a todo svc
// server.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		AddTodoEndpoint:       makeAddTodoEndpoint(s),
		GetTodoEndpoint:       makeGetTodoEndpoint(s),
		UpdateTodoEndpoint:    makeUpdateTodoEndpoint(s),
		DeleteTodoEndpoint:    makeDeleteTodoEndpoint(s),
		GetTodosEndpoint:      makeGetTodosEndpoint(s),
		ServiceStatusEndpoint: makeServiceStatusEndpoint(s),
	}
}

// MakeClientEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the remote instance, via a transport/http.Client.
// Useful in a profilesvc client.
func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	// Note that the request encoders need to modify the request URL, changing
	// the path. That's fine: we simply need to provide specific encoders for
	// each endpoint.

	return Endpoints{
		AddTodoEndpoint:       httptransport.NewClient("POST", tgt, encodeHTTPAddTodoRequest, decodeHTTPAddTodoResponse, options...).Endpoint(),
		GetTodoEndpoint:       httptransport.NewClient("GET", tgt, encodeHTTPGetTodoRequest, decodeHTTPGetTodoResponse, options...).Endpoint(),
		UpdateTodoEndpoint:    httptransport.NewClient("PUT", tgt, encodeHTTPUpdateTodoRequest, decodeHTTPUpdateTodoResponse, options...).Endpoint(),
		DeleteTodoEndpoint:    httptransport.NewClient("DELETE", tgt, encodeHTTPDeleteTodoRequest, decodeHTTPDeleteTodoResponse, options...).Endpoint(),
		GetTodosEndpoint:      httptransport.NewClient("GET", tgt, encodeHTTPGetTodosRequest, decodeHTTPGetTodosResponse, options...).Endpoint(),
		ServiceStatusEndpoint: httptransport.NewClient("GET", tgt, encodeHTTPServiceStatusRequest, decodeHTTPServiceStatusResponse, options...).Endpoint(),
	}, nil
}

// AddTodo implements Service interface. Primarily useful in a client.
func (e Endpoints) AddTodo(ctx context.Context, t Todo) (Todo, error) {
	request := addTodoRequest{Todo: t}
	response, err := e.AddTodoEndpoint(ctx, request)
	if err != nil {
		return Todo{}, err
	}
	resp := response.(addTodoResponse)
	return resp.Todo, resp.Err
}

// GetTodo implements Service interface. Primarily useful in a client.
func (e Endpoints) GetTodo(ctx context.Context, id uuid.UUID) (Todo, error) {
	request := getTodoRequest{ID: id}
	response, err := e.GetTodoEndpoint(ctx, request)
	if err != nil {
		return Todo{}, nil
	}
	resp := response.(getTodoResponse)
	return resp.Todo, resp.Err
}

// UpdateTodo implements Service interface. Primarily useful in a client.
func (e Endpoints) UpdateTodo(ctx context.Context, id uuid.UUID, t Todo) (Todo, error) {
	request := updateTodoRequest{ID: id, Todo: t}
	response, err := e.UpdateTodoEndpoint(ctx, request)
	if err != nil {
		return Todo{}, err
	}
	resp := response.(updateTodoResponse)
	return resp.Todo, resp.Err
}

// DeleteTodo implements Service interface. Primarily useful in a client.
func (e Endpoints) DeleteTodo(ctx context.Context, id uuid.UUID) error {
	request := deleteTodoRequest{ID: id}
	response, err := e.DeleteTodoEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(updateTodoResponse)
	return resp.Err
}

// GetTodos implements Service interface. Primarily useful in a client.
func (e Endpoints) GetTodos(ctx context.Context) ([]Todo, error) {
	request := getTodosRequest{}
	response, err := e.UpdateTodoEndpoint(ctx, request)
	if err != nil {
		return []Todo{}, err
	}
	resp := response.(getTodosResponse)
	return resp.Todos, resp.Err
}

// GetTodosOwned implements Service interface. Primarily useful in a client.
func (e Endpoints) GetTodosOwned(ctx context.Context, user authorization.User) ([]Todo, error) {
	request := getTodosRequest{OwnerID: user.ID}
	response, err := e.UpdateTodoEndpoint(ctx, request)
	if err != nil {
		return []Todo{}, err
	}
	resp := response.(getTodosResponse)
	return resp.Todos, resp.Err
}

// ServiceStatus implements Service interface. Primarily useful in a client.
func (e Endpoints) ServiceStatus(ctx context.Context) (int, error) {
	request := serviceStatusRequest{}
	response, err := e.ServiceStatusEndpoint(ctx, request)
	if err != nil {
		return 0, err
	}
	resp := response.(serviceStatusResponse)
	return resp.Code, resp.Err
}

// makeAddTodoEndpoint returns an enpoint via the passed service.
// Primarily useful in a server
func makeAddTodoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addTodoRequest)
		t, e := s.AddTodo(ctx, req.Todo)
		return addTodoResponse{Todo: t, Err: e}, nil
	}
}

// makeGetTodoEndpoint returns an enpoint via the passed service.
// Primarily useful in a server
func makeGetTodoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getTodoRequest)
		t, e := s.GetTodo(ctx, req.ID)
		return getTodoResponse{Todo: t, Err: e}, nil
	}
}

// makeUpdateTodoEndpoint returns an enpoint via the passed service.
// Primarily useful in a server
func makeUpdateTodoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateTodoRequest)
		t, e := s.UpdateTodo(ctx, req.ID, req.Todo)
		return updateTodoResponse{Todo: t, Err: e}, nil
	}
}

// makedeleteTodoEndpoint returns an enpoint via the passed service.
// Primarily useful in a server
func makeDeleteTodoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteTodoRequest)
		e := s.DeleteTodo(ctx, req.ID)
		return deleteTodoResponse{Err: e}, nil
	}
}

// makeGetTodosEndpoint returns an enpoint via the passed service.
// Primarily useful in a server
func makeGetTodosEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getTodosRequest)

		//get all todos if owner id is empty
		if req.OwnerID == uuid.Nil {
			t, e := s.GetTodos(ctx)
			return getTodosResponse{Todos: t, Err: e}, nil
		}

		t, e := s.GetTodosOwned(ctx, authorization.User{ID: req.OwnerID})
		return getTodosResponse{Todos: t, Err: e}, nil
	}
}

// makeServicestatusEndpoint returns an enpoint via the passed service.
// Primarily useful in a server
func makeServiceStatusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(serviceStatusRequest)
		code, e := s.ServiceStatus(ctx)

		return serviceStatusResponse{Code: code, Err: e}, nil

	}
}
