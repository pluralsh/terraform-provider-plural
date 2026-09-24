package model

import (
	"context"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
)

type Dashboard struct {
	Id          types.String      `tfsdk:"id"`
	WorkbenchID types.String      `tfsdk:"workbench_id"`
	Name        types.String      `tfsdk:"name"`
	Description types.String      `tfsdk:"description"`
	Graphs      []*DashboardGraph `tfsdk:"graphs"`
	Inputs      []*DashboardInput `tfsdk:"inputs"`
}

func (in *Dashboard) Attributes(ctx context.Context, d *diag.Diagnostics) gqlclient.DashboardAttributes {
	return gqlclient.DashboardAttributes{
		WorkbenchID: in.WorkbenchID.ValueStringPointer(),
		Name:        in.Name.ValueStringPointer(),
		Description: in.Description.ValueStringPointer(),
		Graphs:      lo.Map(in.Graphs, func(g *DashboardGraph, _ int) *gqlclient.DashboardGraphAttributes { return g.Attributes() }),
		Inputs:      lo.Map(in.Inputs, func(i *DashboardInput, _ int) *gqlclient.DashboardInputAttributes { return i.Attributes(ctx, d) }),
	}
}

// From maps the dashboard returned by the Console API. The API does not return graph options
// and datasource inputs, so they are kept from the current state, matching graphs by identifier
// and inputs by name.
func (in *Dashboard) From(response *gqlclient.WorkbenchDashboardFragment, ctx context.Context, d *diag.Diagnostics) {
	if in == nil || response == nil {
		return
	}

	in.Id = types.StringValue(response.ID)
	in.Name = types.StringValue(response.Name)
	in.Description = types.StringPointerValue(response.Description)

	if response.Workbench != nil {
		in.WorkbenchID = types.StringValue(response.Workbench.ID)
	}

	// Console returns empty lists for unset graphs and inputs, keep them unset to avoid nil-vs-empty diffs.
	if len(response.Graphs) > 0 || in.Graphs != nil {
		current := lo.KeyBy(in.Graphs, func(g *DashboardGraph) string { return g.Identifier.ValueString() })
		in.Graphs = lo.FilterMap(response.Graphs, func(g *gqlclient.WorkbenchDashboardGraphFragment, _ int) (*DashboardGraph, bool) {
			if g == nil {
				return nil, false
			}

			graph := current[g.Identifier]
			if graph == nil {
				graph = &DashboardGraph{Options: types.StringNull()}
			}
			graph.From(g)
			return graph, true
		})
	}

	if len(response.Inputs) > 0 || in.Inputs != nil {
		current := lo.KeyBy(in.Inputs, func(i *DashboardInput) string { return i.Name.ValueString() })
		in.Inputs = lo.FilterMap(response.Inputs, func(i *gqlclient.WorkbenchDashboardInputFragment, _ int) (*DashboardInput, bool) {
			if i == nil {
				return nil, false
			}

			input := current[i.Name]
			if input == nil {
				input = &DashboardInput{Options: types.ListNull(types.StringType)}
			}
			input.From(i, ctx, d)
			return input, true
		})
	}
}

type DashboardGraph struct {
	Identifier  types.String         `tfsdk:"identifier"`
	Title       types.String         `tfsdk:"title"`
	Description types.String         `tfsdk:"description"`
	Type        types.String         `tfsdk:"type"`
	SectionID   types.String         `tfsdk:"section_id"`
	Markdown    types.String         `tfsdk:"markdown"`
	Options     types.String         `tfsdk:"options"`
	Layout      DashboardGraphLayout `tfsdk:"layout"`
	Datasource  *DashboardDatasource `tfsdk:"datasource"`
}

func (in *DashboardGraph) Attributes() *gqlclient.DashboardGraphAttributes {
	return &gqlclient.DashboardGraphAttributes{
		Identifier:  in.Identifier.ValueString(),
		Title:       in.Title.ValueStringPointer(),
		Description: in.Description.ValueStringPointer(),
		Type:        gqlclient.DashboardGraphType(in.Type.ValueString()),
		SectionID:   in.SectionID.ValueStringPointer(),
		Markdown:    in.Markdown.ValueStringPointer(),
		Options:     in.Options.ValueStringPointer(),
		Layout: gqlclient.DashboardGraphLayoutAttributes{
			X: in.Layout.X.ValueInt64(),
			Y: in.Layout.Y.ValueInt64(),
			W: in.Layout.W.ValueInt64(),
			H: in.Layout.H.ValueInt64(),
		},
		Datasource: in.Datasource.Attributes(),
	}
}

