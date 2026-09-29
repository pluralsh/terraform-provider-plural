package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
)

type Monitor struct {
	Id             types.String       `tfsdk:"id"`
	Name           types.String       `tfsdk:"name"`
	ServiceID      types.String       `tfsdk:"service_id"`
	WorkbenchID    types.String       `tfsdk:"workbench_id"`
	Prompt         types.String       `tfsdk:"prompt"`
	Modes          *WorkbenchJobModes `tfsdk:"modes"`
	Description    types.String       `tfsdk:"description"`
	AlertTemplate  types.String       `tfsdk:"alert_template"`
	Severity       types.String       `tfsdk:"severity"`
	Type           types.String       `tfsdk:"type"`
	EvaluationCron types.String       `tfsdk:"evaluation_cron"`
	Query          *MonitorQuery      `tfsdk:"query"`
	Threshold      *MonitorThreshold  `tfsdk:"threshold"`
}

func (in *Monitor) Attributes(ctx context.Context, d *diag.Diagnostics) gqlclient.MonitorAttributes {
	return gqlclient.MonitorAttributes{
		ServiceID:      in.ServiceID.ValueString(),
		WorkbenchID:    in.WorkbenchID.ValueStringPointer(),
		Prompt:         in.Prompt.ValueStringPointer(),
		Modes:          in.Modes.Attributes(ctx, d),
		Name:           in.Name.ValueString(),
		Description:    in.Description.ValueStringPointer(),
		AlertTemplate:  in.AlertTemplate.ValueStringPointer(),
		Severity:       gqlclient.AlertSeverity(in.Severity.ValueString()),
		Type:           gqlclient.MonitorType(in.Type.ValueString()),
		EvaluationCron: in.EvaluationCron.ValueString(),
		Query:          in.Query.Attributes(),
		Threshold:      in.Threshold.Attributes(),
	}
}

func (in *Monitor) From(response *gqlclient.MonitorFragment, ctx context.Context, d *diag.Diagnostics) {
	if in == nil || response == nil {
		return
	}

	in.Id = types.StringValue(response.ID)
	in.Name = types.StringValue(response.Name)
	in.Prompt = types.StringPointerValue(response.Prompt)
	in.Description = types.StringPointerValue(response.Description)
	in.AlertTemplate = types.StringPointerValue(response.AlertTemplate)
	in.Severity = types.StringValue(string(response.Severity))
	in.Type = types.StringValue(string(response.Type))
	in.EvaluationCron = types.StringValue(response.EvaluationCron)
	ensure(&in.Query)
	ensure(&in.Threshold)
	in.Query.From(&response.Query)
	in.Threshold.From(&response.Threshold)

	if response.Modes != nil {
		ensure(&in.Modes)
		in.Modes.From(response.Modes, ctx, d)
	} else {
		in.Modes = nil
	}

	if response.Service != nil {
		in.ServiceID = types.StringValue(response.Service.ID)
	}

	if response.Workbench != nil {
		in.WorkbenchID = types.StringValue(response.Workbench.ID)
	} else {
		in.WorkbenchID = types.StringNull()
	}
}

type MonitorQuery struct {
	Log     *MonitorLogQuery     `tfsdk:"log"`
	Metrics *MonitorMetricsQuery `tfsdk:"metrics"`
}

func (in *MonitorQuery) Attributes() gqlclient.MonitorQueryAttributes {
	return gqlclient.MonitorQueryAttributes{
		Log:     in.Log.Attributes(),
		Metrics: in.Metrics.Attributes(),
	}
}

func (in *MonitorQuery) From(response *gqlclient.MonitorQueryFragment) {
	if response.Log != nil {
		ensure(&in.Log)
		in.Log.From(response.Log)
	} else {
		in.Log = nil
	}

	if response.Metrics != nil {
		ensure(&in.Metrics)
		in.Metrics.From(response.Metrics)
	} else {
		in.Metrics = nil
	}
}

