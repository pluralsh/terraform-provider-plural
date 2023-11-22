package resource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &providerResource{}
var _ resource.ResourceWithImportState = &providerResource{}

func NewProviderResource() resource.Resource {
	return &providerResource{}
}

// providerResource defines the provider resource implementation.
type providerResource struct {
	client     *client.Client
	consoleUrl string
}

func (r *providerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider"
}

func (r *providerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A representation of a provider you can deploy your clusters to.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal identifier of this provider.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"editable": schema.BoolAttribute{
				MarkdownDescription: "Whether this provider is editable.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable name of this provider. Globally unique.",
				Required:            true,
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: "The namespace the Cluster API resources are deployed into.",
				Optional:            true,
			},
			"cloud": schema.StringAttribute{
				MarkdownDescription: "The name of the cloud service for this provider.",
				Required:            true,
				Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
					model.CloudAWS.String(), model.CloudAzure.String(), model.CloudGCP.String())},
			},
			"cloud_settings": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"aws": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"access_key_id": schema.StringAttribute{
								MarkdownDescription: "",
								Required:            true,
								Sensitive:           true,
							},
							"secret_access_key": schema.StringAttribute{
								MarkdownDescription: "",
								Required:            true,
								Sensitive:           true,
							},
						},
					},
					"azure": schema.SingleNestedAttribute{
						MarkdownDescription: "Azure cloud settings that will be used by this provider to create clusters.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"subscription_id": schema.StringAttribute{
								MarkdownDescription: "GUID of the Azure subscription",
								Required:            true,
								Sensitive:           true,
							},
							"tenant_id": schema.StringAttribute{
								MarkdownDescription: "The unique identifier of the Azure Active Directory instance.",
								Required:            true,
								Sensitive:           true,
							},
							"client_id": schema.StringAttribute{
								MarkdownDescription: "The unique identifier of an application created in the Azure Active Directory.",
								Required:            true,
								Sensitive:           true,
							},
							"client_secret": schema.StringAttribute{
								MarkdownDescription: "A string value your app can use in place of a certificate to identity itself. Sometimes called an application password.",
								Required:            true,
								Sensitive:           true,
							},
						},
					},
					"gcp": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"credentials": schema.StringAttribute{
								MarkdownDescription: "",
								Required:            true,
								Sensitive:           true,
							},
						},
					},
				},
				MarkdownDescription: "Cloud-specific settings for a provider.",
				Required:            true,
			},
		},
	}
}

func (r *providerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*model.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Resource Configure Type",
			fmt.Sprintf("Expected *model.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = data.Client
	r.consoleUrl = data.ConsoleUrl
}

func (r *providerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data model.Provider
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateClusterProvider(ctx, data.Attributes())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create provider, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a provider")

	data.From(result.CreateClusterProvider)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data model.Provider
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetClusterProvider(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read provider, got error: %s", err))
		return
	}

	data.From(result.ClusterProvider)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data model.Provider
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateClusterProvider(ctx, data.Id.ValueString(), data.UpdateAttributes())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update provider, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data model.Provider
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteCluster(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete provider, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted the provider")
}

func (r *providerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
