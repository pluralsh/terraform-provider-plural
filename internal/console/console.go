package console

import (
	"context"
	"fmt"
	"net/http"

	consoleclient "github.com/pluralsh/console-client-go"
)

type consoleClient struct {
	ctx    context.Context
	client *consoleclient.Client
	url    string
	token  string
}

type ConsoleClient interface {
	CreateCluster(attributes consoleclient.ClusterAttributes) (*consoleclient.CreateCluster, error)
	GetCluster(clusterId, clusterName *string) (*consoleclient.ClusterFragment, error)
	UpdateCluster(id string, attr consoleclient.ClusterUpdateAttributes) (*consoleclient.UpdateCluster, error)
	DeleteCluster(id string) error
}

func NewClient(token, url string) (ConsoleClient, error) {
	return &consoleClient{
		url:   url,
		token: token,
		client: consoleclient.NewClient(http.DefaultClient, fmt.Sprintf("%s/gql", url), func(req *http.Request) {
			req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
		}),
		ctx: context.Background(),
	}, nil
}