type MonitorLogQuery struct {
	Tool       types.String       `tfsdk:"tool"`
	Query      types.String       `tfsdk:"query"`
	BucketSize types.String       `tfsdk:"bucket_size"`
	Duration   types.String       `tfsdk:"duration"`
	Operator   types.String       `tfsdk:"operator"`
	Facets     []*MonitorFacet    `tfsdk:"facets"`
	Options    *MonitorLogOptions `tfsdk:"options"`
}

func (in *MonitorLogQuery) Attributes() *gqlclient.MonitorLogQueryAttributes {
	if in == nil {
		return nil
	}

	var operator *gqlclient.MonitorOperator
	if !in.Operator.IsNull() && !in.Operator.IsUnknown() {
		operator = lo.ToPtr(gqlclient.MonitorOperator(in.Operator.ValueString()))
	}

	return &gqlclient.MonitorLogQueryAttributes{
		Tool:       in.Tool.ValueStringPointer(),
		Query:      in.Query.ValueString(),
		BucketSize: in.BucketSize.ValueString(),
		Duration:   in.Duration.ValueStringPointer(),
		Operator:   operator,
		Facets: lo.Map(in.Facets, func(f *MonitorFacet, _ int) *gqlclient.MonitorFacetAttributes {
			return &gqlclient.MonitorFacetAttributes{Key: f.Key.ValueString(), Value: f.Value.ValueString()}
		}),
		Options: in.Options.Attributes(),
	}
}

func (in *MonitorLogQuery) From(response *gqlclient.MonitorQueryFragment_Log) {
	in.Tool = types.StringPointerValue(response.Tool)
	in.Query = types.StringValue(response.Query)
	in.BucketSize = types.StringValue(response.BucketSize)
	in.Duration = types.StringPointerValue(response.Duration)

	if response.Operator != nil {
		in.Operator = types.StringValue(string(*response.Operator))
	} else {
		in.Operator = types.StringNull()
	}

	// Console returns an empty list for unset facets, keep them unset to avoid nil-vs-empty diffs.
	if len(response.Facets) > 0 || in.Facets != nil {
		in.Facets = lo.Map(response.Facets, func(f *gqlclient.MonitorQueryFragment_Log_Facets, _ int) *MonitorFacet {
			return &MonitorFacet{Key: types.StringValue(f.Key), Value: types.StringValue(f.Value)}
		})
	}

	if response.Options != nil {
		ensure(&in.Options)
		in.Options.From(response.Options)
	} else {
		in.Options = nil
	}
}

type MonitorFacet struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type MonitorLogOptions struct {
	Azure *MonitorLogAzureOptions `tfsdk:"azure"`
}

func (in *MonitorLogOptions) Attributes() *gqlclient.MonitorLogOptionsAttributes {
	if in == nil {
		return nil
	}

	var azure *gqlclient.MonitorLogAzureOptionsAttributes
	if in.Azure != nil {
		azure = &gqlclient.MonitorLogAzureOptionsAttributes{ResourceID: in.Azure.ResourceID.ValueStringPointer()}
	}

	return &gqlclient.MonitorLogOptionsAttributes{Azure: azure}
}

func (in *MonitorLogOptions) From(response *gqlclient.MonitorQueryFragment_Log_Options) {
	if response.Azure != nil {
		in.Azure = &MonitorLogAzureOptions{ResourceID: types.StringPointerValue(response.Azure.ResourceID)}
	} else {
		in.Azure = nil
	}
}

type MonitorLogAzureOptions struct {
	ResourceID types.String `tfsdk:"resource_id"`
}

type MonitorMetricsQuery struct {
	Tool     types.String           `tfsdk:"tool"`
	Query    types.String           `tfsdk:"query"`
	Step     types.String           `tfsdk:"step"`
	Duration types.String           `tfsdk:"duration"`
	Options  *MonitorMetricsOptions `tfsdk:"options"`
}

