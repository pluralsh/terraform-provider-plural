package common

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AttributesJson(m map[string]types.String, d diag.Diagnostics) *string {
	stringMap := map[string]string{}
	for key, val := range m {
		stringMap[key] = val.ValueString()
	}

	jsonMap, err := json.Marshal(stringMap)
	if err != nil {
		if d != nil {
			d.AddError("Provider Error", fmt.Sprintf("Cannot marshall labels, got error: %s", err))
		}
		return nil
	}

	result := bytes.NewBuffer(jsonMap).String()
	return &result
}
