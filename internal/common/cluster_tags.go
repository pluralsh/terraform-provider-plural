package common

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	console "github.com/pluralsh/console-client-go"
)

func ClusterTagsFrom(tags []*console.ClusterTags, d diag.Diagnostics) basetypes.MapValue {
	resultMap := map[string]attr.Value{}
	for _, v := range tags {
		resultMap[v.Name] = types.StringValue(v.Value)
	}

	result, tagsDiagnostics := types.MapValue(types.StringType, resultMap)
	d.Append(tagsDiagnostics...)

	return result
}
