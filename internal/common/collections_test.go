package common

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCollectionsFromEmptyResponse(t *testing.T) {
	cases := []struct {
		name string
		list types.List
		set  types.Set
	}{
		{"unset", types.ListNull(types.StringType), types.SetNull(types.StringType)},
		{"empty", types.ListValueMust(types.StringType, nil), types.SetValueMust(types.StringType, nil)},
		{"populated", types.ListValueMust(types.StringType, []attr.Value{types.StringValue("prod")}), types.SetValueMust(types.StringType, []attr.Value{types.StringValue("prod")})},
		{"unknown", types.ListUnknown(types.StringType), types.SetUnknown(types.StringType)},
	}
	for _, tc := range cases {
		for _, response := range []struct {
			name   string
			values []*string
		}{{"nil", nil}, {"empty", []*string{}}} {
			t.Run(tc.name+"/"+response.name, func(t *testing.T) {
				d := diag.Diagnostics{}
				list := ListFrom(response.values, tc.list, context.Background(), &d)
				set := SetFrom(response.values, tc.set, context.Background(), &d)
				if d.HasError() {
					t.Fatal(d)
				}
				if list.IsUnknown() || list.IsNull() != tc.list.IsNull() || len(list.Elements()) != 0 {
					t.Fatalf("expected empty list preserving only prior nullness, got %v", list)
				}
				if set.IsUnknown() || set.IsNull() != tc.set.IsNull() || len(set.Elements()) != 0 {
					t.Fatalf("expected empty set preserving only prior nullness, got %v", set)
				}
			})
		}
	}
}

func TestCollectionsFromEmptyResponseTypesZeroValue(t *testing.T) {
	d := diag.Diagnostics{}
	list := ListFrom(nil, types.List{}, context.Background(), &d)
	set := SetFrom(nil, types.Set{}, context.Background(), &d)
	if d.HasError() {
		t.Fatal(d)
	}
	if !list.Equal(types.ListNull(types.StringType)) || !set.Equal(types.SetNull(types.StringType)) {
		t.Fatalf("expected typed nulls for zero value config, got %v and %v", list, set)
	}
}
