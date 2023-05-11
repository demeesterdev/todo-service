package todo

import (
	"context"
	"errors"
	"time"

	"github.com/demeesterdev/todo-service/pkg/authorization"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Todo presents a single Todo item.
// ID should be globally unique and a valid uuid
type Todo struct {
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primarykey"`
	Title       string         `json:"title"`
	Description string         `json:"description,omitempty"`
	OwnerID     uuid.UUID      `json:"owner_id"`
}

// Before create is a GORM hook
// It makes shure a ToDo has a valid uuid before creation
func (t *Todo) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}

	return
}

// Service is a simple CRUD intreface for user profiles
type Service interface {
	AddTodo(ctx context.Context, t Todo) (Todo, error)
	GetTodo(ctx context.Context, id uuid.UUID) (Todo, error)
	UpdateTodo(ctx context.Context, id uuid.UUID, t Todo) (Todo, error)
	DeleteTodo(ctx context.Context, id uuid.UUID) error
	GetTodos(ctx context.Context) ([]Todo, error)
	GetTodosOwned(ctx context.Context, user authorization.User) ([]Todo, error)
	ServiceStatus(ctx context.Context) (int, error)
}

var (
	ErrPopulatedID     = errors.New("id filled")
	ErrOwnerChanged    = errors.New("owner_id changed")
	ErrOwnerMissing    = errors.New("owner_id missing")
	ErrInconsistentIDs = errors.New("inconsistent ids")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
	ErrInvalidUUID     = errors.New("invalid uuid")
)
