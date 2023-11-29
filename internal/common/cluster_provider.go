package common

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	console "github.com/pluralsh/console-client-go"
)

func ClusterProviderIdFrom(provider *console.ClusterProviderFragment) basetypes.StringValue {
	if provider != nil {
		return types.StringValue(provider.ID)
	}

	return types.StringNull()
}
