package client

import (
	"context"
	"fmt"

	gqlclient "github.com/pluralsh/console/go/client"
)

type getPolicyResponse struct {
	Policy *gqlclient.Policy `json:"policy"`
}

type createPolicyResponse struct {
	CreatePolicy *gqlclient.Policy `json:"createPolicy"`
}

type updatePolicyResponse struct {
	UpdatePolicy *gqlclient.Policy `json:"updatePolicy"`
}

type deletePolicyResponse struct {
	DeletePolicy *gqlclient.Policy `json:"deletePolicy"`
}

type getBindingPolicyResponse struct {
	BindingPolicy *gqlclient.BindingPolicy `json:"bindingPolicy"`
}

type createBindingPolicyResponse struct {
	CreateBindingPolicy *gqlclient.BindingPolicy `json:"createBindingPolicy"`
}

type updateBindingPolicyResponse struct {
	UpdateBindingPolicy *gqlclient.BindingPolicy `json:"updateBindingPolicy"`
}

type deleteBindingPolicyResponse struct {
	DeleteBindingPolicy *gqlclient.BindingPolicy `json:"deleteBindingPolicy"`
}

const policyFragment = `fragment TerraformPolicyFragment on Policy {
  id
  name
  type
  description
  policy
  project {
    id
  }
}`

const bindingPolicyFragment = `fragment TerraformBindingPolicyFragment on BindingPolicy {
  id
  type
  interval
  matches {
    workbench {
      regexes
    }
  }
  policy {
    id
  }
  bindPolicy {
    id
  }
}`

const getPolicyDocument = `query TerraformGetPolicy($id: ID, $name: String) {
  policy(id: $id, name: $name) {
    ...TerraformPolicyFragment
  }
}
` + policyFragment

const createPolicyDocument = `mutation TerraformCreatePolicy($attributes: PolicyAttributes!) {
  createPolicy(attributes: $attributes) {
    ...TerraformPolicyFragment
  }
}
` + policyFragment

const updatePolicyDocument = `mutation TerraformUpdatePolicy($id: ID!, $attributes: PolicyAttributes!) {
  updatePolicy(id: $id, attributes: $attributes) {
    ...TerraformPolicyFragment
  }
}
` + policyFragment

const deletePolicyDocument = `mutation TerraformDeletePolicy($id: ID!) {
  deletePolicy(id: $id) {
    ...TerraformPolicyFragment
  }
}
` + policyFragment

const getBindingPolicyDocument = `query TerraformGetBindingPolicy($id: ID!) {
  bindingPolicy(id: $id) {
    ...TerraformBindingPolicyFragment
  }
}
` + bindingPolicyFragment

const createBindingPolicyDocument = `mutation TerraformCreateBindingPolicy($attributes: BindingPolicyAttributes!) {
  createBindingPolicy(attributes: $attributes) {
    ...TerraformBindingPolicyFragment
  }
}
` + bindingPolicyFragment

const updateBindingPolicyDocument = `mutation TerraformUpdateBindingPolicy($id: ID!, $attributes: BindingPolicyUpdateAttributes!) {
  updateBindingPolicy(id: $id, attributes: $attributes) {
    ...TerraformBindingPolicyFragment
  }
}
` + bindingPolicyFragment

const deleteBindingPolicyDocument = `mutation TerraformDeleteBindingPolicy($id: ID!) {
  deleteBindingPolicy(id: $id) {
    ...TerraformBindingPolicyFragment
  }
}
` + bindingPolicyFragment

func (c *Client) graphqlClient() (*gqlclient.Client, error) {
	client, ok := c.ConsoleClient.(*gqlclient.Client)
	if !ok {
		return nil, fmt.Errorf("policy operations require *client.Client, got %T", c.ConsoleClient)
	}

	return client, nil
}

func (c *Client) GetPolicy(ctx context.Context, id, name *string) (*getPolicyResponse, error) {
	client, err := c.graphqlClient()
	if err != nil {
		return nil, err
	}

	response := new(getPolicyResponse)
	if err := client.Client.Post(ctx, "TerraformGetPolicy", getPolicyDocument, response, map[string]any{"id": id, "name": name}); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) CreatePolicy(ctx context.Context, attributes gqlclient.PolicyAttributes) (*createPolicyResponse, error) {
	client, err := c.graphqlClient()
	if err != nil {
		return nil, err
	}

	response := new(createPolicyResponse)
	if err := client.Client.Post(ctx, "TerraformCreatePolicy", createPolicyDocument, response, map[string]any{"attributes": attributes}); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) UpdatePolicy(ctx context.Context, id string, attributes gqlclient.PolicyAttributes) (*updatePolicyResponse, error) {
	client, err := c.graphqlClient()
	if err != nil {
		return nil, err
	}

	response := new(updatePolicyResponse)
	if err := client.Client.Post(ctx, "TerraformUpdatePolicy", updatePolicyDocument, response, map[string]any{"id": id, "attributes": attributes}); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) DeletePolicy(ctx context.Context, id string) (*deletePolicyResponse, error) {
	client, err := c.graphqlClient()
	if err != nil {
		return nil, err
	}

	response := new(deletePolicyResponse)
	if err := client.Client.Post(ctx, "TerraformDeletePolicy", deletePolicyDocument, response, map[string]any{"id": id}); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) GetBindingPolicy(ctx context.Context, id string) (*getBindingPolicyResponse, error) {
	client, err := c.graphqlClient()
	if err != nil {
		return nil, err
	}

	response := new(getBindingPolicyResponse)
	if err := client.Client.Post(ctx, "TerraformGetBindingPolicy", getBindingPolicyDocument, response, map[string]any{"id": id}); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) CreateBindingPolicy(ctx context.Context, attributes gqlclient.BindingPolicyAttributes) (*createBindingPolicyResponse, error) {
	client, err := c.graphqlClient()
	if err != nil {
		return nil, err
	}

	response := new(createBindingPolicyResponse)
	if err := client.Client.Post(ctx, "TerraformCreateBindingPolicy", createBindingPolicyDocument, response, map[string]any{"attributes": attributes}); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) UpdateBindingPolicy(ctx context.Context, id string, attributes gqlclient.BindingPolicyUpdateAttributes) (*updateBindingPolicyResponse, error) {
	client, err := c.graphqlClient()
	if err != nil {
		return nil, err
	}

	response := new(updateBindingPolicyResponse)
	if err := client.Client.Post(ctx, "TerraformUpdateBindingPolicy", updateBindingPolicyDocument, response, map[string]any{"id": id, "attributes": attributes}); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) DeleteBindingPolicy(ctx context.Context, id string) (*deleteBindingPolicyResponse, error) {
	client, err := c.graphqlClient()
	if err != nil {
		return nil, err
	}

	response := new(deleteBindingPolicyResponse)
	if err := client.Client.Post(ctx, "TerraformDeleteBindingPolicy", deleteBindingPolicyDocument, response, map[string]any{"id": id}); err != nil {
		return nil, err
	}

	return response, nil
}
