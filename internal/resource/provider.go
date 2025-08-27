package resource

import (
	"context"
	"fmt"
	"time"

	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"k8s.io/apimachinery/pkg/util/wait"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"
	"terraform-provider-plural/internal/defaults"
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
				Description:         "Internal identifier of this provider.",
				MarkdownDescription: "Internal identifier of this provider.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"editable": schema.BoolAttribute{
				Description:         "Whether this provider is editable.",
				MarkdownDescription: "Whether this provider is editable.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Human-readable name of this provider. Globally unique.",
				MarkdownDescription: "Human-readable name of this provider. Globally unique.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"namespace": schema.StringAttribute{
				Description:         "The namespace the Cluster API resources are deployed into.",
				MarkdownDescription: "The namespace the Cluster API resources are deployed into.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"cloud": schema.StringAttribute{
				Description:         "The name of the cloud service for this provider.",
				MarkdownDescription: "The name of the cloud service for this provider.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
					common.CloudAWS.String(), common.CloudAzure.String(), common.CloudGCP.String())},
			},
			"cloud_settings": schema.SingleNestedAttribute{
				Description:         "Cloud-specific settings for a provider. See https://docs.plural.sh/reference/configuring-cloud-provider#permissions for more information about required permissions.",
				MarkdownDescription: "Cloud-specific settings for a provider. See [our docs](https://docs.plural.sh/reference/configuring-cloud-provider#permissions) for more information about required permissions.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"aws": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"access_key_id": schema.StringAttribute{
								Description:         "ID of an access key for an IAM user or the AWS account root user. Can be sourced from PLURAL_AWS_ACCESS_KEY_ID.",
								MarkdownDescription: "ID of an access key for an IAM user or the AWS account root user. Can be sourced from `PLURAL_AWS_ACCESS_KEY_ID`.",
								Optional:            true,
								Sensitive:           true,
								Computed:            true,
								Default:             defaults.Env("PLURAL_AWS_ACCESS_KEY_ID", ""),
							},
							"secret_access_key": schema.StringAttribute{
								Description:         "Secret of an access key for an IAM user or the AWS account root user. Can be sourced from PLURAL_AWS_SECRET_ACCESS_KEY.",
								MarkdownDescription: "Secret of an access key for an IAM user or the AWS account root user. Can be sourced from `PLURAL_AWS_SECRET_ACCESS_KEY`.",
								Optional:            true,
								Sensitive:           true,
								Computed:            true,
								Default:             defaults.Env("PLURAL_AWS_SECRET_ACCESS_KEY", ""),
							},
						},
						Validators: []validator.Object{
							objectvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("azure"),
								path.MatchRelative().AtParent().AtName("gcp"),
							),
						},
					},
					"azure": schema.SingleNestedAttribute{
						Description:         "Azure cloud settings that will be used by this provider to create clusters.",
						MarkdownDescription: "Azure cloud settings that will be used by this provider to create clusters.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"subscription_id": schema.StringAttribute{
								Description:         "GUID of the Azure subscription. Can be sourced from PLURAL_AZURE_SUBSCRIPTION_ID.",
								MarkdownDescription: "GUID of the Azure subscription. Can be sourced from `PLURAL_AZURE_SUBSCRIPTION_ID`.",
								Optional:            true,
								Sensitive:           true,
								Computed:            true,
								Default:             defaults.Env("PLURAL_AZURE_SUBSCRIPTION_ID", ""),
							},
							"tenant_id": schema.StringAttribute{
								Description:         "The unique identifier of the Azure Active Directory instance. Can be sourced from PLURAL_AZURE_TENANT_ID.",
								MarkdownDescription: "The unique identifier of the Azure Active Directory instance. Can be sourced from `PLURAL_AZURE_TENANT_ID`.",
								Optional:            true,
								Sensitive:           true,
								Computed:            true,
								Default:             defaults.Env("PLURAL_AZURE_TENANT_ID", ""),
							},
							"client_id": schema.StringAttribute{
								Description:         "The unique identifier of an application created in the Azure Active Directory. Can be sourced from PLURAL_AZURE_CLIENT_ID.",
								MarkdownDescription: "The unique identifier of an application created in the Azure Active Directory. Can be sourced from `PLURAL_AZURE_CLIENT_ID`.",
								Optional:            true,
								Sensitive:           true,
								Computed:            true,
								Default:             defaults.Env("PLURAL_AZURE_CLIENT_ID", ""),
							},
							"client_secret": schema.StringAttribute{
								Description:         "A string value your app can use in place of a certificate to identity itself. Sometimes called an application password. Can be sourced from PLURAL_AZURE_CLIENT_SECRET.",
								MarkdownDescription: "A string value your app can use in place of a certificate to identity itself. Sometimes called an application password. Can be sourced from `PLURAL_AZURE_CLIENT_SECRET`.",
								Optional:            true,
								Sensitive:           true,
								Computed:            true,
								Default:             defaults.Env("PLURAL_AZURE_CLIENT_SECRET", ""),
							},
						},
						Validators: []validator.Object{
							objectvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("aws"),
								path.MatchRelative().AtParent().AtName("gcp"),
							),
						},
					},
					"gcp": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "GCP cloud settings that will be used by this provider to create clusters.",
						MarkdownDescription: "GCP cloud settings that will be used by this provider to create clusters.",
						Attributes: map[string]schema.Attribute{
							"credentials": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Sensitive:           true,
								Default:             defaults.Env("PLURAL_GCP_CREDENTIALS", ""),
								Description:         "Base64 encoded GCP credentials.json file. It's recommended to use custom service account. Can be sourced from PLURAL_GCP_CREDENTIALS.",
								MarkdownDescription: "Base64 encoded GCP credentials.json file. It's recommended to use custom service account. Can be sourced from `PLURAL_GCP_CREDENTIALS`.",
							},
						},
						Validators: []validator.Object{
							objectvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("aws"),
								path.MatchRelative().AtParent().AtName("azure"),
							),
						},
					},
				},
			},
		},
	}
}

func (r *providerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = data.Client
	r.consoleUrl = data.ConsoleUrl
}

func (r *providerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data model.ProviderExtended
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateClusterProvider(ctx, data.Attributes())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create provider, got error: %s", err))
		return
	}

	data.From(result.CreateClusterProvider)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data model.ProviderExtended
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetClusterProvider(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read provider, got error: %s", err))
		return
	}
	if result == nil || result.ClusterProvider == nil {
		// Resource not found, remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	data.From(result.ClusterProvider)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data model.ProviderExtended
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
	var data model.ProviderExtended
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.DeleteClusterProvider(ctx, data.Id.ValueString()); err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete provider, got error: %s", err))
		return
	}

	if err := wait.PollUntilContextTimeout(ctx, 10*time.Second, 10*time.Minute, true, func(ctx context.Context) (bool, error) {
		response, err := r.client.GetClusterProvider(ctx, data.Id.ValueString())
		if client.IsNotFound(err) {
			return true, nil
		}

		if err == nil && (response == nil || response.ClusterProvider == nil) {
			return true, nil
		}

		return false, err
	}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error while watiting for provider to be deleted, got error: %s", err))
		return
	}
}

func (r *providerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
