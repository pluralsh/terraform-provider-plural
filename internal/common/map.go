package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func MapFrom(values map[string]any, ctx context.Context, d *diag.Diagnostics) types.Map {
	if values == nil {
		return types.MapNull(types.StringType)
	}

	mapValue, diags := types.MapValueFrom(ctx, types.StringType, values)
	d.Append(diags...)
	return mapValue
}

func MapFromWithConfig(values map[string]any, config types.Map, ctx context.Context, d *diag.Diagnostics) types.Map {
	if len(values) == 0 {
		// Rewriting config to state to avoid inconsistent result errors.
		// This could happen, for example, when sending "nil" to API and "{}" is returned as a result.
		return config
	}

	mapValue, diags := types.MapValueFrom(ctx, types.StringType, values)
	d.Append(diags...)
	return mapValue
}
