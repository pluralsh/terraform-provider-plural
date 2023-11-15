package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewClusterDataSource() datasource.DataSource {
	return &clusterDataSource{}
}

type clusterDataSource struct {
	Id     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Handle types.String `tfsdk:"handle"`
	Cloud  types.String `tfsdk:"cloud"`
	Tags   types.Map    `tfsdk:"tags"`
}

func (d *clusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (d *clusterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

// TODO: Support read by handle and ID.
func (d *clusterDataSource) Read(_ context.Context, _ datasource.ReadRequest, _ *datasource.ReadResponse) {
}
