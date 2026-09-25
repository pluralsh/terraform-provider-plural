package resource

import (
	"context"
	"testing"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/model"

	"github.com/gqlgo/gqlgenc/clientv2"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
)

type dashboardImportClient struct {
	gqlclient.ConsoleClient
	updated *gqlclient.DashboardAttributes
}

func (c *dashboardImportClient) GetWorkbenchDashboard(_ context.Context, id string, _ ...clientv2.RequestInterceptor) (*gqlclient.GetWorkbenchDashboard, error) {
	return &gqlclient.GetWorkbenchDashboard{WorkbenchDashboard: &gqlclient.WorkbenchDashboardFragment{
		ID:        id,
		Name:      "imported",
		Workbench: &gqlclient.WorkbenchDashboardFragment_Workbench{ID: "workbench-1"},
		Graphs: []*gqlclient.WorkbenchDashboardGraphFragment{{
			Identifier: "cpu",
			Title:      lo.ToPtr("CPU"),
			Type:       gqlclient.DashboardGraphTypeTimeseries,
			Options:    map[string]any{"stacked": true, "legend": map[string]any{"position": "bottom"}},
			Layout:     gqlclient.WorkbenchDashboardGraphFragment_Layout{X: 0, Y: 0, W: 6, H: 4},
			Datasource: &gqlclient.WorkbenchDashboardDatasourceFragment{
				Type:  gqlclient.DashboardDatasourceTypeMetrics,
				Tool:  "plrl_metrics",
				Input: map[string]any{"query": "sum(rate(container_cpu_usage_seconds_total[5m]))"},
			},
		}},
		Inputs: []*gqlclient.WorkbenchDashboardInputFragment{{
			Name: "pod",
			Type: gqlclient.DashboardInputTypeSelect,
			Datasource: &gqlclient.WorkbenchDashboardDatasourceFragment{
				Type:  gqlclient.DashboardDatasourceTypeLabels,
				Tool:  "plrl_metric_label_search",
				Input: map[string]any{"metric": "container_cpu_usage_seconds_total", "label": "pod"},
			},
		}},
	}}, nil
}

func (c *dashboardImportClient) UpdateDashboard(_ context.Context, id string, attributes gqlclient.DashboardAttributes, _ ...clientv2.RequestInterceptor) (*gqlclient.UpdateDashboard, error) {
	c.updated = &attributes
	return &gqlclient.UpdateDashboard{UpdateDashboard: &gqlclient.WorkbenchDashboardFragment{ID: id}}, nil
}

// TestDashboardResourceImportKeepsSettings checks that graph options and datasource inputs of an imported
// dashboard are stored in the state, so that an update of an unrelated field sends them back unchanged.
func TestDashboardResourceImportKeepsSettings(t *testing.T) {
	ctx := context.Background()
	api := &dashboardImportClient{}
	r := &DashboardResource{client: client.NewClient(api)}
	s := resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, &s)

	imported := resource.ImportStateResponse{State: tfsdk.State{Schema: s.Schema, Raw: tftypes.NewValue(s.Schema.Type().TerraformType(ctx), nil)}}
	r.ImportState(ctx, resource.ImportStateRequest{ID: "dashboard-1"}, &imported)
	if imported.Diagnostics.HasError() {
		t.Fatal(imported.Diagnostics)
	}
	refreshed := resource.ReadResponse{State: imported.State}
	r.Read(ctx, resource.ReadRequest{State: imported.State}, &refreshed)
	if refreshed.Diagnostics.HasError() {
		t.Fatal(refreshed.Diagnostics)
	}

	var dashboard model.Dashboard
	if d := refreshed.State.Get(ctx, &dashboard); d.HasError() {
		t.Fatal(d)
	}
	if len(dashboard.Graphs) != 1 || len(dashboard.Inputs) != 1 {
		t.Fatalf("expected imported graphs and inputs, got %+v", dashboard)
	}
	expectStateJSON(t, "graph options", dashboard.Graphs[0].Options, `{"legend":{"position":"bottom"},"stacked":true}`)
	expectStateJSON(t, "graph datasource input", dashboard.Graphs[0].Datasource.Input, `{"query":"sum(rate(container_cpu_usage_seconds_total[5m]))"}`)
	expectStateJSON(t, "input datasource input", dashboard.Inputs[0].Datasource.Input, `{"label":"pod","metric":"container_cpu_usage_seconds_total"}`)

	// Update an unrelated field, the imported settings have to be sent back unchanged.
	dashboard.Name = types.StringValue("renamed")
	plan := tfsdk.Plan{Schema: s.Schema, Raw: tftypes.NewValue(s.Schema.Type().TerraformType(ctx), nil)}
	if d := plan.Set(ctx, &dashboard); d.HasError() {
		t.Fatal(d)
	}
	updated := resource.UpdateResponse{State: refreshed.State}
	r.Update(ctx, resource.UpdateRequest{Plan: plan, State: refreshed.State}, &updated)
	if updated.Diagnostics.HasError() {
		t.Fatal(updated.Diagnostics)
	}
	if api.updated == nil || len(api.updated.Graphs) != 1 || len(api.updated.Inputs) != 1 {
		t.Fatalf("expected the update to send the imported graphs and inputs, got %+v", api.updated)
	}

	graph := api.updated.Graphs[0]
	if graph.Options == nil || *graph.Options != `{"legend":{"position":"bottom"},"stacked":true}` {
		t.Fatalf("expected imported graph options to be sent back, got %v", lo.FromPtr(graph.Options))
	}
	if graph.Datasource.Input != `{"query":"sum(rate(container_cpu_usage_seconds_total[5m]))"}` {
		t.Fatalf("expected imported graph datasource input to be sent back, got %q", graph.Datasource.Input)
	}
	if input := api.updated.Inputs[0]; input.Datasource.Input != `{"label":"pod","metric":"container_cpu_usage_seconds_total"}` {
		t.Fatalf("expected imported input datasource input to be sent back, got %q", input.Datasource.Input)
	}
}

