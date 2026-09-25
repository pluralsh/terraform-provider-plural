package datasource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewServiceDeploymentDataSource() datasource.DataSource {
	return &ServiceDeploymentDataSource{}
}

// ServiceDeploymentDataSource defines the service deployment data source implementation.
type ServiceDeploymentDataSource struct {
	client *client.Client
}

type serviceDeployment struct {
	Id      types.String `tfsdk:"id"`
	Cluster types.String `tfsdk:"cluster"`
	Name    types.String `tfsdk:"name"`
}

func (r *ServiceDeploymentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_deployment"
}

func (r *ServiceDeploymentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Service deployment data source. Looks up a service deployment by the handle of its cluster and its name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Internal identifier of this service deployment.",
				MarkdownDescription: "Internal identifier of this service deployment.",
				Computed:            true,
			},
			"cluster": schema.StringAttribute{
				Description:         "Handle of the cluster this service deployment is deployed to, e.g. mgmt.",
				MarkdownDescription: "Handle of the cluster this service deployment is deployed to, e.g. `mgmt`.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"name": schema.StringAttribute{
				Description:         "Name of this service deployment, e.g. console.",
				MarkdownDescription: "Name of this service deployment, e.g. `console`.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
		},
	}
}

func (r *ServiceDeploymentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Service Deployment Data Source Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = data.Client
}

func (r *ServiceDeploymentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := new(serviceDeployment)
	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetServiceDeploymentTinyByHandle(ctx, data.Cluster.ValueString(), data.Name.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get service deployment, got error: %s", err))
		return
	}
	if response == nil || response.ServiceDeployment == nil || client.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to find service deployment %s in cluster %s", data.Name.ValueString(), data.Cluster.ValueString()))
		return
	}

	data.Id = types.StringValue(response.ServiceDeployment.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
