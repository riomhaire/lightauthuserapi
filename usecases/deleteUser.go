package usecases

func (usecases *Usecases) DeleteUser(user string) LightAuthError {
	// Do some validation here before we save - IE User should not exist
	lerror := NewError(NoError, nil)
	_, err := usecases.Registry.StorageInteractor.LookupUserByName(user)

	if err == nil {
		err = usecases.Registry.StorageInteractor.DeleteUser(user)
		if err != nil {
			lerror = NewError(InternalError, err)
		}

	} else {
		lerror = NewError(Unknown, err)
	}

	return lerror

}
