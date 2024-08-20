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
	console "github.com/pluralsh/console/go/client"
)

func NewInfrastructureStackDataSource() datasource.DataSource {
	return &InfrastructureStackDataSource{}
}

// InfrastructureStackDataSource defines the stack data source implementation.
type InfrastructureStackDataSource struct {
	client *client.Client
}

func (r *InfrastructureStackDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infrastructure_stack"
}

func (r *InfrastructureStackDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Stack data source.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Description:         "Internal identifier of this stack.",
				MarkdownDescription: "Internal identifier of this stack.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("name")),
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Description:         "Human-readable name of this stack.",
				MarkdownDescription: "Human-readable name of this stack.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("id")),
				},
			},
			"type": schema.StringAttribute{
				Description:         "A type for the stack, specifies the tool to use to apply it. ",
				MarkdownDescription: "A type for the stack, specifies the tool to use to apply it. ",
				Computed:            true,
			},
			"approval": schema.BoolAttribute{
				Description:         "Determines whether to require approval.",
				MarkdownDescription: "Determines whether to require approval.",
				Computed:            true,
			},
			"project_id": schema.StringAttribute{
				Description:         "ID of the project that this stack belongs to.",
				MarkdownDescription: "ID of the project that this stack belongs to.",
				Computed:            true,
			},
			"cluster_id": schema.StringAttribute{
				Description:         "The cluster on which the stack is be applied.",
				MarkdownDescription: "The cluster on which the stack is be applied.",
				Computed:            true,
			},
		},
	}
}

func (r *InfrastructureStackDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Infrastructure Stack Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = data.Client
}

func (r *InfrastructureStackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := new(model.InfrastructureStack)
	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var fragment *console.InfrastructureStackFragment
	if !data.Id.IsNull() {
		if c, err := r.client.GetInfrastructureStack(ctx, data.Id.ValueString()); err != nil {
			resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to read stack by ID, got error: %s", err))
		} else {
			fragment = c.InfrastructureStack
		}
	}

	// TODO: Add once API will allow it.
	//if fragment == nil && !data.Name.IsNull() {
	//	if c, err := r.client.GetInfrastructureStack(ctx, data.Name.ValueString()); err != nil {
	//		resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to read stack by name, got error: %s", err))
	//	} else {
	//		fragment = c.InfrastructureStack
	//	}
	//}

	if fragment == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to read stack, see warnings for more information")
		return
	}

	data.From(fragment)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
