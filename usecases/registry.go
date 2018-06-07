package usecases

import (
	"fmt"

	"github.com/riomhaire/lightauthuserapi/frameworks/serviceregistry"
)

// Configuration containing data from the environment which is used to define program behaviour
type Configuration struct {
	Application string
	Version     string
	RoleStore   string
	UserStore   string
	Port        int
	APIKey      string
	Host        string
	Consul      bool
	ConsulHost  string
	ConsulId    string // ID of this client
}

type Registry struct {
	Configuration           Configuration
	Logger                  Logger
	StorageInteractor       StorageInteractor
	Usecases                Usecases
	ExternalServiceRegistry serviceregistry.ServiceRegistry
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
