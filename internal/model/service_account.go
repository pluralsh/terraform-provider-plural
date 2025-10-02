package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
)

type ServiceAccount struct {
	Id    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

func (sa *ServiceAccount) From(response *console.UserFragment) {
	sa.Id = types.StringValue(response.ID)
	sa.Name = types.StringValue(response.Name)
	sa.Email = types.StringValue(response.Email)
}

func (sa *ServiceAccount) Attributes() console.ServiceAccountAttributes {
	return console.ServiceAccountAttributes{
		Name:  sa.Name.ValueStringPointer(),
		Email: sa.Email.ValueStringPointer(),
	}
}
