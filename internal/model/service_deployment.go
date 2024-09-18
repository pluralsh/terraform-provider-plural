package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/pluralsh/polly/algorithms"
	"github.com/samber/lo"
)

type ServiceDeployment struct {
	Id            types.String                 `tfsdk:"id"`
	Name          types.String                 `tfsdk:"name"`
	Namespace     types.String                 `tfsdk:"namespace"`
	Version       types.String                 `tfsdk:"version"`
	DocsPath      types.String                 `tfsdk:"docs_path"`
	Protect       types.Bool                   `tfsdk:"protect"`
	Templated     types.Bool                   `tfsdk:"templated"`
	Kustomize     *ServiceDeploymentKustomize  `tfsdk:"kustomize"`
	Configuration types.Map                    `tfsdk:"configuration"`
	Cluster       *ServiceDeploymentCluster    `tfsdk:"cluster"`
	Repository    *ServiceDeploymentRepository `tfsdk:"repository"`
	Bindings      *ServiceDeploymentBindings   `tfsdk:"bindings"`
	SyncConfig    *ServiceDeploymentSyncConfig `tfsdk:"sync_config"`
	Helm          *ServiceDeploymentHelm       `tfsdk:"helm"`
}

func (this *ServiceDeployment) VersionString() *string {
	result := this.Version.ValueStringPointer()
	if result != nil && len(*result) == 0 {
		result = nil
	}

	return result
}

func (this *ServiceDeployment) FromCreate(response *gqlclient.ServiceDeploymentExtended, d diag.Diagnostics) {
	this.Id = types.StringValue(response.ID)
	this.Name = types.StringValue(response.Name)
	this.Namespace = types.StringValue(response.Namespace)
	this.Protect = types.BoolPointerValue(response.Protect)
	this.Version = types.StringValue(response.Version)
	this.Kustomize.From(response.Kustomize)
	this.Configuration = ToServiceDeploymentConfiguration(response.Configuration, d)
	this.Repository.From(response.Repository, response.Git)
	this.Templated = types.BoolPointerValue(response.Templated)
}

func (this *ServiceDeployment) FromGet(response *gqlclient.ServiceDeploymentExtended, d diag.Diagnostics) {
	this.Id = types.StringValue(response.ID)
	this.Name = types.StringValue(response.Name)
	this.Namespace = types.StringValue(response.Namespace)
	this.Protect = types.BoolPointerValue(response.Protect)
	this.Kustomize.From(response.Kustomize)
	this.Configuration = ToServiceDeploymentConfiguration(response.Configuration, d)
	this.Repository.From(response.Repository, response.Git)
	this.Templated = types.BoolPointerValue(response.Templated)
}

func (this *ServiceDeployment) Attributes(ctx context.Context, d diag.Diagnostics) gqlclient.ServiceDeploymentAttributes {
	if this == nil {
		return gqlclient.ServiceDeploymentAttributes{}
	}

	var repositoryId *string = nil
	if this.Repository != nil && this.Repository.Id.ValueStringPointer() != nil {
		repositoryId = this.Repository.Id.ValueStringPointer()
	}

	return gqlclient.ServiceDeploymentAttributes{
		Name:          this.Name.ValueString(),
		Namespace:     this.Namespace.ValueString(),
		Version:       this.VersionString(),
		DocsPath:      this.DocsPath.ValueStringPointer(),
		SyncConfig:    this.SyncConfig.Attributes(d),
		Protect:       this.Protect.ValueBoolPointer(),
		RepositoryID:  repositoryId,
		Git:           this.Repository.Attributes(),
		Kustomize:     this.Kustomize.Attributes(),
		Configuration: this.ToServiceDeploymentConfigAttributes(ctx, d),
		ReadBindings:  this.Bindings.ReadAttributes(),
		WriteBindings: this.Bindings.WriteAttributes(),
		Helm:          this.Helm.Attributes(),
		Templated:     this.Templated.ValueBoolPointer(),
	}
}

