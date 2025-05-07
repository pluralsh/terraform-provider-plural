package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
)

type SCMWebhook struct {
	Id    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Owner types.String `tfsdk:"owner"`
	Type  types.String `tfsdk:"type"`
	URL   types.String `tfsdk:"url"`
	Hmac  types.String `tfsdk:"hmac"`
}

func (scmw *SCMWebhook) From(response *console.ScmWebhookFragment) {
	scmw.Id = types.StringValue(response.ID)
	scmw.Name = types.StringValue(response.Name)
	scmw.Owner = types.StringValue(response.Owner)
	scmw.Type = types.StringValue(string(response.Type))
	scmw.URL = types.StringValue(response.URL)
}

func (scmw *SCMWebhook) Attributes() console.ScmWebhookAttributes {
	return console.ScmWebhookAttributes{
		Type:  console.ScmType(scmw.Type.ValueString()),
		Owner: scmw.Owner.ValueString(),
		Hmac:  scmw.Hmac.ValueString(),
	}
}
