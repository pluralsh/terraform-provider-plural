package resource

import (
	"context"

	"terraform-provider-plural/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console-client-go"
)

type customStackRun struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Documentation types.String `tfsdk:"documentation"`
	StackId       types.String `tfsdk:"stack_id"`
	Commands      types.Set    `tfsdk:"commands"`
	Configuration types.Set    `tfsdk:"configuration"`
}

func (csr *customStackRun) Attributes(ctx context.Context, d diag.Diagnostics, client *client.Client) (*gqlclient.CustomStackRunAttributes, error) {
	attr := &gqlclient.CustomStackRunAttributes{
		Name:          csr.Name.ValueString(),
		Documentation: csr.Documentation.ValueStringPointer(),
		StackID:       csr.StackId.ValueStringPointer(),
		Commands:      nil, // TODO
		Configuration: nil, // TODO
	}

	return attr, nil
}

func (csr *customStackRun) From(customStackRun *gqlclient.CustomStackRunFragment, ctx context.Context, d diag.Diagnostics) {
	csr.Id = types.StringValue(customStackRun.ID)
	csr.Name = types.StringValue(customStackRun.Name)
	csr.Documentation = types.StringPointerValue(customStackRun.Documentation)
	csr.StackId = types.StringPointerValue(customStackRun.Stack.ID)
	// TODO Commands
	// TODO Configuration
}
