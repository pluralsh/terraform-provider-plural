package resource

import (
	"context"
	"encoding/json"
	"fmt"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
)

type cluster struct {
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
}

func (c *cluster) TagsAttribute(ctx context.Context, d *diag.Diagnostics) []*console.TagAttributes {
	if c.Tags.IsNull() {
		return nil
	}

	result := make([]*console.TagAttributes, 0)
	elements := make(map[string]types.String, len(c.Tags.Elements()))
	d.Append(c.Tags.ElementsAs(ctx, &elements, false)...)

	for k, v := range elements {
		result = append(result, &console.TagAttributes{Name: k, Value: v.ValueString()})
	}

	return result
}

func (c *cluster) Attributes(ctx context.Context, d *diag.Diagnostics) console.ClusterAttributes {
	return console.ClusterAttributes{
		Name:          c.Name.ValueString(),
		Handle:        c.Handle.ValueStringPointer(),
		ProjectID:     c.ProjectId.ValueStringPointer(),
		Protect:       c.Protect.ValueBoolPointer(),
		ReadBindings:  c.Bindings.ReadAttributes(ctx, d),
		WriteBindings: c.Bindings.WriteAttributes(ctx, d),
		Tags:          c.TagsAttribute(ctx, d),
		Metadata:      c.Metadata.ValueStringPointer(),
	}
}

func (c *cluster) UpdateAttributes(ctx context.Context, d *diag.Diagnostics) console.ClusterUpdateAttributes {
	return console.ClusterUpdateAttributes{
		Name:     c.Name.ValueStringPointer(),
		Handle:   c.Handle.ValueStringPointer(),
		Protect:  c.Protect.ValueBoolPointer(),
		Metadata: c.Metadata.ValueStringPointer(),
		Tags:     c.TagsAttribute(ctx, d),
	}
}

func (c *cluster) From(cl *console.ClusterFragment, _ context.Context, d *diag.Diagnostics) {
	metadata, err := json.Marshal(cl.Metadata)
	if err != nil {
		d.AddError("Provider Error", fmt.Sprintf("Cannot marshall metadata, got error: %s", err))
		return
	}

	c.Id = types.StringValue(cl.ID)
	c.InsertedAt = types.StringPointerValue(cl.InsertedAt)
	c.Name = types.StringValue(cl.Name)
	c.Handle = types.StringPointerValue(cl.Handle)
	c.Protect = types.BoolPointerValue(cl.Protect)
	c.Tags = common.TagsFrom(cl.Tags, c.Tags, d)
	c.Metadata = types.StringValue(string(metadata))
}

func (c *cluster) FromCreate(cc *console.CreateCluster, _ context.Context, d *diag.Diagnostics) {
	c.Id = types.StringValue(cc.CreateCluster.ID)
	c.InsertedAt = types.StringPointerValue(cc.CreateCluster.InsertedAt)
	c.Name = types.StringValue(cc.CreateCluster.Name)
	c.Handle = types.StringPointerValue(cc.CreateCluster.Handle)
	c.Protect = types.BoolPointerValue(cc.CreateCluster.Protect)
	c.Tags = common.TagsFrom(cc.CreateCluster.Tags, c.Tags, d)
}

func (c *cluster) ClusterVersionFrom(prov *console.ClusterProviderFragment, version, currentVersion *string) types.String {
	if prov == nil {
		return types.StringValue("unknown")
	}

	if version != nil && len(*version) > 0 {
		return types.StringPointerValue(version)
	}

	if currentVersion != nil && len(*currentVersion) > 0 {
		return types.StringPointerValue(currentVersion)
	}

	return types.StringValue("unknown")
}

func (c *cluster) HasKubeconfig() bool {
	return c.Kubeconfig != nil // || (c.CloudSettings != nil && c.CloudSettings.BYOK != nil && c.CloudSettings.BYOK.Kubeconfig != nil)
}

func (c *cluster) GetKubeconfig() *common.Kubeconfig {
	if !c.HasKubeconfig() {
		return nil
	}

	return c.Kubeconfig
}
