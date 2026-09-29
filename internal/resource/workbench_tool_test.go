package resource

import (
	"context"
	"testing"

	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
)

func TestWorkbenchToolResourceSchemaApprovalDefaultsToFalse(t *testing.T) {
	workbenchToolResource := &WorkbenchToolResource{}
	response := resource.SchemaResponse{}
	workbenchToolResource.Schema(context.Background(), resource.SchemaRequest{}, &response)

	attribute, ok := response.Schema.Attributes["approval"].(schema.BoolAttribute)
	if !ok {
		t.Fatalf("expected approval bool attribute, got %#v", response.Schema.Attributes["approval"])
	}
	if attribute.Required || !attribute.Optional || !attribute.Computed {
		t.Fatalf("unexpected flags for approval: required=%t optional=%t computed=%t", attribute.Required, attribute.Optional, attribute.Computed)
	}
	if attribute.Default == nil {
		t.Fatal("expected approval to have a default")
	}
}

func TestWorkbenchToolAttributesIncludeApproval(t *testing.T) {
	for _, approval := range []bool{true, false} {
		tool := &model.WorkbenchTool{
			Name:       types.StringValue("ops_echo"),
			Tool:       types.StringValue("LAMBDA"),
			Approval:   types.BoolValue(approval),
			Categories: types.SetNull(types.StringType),
		}

		attributes, err := tool.Attributes(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if attributes.Approval == nil || *attributes.Approval != approval {
			t.Fatalf("expected approval %t, got %v", approval, attributes.Approval)
		}
	}
}

func TestWorkbenchToolFromReadsApproval(t *testing.T) {
	for name, testCase := range map[string]struct {
		response *bool
		expected bool
	}{
		"enabled":  {response: lo.ToPtr(true), expected: true},
		"disabled": {response: lo.ToPtr(false), expected: false},
		"unset":    {response: nil, expected: false},
	} {
		t.Run(name, func(t *testing.T) {
			tool := &model.WorkbenchTool{Categories: types.SetNull(types.StringType)}
			diagnostics := diag.Diagnostics{}

			tool.From(&console.WorkbenchToolFragment{
				ID:       "id",
				Name:     "ops_echo",
				Tool:     console.WorkbenchToolTypeLambda,
				Approval: testCase.response,
			}, context.Background(), &diagnostics)

			if diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diagnostics)
			}
			if tool.Approval != types.BoolValue(testCase.expected) {
				t.Fatalf("expected approval %t, got %v", testCase.expected, tool.Approval)
			}
		})
	}
}
