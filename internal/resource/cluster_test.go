package resource

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
)

func TestClusterResourceSchemaIsValid(t *testing.T) {
	response := &resource.SchemaResponse{}
	(&clusterResource{}).Schema(context.Background(), resource.SchemaRequest{}, response)
	if response.Diagnostics.HasError() {
		t.Fatalf("unexpected schema diagnostics: %v", response.Diagnostics)
	}

	if diags := response.Schema.ValidateImplementation(context.Background()); diags.HasError() {
		t.Fatalf("invalid schema implementation: %v", diags)
	}
}

func TestClusterResourceProjectIdIsOptionalAndComputed(t *testing.T) {
	response := &resource.SchemaResponse{}
	(&clusterResource{}).Schema(context.Background(), resource.SchemaRequest{}, response)

	assertStringAttributeFlags(t, response.Schema.Attributes, "project_id", false, true, true)
}

func TestClusterAttributesProjectId(t *testing.T) {
	tests := []struct {
		name      string
		projectId types.String
		expected  *string
	}{
		{name: "unknown project is not sent", projectId: types.StringUnknown(), expected: nil},
		{name: "null project is not sent", projectId: types.StringNull(), expected: nil},
		{name: "configured project is sent", projectId: types.StringValue("project"), expected: lo.ToPtr("project")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := cluster{ProjectId: test.projectId}
			var diags diag.Diagnostics
			attributes := c.Attributes(context.Background(), &diags)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(attributes.ProjectID, test.expected) {
				t.Fatalf("expected %v, got %v", test.expected, attributes.ProjectID)
			}
		})
	}
}

func TestClusterFromCreateSetsProjectId(t *testing.T) {
	c := cluster{ProjectId: types.StringUnknown()}
	var diags diag.Diagnostics
	c.FromCreate(&console.CreateCluster{CreateCluster: &console.CreateCluster_CreateCluster{
		ID:      "id",
		Name:    "name",
		Project: &console.TinyProjectFragment{ID: "project"},
	}}, context.Background(), &diags)

	if !c.ProjectId.Equal(types.StringValue("project")) {
		t.Fatalf("expected project ID to be set from the API response, got %v", c.ProjectId)
	}

	c = cluster{ProjectId: types.StringUnknown()}
	c.FromCreate(&console.CreateCluster{CreateCluster: &console.CreateCluster_CreateCluster{ID: "id", Name: "name"}}, context.Background(), &diags)
	if !c.ProjectId.IsNull() {
		t.Fatalf("expected project ID to be null when the API response has no project, got %v", c.ProjectId)
	}
}
