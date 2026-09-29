package resource

import (
	"context"
	"reflect"
	"testing"

	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
)

func TestWorkbenchToolResourceSchemaApprovalKeepsExistingSetting(t *testing.T) {
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
	// A default would plan false for configurations that don't set approval and disable an
	// existing approval requirement on the next unrelated update.
	if attribute.Default != nil {
		t.Fatal("expected approval to have no default")
	}
	if len(attribute.PlanModifiers) != 1 {
		t.Fatalf("expected a single plan modifier, got %d", len(attribute.PlanModifiers))
	}

	// An existing tool requiring approval, managed by a configuration that doesn't set it.
	// Any non-null prior state marks the resource as existing rather than being created.
	priorState, err := types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{}).ToTerraformValue(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	request := planmodifier.BoolRequest{
		State:       tfsdk.State{Raw: priorState},
		ConfigValue: types.BoolNull(),
		PlanValue:   types.BoolUnknown(),
		StateValue:  types.BoolValue(true),
	}
	modified := &planmodifier.BoolResponse{PlanValue: request.PlanValue}
	attribute.PlanModifiers[0].PlanModifyBool(context.Background(), request, modified)

	if modified.PlanValue != types.BoolValue(true) {
		t.Fatalf("expected the existing approval to be kept, planned %v", modified.PlanValue)
	}
}

func TestWorkbenchToolAttributesApproval(t *testing.T) {
	for name, testCase := range map[string]struct {
		approval types.Bool
		expected *bool
	}{
		"enabled":  {approval: types.BoolValue(true), expected: lo.ToPtr(true)},
		"disabled": {approval: types.BoolValue(false), expected: lo.ToPtr(false)},
		"null":     {approval: types.BoolNull(), expected: nil},
		"unknown":  {approval: types.BoolUnknown(), expected: nil},
	} {
		t.Run(name, func(t *testing.T) {
			tool := &model.WorkbenchTool{
				Name:       types.StringValue("ops_echo"),
				Tool:       types.StringValue("LAMBDA"),
				Approval:   testCase.approval,
				Categories: types.SetNull(types.StringType),
			}

			attributes, err := tool.Attributes(context.Background())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(attributes.Approval, testCase.expected) {
				t.Fatalf("expected approval %v, got %v", lo.FromPtr(testCase.expected), attributes.Approval)
			}
		})
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
