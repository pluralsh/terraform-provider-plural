package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
)

type user struct {
	Id    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

func (c *user) From(cl *console.UserFragment) {
	c.Id = types.StringValue(cl.ID)
	c.Name = types.StringValue(cl.Name)
	c.Email = types.StringValue(cl.Email)
}
