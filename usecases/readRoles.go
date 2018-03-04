package usecases

func (usecases *Usecases) ReadRoles() (roles []string) {
	roles, _ = usecases.Registry.StorageInteractor.LookupRoleNames()
	return roles

}
