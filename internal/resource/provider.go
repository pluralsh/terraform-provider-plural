package resource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/defaults"
	"terraform-provider-plural/internal/model"
	"terraform-provider-plural/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	console "github.com/pluralsh/console-client-go"
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
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable name of this provider. Globally unique.",
				Required:            true,
			},
			"cloud": schema.StringAttribute{
				MarkdownDescription: "The name of the cloud service for this provider.",
				Required:            true,
				Validators:          []validator.String{model.CloudValidator},
			},
			"cloud_settings": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"aws": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"access_key_id": schema.StringAttribute{
								Default:   defaults.Env("PLURAL_AWS_ACCESS_KEY_ID", ""),
								Required:  true, // TODO: Test Default and Required.
								Sensitive: true,
							},
							"secret_access_key": schema.StringAttribute{
								Default:   defaults.Env("PLURAL_AWS_SECRET_ACCESS_KEY", ""),
								Required:  true, // TODO: Test Default and Required.
								Sensitive: true,
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

	data, ok := req.ProviderData.(*provider.ProviderData)
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

	attrs := console.ClusterProviderAttributes{
		Name:  data.Name.ValueString(),
		Cloud: data.Cloud.ValueStringPointer(),
	}
	if model.IsCloud(data.Cloud.ValueString(), model.CloudAWS) {
		attrs.CloudSettings = &console.CloudProviderSettingsAttributes{
			Aws: &console.AwsSettingsAttributes{
				AccessKeyID:     data.CloudSettings.AWS.AccessKeyId.ValueString(),
				SecretAccessKey: data.CloudSettings.AWS.SecretAccessKey.ValueString(),
			},
		}
	}

	result, err := r.client.CreateClusterProvider(ctx, attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create provider, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created a provider")

	data.Name = types.StringValue(result.CreateClusterProvider.Name)
	data.Id = types.StringValue(result.CreateClusterProvider.ID)

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

	data.Id = types.StringValue(result.ClusterProvider.ID)
	data.Name = types.StringValue(result.ClusterProvider.Name)
	data.Cloud = types.StringValue(result.ClusterProvider.Cloud)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data model.Provider
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs := console.ClusterProviderUpdateAttributes{}
	if model.IsCloud(data.Cloud.ValueString(), model.CloudAWS) {
		attrs.CloudSettings = &console.CloudProviderSettingsAttributes{
			Aws: &console.AwsSettingsAttributes{
				AccessKeyID:     data.CloudSettings.AWS.AccessKeyId.ValueString(),
				SecretAccessKey: data.CloudSettings.AWS.SecretAccessKey.ValueString(),
			},
		}
	}

	_, err := r.client.UpdateClusterProvider(ctx, data.Id.ValueString(), attrs)
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
