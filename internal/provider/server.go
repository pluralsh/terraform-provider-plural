package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// NewServer returns the provider server with plural_cluster plans marked as produced by the legacy type system.
func NewServer(version string) func() tfprotov6.ProviderServer {
	return func() tfprotov6.ProviderServer {
		return legacyTypeSystemServer{ProviderServer: providerserver.NewProtocol6(New(version)())()}
	}
}

// legacyTypeSystemServer makes Terraform log a warning instead of failing when the plural_cluster plan does not
// match the config. plural_cluster plans the deprecated kubeconfig as null to keep it out of the plan, as its
// values can change during apply and fail it with "Provider produced inconsistent final plan".
// It can be removed together with plural_cluster kubeconfig.
type legacyTypeSystemServer struct {
	tfprotov6.ProviderServer
}

func (s legacyTypeSystemServer) PlanResourceChange(ctx context.Context, req *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	resp, err := s.ProviderServer.PlanResourceChange(ctx, req)
	if resp != nil && req.TypeName == "plural_cluster" {
		resp.UnsafeToUseLegacyTypeSystem = true //nolint:staticcheck // deliberate, see legacyTypeSystemServer
	}

	return resp, err
}
