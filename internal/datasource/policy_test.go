package datasource

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func TestPolicyDataSourceSchemaLookupAndComputedFields(t *testing.T) {
	policyDataSource := &policyDataSource{}
	response := datasource.SchemaResponse{}
	policyDataSource.Schema(context.Background(), datasource.SchemaRequest{}, &response)

	assertDataSourceStringFlags(t, response.Schema.Attributes, "id", false, true, true)
	assertDataSourceStringFlags(t, response.Schema.Attributes, "name", false, true, true)
	assertDataSourceStringFlags(t, response.Schema.Attributes, "type", false, false, true)
	assertDataSourceStringFlags(t, response.Schema.Attributes, "project_id", false, false, true)
}

func TestBindingPolicyDataSourceSchemaUsesIDLookup(t *testing.T) {
	bindingPolicyDataSource := &bindingPolicyDataSource{}
	response := datasource.SchemaResponse{}
	bindingPolicyDataSource.Schema(context.Background(), datasource.SchemaRequest{}, &response)

	assertDataSourceStringFlags(t, response.Schema.Attributes, "id", true, false, false)
	assertDataSourceStringFlags(t, response.Schema.Attributes, "policy_id", false, false, true)
	assertDataSourceStringFlags(t, response.Schema.Attributes, "bind_policy_id", false, false, true)
}

func assertDataSourceStringFlags(t *testing.T, attributes map[string]schema.Attribute, name string, required, optional, computed bool) {
	t.Helper()
	attribute, ok := attributes[name].(schema.StringAttribute)
	if !ok {
		t.Fatalf("expected %q string attribute, got %#v", name, attributes[name])
	}
	if attribute.Required != required || attribute.Optional != optional || attribute.Computed != computed {
		t.Fatalf("unexpected flags for %q: required=%t optional=%t computed=%t", name, attribute.Required, attribute.Optional, attribute.Computed)
	}
}
