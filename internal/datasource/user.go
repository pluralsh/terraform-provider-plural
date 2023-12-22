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
	console "github.com/pluralsh/console-client-go"
)

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

type userDataSource struct {
	client *client.Client
}

func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A representation of a user to authenticate to your plural console.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Internal identifier of this user.",
				Validators:          []validator.String{stringvalidator.ExactlyOneOf(path.MatchRoot("email"))},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the user.",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Email address of this user.",
				Computed:            true,
				Validators:          []validator.String{stringvalidator.ExactlyOneOf(path.MatchRoot("id"))},
			},
		},
	}
}

func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected User Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = data.Client
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data user
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Email.IsNull() {
		resp.Diagnostics.AddError(
			"Missing user email",
			"The provider could not read user data. Email must be specified.",
		)
	}

	// First try to fetch cluster by ID if it was provided.
	var user *console.UserFragment
	if !data.Id.IsNull() {
		if c, err := d.client.GetUser(ctx, data.Email.ValueString()); err != nil {
			resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to read user by email, got error: %s", err))
		} else {
			user = c.User
		}
	}

	if user == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to read user, see warnings for more information")
		return
	}

	data.From(user, ctx, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
