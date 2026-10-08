package resource

import (
	"context"
	"testing"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestClusterResourceSchemaIsValid(t *testing.T) {
	response := &resource.SchemaResponse{}
	(&clusterResource{}).Schema(context.Background(), resource.SchemaRequest{}, response)
	if response.Diagnostics.HasError() {
		t.Fatalf("unexpected schema diagnostics: %v", response.Diagnostics)
	}

	if diags := response.Schema.ValidateImplementation(context.Background()); diags.HasError() {
		t.Fatalf("invalid schema implementation: %v", diags)
	}
}

func TestClusterResourceKubeconfigIsWriteOnly(t *testing.T) {
	response := &resource.SchemaResponse{}
	(&clusterResource{}).Schema(context.Background(), resource.SchemaRequest{}, response)

	kubeconfig, ok := response.Schema.Attributes["kubeconfig"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("expected kubeconfig to be a single nested attribute")
	}

	if !kubeconfig.IsWriteOnly() {
		t.Fatalf("expected kubeconfig to be write-only")
	}

	for name, attribute := range kubeconfig.Attributes {
		if !attribute.IsWriteOnly() {
			t.Errorf("expected kubeconfig.%s to be write-only", name)
		}
	}
}

type fakePrivateState map[string][]byte

func (f fakePrivateState) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return f[key], nil
}

func (f fakePrivateState) SetKey(_ context.Context, key string, value []byte) diag.Diagnostics {
	f[key] = value
	return nil
}

func TestKubeconfigHostChanged(t *testing.T) {
	ctx := context.Background()
	kubeconfig := func(host string) *common.Kubeconfig {
		return &common.Kubeconfig{Host: types.StringValue(host)}
	}

	tests := []struct {
		name       string
		stored     *common.Kubeconfig
		track      bool
		kubeconfig *common.Kubeconfig
		expected   bool
	}{
		{name: "untracked host is unchanged", track: false, kubeconfig: kubeconfig("a"), expected: false},
		{name: "missing kubeconfig is unchanged", track: true, stored: kubeconfig("a"), kubeconfig: nil, expected: false},
		{name: "same host is unchanged", track: true, stored: kubeconfig("a"), kubeconfig: kubeconfig("a"), expected: false},
		{name: "different host is changed", track: true, stored: kubeconfig("a"), kubeconfig: kubeconfig("b"), expected: true},
		{name: "added kubeconfig is changed", track: true, stored: nil, kubeconfig: kubeconfig(""), expected: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			private := fakePrivateState{}
			if test.track {
				if diags := setKubeconfigHost(ctx, private, test.stored); diags.HasError() {
					t.Fatalf("unexpected diagnostics: %v", diags)
				}
			}

			changed, diags := kubeconfigHostChanged(ctx, private, test.kubeconfig)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if changed != test.expected {
				t.Fatalf("expected %v, got %v", test.expected, changed)
			}
		})
	}
}
