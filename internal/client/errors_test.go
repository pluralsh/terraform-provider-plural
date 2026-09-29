package client

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gqlgo/gqlgenc/clientv2"
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
