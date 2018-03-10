package test

import (
	"errors"

	"github.com/riomhaire/lightauthuserapi/entities"
	"github.com/riomhaire/lightauthuserapi/usecases"
)

// This is a test implementation for test purposes
type InMemoryDBInteractor struct {
	userdb map[string]entities.User
	roledb []entities.Role
	logger usecases.Logger
}

func NewInMemoryDBInteractor(logger usecases.Logger, userdb map[string]entities.User, roledb []entities.Role) *InMemoryDBInteractor {
	d := InMemoryDBInteractor{}
	d.userdb = userdb
	d.roledb = roledb
	d.logger = logger

	return &d
}

func (db *InMemoryDBInteractor) LookupUserByName(username string) (entities.User, error) {
	if val, ok := db.userdb[username]; ok {
		return val, nil
	} else {
		return entities.User{}, errors.New("Unknown user")
	}
}

func (db *InMemoryDBInteractor) CreateUser(user entities.User) error {
	if _, ok := db.userdb[user.Username]; ok {
		return errors.New("User exists")
	}
	db.userdb[user.Username] = user
	return nil
}

func (db *InMemoryDBInteractor) LookupUserNames(search string, page int, pageSize int) ([]string, error) {
	var s []string
	for k := range db.userdb {
		s = append(s, k)
	}

	return s, nil
}

func (db *InMemoryDBInteractor) UpdateUser(user entities.User) error {
	if _, ok := db.userdb[user.Username]; ok {
		db.userdb[user.Username] = user
	} else {
		return errors.New("User Does Not Exists")
	}
	return nil
}

func (db *InMemoryDBInteractor) DeleteUser(user string) error {
	if _, ok := db.userdb[user]; ok {
		delete(db.userdb, user)
	} else {
		return errors.New("User Does Not Exists")
	}
	return nil
}

func (db *InMemoryDBInteractor) LookupRoleNames() ([]string, error) {
	var roles []string
	for _, r := range db.roledb {
		roles = append(roles, r.Name)
	}
	return roles, nil
}
