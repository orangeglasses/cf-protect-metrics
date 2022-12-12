package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pivotal-cf/brokerapi/v7"
)

func CatalogLoad(catalogFilePath string, env brokerConfig) ([]brokerapi.Service, error) {
	var services []brokerapi.Service

	inBuf, err := ioutil.ReadFile(catalogFilePath)
	if err != nil {
		return []brokerapi.Service{}, err
	}

	err = json.Unmarshal(inBuf, &services)
	if err != nil {
		return []brokerapi.Service{}, err
	}

	services[0].Metadata.DocumentationUrl = env.DocsURL
	return services, nil
}
