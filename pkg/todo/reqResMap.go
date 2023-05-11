package todo

import (
	"github.com/google/uuid"
)

type addTodoRequest struct {
	Todo Todo
}

type addTodoResponse struct {
	Todo Todo  `json:"todo,omitempty"`
	Err  error `json:"err,omitempty"`
}

//lint:ignore U1000 used to satisfy error interface in github.com/demeesterdev/todo-service/pkg/todo/transport
func (r addTodoResponse) Error() error { return r.Err }

type getTodoRequest struct {
	ID uuid.UUID
}

type getTodoResponse struct {
	Todo Todo  `json:"todo,omitempty"`
	Err  error `json:"err,omitempty"`
}

//lint:ignore U1000 used to satisfy error interface in github.com/demeesterdev/todo-service/pkg/todo/transport
func (r getTodoResponse) Error() error { return r.Err }

type updateTodoRequest struct {
	ID   uuid.UUID
	Todo Todo
}

type updateTodoResponse struct {
	Todo Todo  `json:"todo,omitempty"`
	Err  error `json:"err,omitempty"`
}

//lint:ignore U1000 used to satisfy error interface in github.com/demeesterdev/todo-service/pkg/todo/transport
func (r updateTodoResponse) Error() error { return r.Err }

type deleteTodoRequest struct {
	ID uuid.UUID
}

type deleteTodoResponse struct {
	Err error `json:"err,omitempty"`
}

//lint:ignore U1000 used to satisfy error interface in github.com/demeesterdev/todo-service/pkg/todo/transport
func (r deleteTodoResponse) Error() error { return r.Err }

type getTodosRequest struct {
	OwnerID uuid.UUID
}

type getTodosResponse struct {
	Todos []Todo `json:"todos,omitempty"`
	Err   error  `json:"err,omitempty"`
}

//lint:ignore U1000 used to satisfy error interface in github.com/demeesterdev/todo-service/pkg/todo/transport
func (r getTodosResponse) Error() error { return r.Err }

type serviceStatusRequest struct{}

type serviceStatusResponse struct {
	Code int   `json:"code"`
	Err  error `json:"err,omitempty"`
}

//lint:ignore U1000 used to satisfy error interface in github.com/demeesterdev/todo-service/pkg/todo/transport
func (r serviceStatusResponse) Error() error { return r.Err }
