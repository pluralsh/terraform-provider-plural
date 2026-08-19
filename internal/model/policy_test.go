package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
)

func TestPolicyAttributesAndFrom(t *testing.T) {
	policy := Policy{
		Name:        types.StringValue("allow-workbench"),
		Type:        types.StringValue("WORKBENCH"),
		Description: types.StringValue("Allows supported workbench actions."),
		Policy:      types.StringValue("package policy\ndefault allow := true"),
		ProjectId:   types.StringValue("project-1"),
	}

	attributes := policy.Attributes()
	if got := *attributes.Name; got != "allow-workbench" {
		t.Fatalf("expected name to map to attributes, got %q", got)
	}
	if got := string(*attributes.Type); got != "WORKBENCH" {
		t.Fatalf("expected type to map to attributes, got %q", got)
	}
	if got := *attributes.ProjectID; got != "project-1" {
		t.Fatalf("expected project ID to map to attributes, got %q", got)
	}

	policy.From(&gqlclient.PolicyFragment{
		ID:          "policy-1",
		Name:        "allow-workbench",
		Type:        gqlclient.PolicyTypeWorkbench,
		Description: lo.ToPtr("Allows supported workbench actions."),
		Policy:      "package policy\ndefault allow := true",
		Project:     &gqlclient.TinyProjectFragment{ID: "project-1"},
	})

	if got := policy.Id.ValueString(); got != "policy-1" {
		t.Fatalf("expected ID from response, got %q", got)
	}
	if got := policy.ProjectId.ValueString(); got != "project-1" {
		t.Fatalf("expected project ID from response, got %q", got)
	}
}

func TestBindingPolicyAttributesAndFrom(t *testing.T) {
	ctx := context.Background()
	matches := BindingPolicyMatches{
		Workbench: &WorkbenchPolicyMatches{
			Regexes: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("^prod-.*$")}),
		},
	}
	bindingPolicy := BindingPolicy{
		PolicyId:     types.StringValue("policy-1"),
		BindPolicyId: types.StringValue("policy-2"),
		Type:         types.StringValue("WORKBENCH"),
		Interval:     types.StringValue("1h"),
		Matches:      &matches,
	}
	diagnostics := diag.Diagnostics{}

	attributes := bindingPolicy.Attributes(ctx, &diagnostics)
	if diagnostics.HasError() {
		t.Fatalf("expected no diagnostics from attributes, got %#v", diagnostics)
	}
	if attributes.PolicyID != "policy-1" || attributes.BindPolicyID != "policy-2" {
		t.Fatalf("expected policy IDs to map, got policy=%q bind=%q", attributes.PolicyID, attributes.BindPolicyID)
	}
	if got := *attributes.Matches.Workbench.Regexes[0]; got != "^prod-.*$" {
		t.Fatalf("expected workbench regex to map, got %q", got)
	}

	bindingPolicy.From(&gqlclient.BindingPolicyFragment{
		ID:         "binding-1",
		Type:       gqlclient.BindingPolicyTypeWorkbench,
		Interval:   "30m",
		Policy:     &gqlclient.TinyPolicyFragment{ID: "policy-1"},
		BindPolicy: &gqlclient.TinyPolicyFragment{ID: "policy-2"},
		Matches: &gqlclient.BindingPolicyFragment_Matches{
			Workbench: &gqlclient.BindingPolicyFragment_Matches_Workbench{Regexes: []*string{lo.ToPtr("^dev-.*$")}},
		},
	}, ctx, &diagnostics)
	if diagnostics.HasError() {
		t.Fatalf("expected no diagnostics from response mapping, got %#v", diagnostics)
	}
	if got := bindingPolicy.Id.ValueString(); got != "binding-1" {
		t.Fatalf("expected ID from response, got %q", got)
	}
	if got := bindingPolicy.Matches.Workbench.Regexes.Elements()[0].(types.String).ValueString(); got != "^dev-.*$" {
		t.Fatalf("expected regex from response, got %q", got)
	}
}

func TestBindingPolicyAttributesReportInvalidRegexes(t *testing.T) {
	bindingPolicy := BindingPolicy{
		PolicyId:     types.StringValue("policy-1"),
		BindPolicyId: types.StringValue("policy-2"),
		Type:         types.StringValue("WORKBENCH"),
		Matches: &BindingPolicyMatches{
			Workbench: &WorkbenchPolicyMatches{Regexes: types.ListUnknown(types.StringType)},
		},
	}
	diagnostics := diag.Diagnostics{}

	attributes := bindingPolicy.Attributes(context.Background(), &diagnostics)
	if !diagnostics.HasError() {
		t.Fatal("expected diagnostics when regexes are unknown")
	}
	if attributes.Matches != nil {
		t.Fatalf("expected no matches attributes after conversion failure, got %#v", attributes.Matches)
	}
}
