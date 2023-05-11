package authorization

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// User presents a single user object
// ID should be globally unique and a valid uuid
type User struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primarykey"`
	Username string    `json:"username" gorm:"unique;index"`
	Password string    `json:"password,omitempty"`
}

// Service is a simple CRUD intreface for userObjects
type Service interface {
	AddUser(ctx context.Context, t User) (User, error)
	GetUser(ctx context.Context, id uuid.UUID) (User, error)
	FindUser(ctx context.Context, username string) (User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, u User) (User, error)
	AuthenticateUser(ctx context.Context, u User) (User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetUsers(ctx context.Context) ([]User, error)
	ServiceStatus(ctx context.Context) (int, error)
}

var (
	ErrInvalidUserObject      = errors.New("invalid user")
	ErrAuthenticationFailed   = errors.New("authentication failed")
	ErrIDMissing              = errors.New("missing user identifier or username")
	ErrInconsistentIDs        = errors.New("inconsistent ids")
	ErrInconsistentIDUserName = errors.New("inconsistent id and username")
	ErrNotFound               = errors.New("not found")
)
