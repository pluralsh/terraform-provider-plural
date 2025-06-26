package model

import (
	"context"
	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	gqlclient "github.com/pluralsh/console/go/client"
)

type CloudConnection struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
	ReadBindings  types.Set    `tfsdk:"read_bindings"`
}

func (c *CloudConnection) Attributes(ctx context.Context, d *diag.Diagnostics) (*gqlclient.CloudConnectionAttributes, error) {
	return &gqlclient.CloudConnectionAttributes{
		Name:         c.Name.ValueString(),
		Provider:     gqlclient.Provider(c.CloudProvider.ValueString()),
		ReadBindings: common.SetToPolicyBindingAttributes(c.ReadBindings, ctx, d),
	}, nil
}

func (c *CloudConnection) From(response *gqlclient.CloudConnectionFragment, ctx context.Context, d *diag.Diagnostics) {
	c.Id = types.StringValue(response.ID)
	c.Name = types.StringValue(response.Name)
	c.CloudProvider = types.StringValue(string(response.Provider))
	c.ReadBindings = cloudConnectionReadBindingsFrom(response.ReadBindings, ctx, d)
}

func cloudConnectionReadBindingsFrom(bindings []*gqlclient.PolicyBindingFragment, ctx context.Context, d *diag.Diagnostics) types.Set {
	if len(bindings) == 0 {
		return types.SetNull(basetypes.ObjectType{AttrTypes: common.PolicyBindingAttrTypes})
	}

	values := make([]attr.Value, len(bindings))
	for i, binding := range bindings {
		value := common.PolicyBinding{
			ID: types.StringPointerValue(binding.ID),
		}

		if binding.User != nil {
			value.UserID = types.StringValue(binding.User.ID)
		}

		if binding.Group != nil {
			value.GroupID = types.StringValue(binding.Group.ID)
		}

		objValue, diags := types.ObjectValueFrom(ctx, common.PolicyBindingAttrTypes, value)
		values[i] = objValue
		d.Append(diags...)
	}

	setValue, diags := types.SetValue(basetypes.ObjectType{AttrTypes: common.PolicyBindingAttrTypes}, values)
	d.Append(diags...)
	return setValue
}
