package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	gqlclient "github.com/pluralsh/console/go/client"
)

func TestServiceContextFromSetsProjectID(t *testing.T) {
	response := &gqlclient.ServiceContextFragment{
		ID:            "sc-1",
		Configuration: map[string]any{"foo": "bar"},
		Project:       &gqlclient.TinyProjectFragment{ID: "proj-1"},
	}

	sc := ServiceContext{}
	diagnostics := diag.Diagnostics{}
	sc.From(response, context.Background(), &diagnostics)

	if diagnostics.HasError() {
		t.Fatalf("expected no diagnostics errors, got: %#v", diagnostics)
	}
	if got := sc.ProjectId.ValueString(); got != "proj-1" {
		t.Fatalf("expected project_id to be set from response project, got %q", got)
	}
}

func TestServiceContextFromSetsNullProjectIDWhenProjectMissing(t *testing.T) {
	response := &gqlclient.ServiceContextFragment{
		ID:            "sc-2",
		Configuration: map[string]any{"foo": "bar"},
		Project:       nil,
	}

	sc := ServiceContext{}
	diagnostics := diag.Diagnostics{}
	sc.From(response, context.Background(), &diagnostics)

	if diagnostics.HasError() {
		t.Fatalf("expected no diagnostics errors, got: %#v", diagnostics)
	}
	if !sc.ProjectId.IsNull() {
		t.Fatalf("expected project_id to be null when project is missing, got %q", sc.ProjectId.ValueString())
	}
}
