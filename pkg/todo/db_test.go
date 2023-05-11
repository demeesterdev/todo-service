package todo

import (
	"context"
	"testing"

	"github.com/demeesterdev/todo-service/pkg/authorization"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddTodo(t *testing.T) {
	type args struct {
		todo Todo
	}

	testCases := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "should add items",
			args: args{
				todo: Todo{
					OwnerID:     uuid.New(),
					ID:          uuid.Nil,
					Title:       "New Item",
					Description: "",
				},
			},
			err: nil,
		},
	}

	s, _ := NewInMemService()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.AddTodo(context.Background(), tc.args.todo)

			assert.Equal(t, tc.err, err)
		})
	}
}

func TestGetTodo(t *testing.T) {
	todoId := uuid.New()
	type args struct {
		id uuid.UUID
	}
	testCases := []struct {
		name     string
		args     args
		expected uuid.UUID
		err      error
	}{
		{
			name: "should get single item",
			args: args{
				id: todoId,
			},
			expected: todoId,

			err: nil,
		},
	}

	for _, tc := range testCases {
		ctx := context.Background()
		s, _ := NewInMemService()
		s.AddTodo(ctx, Todo{ID: todoId})

		t.Run(tc.name, func(t *testing.T) {
			actual, err := s.GetTodo(ctx, tc.args.id)

			assert.NotEmpty(t, actual)
			assert.Equal(t, tc.expected, actual.ID)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	todoId := uuid.New()

	type args struct {
		id   uuid.UUID
		todo Todo
	}
	testCases := []struct {
		name     string
		args     args
		expected Todo
		err      error
	}{
		{
			name: "should get single item",
			args: args{
				id: todoId,
				todo: Todo{
					ID:    todoId,
					Title: "new Title",
				},
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, _ := NewInMemService()
			s.AddTodo(context.Background(), Todo{ID: todoId})

			actual, err := s.UpdateTodo(context.Background(), tc.args.id, tc.args.todo)

			assert.Equal(t, tc.args.todo.Title, actual.Title)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		id   uuid.UUID
		todo Todo
	}

	testCases := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "should delete item",
			args: args{
				id: uuid.MustParse("8f0aca9a-0ab5-496f-9a47-8c82be9c307c"),
				todo: Todo{
					ID: uuid.MustParse("8f0aca9a-0ab5-496f-9a47-8c82be9c307c"),
				},
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, _ := NewInMemService()
			s.AddTodo(context.Background(), tc.args.todo)
			err := s.DeleteTodo(context.Background(), tc.args.id)

			assert.Equal(t, tc.err, err)
		})
	}
}

func TestGetTodos(t *testing.T) {
	testCases := []struct {
		name     string
		expected []Todo
		err      error
	}{
		{
			name:     "should get all items",
			expected: []Todo{},
			err:      nil,
		},
	}

	s, _ := NewInMemService()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := s.GetTodos(context.Background())

			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestGetTodosOwned(t *testing.T) {
	type args struct {
		owner authorization.User
		todos []Todo
	}
	testCases := []struct {
		name     string
		args     args
		expected []Todo
		err      error
	}{
		{
			name: "should get all items owned by user id",
			args: args{
				owner: authorization.User{
					ID: uuid.MustParse("68f54f91-3efe-4742-b357-a653029a002d"),
				},
				todos: []Todo{
					{
						ID:      uuid.MustParse("8f0aca9a-0ab5-496f-9a47-8c82be9c307c"),
						OwnerID: uuid.MustParse("68f54f91-3efe-4742-b357-a653029a002d"),
					},
					{
						ID:      uuid.MustParse("aba881cd-f3d6-4a0d-b5e6-233854692d8a"),
						OwnerID: uuid.MustParse("012b9c6e-d0ae-4f58-a43e-51bbbbea2bb1"),
					},
				},
			},
			expected: []Todo{
				{
					ID:      uuid.MustParse("8f0aca9a-0ab5-496f-9a47-8c82be9c307c"),
					OwnerID: uuid.MustParse("68f54f91-3efe-4742-b357-a653029a002d"),
				},
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, _ := NewInMemService()
			for _, value := range tc.args.todos {
				s.AddTodo(context.Background(), value)
			}

			actual, err := s.GetTodosOwned(context.Background(), tc.args.owner)
			assert.Equal(t, len(tc.expected), len(actual))
			assert.Equal(t, tc.expected[0].ID, actual[0].ID)
			assert.Equal(t, tc.expected[0].OwnerID, actual[0].OwnerID)
			assert.Equal(t, tc.err, err)
		})
	}
}
