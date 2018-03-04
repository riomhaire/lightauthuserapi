package usecases

import "fmt"

// Configuration containing data from the environment which is used to define program behaviour
type Configuration struct {
	Application string
	Version     string
	RoleStore   string
	UserStore   string
	Port        int
	APIKey      string
}

type Registry struct {
	Configuration     Configuration
	Logger            Logger
	StorageInteractor StorageInteractor
	Usecases          Usecases
}

func (c *Configuration) String() string {
	return fmt.Sprintf("\nCONFIGURATION\n\t%15s : '%v'\n\t%15s : '%v'\n\t%15s : '%v'\n\t%15s : '%v'\n\t%15s : '%v'\n\t%15s : '%v'\n",
		"Application",
		c.Application,
		"APIKey",
		c.APIKey,
		"UserStore",
		c.UserStore,
		"RoleStore",
		c.RoleStore,
		"Version",
		c.Version,
		"Port",
		c.Port,
	)
}
