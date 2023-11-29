package resource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"terraform-provider-plural/internal/client"
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
				Description:         "Internal identifier of this GitRepository.",
				MarkdownDescription: "Internal identifier of this GitRepository.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Required:            true,
				Description:         "URL of this GitRepository.",
				MarkdownDescription: "URL of this GitRepository.",
			},
			"private_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         "SSH private key to use with this repo if an ssh url was given.",
				MarkdownDescription: "SSH private key to use with this repo if an ssh url was given.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("username"), path.MatchRoot("password"),
					),
					stringvalidator.AlsoRequires(path.MatchRoot("passphrase")),
				},
			},
			"passphrase": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         "Passphrase to decrypt the given private key.",
				MarkdownDescription: "Passphrase to decrypt the given private key.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("username"), path.MatchRoot("password"),
					),
					stringvalidator.AlsoRequires(path.MatchRoot("private_key")),
				},
			},
			"username": schema.StringAttribute{
				Optional:            true,
				Description:         "HTTP username for authenticated http repos, defaults to apiKey for GitHub.",
				MarkdownDescription: "HTTP username for authenticated http repos, defaults to apiKey for GitHub.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("private_key"), path.MatchRoot("passphrase"),
					),
					stringvalidator.AlsoRequires(path.MatchRoot("password")),
				},
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         "HTTP password for http authenticated repos.",
				MarkdownDescription: "HTTP password for http authenticated repos.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("private_key"), path.MatchRoot("passphrase"),
					),
					stringvalidator.AlsoRequires(path.MatchRoot("username")),
				},
			},
			"url_format": schema.StringAttribute{
				Optional:            true,
				Description:         "Similar to https_Path, a manually supplied url format for custom git. Should be something like {url}/tree/{ref}/{folder}.",
				MarkdownDescription: "Similar to https_Path, a manually supplied url format for custom git. Should be something like {url}/tree/{ref}/{folder}.",
			},
			"https_path": schema.StringAttribute{
				Optional:            true,
				Description:         "Manually supplied https path for non standard git setups. This is auto-inferred in many cases.",
				MarkdownDescription: "Manually supplied https path for non standard git setups. This is auto-inferred in many cases.",
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

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Git Repository Resource Configure Type",
			fmt.Sprintf(
				"Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	r.client = data.Client
}

func (r *GitRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := new(gitRepository)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.CreateGitRepository(ctx, data.Attributes())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create GitRepository, got error: %s", err))
		return
	}

	data.From(response.CreateGitRepository)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *GitRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := new(gitRepository)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetGitRepository(ctx, data.Id.ValueStringPointer(), data.Url.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get GitRepository, got error: %s", err))
		return
	}

	data.From(response.GitRepository)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *GitRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := new(gitRepository)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateGitRepository(ctx, data.Id.String(), data.Attributes())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update GitRepository, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *GitRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := new(gitRepository)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
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
