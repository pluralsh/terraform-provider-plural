package datasource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"
	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCloudConnectionDataSource() datasource.DataSource {
	return &cloudConnectionDataSource{}
}

type cloudConnectionDataSource struct {
	client *client.Client
}

func (d *cloudConnectionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_connection"
}

func (d *cloudConnectionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Internal identifier of this cloud connection",
				MarkdownDescription: "Internal identifier of this cloud connection",
				Optional:            true,
				Computed:            true,
				Validators:          []validator.String{stringvalidator.ExactlyOneOf(path.MatchRoot("name"))},
			},
			"name": schema.StringAttribute{
				Description:         "Human-readable name of this cloud connection.",
				MarkdownDescription: "Human-readable name of this cloud connection.",
				Optional:            true,
				Computed:            true,
				Validators:          []validator.String{stringvalidator.ExactlyOneOf(path.MatchRoot("id"))},
			},
			"provider": schema.StringAttribute{
				Description:         "The cloud provider of this cloud connection.",
				MarkdownDescription: "The cloud provider of this cloud connection.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf("AWS", "GCP", "AZURE")},
			},
			"configuration": schema.SingleNestedAttribute{
				Description:         "Cloud provider configuration",
				MarkdownDescription: "Cloud provider configuration",
				Required:            true,
			},
			"read_bindings": schema.SetAttribute{
				Description:         "The read bindings for this cloud connection.",
				MarkdownDescription: "The read bindings for this cloud connection.",
				Optional:            true,
				ElementType:         types.ObjectType{AttrTypes: common.PolicyBindingAttrTypes},
			},
		},
	}
}

func (d *cloudConnectionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Cloud Connection Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = data.Client
}

func (d *cloudConnectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := new(model.CloudConnection)
	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Cloud Connection ID and Name",
			"The provider could not read cloud connection data. ID or name needs to be specified.",
		)
		return
	}

	response, err := d.client.GetCloudConnection(ctx, data.Id.ValueStringPointer(), data.Name.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cloud connection, got error: %s", err))
		return
	}

	data.From(response.CloudConnection, ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
