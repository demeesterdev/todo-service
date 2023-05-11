package authorization

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/demeesterdev/todo-service/internal/argon2id"
)

// user presents a single user object as stored in the database
// ID should be globally unique and a valid uuid
type storedUser struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;primarykey"`
	Username     string    `gorm:"unique;index"`
	PasswordHash string
}

// TableName overrides the table name used by User to `profiles` (GORM specifics)
func (storedUser) TableName() string {
	return "users"
}

// Before create is a GORM hook
// It makes shure a user has a valid uuid before creation (GORM specifics)
func (u *storedUser) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

func (u storedUser) ToUser() (U User) {
	U.ID = u.ID
	U.Username = u.Username
	return U
}

func (u *storedUser) ComparePassword(password string) (match bool, err error) {
	return argon2id.ComparePasswordAndHash(password, u.PasswordHash)
}

func newStoredUser(U User, passwordHashParameters argon2id.Params) (u storedUser, err error) {
	var passwordHash string
	if U.Password != "" {
		passwordHash, err = argon2id.HashPassword(U.Password, passwordHashParameters)
	}
	return storedUser{
		ID:           U.ID,
		Username:     U.Username,
		PasswordHash: passwordHash,
	}, err
}

type dbSvc struct {
	hashParams argon2id.Params
	db         *gorm.DB
}

// NewService creates a new user service based on a sqlite database with a target file
func NewDBService(dbconnection gorm.Dialector, passwordHashParameters argon2id.Params) (Service, error) {
	db, err := gorm.Open(dbconnection, &gorm.Config{})
	db.AutoMigrate(&storedUser{})
	if err != nil {
		return &dbSvc{}, err
	}

	return &dbSvc{
		db:         db,
		hashParams: passwordHashParameters,
	}, nil
}

// NewSqliteDBService creates a new user service based on a sqlite database with a target file
func NewSqliteDBService(target string, passwordHashParameters argon2id.Params) (Service, error) {
	return NewDBService(sqlite.Open(target), passwordHashParameters)
}

func NewInMemService(passwordHashParameters argon2id.Params) (Service, error) {
	return NewSqliteDBService(":memory:", passwordHashParameters)
}

func (s *dbSvc) AddUser(ctx context.Context, u User) (User, error) {
	if u.Username == "" || u.Password == "" {
		return User{}, ErrInvalidUserObject
	}

	newUser, err := newStoredUser(u, s.hashParams)
	if err != nil {
		return User{}, err
	}

	result := s.db.Create(&newUser)
	if result.Error != nil {
		return User{}, result.Error
	}

	return newUser.ToUser(), nil
}

func (s *dbSvc) GetUser(ctx context.Context, id uuid.UUID) (User, error) {

	// get first user where storedUser.ID = id
	var u storedUser
	result := s.db.Model(&storedUser{ID: id}).First(&u)

	switch result.Error {
	case gorm.ErrRecordNotFound:
		return User{}, ErrNotFound
	default:
		return u.ToUser(), result.Error
	}
}

func (s *dbSvc) FindUser(ctx context.Context, username string) (User, error) {

	// get first user where storedUser.Username = username
	var u storedUser
	result := s.db.Where(&storedUser{Username: username}).First(&u)
	switch result.Error {
	case gorm.ErrRecordNotFound:
		return User{}, ErrNotFound
	default:
		return u.ToUser(), result.Error
	}
}

func (s *dbSvc) UpdateUser(ctx context.Context, id uuid.UUID, U User) (User, error) {
	if U.ID == uuid.Nil {
		U.ID = id
	}

	if U.ID != id {
		return User{}, ErrInconsistentIDs
	}

	// get first user where storedUser.ID = id
	var u storedUser
	result := s.db.Model(&storedUser{ID: id}).First(&u)
	if result.Error != nil {
		return User{}, result.Error
	}

	u, err := newStoredUser(U, s.hashParams)
	if result.Error != nil {
		return User{}, err
	}

	result = s.db.Model(&u).Updates(u)
	if result.Error != nil {
		return User{}, result.Error
	}

	s.db.Model(&storedUser{ID: id}).First(&u)
	return u.ToUser(), nil
}

func (s *dbSvc) AuthenticateUser(ctx context.Context, U User) (User, error) {
	if U.Password == "" && U.Username == "" {
		return User{}, ErrInvalidUserObject
	}

	// get first user where storedUser.Username = username
	var u storedUser
	result := s.db.Model(&storedUser{Username: U.Username}).First(&u)
	if result.Error == gorm.ErrRecordNotFound {
		return User{}, ErrNotFound
	}
	if result.Error != nil {
		return User{}, result.Error
	}

	match, err := u.ComparePassword(U.Password)
	if err != nil {
		return User{}, err
	}

	if match {
		return u.ToUser(), nil
	}
	return User{}, ErrAuthenticationFailed
}

func (s *dbSvc) DeleteUser(ctx context.Context, id uuid.UUID) error {
	u := storedUser{}
	result := s.db.Delete(&u, "id = ?", id.String())

	return result.Error

}
func (s *dbSvc) GetUsers(ctx context.Context) ([]User, error) {

	var users []storedUser
	result := s.db.Find(&users)
	if result.Error != nil {
		return []User{}, result.Error
	}

	Users := make([]User, len(users))
	for i := range users {
		Users[i] = users[i].ToUser()
	}

	return Users, nil
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
