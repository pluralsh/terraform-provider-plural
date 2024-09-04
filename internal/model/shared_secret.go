package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"

	"terraform-provider-plural/internal/common"
)

type SharedSecret struct {
	Name                 types.String `tfsdk:"name"`
	Secret               types.String `tfsdk:"secret"`
	NotificationBindings types.Set    `tfsdk:"notification_bindings"`
}

func (in *SharedSecret) Attributes(ctx context.Context, d diag.Diagnostics) console.SharedSecretAttributes {
	return console.SharedSecretAttributes{
		Name:                 in.Name.ValueString(),
		Secret:               in.Secret.ValueString(),
		NotificationBindings: common.SetToPolicyBindingAttributes(in.NotificationBindings, ctx, d),
	}
}
