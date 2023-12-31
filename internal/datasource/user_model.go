package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console-client-go"
)

type user struct {
	Id    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

func (c *user) From(cl *console.UserFragment, ctx context.Context, d diag.Diagnostics) {
	c.Id = types.StringValue(cl.ID)
	c.Name = types.StringValue(cl.Name)
	c.Email = types.StringValue(cl.Email)
}
