package endpoints

import (
	"context"

	"github.com/demeesterdev/todo-service/pkg/authorization"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type Endpoints struct {
	AddUserEndpoint          endpoint.Endpoint
	GetUserEndpoint          endpoint.Endpoint
	UpdateUserEndpoint       endpoint.Endpoint
	AuthenticateUserEndpoint endpoint.Endpoint
	DeleteUserEndpoint       endpoint.Endpoint
	GetUsersEndpoint         endpoint.Endpoint
	ServiceStatusEndpoint    endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a authorization svc
// server
func MakeServerEndpoints(s authorization.Service) Endpoints {
	return Endpoints{
		AddUserEndpoint:          MakeAddUserEndpoint(s),
		GetUserEndpoint:          MakeGetUserEndpoint(s),
		UpdateUserEndpoint:       MakeUpdateUserEndpoint(s),
		AuthenticateUserEndpoint: MakeAuthenticateUserEndpoint(s),
		DeleteUserEndpoint:       MakeDeleteUserEndpoint(s),
		GetUsersEndpoint:         MakeGetUsersEndpoint(s),
		ServiceStatusEndpoint:    MakeServiceStatusEndpoint(s),
	}
}

// MakeAddUserEndpoint returns an enpoint via the passed service.
// Primarily useful in a server
func MakeAddUserEndpoint(s authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddUserRequest)
		u, e := s.AddUser(ctx, req.User)
		return AddUserResponse{User: u, Err: e}, nil
	}
}

func MakeGetUserEndpoint(s authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserRequest)
		// if neither Id or username is set why bother the service
		if req.User.ID == uuid.Nil && req.User.Username == "" {
			return GetUserResponse{User: authorization.User{}, Err: authorization.ErrIDMissing}, nil
		}

		// if id is set get user from id
		if req.User.ID != uuid.Nil {
			u, e := s.GetUser(ctx, req.User.ID)
			return GetUserResponse{User: u, Err: e}, nil
		}

		// try to find user by username if id fails
		u, e := s.FindUser(ctx, req.User.Username)
		return GetUserResponse{User: u, Err: e}, nil
	}
}
func MakeUpdateUserEndpoint(s authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateUserRequest)
		u, e := s.UpdateUser(ctx, req.ID, req.User)
		return UpdateUserResponse{User: u, Err: e}, nil
	}
}
func MakeAuthenticateUserEndpoint(s authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AuthenticateUserRequest)
		u, e := s.AuthenticateUser(ctx, req.User)
		return AuthenticateUserResponse{User: u, Err: e}, nil
	}
}
func MakeDeleteUserEndpoint(s authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteUserRequest)
		e := s.DeleteUser(ctx, req.ID)
		return DeleteUserResponse{Err: e}, nil
	}
}
func MakeGetUsersEndpoint(s authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(GetUsersRequest)
		us, e := s.GetUsers(ctx)
		return GetUsersResponse{Users: us, Err: e}, nil
	}
}

// MakeServicestatusEndpoint returns an enpoint via the passed service.
// Primarily useful in a server
func MakeServiceStatusEndpoint(s authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(ServiceStatusRequest)
		code, e := s.ServiceStatus(ctx)

		return ServiceStatusResponse{Code: code, Err: e}, nil
	}
}