// From maps the graph returned by the Console API, keeping options and datasource input as they are not returned.
func (in *DashboardGraph) From(response *gqlclient.WorkbenchDashboardGraphFragment) {
	in.Identifier = types.StringValue(response.Identifier)
	in.Title = types.StringPointerValue(response.Title)
	in.Description = types.StringPointerValue(response.Description)
	in.Type = types.StringValue(string(response.Type))
	in.SectionID = types.StringPointerValue(response.SectionID)
	in.Markdown = types.StringPointerValue(response.Markdown)
	in.Layout = DashboardGraphLayout{
		X: types.Int64Value(response.Layout.X),
		Y: types.Int64Value(response.Layout.Y),
		W: types.Int64Value(response.Layout.W),
		H: types.Int64Value(response.Layout.H),
	}
	in.Datasource = in.Datasource.from(response.Datasource)
}

type DashboardGraphLayout struct {
	X types.Int64 `tfsdk:"x"`
	Y types.Int64 `tfsdk:"y"`
	W types.Int64 `tfsdk:"w"`
	H types.Int64 `tfsdk:"h"`
}

type DashboardDatasource struct {
	Type  types.String `tfsdk:"type"`
	Tool  types.String `tfsdk:"tool"`
	Input types.String `tfsdk:"input"`
}

func (in *DashboardDatasource) Attributes() *gqlclient.DashboardDatasourceAttributes {
	if in == nil {
		return nil
	}

	input := "{}"
	if !in.Input.IsNull() && !in.Input.IsUnknown() && in.Input.ValueString() != "" {
		input = in.Input.ValueString()
	}

	return &gqlclient.DashboardDatasourceAttributes{
		Type:  gqlclient.DashboardDatasourceType(in.Type.ValueString()),
		Tool:  in.Tool.ValueString(),
		Input: input,
	}
}

// from maps the datasource returned by the Console API, keeping the input as it is not returned.
func (in *DashboardDatasource) from(response *gqlclient.WorkbenchDashboardDatasourceFragment) *DashboardDatasource {
	if response == nil {
		return nil
	}

	input := types.StringNull()
	if in != nil {
		input = in.Input
	}

	return &DashboardDatasource{
		Type:  types.StringValue(string(response.Type)),
		Tool:  types.StringValue(response.Tool),
		Input: input,
	}
}

type DashboardInput struct {
	Name        types.String         `tfsdk:"name"`
	Label       types.String         `tfsdk:"label"`
	Description types.String         `tfsdk:"description"`
	Type        types.String         `tfsdk:"type"`
	Default     types.String         `tfsdk:"default"`
	Options     types.List           `tfsdk:"options"`
	Required    types.Bool           `tfsdk:"required"`
	Datasource  *DashboardDatasource `tfsdk:"datasource"`
}

func (in *DashboardInput) Attributes(ctx context.Context, d *diag.Diagnostics) *gqlclient.DashboardInputAttributes {
	return &gqlclient.DashboardInputAttributes{
		Name:        in.Name.ValueString(),
		Label:       in.Label.ValueStringPointer(),
		Description: in.Description.ValueStringPointer(),
		Type:        gqlclient.DashboardInputType(in.Type.ValueString()),
		Default:     in.Default.ValueStringPointer(),
		Options:     stringPointersFrom(in.Options, ctx, d),
		Required:    in.Required.ValueBoolPointer(),
		Datasource:  in.Datasource.Attributes(),
	}
}

func (in *DashboardInput) From(response *gqlclient.WorkbenchDashboardInputFragment, ctx context.Context, d *diag.Diagnostics) {
	in.Name = types.StringValue(response.Name)
	in.Label = types.StringPointerValue(response.Label)
	in.Description = types.StringPointerValue(response.Description)
	in.Type = types.StringValue(string(response.Type))
	in.Default = types.StringPointerValue(response.Default)
	in.Options = common.ListFrom(response.Options, in.Options, ctx, d)
	in.Required = types.BoolPointerValue(response.Required)
	in.Datasource = in.Datasource.from(response.Datasource)
}
