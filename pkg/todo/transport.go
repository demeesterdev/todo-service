package todo

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/google/uuid"
)

func MakeHTTPHandler(ep Endpoints, logger log.Logger) http.Handler {
	r := chi.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Get("/status", httptransport.NewServer(
		ep.ServiceStatusEndpoint,
		decodeHTTPServiceStatusRequest,
		encodeResponse,
		options...,
	).ServeHTTP)
	r.Get("/", httptransport.NewServer(
		ep.GetTodosEndpoint,
		decodeHTTPGetTodosRequest,
		encodeResponse,
		options...,
	).ServeHTTP)
	r.Get("/{id}", httptransport.NewServer(
		ep.GetTodoEndpoint,
		DecodeHTTPGetTodoRequest,
		encodeResponse,
		options...,
	).ServeHTTP)
	r.Post("/", httptransport.NewServer(
		ep.AddTodoEndpoint,
		decodeHTTPAddTodoRequest,
		encodeResponse,
		options...,
	).ServeHTTP)
	r.Put("/{id}", httptransport.NewServer(
		ep.UpdateTodoEndpoint,
		decodeHTTPUpdateTodoRequest,
		encodeResponse,
		options...,
	).ServeHTTP)
	r.Delete("/{id}", httptransport.NewServer(
		ep.DeleteTodoEndpoint,
		decodeHTTPDeleteTodoRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	return r
}

// Server functions
func decodeHTTPServiceStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req serviceStatusRequest
	return req, nil
}

func decodeHTTPAddTodoRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req addTodoRequest
	err := json.NewDecoder(r.Body).Decode(&req.Todo)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeHTTPGetTodoRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req getTodoRequest
	var err error
	idRaw := chi.URLParam(r, "id")
	req.ID, err = uuid.Parse(idRaw)
	if err != nil {
		return nil, ErrInvalidUUID
	}
	return req, nil
}

func decodeHTTPGetTodosRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req getTodosRequest
	var err error
	q := r.URL.Query()
	if q.Has("owner") {
		req.OwnerID, err = uuid.Parse(q.Get("owner"))
		if err != nil {
			return nil, ErrInvalidUUID
		}
	}
	return req, nil
}

func decodeHTTPUpdateTodoRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req updateTodoRequest
	var err error
	req.ID, err = uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(r.Body).Decode(&req.Todo)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPDeleteTodoRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req deleteTodoRequest
	var err error
	req.ID, err = uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return nil, ErrInvalidUUID
	}
	return req, nil
}

// client functions
// encode request for server

func encodeHTTPGetTodosRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Get("/", ...)
	req.URL.Path = "/"
	return encodeRequest(ctx, req, request)
}

func encodeHTTPServiceStatusRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Get("/status", ...)
	req.URL.Path = "/status/"
	return encodeRequest(ctx, req, request)
}

func encodeHTTPGetTodoRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Get("/{id}", ...)
	r := request.(getTodoRequest)
	todoID := url.QueryEscape(r.ID.String())
	req.URL.Path = "/" + todoID
	return encodeRequest(ctx, req, request)
}

func encodeHTTPAddTodoRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Post("/", ...)
	req.URL.Path = "/"
	return encodeRequest(ctx, req, request)
}

func encodeHTTPUpdateTodoRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Put("/{id}", ...)
	r := request.(updateTodoRequest)
	todoID := url.QueryEscape(r.ID.String())
	req.URL.Path = "/" + todoID
	return encodeRequest(ctx, req, request)
}

func encodeHTTPDeleteTodoRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Delete("/{id}, ...)
	req.URL.Path = "/status/"
	r := request.(updateTodoRequest)
	todoID := url.QueryEscape(r.ID.String())
	req.URL.Path = "/" + todoID
	return encodeRequest(ctx, req, request)
}

// client functions
// decode response from server

func decodeHTTPAddTodoResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response addTodoResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}
func decodeHTTPGetTodoResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getTodoResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}
func decodeHTTPUpdateTodoResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response updateTodoResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}
func decodeHTTPDeleteTodoResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response deleteTodoResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}
func decodeHTTPGetTodosResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getTodosResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}
func decodeHTTPServiceStatusResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response serviceStatusResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
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

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// profilesvc endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case ErrInconsistentIDs:
		w.WriteHeader(http.StatusBadRequest)
	case ErrOwnerMissing:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
