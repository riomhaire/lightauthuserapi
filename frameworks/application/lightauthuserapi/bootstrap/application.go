package bootstrap

import (
	"fmt"
	"os"
	"strconv"

	"github.com/riomhaire/lightauthuserapi/frameworks"
	"github.com/riomhaire/lightauthuserapi/frameworks/api"
	"github.com/riomhaire/lightauthuserapi/frameworks/serviceregistry/consulagent"
	"github.com/riomhaire/lightauthuserapi/frameworks/serviceregistry/defaultserviceregistry"
	"github.com/riomhaire/lightauthuserapi/usecases"
	"github.com/spf13/cobra"
)

const VERSION = "LightAuthUserAPI Version 1.3.1"

type Application struct {
	registry *usecases.Registry
	restAPI  *api.RestAPI
}

func (a *Application) Initialize(cmd *cobra.Command, args []string) {
	logger := frameworks.ConsoleLogger{}

	logger.Log("INFO", "Initializing")
	// Create Configuration
	configuration := usecases.Configuration{}

	// Populate it
	configuration.Application = "UserAPI"
	configuration.Version = VERSION
	configuration.Port, _ = strconv.Atoi(cmd.Flag("port").Value.String())
	configuration.UserStore = cmd.Flag("usersFile").Value.String()
	configuration.RoleStore = cmd.Flag("rolesFile").Value.String()
	configuration.APIKey = cmd.Flag("key").Value.String()
	hostname, _ := os.Hostname()
	configuration.Host = hostname
	configuration.Consul, _ = strconv.ParseBool(cmd.Flag("consul").Value.String())
	configuration.ConsulHost = cmd.Flag("consulHost").Value.String()

	registry := usecases.Registry{}
	a.registry = &registry
	registry.Configuration = configuration
	registry.Logger = logger
	database := frameworks.NewCSVReaderDatabaseInteractor(&registry)

	registry.StorageInteractor = database
	registry.Usecases = usecases.Usecases{&registry}

	// Do we need external registry
	if configuration.Consul {
		registry.ExternalServiceRegistry = consulagent.NewConsulServiceRegistry(&registry, "/api/v1/user", "/api/v1/user/health")

	} else {
		registry.ExternalServiceRegistry = defaultserviceregistry.NewDefaultServiceRegistry(&registry)
	}

	// Create API
	restAPI := api.NewRestAPI(&registry)
	a.restAPI = &restAPI

}

func (a *Application) Run() {
	a.registry.Logger.Log("INFO", fmt.Sprintf("Running %s", a.registry.Configuration.Version))
	a.registry.Logger.Log("INFO", a.registry.Configuration.String())

	// Register with external service if required ... default does nothing
	a.registry.ExternalServiceRegistry.Register()

	a.restAPI.Negroni.Run(fmt.Sprintf(":%d", a.registry.Configuration.Port))
}

func (a *Application) Stop() {
	a.registry.Logger.Log("INFO", "Shutting Down REST API")
	a.registry.ExternalServiceRegistry.Deregister()
}
