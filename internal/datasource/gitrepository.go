package datasource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	consoleClient "github.com/pluralsh/console-client-go"
	"github.com/samber/lo"

	"terraform-provider-plural/internal/client"
)

func NewGitRepositoryDataSource() datasource.DataSource {
	return &GitRepositoryDataSource{}
}

// GitRepositoryDataSource defines the GitRepository resource implementation.
type GitRepositoryDataSource struct {
	client *client.Client
}

// GitRepositoryModel describes the GitRepository data model.
type GitRepositoryModel struct {
	Id         types.String `tfsdk:"id"`
	Url        types.String `tfsdk:"url"`
	PrivateKey types.String `tfsdk:"private_key"`
	Passphrase types.String `tfsdk:"passphrase"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	UrlFormat  types.String `tfsdk:"url_format"`
	HttpsPath  types.String `tfsdk:"https_path"`
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

	c, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected GitRepository Resource Configure Type",
			fmt.Sprintf("Expected *console.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = c
}

func (r *GitRepositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured HTTP Client",
			"Expected configured HTTP client. Please report this issue to the provider developers.",
		)

		return
	}

	var data GitRepositoryModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	repositories, err := r.client.ListGitRepositories(ctx, nil, nil, lo.ToPtr(int64(999)))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read GitRepository, got error: %s", err))
		return
	}

	var repository *consoleClient.GitRepositoryEdgeFragment
	for _, repo := range repositories.GitRepositories.Edges {
		if repo.Node.URL == data.Url.ValueString() {
			repository = repo
		}
	}

	if repository == nil {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Unable to find GitRepository with ID: %s", data.Id.ValueString()))
		return
	}

	data.Id = types.StringValue(repository.Node.ID)
	data.Url = types.StringValue(repository.Node.URL)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
