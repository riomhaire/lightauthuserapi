package bootstrap

import (
	"fmt"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/riomhaire/lightauthuserapi/frameworks"
	"github.com/riomhaire/lightauthuserapi/frameworks/api"
	"github.com/riomhaire/lightauthuserapi/usecases"
	"github.com/spf13/cobra"
	"github.com/urfave/negroni"
)

const VERSION = "LightAuthUserAPI Version 1.0"

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

	router := mux.NewRouter()
	negroni := negroni.Classic()
	restAPI.Negroni = negroni

	// Add handlers
	router.HandleFunc("/api/v1/user/metrics", restAPI.HandleStatistics).Methods("GET")
	router.HandleFunc("/metrics", restAPI.HandleStatistics).Methods("GET")
	router.HandleFunc("/api/v1/user/health", restAPI.HandleHealth).Methods("GET")
	router.HandleFunc("/health", restAPI.HandleHealth).Methods("GET")

	router.HandleFunc("/api/v1/user/account/{name}", restAPI.HandleSpecificUser).Methods("GET", "PUT", "DELETE")
	router.HandleFunc("/api/v1/user/account", restAPI.HandleGenericUser).Methods("POST", "GET")

	router.HandleFunc("/api/v1/user/roles", restAPI.HandleReadRoles).Methods("GET")

	// Add Middleware
	negroni.Use(restAPI.Statistics)
	negroni.UseFunc(restAPI.RecordCall)       // Calculates per second/minute rates
	negroni.UseFunc(restAPI.AddWorkerHeader)  // Add which instance
	negroni.UseFunc(restAPI.AddWorkerVersion) // Which version
	negroni.UseFunc(restAPI.AddCoorsHeader)   // Add coors
	negroni.UseHandler(router)

}

func (a *Application) Run() {
	a.registry.Logger.Log("INFO", fmt.Sprintf("Running %s", a.registry.Configuration.Version))
	a.registry.Logger.Log("INFO", a.registry.Configuration.String())
	a.restAPI.Negroni.Run(fmt.Sprintf(":%d", a.registry.Configuration.Port))
}
