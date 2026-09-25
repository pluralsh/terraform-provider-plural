package datasource

import (
	"context"
	"testing"

	"terraform-provider-plural/internal/client"

	"github.com/gqlgo/gqlgenc/clientv2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	gqlclient "github.com/pluralsh/console/go/client"
)

type serviceDeploymentClient struct {
	gqlclient.ConsoleClient
	cluster, name string
}

func (c *serviceDeploymentClient) GetServiceDeploymentTinyByHandle(_ context.Context, cluster, name string, _ ...clientv2.RequestInterceptor) (*gqlclient.GetServiceDeploymentTinyByHandle, error) {
	c.cluster, c.name = cluster, name
	if cluster != "mgmt" || name != "console" {
		return &gqlclient.GetServiceDeploymentTinyByHandle{}, nil
	}

	return &gqlclient.GetServiceDeploymentTinyByHandle{
		ServiceDeployment: &gqlclient.GetServiceDeploymentTinyByHandle_ServiceDeployment{ID: "service-1", Name: name},
	}, nil
}

func TestServiceDeploymentDataSourceSchema(t *testing.T) {
	response := datasource.SchemaResponse{}
	(&ServiceDeploymentDataSource{}).Schema(context.Background(), datasource.SchemaRequest{}, &response)

	assertDataSourceStringFlags(t, response.Schema.Attributes, "cluster", true, false, false)
	assertDataSourceStringFlags(t, response.Schema.Attributes, "name", true, false, false)
	assertDataSourceStringFlags(t, response.Schema.Attributes, "id", false, false, true)
}

func TestServiceDeploymentDataSourceRead(t *testing.T) {
	ctx := context.Background()
	api := &serviceDeploymentClient{}
	r := &ServiceDeploymentDataSource{client: client.NewClient(api)}
	s := datasource.SchemaResponse{}
	r.Schema(ctx, datasource.SchemaRequest{}, &s)

	read := func(cluster, name string) datasource.ReadResponse {
		t.Helper()
		config := tfsdk.State{Schema: s.Schema, Raw: tftypes.NewValue(s.Schema.Type().TerraformType(ctx), nil)}
		if d := config.Set(ctx, &serviceDeployment{Id: types.StringNull(), Cluster: types.StringValue(cluster), Name: types.StringValue(name)}); d.HasError() {
			t.Fatal(d)
		}

		response := datasource.ReadResponse{State: tfsdk.State{Schema: s.Schema, Raw: tftypes.NewValue(s.Schema.Type().TerraformType(ctx), nil)}}
		r.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: s.Schema, Raw: config.Raw}}, &response)
		return response
	}

	found := read("mgmt", "console")
	if found.Diagnostics.HasError() {
		t.Fatal(found.Diagnostics)
	}
	if api.cluster != "mgmt" || api.name != "console" {
		t.Fatalf("expected the service to be looked up by cluster handle and name, got %q and %q", api.cluster, api.name)
	}
	var service serviceDeployment
	if d := found.State.Get(ctx, &service); d.HasError() {
		t.Fatal(d)
	}
	if service.Id.ValueString() != "service-1" || service.Cluster.ValueString() != "mgmt" || service.Name.ValueString() != "console" {
		t.Fatalf("expected the found service to be stored, got %+v", service)
	}

	if missing := read("mgmt", "missing"); !missing.Diagnostics.HasError() {
		t.Fatal("expected an error for a service that does not exist")
	}
}
