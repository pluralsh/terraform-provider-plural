package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
)

type group struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (c *group) From(cl *console.GroupFragment, ctx context.Context, d diag.Diagnostics) {
	c.Id = types.StringValue(cl.ID)
	c.Name = types.StringValue(cl.Name)
}
