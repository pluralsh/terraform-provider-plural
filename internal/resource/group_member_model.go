package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
)

type groupMember struct {
	Id      types.String `tfsdk:"id"`
	UserId  types.String `tfsdk:"user_id"`
	GroupId types.String `tfsdk:"group_id"`
}

func (g *groupMember) From(response *gqlclient.GroupMemberFragment) {
	g.Id = types.StringValue(response.ID)
	g.UserId = types.StringValue(response.User.ID)
	g.GroupId = types.StringValue(response.Group.ID)
}
