package todo

import (
	"context"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/demeesterdev/todo-service/pkg/authorization"
	"github.com/google/uuid"
)

type dbSvc struct {
	db *gorm.DB
}

// NewService create a new service based on an sqlite database with a persistent file
func NewDBService(dbconnection gorm.Dialector) (Service, error) {
	db, err := gorm.Open(dbconnection, &gorm.Config{})
	db.AutoMigrate(&Todo{})
	if err != nil {
		return &dbSvc{}, err
	}

	return &dbSvc{
		db: db,
	}, nil
}

func NewSqliteDBService(target string) (Service, error) {
	return NewDBService(sqlite.Open(target))
}

func NewInMemService() (Service, error) {
	return NewSqliteDBService(":memory:")
}

func (s *dbSvc) AddTodo(ctx context.Context, t Todo) (Todo, error) {
	if t.OwnerID == uuid.Nil {
		return Todo{}, ErrOwnerMissing
	}

	result := s.db.Create(&t)

	if result.Error != nil {
		return Todo{}, result.Error
	}

	return t, nil
}

func (s *dbSvc) GetTodo(ctx context.Context, id uuid.UUID) (Todo, error) {

	var t Todo
	result := s.db.First(&t, "id = ?", id.String())

	switch result.Error {
	case gorm.ErrRecordNotFound:
		return Todo{}, ErrNotFound
	default:
		return t, result.Error
	}

}

func (s *dbSvc) UpdateTodo(ctx context.Context, id uuid.UUID, t Todo) (Todo, error) {
	if t.ID == uuid.Nil {
		t.ID = id
	}

	if t.ID != id {
		return Todo{}, ErrInconsistentIDs
	}

	current, err := s.GetTodo(ctx, id)
	if err != nil {
		return Todo{}, err
	}

	if t.OwnerID != uuid.Nil && current.OwnerID != t.OwnerID {
		return Todo{}, ErrOwnerChanged
	}

	current.Title = t.Title
	current.Description = t.Description

	result := s.db.Model(&current).Updates(current)

	if result.Error != nil {
		return Todo{}, nil
	}

	s.db.First(&t, "id = ?", id.String())
	return t, nil
}

func (s *dbSvc) DeleteTodo(ctx context.Context, id uuid.UUID) error {

	t := Todo{}
	result := s.db.Delete(&t, "id = ?", id.String())

	return result.Error
}

func (s *dbSvc) GetTodos(ctx context.Context) ([]Todo, error) {

	var todos []Todo
	result := s.db.Find(&todos)

	if result.Error != nil {
		return []Todo{}, nil
	}

	return todos, nil
}

func (s *dbSvc) GetTodosOwned(ctx context.Context, user authorization.User) ([]Todo, error) {

	var todos []Todo
	result := s.db.Where(&Todo{OwnerID: user.ID}).Find(&todos)

	if result.Error != nil {
		return []Todo{}, nil
	}

	return todos, nil
}

func (s *dbSvc) ServiceStatus(ctx context.Context) (int, error) {
	db, err := s.db.DB()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = db.Ping()
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	return http.StatusOK, nil
}
