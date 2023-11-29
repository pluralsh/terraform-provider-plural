package common

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console-client-go"
)

func ClusterTagsMap(tags []*console.ClusterTags) map[string]attr.Value {
	result := map[string]attr.Value{}
	for _, v := range tags {
		result[v.Name] = types.StringValue(v.Value)
	}

	return result
}
