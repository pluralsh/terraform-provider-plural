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
	Id         types.String `tfsdk:"id"`
	InsertedAt types.String `tfsdk:"inserted_at"`
	Name       types.String `tfsdk:"name"`
	Handle     types.String `tfsdk:"handle"`
	ProjectId  types.String `tfsdk:"project_id"`
	Detach     types.Bool   `tfsdk:"detach"`
	// Version    types.String `tfsdk:"version"`
	// DesiredVersion types.String            `tfsdk:"desired_version"`
	// ProviderId types.String            `tfsdk:"provider_id"`
	// Cloud    types.String            `tfsdk:"cloud"`
	Protect  types.Bool       `tfsdk:"protect"`
	Tags     types.Map        `tfsdk:"tags"`
	Metadata types.String     `tfsdk:"metadata"`
	Bindings *common.Bindings `tfsdk:"bindings"`
	// NodePools      types.Map               `tfsdk:"node_pools"`
	// CloudSettings *ClusterCloudSettings `tfsdk:"cloud_settings"`
	HelmRepoUrl types.String       `tfsdk:"helm_repo_url"`
	HelmValues  types.String       `tfsdk:"helm_values"`
	Kubeconfig  *common.Kubeconfig `tfsdk:"kubeconfig"`
}

// func (c *cluster) NodePoolsAttribute(ctx context.Context, d *diag.Diagnostics) []*console.NodePoolAttributes {
// 	result := make([]*console.NodePoolAttributes, 0, len(c.NodePools.Elements()))
// 	nodePools := make(map[string]common.ClusterNodePool, len(c.NodePools.Elements()))
// 	d.Append(c.NodePools.ElementsAs(ctx, &nodePools, false)...)

// 	for _, nodePool := range nodePools {
// 		var nodePoolCloudSettings *common.NodePoolCloudSettings
// 		d.Append(nodePool.CloudSettings.As(ctx, nodePoolCloudSettings, basetypes.ObjectAsOptions{})...)

// 		result = append(result, &console.NodePoolAttributes{
// 			Name:          nodePool.Name.ValueString(),
// 			MinSize:       nodePool.MinSize.ValueInt64(),
// 			MaxSize:       nodePool.MaxSize.ValueInt64(),
// 			InstanceType:  nodePool.InstanceType.ValueString(),
// 			Labels:        nodePool.LabelsAttribute(ctx, d),
// 			Taints:        nodePool.TaintsAttribute(ctx, d),
// 			CloudSettings: nodePoolCloudSettings.Attributes(),
// 		})
// 	}

// 	return result
// }

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
		Name:      c.Name.ValueString(),
		Handle:    c.Handle.ValueStringPointer(),
		ProjectID: c.ProjectId.ValueStringPointer(),
		// ProviderID:    c.ProviderId.ValueStringPointer(),
		// Version:       c.Version.ValueStringPointer(),
		Protect: c.Protect.ValueBoolPointer(),
		// CloudSettings: c.CloudSettings.Attributes(),
		ReadBindings:  c.Bindings.ReadAttributes(ctx, d),
		WriteBindings: c.Bindings.WriteAttributes(ctx, d),
		Tags:          c.TagsAttribute(ctx, d),
		// NodePools:     c.NodePoolsAttribute(ctx, d),
		Metadata: c.Metadata.ValueStringPointer(),
	}
}

func (c *cluster) UpdateAttributes(ctx context.Context, d *diag.Diagnostics) console.ClusterUpdateAttributes {
	return console.ClusterUpdateAttributes{
		Name: c.Name.ValueStringPointer(),
		// Version: c.Version.ValueStringPointer(),
		Handle:  c.Handle.ValueStringPointer(),
		Protect: c.Protect.ValueBoolPointer(),
		// NodePools: c.NodePoolsAttribute(ctx, d),
		Metadata: c.Metadata.ValueStringPointer(),
		Tags:     c.TagsAttribute(ctx, d),
	}
}

