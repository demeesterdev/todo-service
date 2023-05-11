package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/google/uuid"

	"github.com/demeesterdev/todo-service/pkg/authorization"
	ep "github.com/demeesterdev/todo-service/pkg/authorization/endpoints"
)

func MakeHTTPHandler(ep ep.Endpoints, logger log.Logger) http.Handler {
	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Get("/status", httptransport.NewServer(
		ep.ServiceStatusEndpoint,
		DecodeHTTPServiceStatusRequest,
		encodeResponse,
		options...,
	).ServeHTTP)
	r.Post("/", httptransport.NewServer(
		ep.AddUserEndpoint,
		DecodeHTTPAddUserRequest,
		encodeResponse,
		options...,
	).ServeHTTP)
	r.Get("/", httptransport.NewServer(
		ep.GetUsersEndpoint,
		DecodeHTTPGetUsersRequest,
		encodeResponse,
		options...,
	).ServeHTTP)
	r.Post("/login", httptransport.NewServer(
		ep.AuthenticateUserEndpoint,
		DecodeHTTPAuthenticateUserRequest,
		encodeResponse,
		options...,
	).ServeHTTP)
	r.Get("/{id}", httptransport.NewServer(
		ep.GetUserEndpoint,
		DecodeHTTPGetUserRequest,
		encodeResponse,
		options...,
	).ServeHTTP)
	r.Put("/{id}", httptransport.NewServer(
		ep.UpdateUserEndpoint,
		DecodeHTTPUpdateUserRequest,
		encodeResponse,
		options...,
	).ServeHTTP)
	r.Delete("/{id}", httptransport.NewServer(
		ep.DeleteUserEndpoint,
		DecodeHTTPDeleteUserRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	return r
}

func DecodeHTTPServiceStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ep.ServiceStatusRequest
	return req, nil
}

func DecodeHTTPAddUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ep.AddUserRequest
	err := json.NewDecoder(r.Body).Decode(&req.User)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeHTTPGetUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	userIdRaw := chi.URLParam(r, "id")
	userId, err := uuid.Parse(userIdRaw)
	if err == nil {
		return ep.GetUserRequest{User: authorization.User{ID: userId}}, nil
	}
	return ep.GetUserRequest{User: authorization.User{Username: userIdRaw}}, nil
}

func DecodeHTTPGetUsersRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ep.GetUsersRequest
	return req, nil
}

func DecodeHTTPUpdateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ep.UpdateUserRequest
	var err error
	req.ID, err = uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(r.Body).Decode(&req.User)
	if err != nil {
		fmt.Println("json error: ", err.Error())
		return nil, err
	}

	return req, nil
}

func DecodeHTTPAuthenticateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ep.AuthenticateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req.User)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeHTTPDeleteUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	userIdRaw := chi.URLParam(r, "id")
	userId, err := uuid.Parse(userIdRaw)
	if err != nil {
		return nil, err
	}
	return ep.DeleteUserRequest{ID: userId}, nil
}

// errorer is implemented by all concrete response types that may contain
// errors.
type errorer interface {
	Error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(errorer)
	if ok && e.Error() != nil {
		encodeError(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case authorization.ErrInvalidUserObject:
		w.WriteHeader(http.StatusBadRequest)
	case authorization.ErrAuthenticationFailed:
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
