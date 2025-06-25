package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
)

type Group struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Global      types.Bool   `tfsdk:"global"`
}

func (g *Group) From(response *gqlclient.GroupFragment) {
	g.Id = types.StringValue(response.ID)
	g.Name = types.StringValue(response.Name)
	g.Description = types.StringPointerValue(response.Description)
	g.Global = types.BoolPointerValue(response.Global)
}

func (g *Group) Attributes() gqlclient.GroupAttributes {
	return gqlclient.GroupAttributes{
		Name:        g.Name.ValueString(),
		Description: g.Description.ValueStringPointer(),
		Global:      g.Global.ValueBoolPointer(),
	}
}
