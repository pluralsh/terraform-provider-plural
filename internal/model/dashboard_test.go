package model

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
)

func TestDashboardAttributesAndFrom(t *testing.T) {
	ctx := context.Background()
	d := diag.Diagnostics{}
	dashboard := Dashboard{
		WorkbenchID: types.StringValue("workbench-1"),
		Name:        types.StringValue("overview"),
		Description: types.StringNull(),
		Graphs: []*DashboardGraph{{
			Identifier: types.StringValue("errors"),
			Type:       types.StringValue("TIMESERIES"),
			Options:    types.StringValue(`{"stacked":true}`),
			Layout:     DashboardGraphLayout{X: types.Int64Value(0), Y: types.Int64Value(0), W: types.Int64Value(6), H: types.Int64Value(4)},
			Datasource: &DashboardDatasource{Type: types.StringValue("METRICS"), Tool: types.StringValue("prometheus"), Input: types.StringNull()},
		}},
		Inputs: []*DashboardInput{{
			Name:    types.StringValue("namespace"),
			Type:    types.StringValue("SELECT"),
			Options: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("prod"), types.StringValue("dev")}),
		}},
	}

	attributes := dashboard.Attributes(ctx, &d)
	if d.HasError() {
		t.Fatalf("unexpected diagnostics: %v", d)
	}
	if *attributes.WorkbenchID != "workbench-1" || *attributes.Name != "overview" {
		t.Fatalf("expected workbench ID and name to map to attributes, got %+v", attributes)
	}
	graph := attributes.Graphs[0]
	if graph.Type != gqlclient.DashboardGraphTypeTimeseries || graph.Layout.W != 6 || *graph.Options != `{"stacked":true}` {
		t.Fatalf("expected graph to map to attributes, got %+v", graph)
	}
	if graph.Datasource.Input != "{}" {
		t.Fatalf("expected empty datasource input to default to an empty object, got %q", graph.Datasource.Input)
	}
	if input := attributes.Inputs[0]; input.Type != gqlclient.DashboardInputTypeSelect || *input.Options[1] != "dev" {
		t.Fatalf("expected input to map to attributes, got %+v", input)
	}

	dashboard.Inputs = nil
	encoded, err := json.Marshal(dashboard.Attributes(ctx, &d))
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var sent map[string]any
	if err := json.Unmarshal(encoded, &sent); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if inputs, ok := sent["inputs"].([]any); !ok || len(inputs) != 0 {
		t.Fatalf("expected removed inputs to be sent as an empty list, got %v", sent["inputs"])
	}

	dashboard.Inputs = []*DashboardInput{{
		Name:       types.StringValue("namespace"),
		Type:       types.StringValue("SELECT"),
		Options:    types.ListNull(types.StringType),
		Datasource: &DashboardDatasource{Type: types.StringValue("LABELS"), Tool: types.StringValue("prometheus"), Input: types.StringValue(`{"label":"namespace"}`)},
	}}
	dashboard.From(&gqlclient.WorkbenchDashboardFragment{
		ID:          "dashboard-1",
		Name:        "overview",
		Description: lo.ToPtr("Service overview"),
		Graphs: []*gqlclient.WorkbenchDashboardGraphFragment{
			{
				Identifier: "errors",
				Title:      lo.ToPtr("Errors"),
				Type:       gqlclient.DashboardGraphTypeTimeseries,
				Layout:     gqlclient.WorkbenchDashboardGraphFragment_Layout{X: 0, Y: 0, W: 12, H: 4},
				Datasource: &gqlclient.WorkbenchDashboardDatasourceFragment{Type: gqlclient.DashboardDatasourceTypeMetrics, Tool: "prometheus"},
			},
			{
				Identifier: "notes",
				Type:       gqlclient.DashboardGraphTypeMarkdown,
				Markdown:   lo.ToPtr("# Notes"),
				Layout:     gqlclient.WorkbenchDashboardGraphFragment_Layout{X: 0, Y: 4, W: 12, H: 2},
			},
		},
		Inputs: []*gqlclient.WorkbenchDashboardInputFragment{{
			Name:       "namespace",
			Type:       gqlclient.DashboardInputTypeSelect,
			Options:    []*string{},
			Datasource: &gqlclient.WorkbenchDashboardDatasourceFragment{Type: gqlclient.DashboardDatasourceTypeLabels, Tool: "prometheus"},
		}},
		Workbench: &gqlclient.WorkbenchDashboardFragment_Workbench{ID: "workbench-1"},
	}, ctx, &d)
	if d.HasError() {
		t.Fatalf("unexpected diagnostics: %v", d)
	}

	if got := dashboard.Id.ValueString(); got != "dashboard-1" {
		t.Fatalf("expected ID from response, got %q", got)
	}
	if got := dashboard.Description.ValueString(); got != "Service overview" {
		t.Fatalf("expected description from response, got %q", got)
	}
	if len(dashboard.Graphs) != 2 {
		t.Fatalf("expected graphs from response, got %d", len(dashboard.Graphs))
	}

	errors := dashboard.Graphs[0]
	if errors.Title.ValueString() != "Errors" || errors.Layout.W.ValueInt64() != 12 {
		t.Fatalf("expected graph fields from response, got %+v", errors)
	}
	if errors.Options.ValueString() != `{"stacked":true}` || !errors.Datasource.Input.IsNull() {
		t.Fatalf("expected options and datasource input not returned by the API to be kept, got %q and %q", errors.Options, errors.Datasource.Input)
	}

	notes := dashboard.Graphs[1]
	if notes.Markdown.ValueString() != "# Notes" || !notes.Options.IsNull() || notes.Datasource != nil {
		t.Fatalf("expected graph added in Console to be mapped without kept values, got %+v", notes)
	}

	input := dashboard.Inputs[0]
	if !input.Options.IsNull() {
		t.Fatalf("expected empty input options to stay unset, got %v", input.Options)
	}
	if input.Datasource.Input.ValueString() != `{"label":"namespace"}` {
		t.Fatalf("expected input datasource input not returned by the API to be kept, got %q", input.Datasource.Input)
	}

	dashboard.From(&gqlclient.WorkbenchDashboardFragment{ID: "dashboard-1", Name: "overview"}, ctx, &d)
	if dashboard.Graphs == nil || len(dashboard.Graphs) != 0 || len(dashboard.Inputs) != 0 {
		t.Fatalf("expected graphs and inputs removed in Console to be cleared, got %v and %v", dashboard.Graphs, dashboard.Inputs)
	}

	unset := Dashboard{}
	unset.From(&gqlclient.WorkbenchDashboardFragment{ID: "dashboard-2", Name: "empty"}, ctx, &d)
	if unset.Graphs != nil || unset.Inputs != nil {
		t.Fatalf("expected unset graphs and inputs to stay unset")
	}
}
