package resource

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestPolicyResourceSchemaRequiredAndOptionalFields(t *testing.T) {
	policyResource := &PolicyResource{}
	response := policyResource.SchemaResponse(t)

	assertStringAttributeFlags(t, response.Schema.Attributes, "name", true, false, false)
	assertStringAttributeFlags(t, response.Schema.Attributes, "type", true, false, false)
	assertStringAttributeFlags(t, response.Schema.Attributes, "policy", true, false, false)
	assertStringAttributeFlags(t, response.Schema.Attributes, "description", false, true, false)
	assertStringAttributeFlags(t, response.Schema.Attributes, "project_id", false, true, false)
	assertStringAttributeFlags(t, response.Schema.Attributes, "id", false, false, true)
}

func TestBindingPolicyResourceSchemaRequiredAndOptionalFields(t *testing.T) {
	bindingPolicyResource := &BindingPolicyResource{}
	response := bindingPolicyResource.SchemaResponse(t)

	assertStringAttributeFlags(t, response.Schema.Attributes, "policy_id", true, false, false)
	assertStringAttributeFlags(t, response.Schema.Attributes, "bind_policy_id", true, false, false)
	assertStringAttributeFlags(t, response.Schema.Attributes, "type", true, false, false)
	assertStringAttributeFlags(t, response.Schema.Attributes, "interval", false, true, true)

	matches, ok := response.Schema.Attributes["matches"].(schema.SingleNestedAttribute)
	if !ok || !matches.Optional {
		t.Fatalf("expected optional matches nested attribute, got %#v", response.Schema.Attributes["matches"])
	}
}

func (r *PolicyResource) SchemaResponse(t *testing.T) resource.SchemaResponse {
	t.Helper()
	response := resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, &response)
	return response
}

func (r *BindingPolicyResource) SchemaResponse(t *testing.T) resource.SchemaResponse {
	t.Helper()
	response := resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, &response)
	return response
}

func assertStringAttributeFlags(t *testing.T, attributes map[string]schema.Attribute, name string, required, optional, computed bool) {
	t.Helper()
	attribute, ok := attributes[name].(schema.StringAttribute)
	if !ok {
		t.Fatalf("expected %q string attribute, got %#v", name, attributes[name])
	}
	if attribute.Required != required || attribute.Optional != optional || attribute.Computed != computed {
		t.Fatalf("unexpected flags for %q: required=%t optional=%t computed=%t", name, attribute.Required, attribute.Optional, attribute.Computed)
	}
}
