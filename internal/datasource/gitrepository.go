package datasource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/model"
	"terraform-provider-plural/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-plural/internal/client"
)

func NewGitRepositoryDataSource() datasource.DataSource {
	return &GitRepositoryDataSource{}
}

// GitRepositoryDataSource defines the GitRepository resource implementation.
type GitRepositoryDataSource struct {
	client *client.Client
}

func (r *GitRepositoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_git_repository"
}

func (r *GitRepositoryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "GitRepository resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal identifier of this GitRepository.",
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL of this repository.",
				Required:            true,
			},
			"private_key": schema.StringAttribute{
				MarkdownDescription: "SSH private key to use with this repo if an ssh url was given.",
				Validators:          []validator.String{stringvalidator.ConflictsWith(path.MatchRoot("username"), path.MatchRoot("password"))},
				Optional:            true,
			},
			"passphrase": schema.StringAttribute{
				MarkdownDescription: "Passphrase to decrypt the given private key.",
				Validators:          []validator.String{stringvalidator.ConflictsWith(path.MatchRoot("username"), path.MatchRoot("password"))},
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "HTTP username for authenticated http repos, defaults to apiKey for GitHub.",
				Validators:          []validator.String{stringvalidator.ConflictsWith(path.MatchRoot("private_key"), path.MatchRoot("passphrase"))},
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "HTTP password for http authenticated repos.",
				Validators:          []validator.String{stringvalidator.ConflictsWith(path.MatchRoot("private_key"), path.MatchRoot("passphrase"))},
				Optional:            true,
			},
			"url_format": schema.StringAttribute{
				MarkdownDescription: "Similar to https_Path, a manually supplied url format for custom git. Should be something like {url}/tree/{ref}/{folder}.",
				Optional:            true,
			},
			"https_path": schema.StringAttribute{
				MarkdownDescription: "Manually supplied https path for non standard git setups. This is auto-inferred in many cases.",
				Optional:            true,
			},
		},
	}
}

func (r *GitRepositoryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*provider.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Git Repository Resource Configure Type",
			fmt.Sprintf("Expected *provider.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = data.Client
}

func (r *GitRepositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured HTTP Client",
			"Expected configured HTTP client. Please report this issue to the provider developers.",
		)

		return
	}

	var data model.GitRepository

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	response, err := r.client.GetGitRepository(ctx, nil, data.Url.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get GitRepository, got error: %s", err))
		return
	}

	data.Id = types.StringValue(response.GitRepository.ID)
	data.Url = types.StringValue(response.GitRepository.URL)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
