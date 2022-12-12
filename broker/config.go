package main

import "github.com/kelseyhightower/envconfig"

type brokerConfig struct {
	BrokerUsername  string `envconfig:"broker_username" required:"true"`
	BrokerPassword  string `envconfig:"broker_password" required:"true"`
	LogLevel        string `envconfig:"log_level" default:"INFO"`
	MetricsEndpoint string `envconfig:"metrics_endpoint" default:"/metrics"`
	RouteSvcURL     string `envconfig:"route_svc_url" required:"true"`
	Port            string `envconfig:"port" default:"3000"`
	DocsURL         string `envconfig:"docsurl" default:"default"`
}

func brokerConfigLoad() (brokerConfig, error) {
	var config brokerConfig
	err := envconfig.Process("", &config)
	if err != nil {
		return brokerConfig{}, err
	}

	return config, nil
}
