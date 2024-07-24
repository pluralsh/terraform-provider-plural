package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
)

type group struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (g *group) From(response *gqlclient.GroupFragment) {
	g.Id = types.StringValue(response.ID)
	g.Name = types.StringValue(response.Name)
	g.Name = types.StringPointerValue(response.Description)
}

func (g *group) Attributes() gqlclient.GroupAttributes {
	return gqlclient.GroupAttributes{
		Name:        g.Name.ValueString(),
		Description: g.Description.ValueStringPointer(),
	}
}
