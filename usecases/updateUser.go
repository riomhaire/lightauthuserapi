package usecases

import (
	"errors"

	"github.com/riomhaire/lightauthuserapi/entities"
)

func (usecases *Usecases) UpdateUser(user entities.User) (entities.User, LightAuthError) {
	// Do some validation here before we save - IE User should not exist
	lerror := NewError(NoError, nil)
	_, err := usecases.Registry.StorageInteractor.LookupUserByName(user.Username)

	if err == nil {
		err = usecases.Registry.StorageInteractor.UpdateUser(user)
		if err != nil {
			lerror = NewError(InternalError, err)
		}
	} else {
		lerror = NewError(Unknown, errors.New("No Such User"))
	}

	return user, lerror

}
