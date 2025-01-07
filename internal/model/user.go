package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
)

type User struct {
	Id    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

func (u *User) From(response *console.UserFragment) {
	u.Id = types.StringValue(response.ID)
	u.Name = types.StringValue(response.Name)
	u.Email = types.StringValue(response.Email)
}

func (u *User) Attributes() console.UserAttributes {
	return console.UserAttributes{
		Name:  u.Name.ValueStringPointer(),
		Email: u.Email.ValueStringPointer(),
	}
}
