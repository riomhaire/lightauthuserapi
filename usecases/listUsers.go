package usecases

// List all users, page of users or users matching
// search parameters.
func (usecases *Usecases) ListUsers(search string, page int, pageSize int) (names []string) {
	names, _ = usecases.Registry.StorageInteractor.LookupUserNames(search, page, pageSize)
	return
}
