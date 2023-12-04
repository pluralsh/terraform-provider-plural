package resource

import (
	"context"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console-client-go"
	"github.com/pluralsh/polly/algorithms"
)

type ServiceDeployment struct {
	Id            types.String                      `tfsdk:"id"`
	Name          types.String                      `tfsdk:"name"`
	Namespace     types.String                      `tfsdk:"namespace"`
	Version       types.String                      `tfsdk:"version"`
	DocsPath      types.String                      `tfsdk:"docs_path"`
	Protect       types.Bool                        `tfsdk:"protect"`
	Kustomize     *ServiceDeploymentKustomize       `tfsdk:"kustomize"`
	Configuration []*ServiceDeploymentConfiguration `tfsdk:"configuration"`
	Cluster       *ServiceDeploymentCluster         `tfsdk:"cluster"`
	Repository    *ServiceDeploymentRepository      `tfsdk:"repository"`
	Bindings      *ServiceDeploymentBindings        `tfsdk:"bindings"`
	SyncConfig    *ServiceDeploymentSyncConfig      `tfsdk:"sync_config"`
	Helm          *ServiceDeploymentHelm            `tfsdk:"helm"`
}

func (this *ServiceDeployment) VersionString() *string {
	result := this.Version.ValueStringPointer()
	if result != nil && len(*result) == 0 {
		result = nil
	}

	return result
}

func (this *ServiceDeployment) FromCreate(response *gqlclient.ServiceDeploymentFragment) {
	this.Id = types.StringValue(response.ID)
	this.Name = types.StringValue(response.Name)
	this.Namespace = types.StringValue(response.Namespace)
	this.Protect = types.BoolPointerValue(response.Protect)
	this.Version = types.StringValue(response.Version)
	this.Kustomize.From(response.Kustomize)
	this.Configuration = ToServiceDeploymentConfiguration(response.Configuration)
	this.Repository.From(response.Repository, response.Git)
}

func (this *ServiceDeployment) FromGet(response *gqlclient.ServiceDeploymentExtended) {
	this.Id = types.StringValue(response.ID)
	this.Name = types.StringValue(response.Name)
	this.Namespace = types.StringValue(response.Namespace)
	this.Protect = types.BoolPointerValue(response.Protect)
	this.Kustomize.From(response.Kustomize)
	this.Configuration = ToServiceDeploymentConfiguration(response.Configuration)
	this.Repository.From(response.Repository, response.Git)
}

func (this *ServiceDeployment) Attributes() gqlclient.ServiceDeploymentAttributes {
	if this == nil {
		return gqlclient.ServiceDeploymentAttributes{}
	}

	return gqlclient.ServiceDeploymentAttributes{
		Name:          this.Name.ValueString(),
		Namespace:     this.Namespace.ValueString(),
		Version:       this.VersionString(),
		DocsPath:      this.DocsPath.ValueStringPointer(),
		SyncConfig:    this.SyncConfig.Attributes(),
		Protect:       this.Protect.ValueBoolPointer(),
		RepositoryID:  this.Repository.Id.ValueString(),
		Git:           this.Repository.Attributes(),
		Kustomize:     this.Kustomize.Attributes(),
		Configuration: ToServiceDeploymentConfigAttributes(this.Configuration),
		ReadBindings:  this.Bindings.ReadAttributes(),
		WriteBindings: this.Bindings.WriteAttributes(),
		Helm:          this.Helm.Attributes(),
	}
}

func (this *ServiceDeployment) UpdateAttributes() gqlclient.ServiceUpdateAttributes {
	if this == nil {
		return gqlclient.ServiceUpdateAttributes{}
	}

	return gqlclient.ServiceUpdateAttributes{
		Version:       this.Version.ValueStringPointer(),
		Protect:       this.Protect.ValueBoolPointer(),
		Git:           this.Repository.Attributes(),
		Configuration: ToServiceDeploymentConfigAttributes(this.Configuration),
		Kustomize:     this.Kustomize.Attributes(),
		Helm:          this.Helm.Attributes(),
	}
}