func (c *cluster) From(cl *console.ClusterFragment, ctx context.Context, d *diag.Diagnostics) {
	metadata, err := json.Marshal(cl.Metadata)
	if err != nil {
		d.AddError("Provider Error", fmt.Sprintf("Cannot marshall metadata, got error: %s", err))
		return
	}

	c.Id = types.StringValue(cl.ID)
	c.InsertedAt = types.StringPointerValue(cl.InsertedAt)
	c.Name = types.StringValue(cl.Name)
	c.Handle = types.StringPointerValue(cl.Handle)
	// c.DesiredVersion = c.ClusterVersionFrom(cl.Provider, cl.Version, cl.CurrentVersion)
	c.Protect = types.BoolPointerValue(cl.Protect)
	c.Tags = common.TagsFrom(cl.Tags, c.Tags, d)
	// c.ProviderId = common.ClusterProviderIdFrom(cl.Provider)
	// c.NodePools = common.ClusterNodePoolsFrom(cl.NodePools, c.NodePools, ctx, d)
	c.Metadata = types.StringValue(string(metadata))
}

func (c *cluster) FromCreate(cc *console.CreateCluster, ctx context.Context, d *diag.Diagnostics) {
	c.Id = types.StringValue(cc.CreateCluster.ID)
	c.InsertedAt = types.StringPointerValue(cc.CreateCluster.InsertedAt)
	c.Name = types.StringValue(cc.CreateCluster.Name)
	c.Handle = types.StringPointerValue(cc.CreateCluster.Handle)
	// c.DesiredVersion = c.ClusterVersionFrom(cc.CreateCluster.Provider, cc.CreateCluster.Version, cc.CreateCluster.CurrentVersion)
	c.Protect = types.BoolPointerValue(cc.CreateCluster.Protect)
	c.Tags = common.TagsFrom(cc.CreateCluster.Tags, c.Tags, d)
	// c.ProviderId = common.ClusterProviderIdFrom(cc.CreateCluster.Provider)
	// c.NodePools = common.ClusterNodePoolsFrom(cc.CreateCluster.NodePools, c.NodePools, ctx, d)
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

type ClusterCloudSettings struct {
	AWS   *ClusterCloudSettingsAWS   `tfsdk:"aws"`
	Azure *ClusterCloudSettingsAzure `tfsdk:"azure"`
	GCP   *ClusterCloudSettingsGCP   `tfsdk:"gcp"`
	BYOK  *ClusterCloudSettingsBYOK  `tfsdk:"byok"`
}

func (c *ClusterCloudSettings) Attributes() *console.CloudSettingsAttributes {
	if c == nil {
		return nil
	}

	if c.AWS != nil {
		return &console.CloudSettingsAttributes{AWS: c.AWS.Attributes()}
	}

	if c.Azure != nil {
		return &console.CloudSettingsAttributes{Azure: c.Azure.Attributes()}
	}

	if c.GCP != nil {
		return &console.CloudSettingsAttributes{GCP: c.GCP.Attributes()}
	}

	return nil
}

type ClusterCloudSettingsAWS struct {
	Region types.String `tfsdk:"region"`
}

func (c *ClusterCloudSettingsAWS) Attributes() *console.AWSCloudAttributes {
	return &console.AWSCloudAttributes{
		Region: c.Region.ValueStringPointer(),
	}
}

type ClusterCloudSettingsAzure struct {
	ResourceGroup  types.String `tfsdk:"resource_group"`
	Network        types.String `tfsdk:"network"`
	SubscriptionId types.String `tfsdk:"subscription_id"`
	Location       types.String `tfsdk:"location"`
}

func (c *ClusterCloudSettingsAzure) Attributes() *console.AzureCloudAttributes {
	return &console.AzureCloudAttributes{
		Location:       c.Location.ValueStringPointer(),
		SubscriptionID: c.SubscriptionId.ValueStringPointer(),
		ResourceGroup:  c.ResourceGroup.ValueStringPointer(),
		Network:        c.Network.ValueStringPointer(),
	}
}

type ClusterCloudSettingsGCP struct {
	Region  types.String `tfsdk:"region"`
	Network types.String `tfsdk:"network"`
	Project types.String `tfsdk:"project"`
}

func (c *ClusterCloudSettingsGCP) Attributes() *console.GCPCloudAttributes {
	return &console.GCPCloudAttributes{
		Project: c.Project.ValueStringPointer(),
		Network: c.Network.ValueStringPointer(),
		Region:  c.Region.ValueStringPointer(),
	}
}

type ClusterCloudSettingsBYOK struct {
	Kubeconfig *common.Kubeconfig `tfsdk:"kubeconfig"`
}
