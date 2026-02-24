package resource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/common"
	"terraform-provider-plural/internal/model"
	planmod "terraform-provider-plural/internal/planmodifier"

	"terraform-provider-plural/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
)

var _ resource.Resource = &WorkbenchToolResource{}
var _ resource.ResourceWithImportState = &WorkbenchToolResource{}

func NewWorkbenchToolResource() resource.Resource {
	return &WorkbenchToolResource{}
}

// WorkbenchToolResource defines the workbench tool resource implementation.
type WorkbenchToolResource struct {
	client *client.Client
}

func (r *WorkbenchToolResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_workbench_tool"
}

func (r *WorkbenchToolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Workbench tool resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Internal identifier of this workbench tool.",
				MarkdownDescription: "Internal identifier of this workbench tool.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description:         "Name of this workbench tool.",
				MarkdownDescription: "Name of this workbench tool.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tool": schema.StringAttribute{
				Description:         "Workbench tool type.",
				MarkdownDescription: "Workbench tool type.",
				Required:            true,
				Validators: []validator.String{stringvalidator.OneOf(
					lo.Map(console.AllWorkbenchToolType, func(item console.WorkbenchToolType, index int) string {
						return string(item)
					})...),
				},
			},
			"categories": schema.SetAttribute{
				Description:         "Categories of this workbench tool.",
				MarkdownDescription: "Categories of this workbench tool.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"project_id": schema.StringAttribute{
				Description:         "ID of the project that this workbench belongs to.",
				MarkdownDescription: "ID of the project that this workbench belongs to.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"configuration": schema.SingleNestedAttribute{
				Description:         "Configuration of this workbench tool.",
				MarkdownDescription: "Configuration of this workbench tool.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"http": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "HTTP configuration of this workbench tool.",
						MarkdownDescription: "HTTP configuration of this workbench tool.",
						PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Attributes: map[string]schema.Attribute{
							"url": schema.StringAttribute{
								Description:         "The request URL.",
								MarkdownDescription: "The request URL.",
								Required:            true,
							},
							"method": schema.StringAttribute{
								Description:         "The HTTP method.",
								MarkdownDescription: "The HTTP method.",
								Required:            true,
								PlanModifiers:       []planmodifier.String{planmod.UppercaseString()},
							},
							"headers": schema.MapAttribute{
								Description:         "The request headers.",
								MarkdownDescription: "The request headers.",
								ElementType:         types.StringType,
								Optional:            true,
							},
							"body": schema.StringAttribute{
								Description:         "The request body.",
								MarkdownDescription: "The request body.",
								Optional:            true,
							},
							"input_schema": schema.StringAttribute{
								Description:         "The JSON schema for the tool input.",
								MarkdownDescription: "The JSON schema for the tool input.",
								Required:            true,
							},
						},
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *WorkbenchToolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Workbench Tool Resource Configure Type",
			fmt.Sprintf(
				"Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)
		return
	}

	r.client = data.Client
}

func (r *WorkbenchToolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := new(model.WorkbenchTool)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs, err := data.Attributes(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get attributes, got error: %s", err))
		return
	}

	response, err := r.client.CreateWorkbenchTool(ctx, *attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workbench tool, got error: %s", err))
		return
	}

	data.From(response.CreateWorkbenchTool, ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *WorkbenchToolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := new(model.WorkbenchTool)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetWorkbenchTool(ctx, data.Id.ValueStringPointer(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get workbench tool, got error: %s", err))
		return
	}
	if response == nil || response.WorkbenchTool == nil {
		// Resource not found, remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	data.From(response.WorkbenchTool, ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *WorkbenchToolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := new(model.WorkbenchTool)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs, err := data.Attributes(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get attributes, got error: %s", err))
		return
	}

	_, err = r.client.UpdateWorkbenchTool(ctx, data.Id.ValueString(), *attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workbench tool, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *WorkbenchToolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := new(model.WorkbenchTool)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.DeleteWorkbenchTool(ctx, data.Id.ValueString()); err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete workbench tool, got error: %s", err))
		return
	}
}

func (r *WorkbenchToolResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