type ServiceDeploymentConfiguration struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func ToServiceDeploymentConfiguration(configuration []*struct {
	Name  string "json:\"name\" graphql:\"name\""
	Value string "json:\"value\" graphql:\"value\""
}) []*ServiceDeploymentConfiguration {
	result := make([]*ServiceDeploymentConfiguration, len(configuration))
	for i, c := range configuration {
		result[i] = &ServiceDeploymentConfiguration{
			Name:  types.StringValue(c.Name),
			Value: types.StringValue(c.Value),
		}
	}

	return result
}

func ToServiceDeploymentConfigAttributes(configuration []*ServiceDeploymentConfiguration) []*gqlclient.ConfigAttributes {
	result := make([]*gqlclient.ConfigAttributes, len(configuration))
	for i, c := range configuration {
		result[i] = &gqlclient.ConfigAttributes{
			Name:  c.Name.ValueString(),
			Value: c.Value.ValueStringPointer(),
		}
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
	DiffNormalizer    *ServiceDeploymentDiffNormalizer    `tfsdk:"diff_normalizer"`
	NamespaceMetadata *ServiceDeploymentNamespaceMetadata `tfsdk:"namespace_metadata"`
}

func (this *ServiceDeploymentSyncConfig) Attributes() *gqlclient.SyncConfigAttributes {
	if this == nil {
		return nil
	}

	return &gqlclient.SyncConfigAttributes{
		DiffNormalizer:    this.DiffNormalizer.Attributes(),
		NamespaceMetadata: this.NamespaceMetadata.Attributes(),
	}
}

type ServiceDeploymentDiffNormalizer struct {
	Group       types.String `tfsdk:"group"`
	JsonPatches types.Set    `tfsdk:"json_patches"`
	Kind        types.String `tfsdk:"kind"`
	Name        types.String `tfsdk:"name"`
	Namespace   types.String `tfsdk:"namespace"`
}

func (this *ServiceDeploymentDiffNormalizer) Attributes() *gqlclient.DiffNormalizerAttributes {
	if this == nil {
		return nil
	}

	jsonPatches := make([]types.String, len(this.JsonPatches.Elements()))
	this.JsonPatches.ElementsAs(context.Background(), &jsonPatches, false)

	return &gqlclient.DiffNormalizerAttributes{
		Group:     this.Group.ValueString(),
		Kind:      this.Kind.ValueString(),
		Name:      this.Name.ValueString(),
		Namespace: this.Namespace.ValueString(),
		JSONPatches: algorithms.Map(jsonPatches, func(v types.String) string {
			return v.ValueString()
		}),
	}
}

type ServiceDeploymentNamespaceMetadata struct {
	Annotations types.Map `tfsdk:"annotations"`
	Labels      types.Map `tfsdk:"labels"`
}

func (this *ServiceDeploymentNamespaceMetadata) Attributes() *gqlclient.MetadataAttributes {
	if this == nil {
		return nil
	}

	annotations := make(map[string]types.String, len(this.Annotations.Elements()))
	labels := make(map[string]types.String, len(this.Labels.Elements()))

	this.Annotations.ElementsAs(context.Background(), &annotations, false)
	this.Labels.ElementsAs(context.Background(), &labels, false)

	return &gqlclient.MetadataAttributes{
		Annotations: common.ToAttributesMap(annotations),
		Labels:      common.ToAttributesMap(labels),
	}
}

type ServiceDeploymentHelm struct {
	Chart       types.String                     `tfsdk:"chart"`
	Repository  *ServiceDeploymentNamespacedName `tfsdk:"repository"`
	Values      types.String                     `tfsdk:"values"`
	ValuesFiles types.Set                        `tfsdk:"values_files"`
	Version     types.String                     `tfsdk:"version"`
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
