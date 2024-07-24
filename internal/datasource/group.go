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
	console "github.com/pluralsh/console/go/client"
)

func NewGroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

type groupDataSource struct {
	client *client.Client
}

func (d *groupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *groupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A representation of a group to organize authorization in your plural console.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Internal identifier of this group.",
				Validators:          []validator.String{stringvalidator.ExactlyOneOf(path.MatchRoot("id"))},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the group.",
				Computed:            true,
			},
		},
	}
}

func (d *groupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data group
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Missing group name",
			"The provider could not read group data. Name must be specified.",
		)
	}

	var group *console.GroupFragment
	if !data.Name.IsNull() {
		if c, err := d.client.GetGroup(ctx, data.Name.ValueString()); err != nil {
			resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to read group by name, got error: %s", err))
		} else {
			group = c.Group
		}
	}

	if group == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to read group, see warnings for more information")
		return
	}

	data.From(group, ctx, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
