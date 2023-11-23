package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console-client-go"
)

// Cluster describes the cluster resource and data source model.
type Cluster struct {
	Id            types.String          `tfsdk:"id"`
	InseredAt     types.String          `tfsdk:"inserted_at"`
	Name          types.String          `tfsdk:"name"`
	Handle        types.String          `tfsdk:"handle"`
	Version       types.String          `tfsdk:"version"`
	ProviderId    types.String          `tfsdk:"provider_id"`
	Cloud         types.String          `tfsdk:"cloud"`
	CloudSettings *ClusterCloudSettings `tfsdk:"cloud_settings"`
	NodePools     []*ClusterNodePool    `tfsdk:"node_pools"`
	Protect       types.Bool            `tfsdk:"protect"`
	Tags          types.Map             `tfsdk:"tags"`
	Bindings      *ClusterBindings      `tfsdk:"bindings"`
}

func (c *Cluster) CloudSettingsAttributes() *console.CloudSettingsAttributes {
	if IsCloud(c.Cloud.ValueString(), CloudAWS) {
		return &console.CloudSettingsAttributes{Aws: c.CloudSettings.AWS.Attributes()}
	}

	if IsCloud(c.Cloud.ValueString(), CloudAzure) {
		return &console.CloudSettingsAttributes{Azure: c.CloudSettings.Azure.Attributes()}
	}

	if IsCloud(c.Cloud.ValueString(), CloudGCP) {
		return &console.CloudSettingsAttributes{Gcp: c.CloudSettings.GCP.Attributes()}
	}

	return nil
}

type ClusterNodePool struct {
	Name          types.String           `tfsdk:"name"`
	MinSize       types.Int64            `tfsdk:"min_size"`
	MaxSize       types.Int64            `tfsdk:"max_size"`
	InstanceType  types.String           `tfsdk:"instance_type"`
	Labels        types.Map              `tfsdk:"labels"`
	Taints        types.List             `tfsdk:"taints"`
	CloudSettings *NodePoolCloudSettings `tfsdk:"cloud_settings"`
}

type NodePoolCloudSettings struct {
	AWS *NodePoolCloudSettingsAWS `tfsdk:"aws"`
}

type NodePoolCloudSettingsAWS struct {
	LaunchTemplateId types.String `tfsdk:"launch_template_id"`
}

func (c *Cluster) NodePoolsAttribute() (result []*console.NodePoolAttributes) {
	return nil // TODO
}

func (c *Cluster) TagsAttribute(ctx context.Context, d diag.Diagnostics) (result []*console.TagAttributes) {
	elements := make(map[string]types.String, len(c.Tags.Elements()))
	d.Append(c.Tags.ElementsAs(context.TODO(), &elements, false)...)
	for k, v := range elements {
		result = append(result, &console.TagAttributes{Name: k, Value: v.ValueString()})
	}

	return
}

func (c *Cluster) CreateAttributes(ctx context.Context, d diag.Diagnostics) console.ClusterAttributes {
	return console.ClusterAttributes{
		Name:          c.Name.ValueString(),
		Handle:        c.Handle.ValueStringPointer(),
		ProviderID:    c.ProviderId.ValueStringPointer(),
		Version:       c.Version.ValueStringPointer(),
		Protect:       c.Protect.ValueBoolPointer(),
		CloudSettings: c.CloudSettingsAttributes(),
		NodePools:     c.NodePoolsAttribute(),
		ReadBindings:  c.Bindings.ReadAttributes(),
		WriteBindings: c.Bindings.WriteAttributes(),
		Tags:          c.TagsAttribute(ctx, d),
	}
}

func (c *Cluster) UpdateAttributes() console.ClusterUpdateAttributes {
	return console.ClusterUpdateAttributes{
		Version:   c.Version.ValueStringPointer(),
		Handle:    c.Handle.ValueStringPointer(),
		Protect:   c.Protect.ValueBoolPointer(),
		NodePools: c.NodePoolsAttribute(),
	}
}

func (c *Cluster) ProviderFrom(provider *console.ClusterProviderFragment) {
	if provider != nil {
		c.ProviderId = types.StringValue(provider.ID)
	}
}

func (c *Cluster) TagsFrom(tags []*console.ClusterTags, d diag.Diagnostics) {
	elements := map[string]attr.Value{}
	for _, v := range tags {
		elements[v.Name] = types.StringValue(v.Value)
	}

	tagsValue, tagsDiagnostics := types.MapValue(types.StringType, elements)
	c.Tags = tagsValue
	d.Append(tagsDiagnostics...)
}

func (c *Cluster) From(cl *console.ClusterFragment, d diag.Diagnostics) {
	c.Id = types.StringValue(cl.ID)
	c.InseredAt = types.StringPointerValue(cl.InsertedAt)
	c.Name = types.StringValue(cl.Name)
	c.Handle = types.StringPointerValue(cl.Handle)
	c.Version = types.StringPointerValue(cl.Version)
	c.Protect = types.BoolPointerValue(cl.Protect)
	c.ProviderFrom(cl.Provider)
	c.NodePoolsFrom(cl.NodePools, d)
	c.TagsFrom(cl.Tags, d)
}

func (c *Cluster) FromCreate(cc *console.CreateCluster, d diag.Diagnostics) {
	c.Id = types.StringValue(cc.CreateCluster.ID)
	c.InseredAt = types.StringPointerValue(cc.CreateCluster.InsertedAt)
	c.Name = types.StringValue(cc.CreateCluster.Name)
	c.Handle = types.StringPointerValue(cc.CreateCluster.Handle)
	c.Version = types.StringPointerValue(cc.CreateCluster.Version)
	c.Protect = types.BoolPointerValue(cc.CreateCluster.Protect)
	c.ProviderFrom(cc.CreateCluster.Provider)
	c.NodePoolsFrom(cc.CreateCluster.NodePools, d)
	c.TagsFrom(cc.CreateCluster.Tags, d)
}

func (c *Cluster) NodePoolsFrom(nodepools []*console.NodePoolFragment, d diag.Diagnostics) {
	// TODO
}
