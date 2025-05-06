package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
)

type ObservabilityWebhook struct {
	Id     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Type   types.String `tfsdk:"type"`
	URL    types.String `tfsdk:"url"`
	Secret types.String `tfsdk:"hmac"`
}

func (ow *ObservabilityWebhook) From(response *console.ObservabilityWebhookFragment) {
	ow.Id = types.StringValue(response.ID)
	ow.Name = types.StringValue(response.Name)
	ow.Type = types.StringValue(string(response.Type))
	ow.URL = types.StringValue(response.URL)
}

func (ow *ObservabilityWebhook) Attributes() console.ObservabilityWebhookAttributes {
	return console.ObservabilityWebhookAttributes{
		Type:   console.ObservabilityWebhookType(ow.Type.ValueString()),
		Name:   ow.Name.ValueString(),
		Secret: ow.Secret.ValueStringPointer(),
	}
}
