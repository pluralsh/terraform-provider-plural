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
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/pluralsh/polly/algorithms"
)

var _ resource.Resource = &OIDCProviderResource{}
var _ resource.ResourceWithImportState = &OIDCProviderResource{}

func NewOIDCProviderResourceResource() resource.Resource {
	return &OIDCProviderResource{}
}

// OIDCProviderResource defines the OIDC provider resource implementation.
type OIDCProviderResource struct {
	client *client.Client
}

func (r *OIDCProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oidc_provider"
}

func (r *OIDCProviderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Internal identifier of this OIDC provider.",
				MarkdownDescription: "Internal identifier of this OIDC provider.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description:         "Human-readable name of this OIDC provider.",
				MarkdownDescription: "Human-readable name of this OIDC provider.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						algorithms.Map(gqlclient.AllOidcProviderType,
							func(t gqlclient.OidcProviderType) string { return string(t) })...),
				},
			},
			"description": schema.StringAttribute{
				Description:         "Description of this OIDC provider.",
				MarkdownDescription: "Description of this OIDC provider.",
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"client_id": schema.StringAttribute{
				Computed:      true,
				Sensitive:     true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"client_secret": schema.StringAttribute{
				Computed:      true,
				Sensitive:     true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"auth_method": schema.StringAttribute{
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						algorithms.Map(gqlclient.AllOidcAuthMethod,
							func(t gqlclient.OidcAuthMethod) string { return string(t) })...),
				},
			},
			"redirect_uris": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *OIDCProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected OIDC Provider Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = data.Client
}

func (r *OIDCProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := new(model.OIDCProvider)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateOIDCProvider(ctx, data.TypeAttribute(), data.Attributes(ctx, resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create OIDC provider, got error: %s", err))
		return
	}

	data.From(result.CreateOidcProvider, ctx, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *OIDCProviderResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// Ignore.
}

func (r *OIDCProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := new(model.OIDCProvider)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.UpdateOIDCProvider(ctx, data.ID.ValueString(), data.TypeAttribute(), data.Attributes(ctx, resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update OIDC provider, got error: %s", err))
		return
	}

	data.From(result.UpdateOidcProvider, ctx, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *OIDCProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := new(model.OIDCProvider)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteOIDCProvider(ctx, data.ID.ValueString(), data.TypeAttribute())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete OIDC provider, got error: %s", err))
		return
	}
}

func (r *OIDCProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
