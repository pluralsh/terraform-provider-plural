package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gqlgo/gqlgenc/clientv2"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func TestIsNotFound(t *testing.T) {
	notFound := &clientv2.ErrorResponse{GqlErrors: &gqlerror.List{{Message: string(ErrorNotFound)}}}
	otherError := &clientv2.ErrorResponse{GqlErrors: &gqlerror.List{{Message: "forbidden"}}}

	cases := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "not found", err: notFound, want: true},
		{name: "wrapped not found", err: fmt.Errorf("get monitor: %w", notFound), want: true},
		{name: "other graphql error", err: otherError, want: false},
		{name: "network error", err: &clientv2.ErrorResponse{NetworkError: &clientv2.HTTPError{Code: 500}}, want: false},
		{name: "other error", err: errors.New(string(ErrorNotFound)), want: false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := IsNotFound(c.err); got != c.want {
				t.Fatalf("expected IsNotFound to return %t, got %t", c.want, got)
			}
		})
	}
}

// Goes through the real console client so a gqlgenc module mismatch between it and errors.go fails here.
func TestIsNotFoundFromConsoleClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"data":null,"errors":[{"message":%q,"path":["detachCluster"]}]}`, ErrorNotFound)
	}))
	defer server.Close()

	_, err := gqlclient.NewClient(server.Client(), server.URL, nil).DetachCluster(context.Background(), "id")
	if err == nil {
		t.Fatal("expected an error from the console client")
	}

	if !IsNotFound(err) {
		t.Fatalf("expected IsNotFound to match console client error %T: %v", err, err)
	}
}
