package resource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/common"
	"terraform-provider-plural/internal/model"

	"terraform-provider-plural/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &WorkbenchResource{}
var _ resource.ResourceWithImportState = &WorkbenchResource{}

func NewWorkbenchResource() resource.Resource {
	return &WorkbenchResource{}
}

// WorkbenchResource defines the workbench resource implementation.
type WorkbenchResource struct {
	client *client.Client
}

func (r *WorkbenchResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_workbench"
}

func (r *WorkbenchResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Workbench resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Internal identifier of this workbench.",
				MarkdownDescription: "Internal identifier of this workbench.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description:         "Name of this workbench.",
				MarkdownDescription: "Name of this workbench.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"description": schema.StringAttribute{
				Description:         "Description of this workbench.",
				MarkdownDescription: "Description of this workbench.",
				Optional:            true,
			},
			"system_prompt": schema.StringAttribute{
				Description:         "System prompt for this workbench.",
				MarkdownDescription: "System prompt for this workbench.",
				Optional:            true,
			},
			"project_id": schema.StringAttribute{
				Description:         "ID of the project that this workbench belongs to.",
				MarkdownDescription: "ID of the project that this workbench belongs to.",
				Optional:            true,
				Computed:            true,
			},
			"repository_id": schema.StringAttribute{
				Description:         "The Git repository for this workbench.",
				MarkdownDescription: "The Git repository for this workbench.",
				Optional:            true,
			},
			"agent_runtime": schema.StringAttribute{
				Description:         "The runtime for the agent to use in the '<cluster-handle>/<agent-runtime>' format.",
				MarkdownDescription: "The runtime for the agent to use in the `<cluster-handle>/<agent-runtime>` format.",
				Optional:            true,
			},
			"configuration": schema.SingleNestedAttribute{
				Description:         "Configuration for this workbench.",
				MarkdownDescription: "Configuration for this workbench.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"coding": schema.SingleNestedAttribute{
						Description:         "Coding capabilities.",
						MarkdownDescription: "Coding capabilities.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"mode": schema.StringAttribute{
								Description:         "The coding mode to use. Valid values: ANALYZE or WRITE.",
								MarkdownDescription: "The coding mode to use. Valid values: `ANALYZE` or `WRITE`.",
								Optional:            true,
								Validators:          []validator.String{stringvalidator.OneOf("ANALYZE", "WRITE")},
							},
							"repositories": schema.SetAttribute{
								Description:         "Allowed repository identifiers.",
								MarkdownDescription: "Allowed repository identifiers.",
								Optional:            true,
								ElementType:         types.StringType,
							},
						},
					},
					"infrastructure": schema.SingleNestedAttribute{
						Description:         "Infrastructure capabilities.",
						MarkdownDescription: "Infrastructure capabilities.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"services": schema.BoolAttribute{
								Description:         "Whether to enable services capability.",
								MarkdownDescription: "Whether to enable services capability.",
								Optional:            true,
							},
							"stacks": schema.BoolAttribute{
								Description:         "Whether to enable stacks capability.",
								MarkdownDescription: "Whether to enable stacks capability.",
								Optional:            true,
							},
							"kubernetes": schema.BoolAttribute{
								Description:         "Whether to enable Kubernetes capability.",
								MarkdownDescription: "Whether to enable Kubernetes capability.",
								Optional:            true,
							},
						},
					},
				},
			},
			"skills": schema.SingleNestedAttribute{
				Description:         "Skills for this workbench.",
				MarkdownDescription: "Skills for this workbench.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"ref": schema.SingleNestedAttribute{
						Description:         "Git reference to the skill.",
						MarkdownDescription: "Git reference to the skill.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"ref": schema.StringAttribute{
								Description:         "Git reference to use.",
								MarkdownDescription: "Git reference to use.",
								Required:            true,
							},
							"folder": schema.StringAttribute{
								Description:         "The subdirectory in the Git repository to use.",
								MarkdownDescription: "The subdirectory in the Git repository to use.",
								Required:            true,
							},
							"files": schema.SetAttribute{
								Description:         "Files to include.",
								MarkdownDescription: "Files to include.",
								Optional:            true,
								ElementType:         types.StringType,
							},
						},
					},
					"files": schema.SetAttribute{
						Description:         "Files to include.",
						MarkdownDescription: "Files to include.",
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"tool_ids": schema.SetAttribute{
				Description:         "Tools for this workbench.",
				MarkdownDescription: "Tools for this workbench.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *WorkbenchResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Workbench Resource Configure Type",
			fmt.Sprintf(
				"Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)
		return
	}

	r.client = data.Client
}

func (r *WorkbenchResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := new(model.Workbench)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs, err := data.Attributes(r.client, ctx, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get attributes, got error: %s", err))
		return
	}

	response, err := r.client.CreateWorkbench(ctx, *attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workbench, got error: %s", err))
		return
	}

	data.From(response.CreateWorkbench, ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *WorkbenchResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := new(model.Workbench)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetWorkbench(ctx, data.Id.ValueStringPointer(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get workbench, got error: %s", err))
		return
	}
	if response == nil || response.Workbench == nil {
		// Resource not found, remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	data.From(response.Workbench, ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *WorkbenchResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := new(model.Workbench)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs, err := data.Attributes(r.client, ctx, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get attributes, got error: %s", err))
		return
	}

	_, err = r.client.UpdateWorkbench(ctx, data.Id.ValueString(), *attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workbench, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *WorkbenchResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := new(model.Workbench)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.DeleteWorkbench(ctx, data.Id.ValueString()); err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete workbench, got error: %s", err))
		return
	}
}

func (r *WorkbenchResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
