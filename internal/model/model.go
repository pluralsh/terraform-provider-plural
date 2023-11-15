package model

import "github.com/hashicorp/terraform-plugin-framework/types"

// Cluster describes the cluster resource and data source model.
type Cluster struct {
	Id        types.String `tfsdk:"id"`
	InseredAt types.String `tfsdk:"inserted_at"`
	Name      types.String `tfsdk:"name"`
	Handle    types.String `tfsdk:"handle"`
	Cloud     types.String `tfsdk:"cloud"`
	Protect   types.Bool   `tfsdk:"protect"`
	Tags      types.Map    `tfsdk:"tags"`
}
