package model

import (
	"context"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	console "github.com/pluralsh/console/go/client"
)

type CloudConnection struct {
	Id            types.String                  `tfsdk:"id"`
	Name          types.String                  `tfsdk:"name"`
	CloudProvider types.String                  `tfsdk:"cloud_provider"`
	Configuration *CloudConnectionConfiguration `tfsdk:"configuration"`
	ReadBindings  types.Set                     `tfsdk:"read_bindings"`
}

type CloudConnectionConfiguration struct {
	AWS   *AwsCloudConnectionAttributes   `tfsdk:"aws"`
	GCP   *GcpCloudConnectionAttributes   `tfsdk:"gcp"`
	Azure *AzureCloudConnectionAttributes `tfsdk:"azure"`
}

func (c *CloudConnectionConfiguration) Attributes() *console.CloudConnectionConfigurationAttributes {
	if c == nil {
		return nil
	}

	if c.AWS != nil {
		return &console.CloudConnectionConfigurationAttributes{AWS: c.AWS.Attributes()}
	}

	if c.Azure != nil {
		return &console.CloudConnectionConfigurationAttributes{Azure: c.Azure.Attributes()}
	}

	if c.GCP != nil {
		return &console.CloudConnectionConfigurationAttributes{GCP: c.GCP.Attributes()}
	}

	return nil
}

type AwsCloudConnectionAttributes struct {
	AccessKeyID     types.String `tfsdk:"access_key_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
	Region          types.String `tfsdk:"region"`
}

func (c *AwsCloudConnectionAttributes) Attributes() *console.AWSCloudConnectionAttributes {
	return &console.AWSCloudConnectionAttributes{
		AccessKeyID:     c.AccessKeyID.ValueString(),
		SecretAccessKey: c.SecretAccessKey.ValueString(),
		Region:          c.Region.ValueString(),
	}
}

type GcpCloudConnectionAttributes struct {
	ServiceAccountKey types.String `tfsdk:"service_account_key"`
	ProjectID         types.String `tfsdk:"project_id"`
}

func (c *GcpCloudConnectionAttributes) Attributes() *console.GCPCloudConnectionAttributes {
	return &console.GCPCloudConnectionAttributes{
		ServiceAccountKey: c.ServiceAccountKey.ValueString(),
		ProjectID:         c.ProjectID.ValueString(),
	}
}

type AzureCloudConnectionAttributes struct {
	SubscriptionID types.String `tfsdk:"subscription_id"`
	TenantID       types.String `tfsdk:"tenant_id"`
	ClientID       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
}

func (c *AzureCloudConnectionAttributes) Attributes() *console.AzureCloudConnectionAttributes {
	return &console.AzureCloudConnectionAttributes{
		SubscriptionID: c.SubscriptionID.ValueString(),
		TenantID:       c.TenantID.ValueString(),
		ClientID:       c.ClientID.ValueString(),
		ClientSecret:   c.ClientSecret.ValueString(),
	}
}

func (c *CloudConnection) Attributes(ctx context.Context, d *diag.Diagnostics) console.CloudConnectionAttributes {
	return console.CloudConnectionAttributes{
		Name:          c.Name.ValueString(),
		Provider:      console.Provider(c.CloudProvider.ValueString()),
		Configuration: *c.Configuration.Attributes(),
		ReadBindings:  common.SetToPolicyBindingAttributes(c.ReadBindings, ctx, d),
	}
}

func (c *CloudConnection) From(cc *console.CloudConnectionFragment, ctx context.Context, d *diag.Diagnostics) {
	c.Id = types.StringValue(cc.ID)
	c.Name = types.StringValue(cc.Name)
	c.CloudProvider = types.StringValue(string(cc.Provider))
}
