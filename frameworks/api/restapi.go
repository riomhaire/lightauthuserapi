package api

import (
	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	"github.com/riomhaire/lightauthuserapi/usecases"
	"github.com/thoas/stats"
	"github.com/urfave/negroni"
)

var bearerPrefix = "bearer "

type RestAPI struct {
	Registry         *usecases.Registry
	Statistics       *stats.Stats
	Negroni          *negroni.Negroni
	Producer         sarama.SyncProducer
	KafkaInitialized bool
	IPAddress        string
	MetricsRegistry  MetricsRegistry
}

func NewRestAPI(registry *usecases.Registry) RestAPI {
	api := RestAPI{}
	api.Registry = registry
	api.Statistics = stats.New()
	api.MetricsRegistry = MetricsRegistry{}
	router := mux.NewRouter()
	negroni := negroni.Classic()
	api.Negroni = negroni

	// Add handlers
	router.HandleFunc("/api/v1/user/metrics", api.HandleStatistics).Methods("GET")
	router.HandleFunc("/metrics", api.HandleStatistics).Methods("GET")
	router.HandleFunc("/api/v1/user/health", api.HandleHealth).Methods("GET")
	router.HandleFunc("/health", api.HandleHealth).Methods("GET")

	router.HandleFunc("/api/v1/user/account/{name}", api.HandleSpecificUser).Methods("GET", "PUT", "DELETE")
	router.HandleFunc("/api/v1/user/account", api.HandleGenericUser).Methods("POST", "GET")

	router.HandleFunc("/api/v1/user/roles", api.HandleReadRoles).Methods("GET")

	// This is for options call
	router.HandleFunc("/api/v1/user/metrics", api.HandleOptions).Methods("OPTIONS")
	router.HandleFunc("/metrics", api.HandleOptions).Methods("OPTIONS")
	router.HandleFunc("/api/v1/user/health", api.HandleOptions).Methods("OPTIONS")
	router.HandleFunc("/health", api.HandleOptions).Methods("OPTIONS")

	router.HandleFunc("/api/v1/user/account/{name}", api.HandleOptions).Methods("OPTIONS")
	router.HandleFunc("/api/v1/user/account", api.HandleOptions).Methods("OPTIONS")

	router.HandleFunc("/api/v1/user/roles", api.HandleOptions).Methods("OPTIONS")

	// Add Middleware
	negroni.Use(api.Statistics)
	negroni.UseFunc(api.RecordCall)       // Calculates per second/minute rates
	negroni.UseFunc(api.AddWorkerHeader)  // Add which instance
	negroni.UseFunc(api.AddWorkerVersion) // Which version
	negroni.UseFunc(api.AddCoorsHeader)   // Add coors
	negroni.UseHandler(router)
	return api
}
