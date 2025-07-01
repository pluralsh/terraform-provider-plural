package model

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

func (c *CloudConnection) FromUpsert(cc *console.UpsertCloudConnection, ctx context.Context, d *diag.Diagnostics) {
	c.Id = types.StringValue(cc.UpsertCloudConnection.ID)
	c.Name = types.StringValue(cc.UpsertCloudConnection.Name)
	c.CloudProvider = types.StringValue(string(cc.UpsertCloudConnection.Provider))
	c.ReadBindings = cloudConnectionReadBindingsFrom(cc.UpsertCloudConnection.ReadBindings, ctx, d)
}

func (c *CloudConnection) From(cc *console.CloudConnectionFragment, ctx context.Context, d *diag.Diagnostics) {
	c.Id = types.StringValue(cc.ID)
	c.Name = types.StringValue(cc.Name)
	c.CloudProvider = types.StringValue(string(cc.Provider))
	c.ReadBindings = cloudConnectionReadBindingsFrom(cc.ReadBindings, ctx, d)
}

func cloudConnectionReadBindingsFrom(bindings []*console.PolicyBindingFragment, ctx context.Context, d *diag.Diagnostics) types.Set {
	if len(bindings) == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: common.PolicyBindingAttrTypes})
	}

	values := make([]attr.Value, 0, len(bindings))
	for _, binding := range bindings {
		var userID, groupID types.String

		if binding.User != nil && binding.Group == nil {
			userID = types.StringValue(binding.User.ID)
			groupID = types.StringNull()
		} else if binding.Group != nil && binding.User == nil {
			groupID = types.StringValue(binding.Group.ID)
			userID = types.StringNull()
		} else {
			// Skip invalid binding (either both or neither set)
			continue
		}

		id := types.StringNull()
		if !userID.IsNull() {
			id = types.StringValue(fmt.Sprintf("user:%s", userID.ValueString()))
		} else if !groupID.IsNull() {
			id = types.StringValue(fmt.Sprintf("group:%s", groupID.ValueString()))
		}

		objValue, diags := types.ObjectValueFrom(ctx, common.PolicyBindingAttrTypes, common.PolicyBinding{
			ID:      id,
			UserID:  userID,
			GroupID: groupID,
		})
		values = append(values, objValue)
		d.Append(diags...)
	}

	setValue, diags := types.SetValue(
		types.ObjectType{AttrTypes: common.PolicyBindingAttrTypes},
		values,
	)
	d.Append(diags...)
	return setValue
}