func (this *ServiceDeployment) UpdateAttributes(ctx context.Context, d diag.Diagnostics) gqlclient.ServiceUpdateAttributes {
	if this == nil {
		return gqlclient.ServiceUpdateAttributes{}
	}

	return gqlclient.ServiceUpdateAttributes{
		Version:       this.Version.ValueStringPointer(),
		Protect:       this.Protect.ValueBoolPointer(),
		Git:           this.Repository.Attributes(),
		Configuration: this.ToServiceDeploymentConfigAttributes(ctx, d),
		Kustomize:     this.Kustomize.Attributes(),
		Helm:          this.Helm.Attributes(),
		Templated:     this.Templated.ValueBoolPointer(),
	}
}

type ServiceDeploymentConfiguration struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func ToServiceDeploymentConfiguration(configuration []*gqlclient.ServiceDeploymentExtended_ServiceDeploymentFragment_Configuration, d diag.Diagnostics) basetypes.MapValue {
	resultMap := make(map[string]attr.Value, len(configuration))
	for _, c := range configuration {
		resultMap[c.Name] = types.StringValue(c.Value)
	}

	result, tagsDiagnostics := types.MapValue(types.StringType, resultMap)
	d.Append(tagsDiagnostics...)

	return result
}

func (this *ServiceDeployment) ToServiceDeploymentConfigAttributes(ctx context.Context, d diag.Diagnostics) []*gqlclient.ConfigAttributes {
	result := make([]*gqlclient.ConfigAttributes, 0)
	elements := make(map[string]types.String, len(this.Configuration.Elements()))
	d.Append(this.Configuration.ElementsAs(ctx, &elements, false)...)

	for k, v := range elements {
		result = append(result, &gqlclient.ConfigAttributes{Name: k, Value: lo.ToPtr(v.ValueString())})
	}

	return result
}

type ServiceDeploymentCluster struct {
	Id     types.String `tfsdk:"id"`
	Handle types.String `tfsdk:"handle"`
}

func (this *ServiceDeploymentCluster) From(cluster *gqlclient.BaseClusterFragment) {
	if this == nil {
		return
	}

	this.Id = types.StringValue(cluster.ID)
	this.Handle = types.StringPointerValue(cluster.Handle)
}

type ServiceDeploymentRepository struct {
	Id     types.String `tfsdk:"id"`
	Ref    types.String `tfsdk:"ref"`
	Folder types.String `tfsdk:"folder"`
}

func (this *ServiceDeploymentRepository) From(repository *gqlclient.GitRepositoryFragment, git *gqlclient.GitRefFragment) {
	if this == nil {
		return
	}

	this.Id = types.StringValue(repository.ID)

	if git == nil {
		return
	}

	this.Ref = types.StringValue(git.Ref)
	this.Folder = types.StringValue(git.Folder)
}

func (this *ServiceDeploymentRepository) Attributes() *gqlclient.GitRefAttributes {
	if this == nil {
		return nil
	}

	if len(this.Ref.ValueString()) == 0 && len(this.Folder.ValueString()) == 0 {
		return nil
	}

	return &gqlclient.GitRefAttributes{
		Ref:    this.Ref.ValueString(),
		Folder: this.Folder.ValueString(),
	}
}

type ServiceDeploymentKustomize struct {
	Path types.String `tfsdk:"path"`
}

func (this *ServiceDeploymentKustomize) From(kustomize *gqlclient.KustomizeFragment) {
	if this == nil {
		return
	}

	this.Path = types.StringValue(kustomize.Path)
}

func (this *ServiceDeploymentKustomize) Attributes() *gqlclient.KustomizeAttributes {
	if this == nil {
		return nil
	}

	return &gqlclient.KustomizeAttributes{
		Path: this.Path.ValueString(),
	}
}

type ServiceDeploymentBindings struct {
	Read  []*ServiceDeploymentPolicyBinding `tfsdk:"read"`
	Write []*ServiceDeploymentPolicyBinding `tfsdk:"write"`
}

func (this *ServiceDeploymentBindings) ReadAttributes() []*gqlclient.PolicyBindingAttributes {
	if this == nil {
		return []*gqlclient.PolicyBindingAttributes{}
	}

	return this.attributes(this.Read)
}

