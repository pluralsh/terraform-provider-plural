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

var _ resource.Resource = &ObservabilityWebhookResource{}
var _ resource.ResourceWithImportState = &ObservabilityWebhookResource{}

func NewObservabilityWebhookResource() resource.Resource {
	return &ObservabilityWebhookResource{}
}

// ObservabilityWebhookResource defines the observability webhook resource implementation.
type ObservabilityWebhookResource struct {
	client *client.Client
}

func (r *ObservabilityWebhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_observability_webhook"
}

func (r *ObservabilityWebhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Internal identifier of this observability webhook.",
				MarkdownDescription: "Internal identifier of this observability webhook.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description:         "Human-readable name of this observability webhook.",
				MarkdownDescription: "Human-readable name of this observability webhook.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"type": schema.StringAttribute{
				Description:         "Observability webhook type.",
				MarkdownDescription: "Observability webhook type.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators: []validator.String{
					stringvalidator.OneOf(
						algorithms.Map(gqlclient.AllObservabilityWebhookType,
							func(t gqlclient.ObservabilityWebhookType) string { return string(t) })...),
				},
			},
			"url": schema.StringAttribute{
				Description:         "Observability webhook URL.",
				MarkdownDescription: "Observability webhook URL.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"secret": schema.StringAttribute{
				Description:         "Observability webhook secret.",
				MarkdownDescription: "Observability webhook secret.",
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *ObservabilityWebhookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Observability Webhook Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = data.Client
}

func (r *ObservabilityWebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := new(model.ObservabilityWebhook)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sc, err := r.client.UpsertObservabilityWebhook(ctx, data.Attributes())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create observability webhook, got error: %s", err))
		return
	}

	data.From(sc.UpsertObservabilityWebhook)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *ObservabilityWebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := new(model.ObservabilityWebhook)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetObservabilityWebhook(ctx, nil, data.Name.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read observability webhook, got error: %s", err))
		return
	}
	if response == nil || response.ObservabilityWebhook == nil {
		// Resource not found, remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	data.From(response.ObservabilityWebhook)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *ObservabilityWebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := new(model.ObservabilityWebhook)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpsertObservabilityWebhook(ctx, data.Attributes())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update observability webhook, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *ObservabilityWebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := new(model.ObservabilityWebhook)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.DeleteObservabilityWebhook(ctx, data.Id.ValueString()); err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete observability webhook, got error: %s", err))
		return
	}
}

func (r *ObservabilityWebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
