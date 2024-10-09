package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SetFrom(values []*string, ctx context.Context, d diag.Diagnostics) types.Set {
	if values == nil {
		return types.SetNull(types.StringType)
	}

	setValue, diags := types.SetValueFrom(ctx, types.StringType, values)
	d.Append(diags...)
	return setValue
}
