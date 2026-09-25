package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ListFrom(values []*string, config types.List, ctx context.Context, d *diag.Diagnostics) types.List {
	if len(values) == 0 {
		// Preserve null only when already unset; populated state must reflect remote removals.
		// Typed null, as config can be an untyped zero value, e.g. on import.
		if config.IsNull() {
			return types.ListNull(types.StringType)
		}
		return types.ListValueMust(types.StringType, nil)
	}

	listValue, diags := types.ListValueFrom(ctx, types.StringType, values)
	d.Append(diags...)
	return listValue
}
