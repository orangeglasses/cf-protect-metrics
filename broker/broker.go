package main

import (
	"context"
	"fmt"
	"strings"

	brokerapi "github.com/pivotal-cf/brokerapi/v8"
	domain "github.com/pivotal-cf/brokerapi/v8/domain"
)

type broker struct {
	services []brokerapi.Service
	env      brokerConfig
}

func (b *broker) Services(context context.Context) ([]brokerapi.Service, error) {
	return b.services, nil
}

func (b *broker) Provision(context context.Context, instanceID string, details domain.ProvisionDetails, asyncAllowed bool) (domain.ProvisionedServiceSpec, error) {
	return domain.ProvisionedServiceSpec{
		IsAsync:       false,
		AlreadyExists: false,
		DashboardURL:  "",
		OperationData: "",
	}, nil
}

func (b *broker) GetInstance(ctx context.Context, instanceID string, details domain.FetchInstanceDetails) (domain.GetInstanceDetailsSpec, error) {
	return domain.GetInstanceDetailsSpec{}, fmt.Errorf("Instances are not retrievable")
}

func (b *broker) Deprovision(context context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (domain.DeprovisionServiceSpec, error) {
	return domain.DeprovisionServiceSpec{
		IsAsync:       false,
		OperationData: "",
	}, nil
}

func (b *broker) Bind(context context.Context, instanceID, bindingID string, details domain.BindDetails, asyncAllowed bool) (domain.Binding, error) {
	route := details.BindResource.Route
	if !strings.HasSuffix(route, b.env.MetricsEndpoint) {
		return domain.Binding{}, fmt.Errorf("Can only bind to %v routes", b.env.MetricsEndpoint)
	}

	return domain.Binding{
		RouteServiceURL: b.env.RouteSvcURL,
	}, nil
}

func (b *broker) GetBinding(ctx context.Context, instanceID, bindingID string, details domain.FetchBindingDetails) (domain.GetBindingSpec, error) {
	return domain.GetBindingSpec{}, fmt.Errorf("Bindings are not retrievable")
}

func (b *broker) Unbind(context context.Context, instanceID, bindingID string, details domain.UnbindDetails, asyncAllowed bool) (domain.UnbindSpec, error) {
	return domain.UnbindSpec{
		IsAsync:       false,
		OperationData: "",
	}, nil
}

func (b *broker) Update(context context.Context, instanceID string, details domain.UpdateDetails, asyncAllowed bool) (domain.UpdateServiceSpec, error) {

	return domain.UpdateServiceSpec{}, nil
}

func (b *broker) LastOperation(context context.Context, instanceID string, details domain.PollDetails) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}

func (b *broker) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details domain.PollDetails) (domain.LastOperation, error) {
	return domain.LastOperation{}, nil
}
