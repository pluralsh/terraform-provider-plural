package resource

import (
	"context"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	gqlclient "github.com/pluralsh/console-client-go"
)

type infrastructureStack struct {
	Id            types.String                      `tfsdk:"id"`
	Name          types.String                      `tfsdk:"name"`
	Type          types.String                      `tfsdk:"type"`
	Approval      types.Bool                        `tfsdk:"protect"`
	ClusterId     types.String                      `tfsdk:"cluster_id"`
	Repository    *InfrastructureStackRepository    `tfsdk:"repository"`
	Configuration *InfrastructureStackConfiguration `tfsdk:"configuration"`
	Files         types.Map                         `tfsdk:"files"`
	Environment   types.Set                         `tfsdk:"environment"`
	JobSpec       *InfrastructureStackJobSpec       `tfsdk:"job_spec"`
	Bindings      *common.ClusterBindings           `tfsdk:"bindings"`
}

func (is *infrastructureStack) Attributes(ctx context.Context, d diag.Diagnostics) gqlclient.StackAttributes {
	return gqlclient.StackAttributes{
		Name:          is.Name.ValueString(),
		Type:          gqlclient.StackType(is.Type.ValueString()),
		RepositoryID:  is.Repository.Id.ValueString(),
		ClusterID:     is.ClusterId.ValueString(),
		Git:           is.Repository.Attributes(),
		JobSpec:       is.JobSpec.Attributes(),
		Configuration: is.Configuration.Attributes(),
		Approval:      is.Approval.ValueBoolPointer(),
		ReadBindings:  is.Bindings.ReadAttributes(ctx, d),
		WriteBindings: is.Bindings.WriteAttributes(ctx, d),
		Files:         nil,
		Environemnt:   nil,
	}
}

func (is *infrastructureStack) From(stack *gqlclient.InfrastructureStackFragment, ctx context.Context, d diag.Diagnostics) {
	is.Id = types.StringPointerValue(stack.ID)
	is.Name = types.StringValue(stack.Name)
	is.Type = types.StringValue(string(stack.Type))
	is.Approval = types.BoolPointerValue(stack.Approval)
	is.ClusterId = types.StringValue(stack.Cluster.ID)
	is.Repository.From(stack.Repository, stack.Git)
	is.Configuration.From(stack.Configuration)
	is.Files = infrastructureStackFilesFrom(stack.Files, d)
	is.Environment = infrastructureStackEnvironmentsFrom(stack.Environment, ctx, d)
	is.Bindings.From(stack.ReadBindings, stack.WriteBindings, ctx, d)
	is.JobSpec.From(stack.JobSpec)
}

func infrastructureStackFilesFrom(files []*gqlclient.StackFileFragment, d diag.Diagnostics) basetypes.MapValue {
	resultMap := make(map[string]attr.Value, len(files))
	for _, file := range files {
		resultMap[file.Path] = types.StringValue(file.Content)
	}

	result, tagsDiagnostics := types.MapValue(types.StringType, resultMap)
	d.Append(tagsDiagnostics...)

	return result
}

type InfrastructureStackRepository struct {
	Id     types.String `tfsdk:"id"`
	Ref    types.String `tfsdk:"ref"`
	Folder types.String `tfsdk:"folder"`
}

func (isr *InfrastructureStackRepository) Attributes() gqlclient.GitRefAttributes {
	if isr == nil {
		return gqlclient.GitRefAttributes{}
	}

	return gqlclient.GitRefAttributes{
		Ref:    isr.Ref.ValueString(),
		Folder: isr.Folder.ValueString(),
	}
}

func (isr *InfrastructureStackRepository) From(repository *gqlclient.GitRepositoryFragment, ref *gqlclient.GitRefFragment) {
	if isr == nil {
		return
	}

	isr.Id = types.StringValue(repository.ID)

	if ref == nil {
		return
	}

	isr.Ref = types.StringValue(ref.Ref)
	isr.Folder = types.StringValue(ref.Folder)
}

type InfrastructureStackConfiguration struct {
	Image   types.String `tfsdk:"image"`
	Version types.String `tfsdk:"version"`
}

func (isc *InfrastructureStackConfiguration) Attributes() gqlclient.StackConfigurationAttributes {
	if isc == nil {
		return gqlclient.StackConfigurationAttributes{}
	}

	return gqlclient.StackConfigurationAttributes{
		Image:   isc.Image.ValueStringPointer(),
		Version: isc.Version.ValueString(),
	}
}

func (isc *InfrastructureStackConfiguration) From(configuration *gqlclient.StackConfigurationFragment) {
	if isc == nil || configuration == nil {
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

var InfrastructureStackEnvironmentAttrTypes = map[string]attr.Type{
	"name":   types.StringType,
	"value":  types.StringType,
	"secret": types.BoolType,
}

func infrastructureStackEnvironmentsFrom(envs []*gqlclient.StackEnvironmentFragment, ctx context.Context, d diag.Diagnostics) types.Set {
	if len(envs) == 0 {
		return types.SetNull(basetypes.ObjectType{AttrTypes: InfrastructureStackEnvironmentAttrTypes})
	}

	values := make([]attr.Value, len(envs))
	for i, file := range envs {
		objValue, diags := types.ObjectValueFrom(ctx, InfrastructureStackEnvironmentAttrTypes, InfrastructureStackEnvironment{
			Name:   types.StringValue(file.Name),
			Value:  types.StringValue(file.Value),
			Secret: types.BoolPointerValue(file.Secret),
		})
		values[i] = objValue
		d.Append(diags...)
	}

	setValue, diags := types.SetValue(basetypes.ObjectType{AttrTypes: InfrastructureStackEnvironmentAttrTypes}, values)
	d.Append(diags...)
	return setValue
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
	Namespace      types.String `tfsdk:"namespace"`
	Raw            types.String `tfsdk:"raw"`
	Containers     types.Set    `tfsdk:"containers"`
	Labels         types.Map    `tfsdk:"labels"`
	Annotations    types.Map    `tfsdk:"annotations"`
	ServiceAccount types.String `tfsdk:"service_account"`
}

func (isjs *InfrastructureStackJobSpec) Attributes() *gqlclient.GateJobAttributes {
	if isjs == nil {
		return nil
	}

	return &gqlclient.GateJobAttributes{} // TODO
}

func (isjs *InfrastructureStackJobSpec) From(spec *gqlclient.JobGateSpecFragment) {
	if isjs == nil {
		return
	}

	isjs.Namespace = types.StringValue(spec.Namespace)
	isjs.Raw = types.StringPointerValue(spec.Raw)
	// TODO: Containers
	// TODO: Labels
	// TODO: Annotations
	isjs.ServiceAccount = types.StringPointerValue(spec.ServiceAccount)
}

type InfrastructureStackContainerSpec struct {
	Image   types.String `tfsdk:"image"`
	Args    types.Set    `tfsdk:"args"`
	Env     types.Map    `tfsdk:"env"`
	EnvFrom types.Set    `tfsdk:"env_from"`
}

type InfrastructureStackContainerEnvFrom struct {
	Secret    types.String `tfsdk:"secret"`
	ConfigMap types.String `tfsdk:"config_map"`
}
