package usecases

import "github.com/riomhaire/lightauthuserapi/entities"

func (usecases *Usecases) ReadUser(name string) (entities.User, LightAuthError) {
	user, err := usecases.Registry.StorageInteractor.LookupUserByName(name)
	lerror := NewError(NoError, nil)
	if err != nil {
		lerror = NewError(Unknown, err)
	}

	return user, lerror

}
