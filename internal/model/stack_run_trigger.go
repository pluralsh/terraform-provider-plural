package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StackRunTrigger struct {
	ID types.String `tfsdk:"id"`
}
