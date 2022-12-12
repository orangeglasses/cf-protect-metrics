package main

import (
	"fmt"
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	brokerapi "github.com/pivotal-cf/brokerapi/v8"
)

func main() {
	var logLevels = map[string]lager.LogLevel{
		"DEBUG": lager.DEBUG,
		"INFO":  lager.INFO,
		"ERROR": lager.ERROR,
		"FATAL": lager.FATAL,
	}

	config, err := brokerConfigLoad()
	if err != nil {
		panic(err)
	}

	brokerCredentials := brokerapi.BrokerCredentials{
		Username: config.BrokerUsername,
		Password: config.BrokerPassword,
	}

	services, err := CatalogLoad("./catalog.json", config)
	if err != nil {
		panic(err)
	}

	logger := lager.NewLogger("cf-protectmetrics-broker")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, logLevels[config.LogLevel]))

	serviceBroker := &broker{
		services: services,
		env:      config,
	}

	brokerHandler := brokerapi.New(serviceBroker, logger, brokerCredentials)
	fmt.Println("Starting service")
	http.Handle("/", brokerHandler)
	http.ListenAndServe(":"+config.Port, nil)
}
