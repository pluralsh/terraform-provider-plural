package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	consoleClient "github.com/pluralsh/console-client-go"
	"github.com/samber/lo"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/model"
)

var _ resource.Resource = &GitRepositoryResource{}
var _ resource.ResourceWithImportState = &GitRepositoryResource{}

func NewGitRepositoryResource() resource.Resource {
	return &GitRepositoryResource{}
}

// GitRepositoryResource defines the GitRepository resource implementation.
type GitRepositoryResource struct {
	client *client.Client
}

func (r *GitRepositoryResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_git_repository"
}

func (r *GitRepositoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "GitRepository resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal identifier of this GitRepository.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL of this repository.",
				Required:            true,
			},
			"private_key": schema.StringAttribute{
				MarkdownDescription: "SSH private key to use with this repo if an ssh url was given.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("username"), path.MatchRoot("password"),
					),
					stringvalidator.AlsoRequires(path.MatchRoot("passphrase")),
				},
				Optional:  true,
				Sensitive: true,
			},
			"passphrase": schema.StringAttribute{
				MarkdownDescription: "Passphrase to decrypt the given private key.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("username"), path.MatchRoot("password"),
					),
					stringvalidator.AlsoRequires(path.MatchRoot("private_key")),
				},
				Optional:  true,
				Sensitive: true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "HTTP username for authenticated http repos, defaults to apiKey for GitHub.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("private_key"), path.MatchRoot("passphrase"),
					),
					stringvalidator.AlsoRequires(path.MatchRoot("password")),
				},
				Optional: true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "HTTP password for http authenticated repos.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("private_key"), path.MatchRoot("passphrase"),
					),
					stringvalidator.AlsoRequires(path.MatchRoot("username")),
				},
				Optional:  true,
				Sensitive: true,
			},
			"url_format": schema.StringAttribute{
				MarkdownDescription: "Similar to `https_path`, a manually supplied url format for custom Git. Should be something like `{url}/tree/{ref}/{folder}`.",
				Optional:            true,
			},
			"https_path": schema.StringAttribute{
				MarkdownDescription: "Manually supplied https path for non standard Git setups. This is auto-inferred in many cases.",
				Optional:            true,
			},
		},
	}
}

func (r *GitRepositoryResource) Configure(
	_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*model.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Git Repository Resource Configure Type",
			fmt.Sprintf(
				"Expected *model.ProviderData, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	r.client = data.Client
}

func (r *GitRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data model.GitRepository
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs := consoleClient.GitAttributes{
		URL:        data.Url.ValueString(),
		PrivateKey: lo.ToPtr(data.PrivateKey.ValueString()),
		Passphrase: lo.ToPtr(data.Passphrase.ValueString()),
		Username:   lo.ToPtr(data.Username.ValueString()),
		Password:   lo.ToPtr(data.Password.ValueString()),
		HTTPSPath:  lo.ToPtr(data.HttpsPath.ValueString()),
		URLFormat:  lo.ToPtr(data.UrlFormat.ValueString()),
	}

	result, err := r.client.CreateGitRepository(ctx, attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create GitRepository, got error: %s", err))
		return
	}

	// TODO: figure out if we need to read response and update state
	data.Id = types.StringValue(result.CreateGitRepository.ID)

	tflog.Trace(ctx, "created a GitRepository")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GitRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data model.GitRepository
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetGitRepository(ctx, data.Id.ValueStringPointer(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get GitRepository, got error: %s", err))
		return
	}

	data.Id = types.StringValue(result.GitRepository.ID)
	data.Url = types.StringValue(result.GitRepository.URL)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GitRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data model.GitRepository
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs := consoleClient.GitAttributes{
		URL:        data.Url.ValueString(),
		Username:   data.Username.ValueStringPointer(),
		Password:   data.Password.ValueStringPointer(),
		PrivateKey: data.PrivateKey.ValueStringPointer(),
		Passphrase: data.Passphrase.ValueStringPointer(),
		HTTPSPath:  data.HttpsPath.ValueStringPointer(),
		URLFormat:  data.UrlFormat.ValueStringPointer(),
	}

	result, err := r.client.UpdateGitRepository(ctx, data.Id.String(), attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update GitRepository, got error: %s", err))
		return
	}

	data.Id = types.StringValue(result.UpdateGitRepository.ID)
	data.Url = types.StringValue(result.UpdateGitRepository.URL)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GitRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data model.GitRepository
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteGitRepository(ctx, data.Id.String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete GitRepository, got error: %s", err))
		return
	}
}

func (r *GitRepositoryResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
