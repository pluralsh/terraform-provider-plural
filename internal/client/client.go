package client

import (
	"context"
	"fmt"

	gqlclient "github.com/pluralsh/console/go/client"
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

	res, err := c.CreateServiceDeploymentWithHandle(ctx, *handle, attrs)
	if err != nil {
		return nil, err
	}

	return res.CreateServiceDeployment, err
}

func (c *Client) GetDeploymentSettings(ctx context.Context) (*gqlclient.GetDeploymentSettings, error) {
	res, err := c.ConsoleClient.GetDeploymentSettings(ctx)
	if err == nil && res != nil && res.DeploymentSettings != nil {
		return res, nil
	}

	minimal, minimalErr := c.GetDeploymentSettingsMinimal(ctx)
	if minimalErr != nil {
		if err != nil {
			return nil, err
		}
		return nil, minimalErr
	}

	return &gqlclient.GetDeploymentSettings{
		DeploymentSettings: toDeploymentSettingsFragment(minimal.DeploymentSettings),
	}, nil
}

func toDeploymentSettingsFragment(minimal *gqlclient.DeploymentSettingsMinimalFragment) *gqlclient.DeploymentSettingsFragment {
	if minimal == nil {
		return nil
	}

	return &gqlclient.DeploymentSettingsFragment{
		AgentHelmValues: minimal.AgentHelmValues,
		AgentVsn:        minimal.AgentVsn,
	}
}

func NewClient(client gqlclient.ConsoleClient) *Client {
	return &Client{
		ConsoleClient: client,
	}
}
