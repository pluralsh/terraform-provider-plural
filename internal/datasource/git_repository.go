package datasource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
				Optional:            true,
				Description:         "Internal identifier of this GitRepository.",
				MarkdownDescription: "Internal identifier of this GitRepository.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("url")),
				},
			},
			"url": schema.StringAttribute{
				Optional:            true,
				Description:         "URL of this GitRepository.",
				MarkdownDescription: "URL of this GitRepository.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("id")),
				},
			},
		},
	}
}

func (r *GitRepositoryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Git Repository Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = data.Client
}

func (r *GitRepositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := new(gitRepository)
	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetGitRepository(ctx, data.Id.ValueStringPointer(), data.Url.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get GitRepository, got error: %s", err))
		return
	}

	if response == nil || response.GitRepository == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to find GitRepository"))
		return
	}

	data.From(response.GitRepository)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
