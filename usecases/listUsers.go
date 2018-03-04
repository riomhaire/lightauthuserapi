package usecases

func (usecases *Usecases) ListUsers() (names []string) {
	names, _ = usecases.Registry.StorageInteractor.LookupUserNames()
	return names

}
