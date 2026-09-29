package model

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestInfrastructureStackJobSpecContainersAttributes(t *testing.T) {
	containerType := basetypes.ObjectType{AttrTypes: InfrastructureStackContainerSpecAttrTypes}

	t.Run("unset containers are omitted", func(t *testing.T) {
		spec := &InfrastructureStackJobSpec{Containers: types.SetNull(containerType)}
		diagnostics := diag.Diagnostics{}

		if containers := spec.ContainersAttributes(context.Background(), &diagnostics); containers != nil {
			t.Fatalf("expected nil containers, got %v", *containers)
		}
	})

	t.Run("an empty list is kept to clear containers", func(t *testing.T) {
		spec := &InfrastructureStackJobSpec{Containers: types.SetValueMust(containerType, []attr.Value{})}
		diagnostics := diag.Diagnostics{}

		containers := spec.ContainersAttributes(context.Background(), &diagnostics)
		if diagnostics.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diagnostics)
		}
		if containers == nil || len(*containers) != 0 {
			t.Fatalf("expected an empty container list, got %v", containers)
		}
	})
}
