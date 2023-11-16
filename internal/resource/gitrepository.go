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

// GitRepositoryResourceModel describes the GitRepository resource data model.
type GitRepositoryResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Url        types.String `tfsdk:"url"`
	PrivateKey types.String `tfsdk:"private_key"`
	Passphrase types.String `tfsdk:"passphrase"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	UrlFormat  types.String `tfsdk:"url_format"`
	HttpsPath  types.String `tfsdk:"https_path"`
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

func (r *GitRepositoryResource) Configure(
	_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected GitRepository Resource Configure Type",
			fmt.Sprintf(
				"Expected *client.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	r.client = c
}

func (r *GitRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GitRepositoryResourceModel
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

	repository, err := r.client.CreateGitRepository(ctx, attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create GitRepository, got error: %s", err))
		return
	}

	// TODO: figure out if we need to read response and update state
	data.Id = types.StringValue(repository.CreateGitRepository.ID)

	tflog.Trace(ctx, "created a GitRepository")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GitRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GitRepositoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	repositories, err := r.client.ListGitRepositories(ctx, nil, nil, lo.ToPtr(int64(999)))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read GitRepository, got error: %s", err))
		return
	}

	var repository *consoleClient.GitRepositoryEdgeFragment
	for _, repo := range repositories.GitRepositories.Edges {
		if repo.Node.ID == data.Id.ValueString() {
			repository = repo
		}
	}

	if repository == nil {
		resp.Diagnostics.AddError(
			"Not Found", fmt.Sprintf("Unable to find GitRepository with ID: %s", data.Id.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GitRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data GitRepositoryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: figure out what can be updated
	//attrs := consoleClient.GitRepositoryUpdateAttributes{
	//	Handle: lo.ToPtr(data.Handle.String()),
	//}
	//GitRepository, err := r.client.UpdateGitRepository(ctx, data.Id.String(), attrs)
	//if err != nil {
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update GitRepository, got error: %s", err))
	//	return
	//}
	//
	//data.Handle = types.StringValue(*GitRepository.UpdateGitRepository.Handle)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GitRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GitRepositoryResourceModel
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
