package resource

import (
	"context"
	"testing"

	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestMonitorResourceSchemaMatchesModel(t *testing.T) {
	monitor := &model.Monitor{
		Id:             types.StringValue("monitor-1"),
		Name:           types.StringValue("error-rate"),
		ServiceID:      types.StringValue("service-1"),
		WorkbenchID:    types.StringValue("workbench-1"),
		Prompt:         types.StringValue("Investigate the error rate."),
		Description:    types.StringNull(),
		AlertTemplate:  types.StringNull(),
		Severity:       types.StringValue("HIGH"),
		Type:           types.StringValue("LOG"),
		EvaluationCron: types.StringValue("*/5 * * * *"),
		Modes: &model.WorkbenchJobModes{
			Plan:         types.BoolValue(true),
			Verification: types.BoolNull(),
			Model:        &model.WorkbenchJobModel{Provider: types.StringValue("ANTHROPIC"), Model: types.StringValue("claude")},
			Coding:       &model.WorkbenchJobCodingModes{Babysit: types.BoolNull(), Approval: types.BoolValue(true), Review: types.BoolNull()},
			Budget:       &model.WorkbenchJobBudget{Cost: types.Float64Value(12.5), Tokens: types.Int64Value(1000)},
			Kubernetes: &model.WorkbenchJobKubernetesModes{
				Update:            types.BoolValue(false),
				Delete:            types.BoolValue(false),
				Exec:              types.BoolValue(false),
				Drain:             types.BoolValue(false),
				ExcludeNamespaces: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("kube-system")}),
				RequireNamespaces: types.SetNull(types.StringType),
			},
		},
		Query: model.MonitorQuery{
			Log: &model.MonitorLogQuery{
				Tool:       types.StringNull(),
				Query:      types.StringValue("level:error"),
				BucketSize: types.StringValue("5m"),
				Duration:   types.StringValue("1h"),
				Operator:   types.StringValue("OR"),
				Facets:     []*model.MonitorFacet{{Key: types.StringValue("namespace"), Value: types.StringValue("prod")}},
				Options:    &model.MonitorLogOptions{Azure: &model.MonitorLogAzureOptions{ResourceID: types.StringValue("resource-1")}},
			},
		},
		Threshold: model.MonitorThreshold{Aggregate: types.StringValue("MAX"), Value: types.Float64Value(0.95)},
	}

	assertSchemaMatchesModel(t, NewMonitorResource(), monitor, new(model.Monitor))

	metrics := &model.MonitorMetricsQuery{
		Query: types.StringValue("rate(errors[5m])"),
		Options: &model.MonitorMetricsOptions{Azure: &model.MonitorMetricsAzureOptions{
			ResourceID:       types.StringValue("resource-1"),
			MetricsNamespace: types.StringNull(),
			Aggregation:      types.StringValue("Average"),
			Filter:           types.StringNull(),
			OrderBy:          types.StringNull(),
			RollUpBy:         types.StringNull(),
			MetricsEndpoint:  types.StringNull(),
		}},
	}
	monitor.Type = types.StringValue("METRICS")
	monitor.Query = model.MonitorQuery{Metrics: metrics}
	monitor.Modes = nil
	assertSchemaMatchesModel(t, NewMonitorResource(), monitor, new(model.Monitor))
}

// assertSchemaMatchesModel checks that the resource schema is valid and that the model can be stored in
// and read back from a state using that schema, which fails when schema attributes and model fields differ.
func assertSchemaMatchesModel[T any](t *testing.T, r resource.Resource, in *T, out *T) {
	t.Helper()
	ctx := context.Background()

	response := resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, &response)
	if response.Diagnostics.HasError() {
		t.Fatalf("unexpected schema diagnostics: %v", response.Diagnostics)
	}
	if diagnostics := response.Schema.ValidateImplementation(ctx); diagnostics.HasError() {
		t.Fatalf("invalid schema: %v", diagnostics)
	}

	state := tfsdk.State{Schema: response.Schema, Raw: tftypes.NewValue(response.Schema.Type().TerraformType(ctx), nil)}
	if diagnostics := state.Set(ctx, in); diagnostics.HasError() {
		t.Fatalf("unable to set model in state: %v", diagnostics)
	}
	if diagnostics := state.Get(ctx, out); diagnostics.HasError() {
		t.Fatalf("unable to get model from state: %v", diagnostics)
	}
}