func expectStateJSON(t *testing.T, name string, actual jsontypes.Normalized, expected string) {
	t.Helper()

	if actual.IsNull() || actual.IsUnknown() || actual.ValueString() != expected {
		t.Fatalf("expected imported %s to be %s, got %v", name, expected, actual)
	}
}

func TestDashboardResourceSchemaMatchesModel(t *testing.T) {
	datasource := &model.DashboardDatasource{
		Type:  types.StringValue("METRICS"),
		Tool:  types.StringValue("prometheus"),
		Input: jsontypes.NewNormalizedValue(`{"query":"up"}`),
	}
	dashboard := &model.Dashboard{
		Id:          types.StringValue("dashboard-1"),
		WorkbenchID: types.StringValue("workbench-1"),
		Name:        types.StringValue("overview"),
		Description: types.StringNull(),
		Graphs: []*model.DashboardGraph{
			{
				Identifier:  types.StringValue("section"),
				Title:       types.StringValue("Errors"),
				Description: types.StringNull(),
				Type:        types.StringValue("SECTION"),
				SectionID:   types.StringNull(),
				Markdown:    types.StringNull(),
				Options:     jsontypes.NewNormalizedValue(`{"collapsed":false}`),
				Layout:      model.DashboardGraphLayout{X: types.Int64Value(0), Y: types.Int64Value(0), W: types.Int64Value(12), H: types.Int64Value(1)},
			},
			{
				Identifier:  types.StringValue("errors"),
				Title:       types.StringNull(),
				Description: types.StringNull(),
				Type:        types.StringValue("TIMESERIES"),
				SectionID:   types.StringValue("section"),
				Markdown:    types.StringNull(),
				Options:     jsontypes.NewNormalizedNull(),
				Layout:      model.DashboardGraphLayout{X: types.Int64Value(0), Y: types.Int64Value(1), W: types.Int64Value(6), H: types.Int64Value(4)},
				Datasource:  datasource,
			},
		},
		Inputs: []*model.DashboardInput{{
			Name:        types.StringValue("namespace"),
			Label:       types.StringValue("Namespace"),
			Description: types.StringNull(),
			Type:        types.StringValue("SELECT"),
			Default:     types.StringValue("prod"),
			Options:     types.ListValueMust(types.StringType, []attr.Value{types.StringValue("prod"), types.StringValue("dev")}),
			Required:    types.BoolValue(false),
			Datasource:  datasource,
		}},
	}

	assertSchemaMatchesModel(t, NewDashboardResource(), dashboard, new(model.Dashboard))
	assertSchemaMatchesModel(t, NewDashboardResource(), &model.Dashboard{
		Id:          types.StringValue("dashboard-2"),
		WorkbenchID: types.StringValue("workbench-1"),
		Name:        types.StringValue("empty"),
		Description: types.StringNull(),
	}, new(model.Dashboard))
}
