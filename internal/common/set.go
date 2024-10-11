package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SetFrom(values []*string, config types.Set, ctx context.Context, d diag.Diagnostics) types.Set {
	if len(values) == 0 {
		// Rewriting config to state to avoid inconsistent result errors.
		// This could happen, for example, when sending "nil" to API and "[]" is returned as a result.
		return config
	}

	setValue, diags := types.SetValueFrom(ctx, types.StringType, values)
	d.Append(diags...)
	return setValue
}
