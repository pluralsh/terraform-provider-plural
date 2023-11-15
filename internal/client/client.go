package client

import (
	"context"
	"fmt"

	gqlclient "github.com/pluralsh/console-client-go"
)

type Client struct {
	*gqlclient.Client
}

func (c *Client) CreateServiceDeployment(ctx context.Context, id, handle *string, attrs gqlclient.ServiceDeploymentAttributes) (*gqlclient.ServiceDeploymentFragment, error) {
	if id == nil && handle == nil {
		return nil, fmt.Errorf("could not create ServiceDeployment: id and handle not provided")
	}

	if id != nil {
		res, err := c.Client.CreateServiceDeployment(ctx, *id, attrs)
		return res.CreateServiceDeployment, err
	}

	res, err := c.Client.CreateServiceDeploymentWithHandle(ctx, *handle, attrs)
	return res.CreateServiceDeployment, err
}

func NewClient(client *gqlclient.Client) *Client {
	return &Client{
		Client: client,
	}
}
