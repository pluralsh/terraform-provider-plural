package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

type fakeProviderServer struct {
	tfprotov6.ProviderServer
}

func (fakeProviderServer) PlanResourceChange(context.Context, *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	return &tfprotov6.PlanResourceChangeResponse{}, nil
}

func TestLegacyTypeSystemServerPlanResourceChange(t *testing.T) {
	tests := []struct {
		typeName string
		expected bool
	}{
		{typeName: "plural_cluster", expected: true},
		{typeName: "plural_service_deployment", expected: false},
	}

	for _, test := range tests {
		t.Run(test.typeName, func(t *testing.T) {
			server := legacyTypeSystemServer{ProviderServer: fakeProviderServer{}}
			resp, err := server.PlanResourceChange(context.Background(), &tfprotov6.PlanResourceChangeRequest{TypeName: test.typeName})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if legacy := resp.UnsafeToUseLegacyTypeSystem; legacy != test.expected { //nolint:staticcheck // tested on purpose
				t.Fatalf("expected UnsafeToUseLegacyTypeSystem %v, got %v", test.expected, legacy)
			}
		})
	}
}
