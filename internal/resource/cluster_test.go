package resource

import (
	"context"
	"reflect"
	"testing"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	console "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
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

func TestClusterResourceProjectIdIsOptionalAndComputed(t *testing.T) {
	response := &resource.SchemaResponse{}
	(&clusterResource{}).Schema(context.Background(), resource.SchemaRequest{}, response)

	assertStringAttributeFlags(t, response.Schema.Attributes, "project_id", false, true, true)
}

func TestClusterAttributesProjectId(t *testing.T) {
	tests := []struct {
		name      string
		projectId types.String
		expected  *string
	}{
		{name: "unknown project is not sent", projectId: types.StringUnknown(), expected: nil},
		{name: "null project is not sent", projectId: types.StringNull(), expected: nil},
		{name: "configured project is sent", projectId: types.StringValue("project"), expected: lo.ToPtr("project")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := cluster{ProjectId: test.projectId}
			var diags diag.Diagnostics
			attributes := c.Attributes(context.Background(), &diags)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if !reflect.DeepEqual(attributes.ProjectID, test.expected) {
				t.Fatalf("expected %v, got %v", test.expected, attributes.ProjectID)
			}
		})
	}
}

func TestClusterFromCreateSetsProjectId(t *testing.T) {
	c := cluster{ProjectId: types.StringUnknown()}
	var diags diag.Diagnostics
	c.FromCreate(&console.CreateCluster{CreateCluster: &console.CreateCluster_CreateCluster{
		ID:      "id",
		Name:    "name",
		Project: &console.TinyProjectFragment{ID: "project"},
	}}, context.Background(), &diags)

	if !c.ProjectId.Equal(types.StringValue("project")) {
		t.Fatalf("expected project ID to be set from the API response, got %v", c.ProjectId)
	}

	c = cluster{ProjectId: types.StringUnknown()}
	c.FromCreate(&console.CreateCluster{CreateCluster: &console.CreateCluster_CreateCluster{ID: "id", Name: "name"}}, context.Background(), &diags)
	if !c.ProjectId.IsNull() {
		t.Fatalf("expected project ID to be null when the API response has no project, got %v", c.ProjectId)
	}
}

func TestClusterResourceModifyPlanRemovesKubeconfig(t *testing.T) {
	ctx := context.Background()
	r := &clusterResource{}
	s := r.schema()
	plan := tfsdk.Plan{Schema: s, Raw: configWithKubeconfig(ctx, s, kubeconfigValue(ctx, s))}

	resp := &resource.ModifyPlanResponse{Plan: plan}
	r.ModifyPlan(ctx, resource.ModifyPlanRequest{Plan: plan}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	var kubeconfig types.Object
	if diags := resp.Plan.GetAttribute(ctx, path.Root("kubeconfig"), &kubeconfig); diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if !kubeconfig.IsNull() {
		t.Fatalf("expected kubeconfig to be removed from the plan, got %v", kubeconfig)
	}
}

func TestClusterResourceModifyPlanIgnoresDestroy(t *testing.T) {
	ctx := context.Background()
	r := &clusterResource{}
	s := r.schema()
	plan := tfsdk.Plan{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}

	resp := &resource.ModifyPlanResponse{Plan: plan}
	r.ModifyPlan(ctx, resource.ModifyPlanRequest{Plan: plan}, resp)
	if resp.Diagnostics.HasError() || !resp.Plan.Raw.IsNull() {
		t.Fatalf("expected destroy plan to stay null, got %v, diagnostics: %v", resp.Plan.Raw, resp.Diagnostics)
	}
}

func kubeconfigType(ctx context.Context, s schema.Schema) tftypes.Type {
	return s.Attributes["kubeconfig"].GetType().TerraformType(ctx)
}

// kubeconfigValue returns a kubeconfig object with only the host set.
func kubeconfigValue(ctx context.Context, s schema.Schema) tftypes.Value {
	value := nullObject(kubeconfigType(ctx, s).(tftypes.Object))
	value["host"] = tftypes.NewValue(tftypes.String, "host")
	return tftypes.NewValue(kubeconfigType(ctx, s), value)
}

func configWithKubeconfig(ctx context.Context, s schema.Schema, kubeconfig tftypes.Value) tftypes.Value {
	objectType := s.Type().TerraformType(ctx).(tftypes.Object)
	attributes := nullObject(objectType)
	attributes["kubeconfig"] = kubeconfig
	return tftypes.NewValue(objectType, attributes)
}

func nullObject(objectType tftypes.Object) map[string]tftypes.Value {
	attributes := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attributeType := range objectType.AttributeTypes {
		attributes[name] = tftypes.NewValue(attributeType, nil)
	}
	return attributes
}

func TestClusterResourceUpgradeStateRemovesKubeconfig(t *testing.T) {
	ctx := context.Background()
	r := &clusterResource{}

	upgrader, ok := r.UpgradeState(ctx)[0]
	if !ok {
		t.Fatalf("expected state upgrader from version 0")
	}

	priorState := tfsdk.State{Schema: *upgrader.PriorSchema}
	diags := priorState.Set(ctx, struct {
		Id          types.String       `tfsdk:"id"`
		InsertedAt  types.String       `tfsdk:"inserted_at"`
		Name        types.String       `tfsdk:"name"`
		Handle      types.String       `tfsdk:"handle"`
		ProjectId   types.String       `tfsdk:"project_id"`
		Detach      types.Bool         `tfsdk:"detach"`
		Protect     types.Bool         `tfsdk:"protect"`
		Tags        types.Map          `tfsdk:"tags"`
		Metadata    types.String       `tfsdk:"metadata"`
		Bindings    *common.Bindings   `tfsdk:"bindings"`
		HelmRepoUrl types.String       `tfsdk:"helm_repo_url"`
		HelmValues  types.String       `tfsdk:"helm_values"`
		Kubeconfig  *common.Kubeconfig `tfsdk:"kubeconfig"`
	}{
		Id:         types.StringValue("id"),
		Name:       types.StringValue("name"),
		ProjectId:  types.StringValue("project"),
		Tags:       types.MapNull(types.StringType),
		Kubeconfig: &common.Kubeconfig{Host: types.StringValue("host")},
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	resp := &resource.UpgradeStateResponse{State: tfsdk.State{Schema: r.schema()}}
	upgrader.StateUpgrader(ctx, resource.UpgradeStateRequest{State: &priorState}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	var upgraded cluster
	if diags := resp.State.Get(ctx, &upgraded); diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if upgraded.Kubeconfig != nil {
		t.Fatalf("expected kubeconfig to be removed from the state, got %+v", upgraded.Kubeconfig)
	}

	if !upgraded.Id.Equal(types.StringValue("id")) || !upgraded.ProjectId.Equal(types.StringValue("project")) {
		t.Fatalf("expected other attributes to be preserved, got id=%v project_id=%v", upgraded.Id, upgraded.ProjectId)
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
		{name: "untracked host is changed", track: false, kubeconfig: kubeconfig("a"), expected: true},
		{name: "untracked host without kubeconfig is unchanged", track: false, kubeconfig: nil, expected: false},
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
