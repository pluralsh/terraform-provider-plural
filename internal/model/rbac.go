package model

import (
	"context"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
)

type RBAC struct {
	ClusterId types.String     `tfsdk:"cluster_id"`
	ServiceId types.String     `tfsdk:"service_id"`
	Bindings  *common.Bindings `tfsdk:"bindings"`
}

func (rbac *RBAC) Attributes(ctx context.Context, d diag.Diagnostics) gqlclient.RbacAttributes {
	return gqlclient.RbacAttributes{
		ReadBindings:  rbac.Bindings.ReadAttributes(ctx, d),
		WriteBindings: rbac.Bindings.WriteAttributes(ctx, d),
	}
}
