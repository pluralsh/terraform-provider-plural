package resource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"
	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ resource.Resource = &WorkbenchWebhookResource{}
var _ resource.ResourceWithImportState = &WorkbenchWebhookResource{}

func NewWorkbenchWebhookResource() resource.Resource {
	return &WorkbenchWebhookResource{}
}

type WorkbenchWebhookResource struct {
	client *client.Client
}

func (r *WorkbenchWebhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workbench_webhook"
}

func (r *WorkbenchWebhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Workbench webhook resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Internal identifier of this workbench webhook.",
				MarkdownDescription: "Internal identifier of this workbench webhook.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"workbench_id": schema.StringAttribute{
				Description:         "ID of the workbench this webhook belongs to.",
				MarkdownDescription: "ID of the workbench this webhook belongs to.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Description:         "Unique name of this webhook trigger within a workbench.",
				MarkdownDescription: "Unique name of this webhook trigger within a workbench.",
				Required:            true,
			},
			"prompt": schema.StringAttribute{
				Description:         "Prompt to run when this webhook trigger matches.",
				MarkdownDescription: "Prompt to run when this webhook trigger matches.",
				Optional:            true,
			},
			"webhook_id": schema.StringAttribute{
				Description:         "Observability webhook ID that sends events to this trigger.",
				MarkdownDescription: "Observability webhook ID that sends events to this trigger.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("issue_webhook_id")),
				},
			},
			"issue_webhook_id": schema.StringAttribute{
				Description:         "Issue webhook ID that sends events to this trigger.",
				MarkdownDescription: "Issue webhook ID that sends events to this trigger.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("webhook_id")),
				},
			},
			"matches": schema.SingleNestedAttribute{
				Description:         "Webhook payload matching rules.",
				MarkdownDescription: "Webhook payload matching rules.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"regex": schema.StringAttribute{
						Description:         "Regex expression to match webhook payloads.",
						MarkdownDescription: "Regex expression to match webhook payloads.",
						Optional:            true,
					},
					"substring": schema.StringAttribute{
						Description:         "Substring expression to match webhook payloads.",
						MarkdownDescription: "Substring expression to match webhook payloads.",
						Optional:            true,
					},
					"case_insensitive": schema.BoolAttribute{
						Description:         "Whether matching should be case-insensitive.",
						MarkdownDescription: "Whether matching should be case-insensitive.",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (r *WorkbenchWebhookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Workbench Webhook Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = data.Client
}

func (r *WorkbenchWebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := new(model.WorkbenchWebhook)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.CreateWorkbenchWebhook(ctx, data.WorkbenchID.ValueString(), data.Attributes())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workbench webhook, got error: %s", err))
		return
	}

	data.From(response.CreateWorkbenchWebhook)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *WorkbenchWebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := new(model.WorkbenchWebhook)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetWorkbenchWebhook(ctx, data.Id.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get workbench webhook, got error: %s", err))
		return
	}
	if response == nil || response.GetWorkbenchWebhook == nil || client.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}

	data.From(response.GetWorkbenchWebhook)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *WorkbenchWebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := new(model.WorkbenchWebhook)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateWorkbenchWebhook(ctx, data.Id.ValueString(), data.Attributes())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workbench webhook, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *WorkbenchWebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := new(model.WorkbenchWebhook)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.DeleteWorkbenchWebhook(ctx, data.Id.ValueString()); err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete workbench webhook, got error: %s", err))
		return
	}
}

func (r *WorkbenchWebhookResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
