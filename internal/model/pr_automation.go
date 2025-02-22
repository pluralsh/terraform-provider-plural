package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
)

type PRAutomation struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Identifier types.String `tfsdk:"identifier"`
	Title      types.String `tfsdk:"title"`
	Message    types.String `tfsdk:"message"`
}

func (pra *PRAutomation) From(response *gqlclient.PrAutomationFragment) {
	pra.Id = types.StringValue(response.ID)
	pra.Name = types.StringValue(response.Name)
	pra.Message = types.StringValue(response.Message)
	pra.Identifier = types.StringPointerValue(response.Identifier)
	pra.Title = types.StringValue(response.Title)
}