func (in *MonitorMetricsQuery) Attributes() *gqlclient.MonitorMetricsQueryAttributes {
	if in == nil {
		return nil
	}

	return &gqlclient.MonitorMetricsQueryAttributes{
		Tool:     in.Tool.ValueStringPointer(),
		Query:    in.Query.ValueString(),
		Step:     in.Step.ValueStringPointer(),
		Duration: in.Duration.ValueStringPointer(),
		Options:  in.Options.Attributes(),
	}
}

func (in *MonitorMetricsQuery) From(response *gqlclient.MonitorQueryFragment_Metrics) {
	in.Tool = types.StringPointerValue(response.Tool)
	in.Query = types.StringValue(response.Query)
	in.Step = types.StringPointerValue(response.Step)
	in.Duration = types.StringPointerValue(response.Duration)

	if response.Options != nil {
		ensure(&in.Options)
		in.Options.From(response.Options)
	} else {
		in.Options = nil
	}
}

type MonitorMetricsOptions struct {
	Azure *MonitorMetricsAzureOptions `tfsdk:"azure"`
}

func (in *MonitorMetricsOptions) Attributes() *gqlclient.MonitorMetricsOptionsAttributes {
	if in == nil {
		return nil
	}

	var azure *gqlclient.MonitorMetricsAzureOptionsAttributes
	if in.Azure != nil {
		azure = &gqlclient.MonitorMetricsAzureOptionsAttributes{
			ResourceID:       in.Azure.ResourceID.ValueStringPointer(),
			MetricsNamespace: in.Azure.MetricsNamespace.ValueStringPointer(),
			Aggregation:      in.Azure.Aggregation.ValueStringPointer(),
			Filter:           in.Azure.Filter.ValueStringPointer(),
			OrderBy:          in.Azure.OrderBy.ValueStringPointer(),
			RollUpBy:         in.Azure.RollUpBy.ValueStringPointer(),
			MetricsEndpoint:  in.Azure.MetricsEndpoint.ValueStringPointer(),
		}
	}

	return &gqlclient.MonitorMetricsOptionsAttributes{Azure: azure}
}

func (in *MonitorMetricsOptions) From(response *gqlclient.MonitorQueryFragment_Metrics_Options) {
	if response.Azure != nil {
		in.Azure = &MonitorMetricsAzureOptions{
			ResourceID:       types.StringPointerValue(response.Azure.ResourceID),
			MetricsNamespace: types.StringPointerValue(response.Azure.MetricsNamespace),
			Aggregation:      types.StringPointerValue(response.Azure.Aggregation),
			Filter:           types.StringPointerValue(response.Azure.Filter),
			OrderBy:          types.StringPointerValue(response.Azure.OrderBy),
			RollUpBy:         types.StringPointerValue(response.Azure.RollUpBy),
			MetricsEndpoint:  types.StringPointerValue(response.Azure.MetricsEndpoint),
		}
	} else {
		in.Azure = nil
	}
}

type MonitorMetricsAzureOptions struct {
	ResourceID       types.String `tfsdk:"resource_id"`
	MetricsNamespace types.String `tfsdk:"metrics_namespace"`
	Aggregation      types.String `tfsdk:"aggregation"`
	Filter           types.String `tfsdk:"filter"`
	OrderBy          types.String `tfsdk:"order_by"`
	RollUpBy         types.String `tfsdk:"roll_up_by"`
	MetricsEndpoint  types.String `tfsdk:"metrics_endpoint"`
}

type MonitorThreshold struct {
	Aggregate types.String  `tfsdk:"aggregate"`
	Value     types.Float64 `tfsdk:"value"`
}

func (in *MonitorThreshold) Attributes() gqlclient.MonitorThresholdAttributes {
	return gqlclient.MonitorThresholdAttributes{
		Aggregate: gqlclient.MonitorAggregate(in.Aggregate.ValueString()),
		Value:     in.Value.ValueFloat64(),
	}
}

func (in *MonitorThreshold) From(response *gqlclient.MonitorFragment_Threshold) {
	in.Aggregate = types.StringValue(string(response.Aggregate))
	in.Value = types.Float64Value(response.Value)
}
