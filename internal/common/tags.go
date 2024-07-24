package common

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	console "github.com/pluralsh/console/go/client"
)

func TagsFrom(tags []*console.ClusterTags, config types.Map, d diag.Diagnostics) basetypes.MapValue {
	if len(tags) == 0 {
		// Rewriting config to state to avoid inconsistent result errors.
		// This could happen, for example, when sending "nil" to API and "{}" is returned as a result.
		return config
	}

	resultMap := make(map[string]attr.Value, len(tags))
	for _, tag := range tags {
		resultMap[tag.Name] = types.StringValue(tag.Value)
	}

	result, tagsDiagnostics := types.MapValue(types.StringType, resultMap)
	d.Append(tagsDiagnostics...)

	return result
}
