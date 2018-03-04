package usecases

import "github.com/riomhaire/lightauthuserapi/entities"

const (
	NoError        = 0
	AlreadyExists  = 1
	NotImplemented = 2
	Unknown        = 3
	Invalid        = 4
	NotAuthorized  = 5
	InternalError  = 6
)

type LightAuthError struct {
	Code  int
	Error error
}

func NewError(code int, err error) LightAuthError {
	return LightAuthError{code, err}
}

type Logger interface {
	Log(level, message string)
}

type StorageInteractor interface {
	LookupUserByName(username string) (entities.User, error)
	LookupUserNames() ([]string, error)
	CreateUser(user entities.User) error
	UpdateUser(user entities.User) error
	DeleteUser(user string) error

	LookupRoleNames() ([]string, error)
}

type Usecases struct {
	Registry *Registry
}
