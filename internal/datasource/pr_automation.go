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

func NewPRAutomationDataSource() datasource.DataSource {
	return &PRAutomationDataSource{}
}

// PRAutomationDataSource defines the PR automation data source implementation.
type PRAutomationDataSource struct {
	client *client.Client
}

func (r *PRAutomationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pr_automation"
}

func (r *PRAutomationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "PR automation data source.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Description:         "Internal identifier of this PR automation.",
				MarkdownDescription: "Internal identifier of this PR automation.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("name")),
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Description:         "Name of this PR automation.",
				MarkdownDescription: "Name of this PR automation.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("id")),
				},
			},
			"identifier": schema.StringAttribute{
				Computed:            true,
				Description:         "Identifier of this PR automation.",
				MarkdownDescription: "Identifier of this PR automation.",
			},
			"title": schema.StringAttribute{
				Computed:            true,
				Description:         "Title of this PR automation.",
				MarkdownDescription: "Title of this PR automation.",
			},
			"message": schema.StringAttribute{
				Computed:            true,
				Description:         "Message of this PR automation.",
				MarkdownDescription: "Message of this PR automation.",
			},
		},
	}
}

func (r *PRAutomationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected PR Automation Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = data.Client
}

func (r *PRAutomationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := new(model.PRAutomation)
	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var fragment *console.PrAutomationFragment
	if !data.Id.IsNull() {
		if c, err := r.client.GetPrAutomation(ctx, data.Id.ValueString()); err != nil {
			resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to read PR automation by ID, got error: %s", err))
		} else {
			fragment = c.PrAutomation
		}
	}

	if fragment == nil && !data.Name.IsNull() {
		if c, err := r.client.GetPrAutomationByName(ctx, data.Name.ValueString()); err != nil {
			resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to read PR automation by handle, got error: %s", err))
		} else {
			fragment = c.PrAutomation
		}
	}

	if fragment == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to read PR automation, see warnings for more information")
		return
	}

	data.From(fragment)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
