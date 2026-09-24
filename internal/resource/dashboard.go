package resource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"
	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
)

var _ resource.Resource = &DashboardResource{}
var _ resource.ResourceWithImportState = &DashboardResource{}

func NewDashboardResource() resource.Resource {
	return &DashboardResource{}
}

type DashboardResource struct {
	client *client.Client
}

func (r *DashboardResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboard"
}

func (r *DashboardResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Workbench dashboard resource. A dashboard arranges graphs, backed by workbench observability tools, " +
			"on a grid. Graph `options` and datasource `input` are not returned by the Console API, so changes to them " +
			"made outside of Terraform are not detected, and they are not set on import.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Internal identifier of this dashboard.",
				MarkdownDescription: "Internal identifier of this dashboard.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"workbench_id": schema.StringAttribute{
				Description:         "ID of the workbench that owns this dashboard. A dashboard cannot be moved to a different workbench.",
				MarkdownDescription: "ID of the workbench that owns this dashboard. A dashboard cannot be moved to a different workbench.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:         "Dashboard name, unique within its workbench.",
				MarkdownDescription: "Dashboard name, unique within its workbench.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"description": schema.StringAttribute{
				Description:         "Dashboard description.",
				MarkdownDescription: "Dashboard description.",
				Optional:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"graphs": schema.ListNestedAttribute{
				Description:         "Graphs arranged on the dashboard grid. Graph identifiers have to be unique within the dashboard.",
				MarkdownDescription: "Graphs arranged on the dashboard grid. Graph identifiers have to be unique within the dashboard.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"identifier": schema.StringAttribute{
							Description:         "Stable identifier of the graph, unique within the dashboard.",
							MarkdownDescription: "Stable identifier of the graph, unique within the dashboard.",
							Required:            true,
							Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
						},
						"title": schema.StringAttribute{
							Description:         "Graph title.",
							MarkdownDescription: "Graph title.",
							Optional:            true,
							Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
						},
						"description": schema.StringAttribute{
							Description:         "Graph description.",
							MarkdownDescription: "Graph description.",
							Optional:            true,
							Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
						},
						"type": schema.StringAttribute{
							Description:         withAllowedValues("Graph visualization type. SECTION graphs group other graphs and MARKDOWN graphs require markdown.", enumValues(gqlclient.AllDashboardGraphType), false),
							MarkdownDescription: withAllowedValues("Graph visualization type. `SECTION` graphs group other graphs and `MARKDOWN` graphs require `markdown`.", enumValues(gqlclient.AllDashboardGraphType), true),
							Required:            true,
							Validators:          []validator.String{stringvalidator.OneOf(enumValues(gqlclient.AllDashboardGraphType)...)},
						},
						"section_id": schema.StringAttribute{
							Description:         "Identifier of the SECTION graph containing this graph. Sections cannot be nested.",
							MarkdownDescription: "Identifier of the `SECTION` graph containing this graph. Sections cannot be nested.",
							Optional:            true,
							Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
						},
						"markdown": schema.StringAttribute{
							Description:         "Markdown content for MARKDOWN graphs.",
							MarkdownDescription: "Markdown content for `MARKDOWN` graphs.",
							Optional:            true,
							Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
						},
						"options": schema.StringAttribute{
							Description:         "JSON-encoded visualization-specific display options, e.g. jsonencode({ collapsed = true }) for sections.",
							MarkdownDescription: "JSON-encoded visualization-specific display options, e.g. `jsonencode({ collapsed = true })` for sections.",
							Optional:            true,
						},
						"layout": schema.SingleNestedAttribute{
							Description:         "Grid position and size of the graph.",
							MarkdownDescription: "Grid position and size of the graph.",
							Required:            true,
							Attributes: map[string]schema.Attribute{
								"x": gridAttribute("Zero-based horizontal grid coordinate.", 0),
								"y": gridAttribute("Zero-based vertical grid coordinate.", 0),
								"w": gridAttribute("Width in grid columns.", 1),
								"h": gridAttribute("Height in grid rows.", 1),
							},
						},
						"datasource": dashboardDatasourceSchema("Observability tool call used to fetch data for the graph."),
					},
				},
			},
			"inputs": schema.ListNestedAttribute{
				Description:         "User-configurable dashboard variables. Input names have to be unique within the dashboard.",
				MarkdownDescription: "User-configurable dashboard variables. Input names have to be unique within the dashboard.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "Variable name referenced by graph datasource inputs.",
							MarkdownDescription: "Variable name referenced by graph datasource inputs.",
							Required:            true,
							Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
						},
						"label": schema.StringAttribute{
							Description:         "Human-readable input label.",
							MarkdownDescription: "Human-readable input label.",
							Optional:            true,
							Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
						},
						"description": schema.StringAttribute{
							Description:         "Input description.",
							MarkdownDescription: "Input description.",
							Optional:            true,
							Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
						},
						"type": schema.StringAttribute{
							Description:         withAllowedValues("Input control type.", enumValues(gqlclient.AllDashboardInputType), false),
							MarkdownDescription: withAllowedValues("Input control type.", enumValues(gqlclient.AllDashboardInputType), true),
							Required:            true,
							Validators:          []validator.String{stringvalidator.OneOf(enumValues(gqlclient.AllDashboardInputType)...)},
						},
						"default": schema.StringAttribute{
							Description:         "Default input value.",
							MarkdownDescription: "Default input value.",
							Optional:            true,
							Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
						},
						"options": schema.ListAttribute{
							Description:         "Allowed values for SELECT inputs.",
							MarkdownDescription: "Allowed values for `SELECT` inputs.",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"required": schema.BoolAttribute{
							Description:         "Whether a value is required when rendering. Defaults to false.",
							MarkdownDescription: "Whether a value is required when rendering. Defaults to `false`.",
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
						"datasource": dashboardDatasourceSchema("Tool query used to populate input options, e.g. a metric label search."),
					},
				},
			},
		},
	}
}

