package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ServiceDeployment describes the service deployment resource and data source model.
type ServiceDeployment struct {
	Id            types.String                     `tfsdk:"id"`
	Name          types.String                     `tfsdk:"name"`
	Namespace     types.String                     `tfsdk:"namespace"`
	Protect       types.Bool                       `tfsdk:"protect"`
	Configuration []ServiceDeploymentConfiguration `tfsdk:"configuration"`
	Cluster       ServiceDeploymentCluster         `tfsdk:"cluster"`
	Repository    ServiceDeploymentRepository      `tfsdk:"repository"`
}

type ServiceDeploymentConfiguration struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type ServiceDeploymentCluster struct {
	Id     types.String `tfsdk:"id"`
	Handle types.String `tfsdk:"handle"`
}

type ServiceDeploymentRepository struct {
	Id     types.String `tfsdk:"id"`
	Ref    types.String `tfsdk:"ref"`
	Folder types.String `tfsdk:"folder"`
}
