package resource

import (
	"context"

	"terraform-provider-plural/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
	csr.Commands = customStackRunCommandsFrom(customStackRun.Commands, csr.Commands, ctx, d)
	// TODO Configuration
}

type CustomStackRunCommand struct {
	Cmd  types.String `tfsdk:"cmd"`
	Args types.Set    `tfsdk:"args"`
	Dir  types.String `tfsdk:"dir"`
}

var CustomStackRunCommandAttrTypes = map[string]attr.Type{
	"cmd":  types.StringType,
	"args": types.SetType{ElemType: types.StringType},
	"dir":  types.StringType,
}

func customStackRunCommandsFrom(commands []*gqlclient.StackCommandFragment, config types.Set, ctx context.Context, d diag.Diagnostics) types.Set {
	if len(commands) == 0 {
		// Rewriting config to state to avoid inconsistent result errors.
		// This could happen, for example, when sending "nil" to API and "[]" is returned as a result.
		return config
	}

	values := make([]attr.Value, len(commands))
	for i, command := range commands {
		objValue, diags := types.ObjectValueFrom(ctx, CustomStackRunCommandAttrTypes, CustomStackRunCommand{
			Cmd:  types.StringValue(command.Cmd),
			Args: customStackRunCommandArgsFrom(command.Args, ctx, d),
			Dir:  types.StringPointerValue(command.Dir),
		})
		values[i] = objValue
		d.Append(diags...)
	}

	setValue, diags := types.SetValue(basetypes.ObjectType{AttrTypes: CustomStackRunCommandAttrTypes}, values)
	d.Append(diags...)
	return setValue
}

func customStackRunCommandArgsFrom(values []*string, ctx context.Context, d diag.Diagnostics) types.Set {
	if values == nil {
		return types.SetNull(types.StringType)
	}

	setValue, diags := types.SetValueFrom(ctx, types.StringType, values)
	d.Append(diags...)
	return setValue
}
