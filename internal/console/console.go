package console

import (
	"fmt"
	"net/http"

	client "github.com/pluralsh/console-client-go"
)

func NewClient(token, url string) *client.Client {
	return client.NewClient(http.DefaultClient, fmt.Sprintf("%s/gql", url), func(req *http.Request) {
		req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
	})
}
