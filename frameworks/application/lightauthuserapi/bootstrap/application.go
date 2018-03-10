package bootstrap

import (
	"fmt"
	"strconv"

	"github.com/riomhaire/lightauthuserapi/frameworks"
	"github.com/riomhaire/lightauthuserapi/frameworks/api"
	"github.com/riomhaire/lightauthuserapi/usecases"
	"github.com/spf13/cobra"
)

const VERSION = "LightAuthUserAPI Version 1.2"

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
	registry := usecases.Registry{}
	a.registry = &registry
	registry.Configuration = configuration
	registry.Logger = logger
	database := frameworks.NewCSVReaderDatabaseInteractor(&registry)

	registry.StorageInteractor = database
	registry.Usecases = usecases.Usecases{&registry}

	// Create API
	restAPI := api.NewRestAPI(&registry)
	a.restAPI = &restAPI

}

func (a *Application) Run() {
	a.registry.Logger.Log("INFO", fmt.Sprintf("Running %s", a.registry.Configuration.Version))
	a.registry.Logger.Log("INFO", a.registry.Configuration.String())
	a.restAPI.Negroni.Run(fmt.Sprintf(":%d", a.registry.Configuration.Port))
}
