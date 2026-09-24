package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SetFrom(values []*string, config types.Set, ctx context.Context, d *diag.Diagnostics) types.Set {
	if len(values) == 0 {
		// Preserve null only when already unset; populated state must reflect remote removals.
		// Typed null, as config can be an untyped zero value, e.g. on import.
		if config.IsNull() {
			return types.SetNull(types.StringType)
		}
		return types.SetValueMust(types.StringType, nil)
	}

	setValue, diags := types.SetValueFrom(ctx, types.StringType, values)
	d.Append(diags...)
	return setValue
}
