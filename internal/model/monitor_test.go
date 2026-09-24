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

func TestMonitorAttributesAndFrom(t *testing.T) {
	ctx := context.Background()
	d := diag.Diagnostics{}
	monitor := Monitor{
		Name:           types.StringValue("error-rate"),
		ServiceID:      types.StringValue("service-1"),
		WorkbenchID:    types.StringValue("workbench-1"),
		Prompt:         types.StringValue("Investigate the error rate."),
		Severity:       types.StringValue("HIGH"),
		Type:           types.StringValue("LOG"),
		EvaluationCron: types.StringValue("*/5 * * * *"),
		Modes: &WorkbenchJobModes{
			Plan:   types.BoolValue(true),
			Model:  &WorkbenchJobModel{Provider: types.StringValue("ANTHROPIC"), Model: types.StringValue("claude")},
			Budget: &WorkbenchJobBudget{Cost: types.Float64Value(12.5), Tokens: types.Int64Null()},
			Kubernetes: &WorkbenchJobKubernetesModes{
				ExcludeNamespaces: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("kube-system")}),
				RequireNamespaces: types.SetNull(types.StringType),
			},
		},
		Query: MonitorQuery{
			Log: &MonitorLogQuery{
				Query:      types.StringValue("level:error"),
				BucketSize: types.StringValue("5m"),
				Operator:   types.StringValue("AND"),
				Facets:     []*MonitorFacet{{Key: types.StringValue("namespace"), Value: types.StringValue("prod")}},
			},
		},
		Threshold: MonitorThreshold{Aggregate: types.StringValue("MAX"), Value: types.Float64Value(0.95)},
	}

	attributes := monitor.Attributes(ctx, &d)
	if d.HasError() {
		t.Fatalf("unexpected diagnostics: %v", d)
	}
	if attributes.ServiceID != "service-1" || *attributes.WorkbenchID != "workbench-1" {
		t.Fatalf("expected service and workbench IDs to map to attributes, got %+v", attributes)
	}
	if attributes.Severity != gqlclient.AlertSeverityHigh || attributes.Type != gqlclient.MonitorTypeLog {
		t.Fatalf("expected enums to map to attributes, got %q and %q", attributes.Severity, attributes.Type)
	}
	if attributes.Threshold.Aggregate != gqlclient.MonitorAggregateMax || attributes.Threshold.Value != 0.95 {
		t.Fatalf("expected threshold to map to attributes, got %+v", attributes.Threshold)
	}
	if *attributes.Query.Log.Operator != gqlclient.MonitorOperatorAnd || attributes.Query.Log.Facets[0].Value != "prod" {
		t.Fatalf("expected log query to map to attributes, got %+v", attributes.Query.Log)
	}
	if attributes.Query.Metrics != nil {
		t.Fatalf("expected metrics query to be nil")
	}
	if *attributes.Modes.Budget.Cost != 12.5 || attributes.Modes.Budget.Tokens != nil {
		t.Fatalf("expected budget to map to attributes, got %+v", attributes.Modes.Budget)
	}
	if got := attributes.Modes.Kubernetes; *got.ExcludeNamespaces[0] != "kube-system" || got.RequireNamespaces != nil {
		t.Fatalf("expected kubernetes namespaces to map to attributes, got %+v", got)
	}

	// Monitor attributes are sent as the full desired state, so unset fields have to be cleared explicitly.
	encoded, err := json.Marshal(attributes)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var sent map[string]any
	if err := json.Unmarshal(encoded, &sent); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if value, ok := sent["description"]; !ok || value != nil {
		t.Fatalf("expected unset description to be sent as null, got %v", sent["description"])
	}

	monitor.From(&gqlclient.MonitorFragment{
		ID:             "monitor-1",
		Name:           "error-rate",
		Severity:       gqlclient.AlertSeverityCritical,
		Type:           gqlclient.MonitorTypeMetrics,
		EvaluationCron: "*/10 * * * *",
		Modes: &gqlclient.WorkbenchJobModesFragment{
			Plan:       lo.ToPtr(false),
			Kubernetes: &gqlclient.WorkbenchJobModesFragment_Kubernetes{ExcludeNamespaces: []*string{lo.ToPtr("default")}},
		},
		Query: gqlclient.MonitorQueryFragment{
			Metrics: &gqlclient.MonitorQueryFragment_Metrics{Query: "rate(errors[5m])", Step: lo.ToPtr("1m")},
		},
		Threshold: gqlclient.MonitorFragment_Threshold{Aggregate: gqlclient.MonitorAggregateAvg, Value: 1.5},
		Service:   &gqlclient.MonitorFragment_Service{ID: "service-1"},
	}, ctx, &d)
	if d.HasError() {
		t.Fatalf("unexpected diagnostics: %v", d)
	}

	if got := monitor.Id.ValueString(); got != "monitor-1" {
		t.Fatalf("expected ID from response, got %q", got)
	}
	if !monitor.WorkbenchID.IsNull() || !monitor.Prompt.IsNull() {
		t.Fatalf("expected workbench ID and prompt to be cleared when not returned")
	}
	if monitor.Severity.ValueString() != "CRITICAL" || monitor.Type.ValueString() != "METRICS" || monitor.EvaluationCron.ValueString() != "*/10 * * * *" {
		t.Fatalf("expected scalar fields from response, got %q, %q and %q", monitor.Severity, monitor.Type, monitor.EvaluationCron)
	}
	if monitor.Query.Log != nil || monitor.Query.Metrics.Step.ValueString() != "1m" {
		t.Fatalf("expected query from response, got %+v", monitor.Query)
	}
	if monitor.Threshold.Aggregate.ValueString() != "AVG" || monitor.Threshold.Value.ValueFloat64() != 1.5 {
		t.Fatalf("expected threshold from response, got %+v", monitor.Threshold)
	}
	if monitor.Modes.Plan.ValueBool() || monitor.Modes.Model != nil || monitor.Modes.Budget != nil {
		t.Fatalf("expected modes from response, got %+v", monitor.Modes)
	}
	if !monitor.Modes.Kubernetes.RequireNamespaces.IsNull() || len(monitor.Modes.Kubernetes.ExcludeNamespaces.Elements()) != 1 {
		t.Fatalf("expected kubernetes namespaces from response, got %+v", monitor.Modes.Kubernetes)
	}

	monitor.From(&gqlclient.MonitorFragment{ID: "monitor-1", Name: "error-rate"}, ctx, &d)
	if monitor.Modes != nil {
		t.Fatalf("expected modes to be cleared when not returned")
	}
}

func TestMonitorLogFacetsFrom(t *testing.T) {
	query := MonitorLogQuery{}
	query.From(&gqlclient.MonitorQueryFragment_Log{Query: "level:error", BucketSize: "5m", Facets: []*gqlclient.MonitorQueryFragment_Log_Facets{}})
	if query.Facets != nil {
		t.Fatalf("expected empty facets to stay unset, got %v", query.Facets)
	}

	query.Facets = []*MonitorFacet{{Key: types.StringValue("namespace"), Value: types.StringValue("prod")}}
	query.From(&gqlclient.MonitorQueryFragment_Log{Query: "level:error", BucketSize: "5m"})
	if query.Facets == nil || len(query.Facets) != 0 {
		t.Fatalf("expected facets removed in Console to be cleared, got %v", query.Facets)
	}
}
