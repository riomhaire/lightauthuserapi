package api

import (
	"github.com/Shopify/sarama"
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

	return api
}
