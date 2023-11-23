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
	Id            types.String         `tfsdk:"id"`
	InseredAt     types.String         `tfsdk:"inserted_at"`
	Name          types.String         `tfsdk:"name"`
	Handle        types.String         `tfsdk:"handle"`
	Version       types.String         `tfsdk:"version"`
	ProviderId    types.String         `tfsdk:"provider_id"`
	Cloud         types.String         `tfsdk:"cloud"`
	CloudSettings ClusterCloudSettings `tfsdk:"cloud_settings"`
	NodePools     types.List           `tfsdk:"node_pools"`
	Protect       types.Bool           `tfsdk:"protect"`
	Tags          types.Map            `tfsdk:"tags"`
}

type ClusterCloudSettings struct {
	AWS   *ClusterCloudSettingsAWS   `tfsdk:"aws"`
	Azure *ClusterCloudSettingsAzure `tfsdk:"azure"`
	GCP   *ClusterCloudSettingsGCP   `tfsdk:"gcp"`
	BYOK  *ClusterCloudSettingsBYOK  `tfsdk:"byok"`
}

type ClusterCloudSettingsAWS struct {
	Region types.String `tfsdk:"region"`
}

type ClusterCloudSettingsAzure struct {
	ResourceGroup  types.String `tfsdk:"resource_group"`
	Network        types.String `tfsdk:"network"`
	SubscriptionId types.String `tfsdk:"subscription_id"`
	Location       types.String `tfsdk:"location"`
}

type ClusterCloudSettingsGCP struct {
	Region  types.String `tfsdk:"region"`
	Network types.String `tfsdk:"network"`
	Project types.String `tfsdk:"project"`
}

type ClusterCloudSettingsBYOK struct {
	Kubeconfig Kubeconfig `tfsdk:"kubeconfig"`
}

type Kubeconfig struct {
	Host                  types.String    `tfsdk:"host"`
	Username              types.String    `tfsdk:"username"`
	Password              types.String    `tfsdk:"password"`
	Insecure              types.Bool      `tfsdk:"insecure"`
	TlsServerName         types.String    `tfsdk:"tls_server_name"`
	ClientCertificate     types.String    `tfsdk:"client_certificate"`
	ClientKey             types.String    `tfsdk:"client_key"`
	ClusterCACertificate  types.String    `tfsdk:"cluster_ca_certificate"`
	ConfigPath            types.String    `tfsdk:"config_path"`
	ConfigContext         types.String    `tfsdk:"config_context"`
	ConfigContextAuthInfo types.String    `tfsdk:"config_context_auth_info"`
	ConfigContextCluster  types.String    `tfsdk:"config_context_cluster"`
	Token                 types.String    `tfsdk:"token"`
	ProxyURL              types.String    `tfsdk:"proxy_url"`
	Exec                  *KubeconfigExec `tfsdk:"exec"`
}

type KubeconfigExec struct {
	Command    types.String `tfsdk:"command"`
	Args       types.String `tfsdk:"args"`
	Env        types.String `tfsdk:"env"`
	APIVersion types.String `tfsdk:"api_version"`
}

type NodePoolAttributes struct {
	Name          types.String          `tfsdk:"name"`
	MinSize       types.Int64           `tfsdk:"min_size"`
	MaxSize       types.Int64           `tfsdk:"max_size"`
	InstanceType  types.String          `tfsdk:"instance_type"`
	Labels        types.Map             `tfsdk:"labels"`
	Taints        types.List            `tfsdk:"taints"`
	CloudSettings NodePoolCloudSettings `tfsdk:"cloud_settings"`
}

type NodePoolCloudSettings struct {
	AWS *NodePoolCloudSettingsAWS `tfsdk:"aws"`
}

type NodePoolCloudSettingsAWS struct {
	LaunchTemplateId types.String `tfsdk:"launch_template_id"`
}

func (c *Cluster) CloudSettingsAttributes() *console.CloudSettingsAttributes {
	if IsCloud(c.Cloud.ValueString(), CloudAWS) {
		return &console.CloudSettingsAttributes{
			Aws: &console.AwsCloudAttributes{
				Region: c.CloudSettings.AWS.Region.ValueStringPointer(),
			},
		}
	}

	if IsCloud(c.Cloud.ValueString(), CloudAzure) {
		return &console.CloudSettingsAttributes{
			Azure: &console.AzureCloudAttributes{
				Location:       c.CloudSettings.Azure.Location.ValueStringPointer(),
				SubscriptionID: c.CloudSettings.Azure.SubscriptionId.ValueStringPointer(),
				ResourceGroup:  c.CloudSettings.Azure.ResourceGroup.ValueStringPointer(),
				Network:        c.CloudSettings.Azure.Network.ValueStringPointer(),
			},
		}
	}

	if IsCloud(c.Cloud.ValueString(), CloudGCP) {
		return &console.CloudSettingsAttributes{
			Gcp: &console.GcpCloudAttributes{
				Project: c.CloudSettings.GCP.Project.ValueStringPointer(),
				Network: c.CloudSettings.GCP.Network.ValueStringPointer(),
				Region:  c.CloudSettings.GCP.Region.ValueStringPointer(),
			},
		}
	}

	return nil
}

func (c *Cluster) TagsAttribute(d diag.Diagnostics) (result []*console.TagAttributes) {
	elements := make(map[string]types.String, len(c.Tags.Elements()))
	d.Append(c.Tags.ElementsAs(context.TODO(), &elements, false)...)
	for k, v := range elements {
		result = append(result, &console.TagAttributes{Name: k, Value: v.ValueString()})
	}

	return
}

func (c *Cluster) CreateAttributes(d diag.Diagnostics) console.ClusterAttributes {
	return console.ClusterAttributes{
		Name:          c.Name.ValueString(),
		Handle:        c.Handle.ValueStringPointer(),
		Version:       c.Version.ValueStringPointer(),
		ProviderID:    c.ProviderId.ValueStringPointer(),
		CloudSettings: c.CloudSettingsAttributes(),
		Tags:          c.TagsAttribute(d),
		Protect:       c.Protect.ValueBoolPointer(),
	}
}

func (c *Cluster) UpdateAttributes() console.ClusterUpdateAttributes {
	return console.ClusterUpdateAttributes{
		Handle:  c.Handle.ValueStringPointer(),
		Version: c.Version.ValueStringPointer(),
		Protect: c.Protect.ValueBoolPointer(),
	}
}

func (c *Cluster) ProviderFrom(provider *console.ClusterProviderFragment) {
	if provider != nil {
		c.ProviderId = types.StringValue(provider.ID)
	}
}

func (c *Cluster) NodePoolsFrom(nodepools []*console.NodePoolFragment, d diag.Diagnostics) {
	// TODO
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