func gridAttribute(description string, minimum int64) schema.Int64Attribute {
	return schema.Int64Attribute{
		Description:         description,
		MarkdownDescription: description,
		Required:            true,
		Validators:          []validator.Int64{int64validator.AtLeast(minimum)},
	}
}

func dashboardDatasourceSchema(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description:         description,
		MarkdownDescription: description,
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description:         withAllowedValues("Kind of data returned by the datasource.", enumValues(gqlclient.AllDashboardDatasourceType), false),
				MarkdownDescription: withAllowedValues("Kind of data returned by the datasource.", enumValues(gqlclient.AllDashboardDatasourceType), true),
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf(enumValues(gqlclient.AllDashboardDatasourceType)...)},
			},
			"tool": schema.StringAttribute{
				Description:         "Name of the workbench observability tool used to fetch the data.",
				MarkdownDescription: "Name of the workbench observability tool used to fetch the data.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"input": schema.StringAttribute{
				Description:         "JSON-encoded input passed to the tool, e.g. jsonencode({ query = \"up\" }). Defaults to an empty object.",
				MarkdownDescription: "JSON-encoded input passed to the tool, e.g. `jsonencode({ query = \"up\" })`. Defaults to an empty object.",
				Optional:            true,
			},
		},
	}
}

func (r *DashboardResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Dashboard Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = data.Client
}

func (r *DashboardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := new(model.Dashboard)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs := data.Attributes(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.CreateDashboard(ctx, attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create dashboard, got error: %s", err))
		return
	}

	data.From(response.CreateDashboard, ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *DashboardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := new(model.Dashboard)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetWorkbenchDashboard(ctx, data.Id.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get dashboard, got error: %s", err))
		return
	}
	if response == nil || response.WorkbenchDashboard == nil || client.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}

	data.From(response.WorkbenchDashboard, ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *DashboardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := new(model.Dashboard)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs := data.Attributes(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateDashboard(ctx, data.Id.ValueString(), attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update dashboard, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *DashboardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := new(model.Dashboard)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.DeleteDashboard(ctx, data.Id.ValueString()); err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete dashboard, got error: %s", err))
		return
	}
}

func (r *DashboardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
