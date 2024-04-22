package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func MapFrom(values map[string]any, ctx context.Context, d diag.Diagnostics) types.Map {
	if len(values) == 0 {
		return types.MapNull(types.StringType)
	}

	mapValue, diags := types.MapValueFrom(ctx, types.StringType, values)
	d.Append(diags...)
	return mapValue
}
