package client

import (
	"context"
	"fmt"

	gqlclient "github.com/pluralsh/console-client-go"
	"github.com/samber/lo"
)

type Client struct {
	gqlclient.ConsoleClient
}

func (c *Client) CreateServiceDeployment(ctx context.Context, id, handle *string, attrs gqlclient.ServiceDeploymentAttributes) (*gqlclient.ServiceDeploymentExtended, error) {
	if len(lo.FromPtr(id)) == 0 && len(lo.FromPtr(handle)) == 0 {
		return nil, fmt.Errorf("could not create service deployment: id or handle not provided")
	}

	if len(lo.FromPtr(id)) > 0 {
		res, err := c.ConsoleClient.CreateServiceDeployment(ctx, *id, attrs)
		if err != nil {
			return nil, err
		}

		return res.CreateServiceDeployment, err
	}

	res, err := c.ConsoleClient.CreateServiceDeploymentWithHandle(ctx, *handle, attrs)
	if err != nil {
		return nil, err
	}

	return res.CreateServiceDeployment, err
}

func NewClient(client gqlclient.ConsoleClient) *Client {
	return &Client{
		ConsoleClient: client,
	}
}
