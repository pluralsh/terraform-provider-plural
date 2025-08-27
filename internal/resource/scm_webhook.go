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
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/pluralsh/polly/algorithms"
)

var _ resource.Resource = &SCMWebhookResource{}
var _ resource.ResourceWithImportState = &SCMWebhookResource{}

func NewSCMWebhookResource() resource.Resource {
	return &SCMWebhookResource{}
}

// SCMWebhookResource defines the SCM webhook resource implementation.
type SCMWebhookResource struct {
	client *client.Client
}

func (r *SCMWebhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scm_webhook"
}

func (r *SCMWebhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Internal identifier of this SCM webhook.",
				MarkdownDescription: "Internal identifier of this SCM webhook.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description:         "Human-readable name of this SCM webhook.",
				MarkdownDescription: "Human-readable name of this SCM webhook.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"owner": schema.StringAttribute{
				Description:         "Owner for this webhook in your SCM, eg. a GitHub org or GitLab group.",
				MarkdownDescription: "Owner for this webhook in your SCM, eg. a GitHub org or GitLab group.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"type": schema.StringAttribute{
				Description:         "SCM webhook type.",
				MarkdownDescription: "SCM webhook type.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					stringvalidator.OneOf(
						algorithms.Map(gqlclient.AllScmType,
							func(t gqlclient.ScmType) string { return string(t) })...),
				},
			},
			"url": schema.StringAttribute{
				Description:         "SCM webhook URL.",
				MarkdownDescription: "SCM webhook URL.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"hmac": schema.StringAttribute{
				Description:         "Secret token for authenticating this webhook via hmac signature.",
				MarkdownDescription: "Secret token for authenticating this webhook via hmac signature.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *SCMWebhookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected SCM Webhook Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = data.Client
}

func (r *SCMWebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := new(model.SCMWebhook)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sc, err := r.client.CreateScmWebhookPointer(ctx, data.Attributes())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create SCM webhook, got error: %s", err))
		return
	}

	data.From(sc.CreateScmWebhookPointer)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *SCMWebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := new(model.SCMWebhook)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetScmWebhook(ctx, data.Id.ValueStringPointer(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SCM webhook, got error: %s", err))
		return
	}
	if response == nil || response.ScmWebhook == nil {
		// Resource not found, remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	data.From(response.ScmWebhook)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *SCMWebhookResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// Ignore.
}

func (r *SCMWebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := new(model.SCMWebhook)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.DeleteScmWebhook(ctx, data.Id.ValueString()); err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete SCM webhook, got error: %s", err))
		return
	}
}

func (r *SCMWebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
