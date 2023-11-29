package common

import "github.com/hashicorp/terraform-plugin-framework/types"

func ToAttributesMap(m map[string]types.String) map[string]interface{} {
	result := map[string]interface{}{}
	for key, val := range m {
		result[key] = val.ValueString()
	}

	return result
}
