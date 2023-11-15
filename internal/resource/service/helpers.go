package service

import (
	"context"
	"fmt"

	gqlclient "github.com/pluralsh/console-client-go"
)

func CreateServiceDeployment(ctx context.Context, client *gqlclient.Client, id, handle *string, attrs gqlclient.ServiceDeploymentAttributes) (*gqlclient.ServiceDeploymentFragment, error) {
	if id == nil && handle == nil {
		return nil, fmt.Errorf("could not create ServiceDeployment: id and handle not provided")
	}

	if id != nil {
		res, err := client.CreateServiceDeployment(ctx, *id, attrs)
		return res.CreateServiceDeployment, err
	}

	res, err := client.CreateServiceDeploymentWithHandle(ctx, *handle, attrs)
	return res.CreateServiceDeployment, err
}
