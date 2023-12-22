package resource

import (
	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console-client-go"
)

type rbac struct {
	ClusterId types.String            `tfsdk:"cluster_id"`
	ServiceId types.String            `tfsdk:"service_id"`
	Bindings  *common.ClusterBindings `tfsdk:"rbac"`
}

func (g *rbac) Attributes() gqlclient.RbacAttributes {
	return gqlclient.RbacAttributes{
		ReadBindings:  g.Bindings.ReadAttributes(),
		WriteBindings: g.Bindings.WriteAttributes(),
	}
}
