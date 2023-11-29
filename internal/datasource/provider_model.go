package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console-client-go"
)

type provider struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Namespace types.String `tfsdk:"namespace"`
	Editable  types.Bool   `tfsdk:"editable"`
	Cloud     types.String `tfsdk:"cloud"`
}

func (p *provider) From(cp *console.ClusterProviderFragment) {
	p.Id = types.StringValue(cp.ID)
	p.Name = types.StringValue(cp.Name)
	p.Namespace = types.StringValue(cp.Namespace)
	p.Editable = types.BoolPointerValue(cp.Editable)
	p.Cloud = types.StringValue(cp.Cloud)
}
