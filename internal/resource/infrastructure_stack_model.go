package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console-client-go"
)

type InfrastructureStack struct {
	Id            types.String                      `tfsdk:"id"`
	Name          types.String                      `tfsdk:"name"`
	Type          types.String                      `tfsdk:"type"`
	ClusterId     types.String                      `tfsdk:"cluster_id"`
	Repository    *InfrastructureStackRepository    `tfsdk:"repository"`
	Approval      types.Bool                        `tfsdk:"protect"`
	Configuration *InfrastructureStackConfiguration `tfsdk:"configuration"`
	Files         []*InfrastructureStackFile        `tfsdk:"environment"`
	Environemnt   []*InfrastructureStackEnvironment `tfsdk:"files"`
	Bindings      *InfrastructureStackBindings      `tfsdk:"bindings"`
	JobSpec       *InfrastructureStackJobSpec       `tfsdk:"job_spec"`
}

type InfrastructureStackRepository struct {
	Id     types.String `tfsdk:"id"`
	Ref    types.String `tfsdk:"ref"`
	Folder types.String `tfsdk:"folder"`
}

func (isr *InfrastructureStackRepository) From(repository *gqlclient.GitRepository, ref gqlclient.GitRef) {
	if isr == nil {
		return
	}

	isr.Id = types.StringValue(repository.ID)
	isr.Ref = types.StringValue(ref.Ref)
	isr.Folder = types.StringValue(ref.Folder)
}

type InfrastructureStackConfiguration struct {
	Image   types.String `tfsdk:"image"`
	Version types.String `tfsdk:"version"`
}

func (isc *InfrastructureStackConfiguration) From(configuration gqlclient.StackConfiguration) {
	if isc == nil {
		return
	}

	isc.Image = types.StringPointerValue(configuration.Image)
	isc.Version = types.StringValue(configuration.Version)
}

type InfrastructureStackEnvironment struct {
	Name   types.String `tfsdk:"name"`
	Value  types.String `tfsdk:"value"`
	Secret types.Bool   `tfsdk:"secret"`
}

type InfrastructureStackFile struct {
	Path    types.String `tfsdk:"path"`
	Content types.String `tfsdk:"content"`
}

type InfrastructureStackBindings struct {
	Read  []*InfrastructureStackPolicyBinding `tfsdk:"read"`
	Write []*InfrastructureStackPolicyBinding `tfsdk:"write"`
}

type InfrastructureStackPolicyBinding struct {
	GroupID types.String `tfsdk:"group_id"`
	ID      types.String `tfsdk:"id"`
	UserID  types.String `tfsdk:"user_id"`
}

type InfrastructureStackJobSpec struct {
	Namespace      types.String                        `tfsdk:"namespace"`
	Raw            types.String                        `tfsdk:"raw"`
	Containers     []*InfrastructureStackContainerSpec `tfsdk:"containers"`
	Labels         types.String                        `tfsdk:"labels"`
	Annotations    types.String                        `tfsdk:"annotations"`
	ServiceAccount types.String                        `tfsdk:"serviceAccount"`
}

type InfrastructureStackContainerSpec struct {
	Image   types.String                           `tfsdk:"image"`
	Args    types.List                             `tfsdk:"args"`
	Env     types.Map                              `tfsdk:"env"`
	EnvFrom []*InfrastructureStackContainerEnvFrom `tfsdk:"envFrom"`
}

type InfrastructureStackContainerEnvFrom struct {
	Secret    types.String `tfsdk:"secret"`
	ConfigMap types.String `tfsdk:"configMap"`
}

func (is *InfrastructureStack) FromCreate(stack *gqlclient.InfrastructureStack, d diag.Diagnostics) {
	is.Id = types.StringPointerValue(stack.ID)
	is.Name = types.StringValue(stack.Name)
	is.Type = types.StringValue(string(stack.Type))
	is.ClusterId = types.StringValue(stack.Cluster.ID)
	is.Repository.From(stack.Repository, stack.Git)
	is.Approval = types.BoolPointerValue(stack.Approval)
	is.Configuration.From(stack.Configuration)
	// TODO ...
}
