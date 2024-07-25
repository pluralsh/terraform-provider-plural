package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
)

type group struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (g *group) From(gf *console.GroupFragment) {
	g.Id = types.StringValue(gf.ID)
	g.Name = types.StringValue(gf.Name)
	g.Description = types.StringPointerValue(gf.Description)
}
