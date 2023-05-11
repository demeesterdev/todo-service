package endpoints

import (
	"github.com/demeesterdev/todo-service/pkg/authorization"
	"github.com/google/uuid"
)

// AddUserRequest and AddUserResponse
// returns user without password
type AddUserRequest struct {
	User authorization.User
}

type AddUserResponse struct {
	User authorization.User `json:"user,omitempty"`
	Err  error              `json:"err,omitempty"`
}

//lint:ignore U1000 used to satisfy error interface in transport
func (r AddUserResponse) Error() error { return r.Err }

// GetUserRequest and GetUserResponse
// no filter so an empty request
type GetUserRequest struct {
	User authorization.User
}

type GetUserResponse struct {
	User authorization.User `json:"user,omitempty"`
	Err  error              `json:"err,omitempty"`
}

//lint:ignore U1000 used to satisfy error interface in transport
func (r GetUserResponse) Error() error { return r.Err }

// UpdateUsersRequest and UpdateUsersResponse
// no filter so an empty request
type UpdateUserRequest struct {
	ID   uuid.UUID
	User authorization.User
}

type UpdateUserResponse struct {
	User authorization.User `json:"user,omitempty"`
	Err  error              `json:"err,omitempty"`
}

//lint:ignore U1000 used to satisfy error interface in transport
func (r UpdateUserResponse) Error() error { return r.Err }

// GetUsersRequest and GetUsersResponse
// no filter so an empty request
type GetUsersRequest struct {
}

type GetUsersResponse struct {
	Users []authorization.User `json:"users,omitempty"`
	Err   error                `json:"err,omitempty"`
}

//lint:ignore U1000 used to satisfy error interface in transport
func (r GetUsersResponse) Error() error { return r.Err }

// AuthenticateUserRequestRequest and AuthenticateUserRequestResponse
// no filter so an empty request
type AuthenticateUserRequest struct {
	User authorization.User
}

type AuthenticateUserResponse struct {
	User authorization.User `json:"user,omitempty"`
	Err  error              `json:"err,omitempty"`
}

//lint:ignore U1000 used to satisfy error interface in transport
func (r AuthenticateUserResponse) Error() error { return r.Err }

// DeleteUserRequestRequest and DeleteUserRequestRequestResponse
// only id needede to delete
type DeleteUserRequest struct {
	ID uuid.UUID
}

type DeleteUserResponse struct {
	Err error `json:"err,omitempty"`
}

//lint:ignore U1000 used to satisfy error interface in transport
func (r DeleteUserResponse) Error() error { return r.Err }

// ServiceStatusRequest and ServiceStatusResponse
// empty request and status response
type ServiceStatusRequest struct{}

type ServiceStatusResponse struct {
	Code int   `json:"code"`
	Err  error `json:"err,omitempty"`
}

//lint:ignore U1000 used to satisfy error interface in transport
func (r ServiceStatusResponse) Error() error { return r.Err }