func (this *ServiceDeploymentBindings) WriteAttributes() []*gqlclient.PolicyBindingAttributes {
	if this == nil {
		return []*gqlclient.PolicyBindingAttributes{}
	}

	return this.attributes(this.Write)
}

func (this *ServiceDeploymentBindings) attributes(bindings []*ServiceDeploymentPolicyBinding) []*gqlclient.PolicyBindingAttributes {
	result := make([]*gqlclient.PolicyBindingAttributes, len(bindings))
	for i, b := range bindings {
		result[i] = &gqlclient.PolicyBindingAttributes{
			ID:      b.ID.ValueStringPointer(),
			UserID:  b.UserID.ValueStringPointer(),
			GroupID: b.GroupID.ValueStringPointer(),
		}
	}

	return result
}

type ServiceDeploymentPolicyBinding struct {
	GroupID types.String `tfsdk:"group_id"`
	ID      types.String `tfsdk:"id"`
	UserID  types.String `tfsdk:"user_id"`
}

type ServiceDeploymentSyncConfig struct {
	NamespaceMetadata *ServiceDeploymentNamespaceMetadata `tfsdk:"namespace_metadata"`
}

func (this *ServiceDeploymentSyncConfig) Attributes(d diag.Diagnostics) *gqlclient.SyncConfigAttributes {
	if this == nil {
		return nil
	}

	return &gqlclient.SyncConfigAttributes{
		NamespaceMetadata: this.NamespaceMetadata.Attributes(d),
	}
}

type ServiceDeploymentNamespaceMetadata struct {
	Annotations types.Map `tfsdk:"annotations"`
	Labels      types.Map `tfsdk:"labels"`
}

func (this *ServiceDeploymentNamespaceMetadata) Attributes(d diag.Diagnostics) *gqlclient.MetadataAttributes {
	if this == nil {
		return nil
	}

	annotations := make(map[string]types.String, len(this.Annotations.Elements()))
	labels := make(map[string]types.String, len(this.Labels.Elements()))

	this.Annotations.ElementsAs(context.Background(), &annotations, false)
	this.Labels.ElementsAs(context.Background(), &labels, false)

	return &gqlclient.MetadataAttributes{
		Annotations: common.AttributesJson(annotations, d),
		Labels:      common.AttributesJson(labels, d),
	}
}

type ServiceDeploymentHelm struct {
	Chart       types.String                     `tfsdk:"chart"`
	Repository  *ServiceDeploymentNamespacedName `tfsdk:"repository"`
	Values      types.String                     `tfsdk:"values"`
	ValuesFiles types.Set                        `tfsdk:"values_files"`
	Version     types.String                     `tfsdk:"version"`
	URL         types.String                     `tfsdk:"url"`
}

func (this *ServiceDeploymentHelm) Attributes() *gqlclient.HelmConfigAttributes {
	if this == nil {
		return nil
	}

	valuesFiles := make([]types.String, len(this.ValuesFiles.Elements()))
	this.ValuesFiles.ElementsAs(context.Background(), &valuesFiles, false)

	return &gqlclient.HelmConfigAttributes{
		Values: this.Values.ValueStringPointer(),
		ValuesFiles: algorithms.Map(valuesFiles, func(v types.String) *string {
			return v.ValueStringPointer()
		}),
		Chart:      this.Chart.ValueStringPointer(),
		Version:    this.Version.ValueStringPointer(),
		Repository: this.Repository.Attributes(),
		URL:        this.URL.ValueStringPointer(),
	}
}

type ServiceDeploymentNamespacedName struct {
	Name      types.String `tfsdk:"name"`
	Namespace types.String `tfsdk:"namespace"`
}

func (this *ServiceDeploymentNamespacedName) Attributes() *gqlclient.NamespacedName {
	if this == nil {
		return nil
	}

	return &gqlclient.NamespacedName{
		Name:      this.Name.ValueString(),
		Namespace: this.Namespace.ValueString(),
	}
}
