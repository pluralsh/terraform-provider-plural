package resource

import (
	"context"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console-client-go"
)

type rbac struct {
	ClusterId types.String     `tfsdk:"cluster_id"`
	ServiceId types.String     `tfsdk:"service_id"`
	Bindings  *common.Bindings `tfsdk:"rbac"`
}

func (g *rbac) Attributes(ctx context.Context, d diag.Diagnostics) gqlclient.RbacAttributes {
	return gqlclient.RbacAttributes{
		ReadBindings:  g.Bindings.ReadAttributes(ctx, d),
		WriteBindings: g.Bindings.WriteAttributes(ctx, d),
	}
}
