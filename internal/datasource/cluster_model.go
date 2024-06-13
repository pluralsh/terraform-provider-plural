package datasource

import (
	"context"
	"encoding/json"
	"fmt"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console-client-go"
)

type cluster struct {
	Id             types.String `tfsdk:"id"`
	InsertedAt     types.String `tfsdk:"inserted_at"`
	Name           types.String `tfsdk:"name"`
	Handle         types.String `tfsdk:"handle"`
	ProjectId      types.String `tfsdk:"project_id"`
	DesiredVersion types.String `tfsdk:"desired_version"`
	ProviderId     types.String `tfsdk:"provider_id"`
	Cloud          types.String `tfsdk:"cloud"`
	Protect        types.Bool   `tfsdk:"protect"`
	Tags           types.Map    `tfsdk:"tags"`
	Metadata       types.String `tfsdk:"metadata"`
	NodePools      types.Map    `tfsdk:"node_pools"`
}

func (c *cluster) From(cl *console.ClusterFragment, ctx context.Context, d diag.Diagnostics) {
	metadata, err := json.Marshal(cl.Metadata)
	if err != nil {
		d.AddError("Provider Error", fmt.Sprintf("Cannot marshall metadata, got error: %s", err))
		return
	}

	c.Id = types.StringValue(cl.ID)
	c.InsertedAt = types.StringPointerValue(cl.InsertedAt)
	c.Name = types.StringValue(cl.Name)
	c.Handle = types.StringPointerValue(cl.Handle)
	c.ProjectId = common.ProjectFrom(cl.Project)
	c.DesiredVersion = types.StringPointerValue(cl.Version)
	c.Protect = types.BoolPointerValue(cl.Protect)
	c.Tags = common.TagsFrom(cl.Tags, c.Tags, d)
	c.Metadata = types.StringValue(string(metadata))
	c.ProviderId = common.ClusterProviderIdFrom(cl.Provider)
	c.NodePools = common.ClusterNodePoolsFrom(cl.NodePools, c.NodePools, ctx, d)
}
