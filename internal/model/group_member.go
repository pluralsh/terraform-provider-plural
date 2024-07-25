package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
)

type GroupMember struct {
	Id      types.String `tfsdk:"id"`
	UserId  types.String `tfsdk:"user_id"`
	GroupId types.String `tfsdk:"group_id"`
}

func (gm *GroupMember) From(response *gqlclient.GroupMemberFragment) {
	gm.Id = types.StringValue(response.ID)
	gm.UserId = types.StringValue(response.User.ID)
	gm.GroupId = types.StringValue(response.Group.ID)
}
