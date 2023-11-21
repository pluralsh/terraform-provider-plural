package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console-client-go"
	"github.com/pluralsh/polly/algorithms"
	"github.com/samber/lo"
)

// ServiceDeployment describes the service deployment resource and data source model.
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
}

func NewServiceDeployment() *ServiceDeployment {
	return &ServiceDeployment{
		Id:            types.String{},
		Name:          types.String{},
		Namespace:     types.String{},
		Version:       types.String{},
		DocsPath:      types.String{},
		Protect:       types.Bool{},
		Kustomize:     &ServiceDeploymentKustomize{},
		Configuration: []*ServiceDeploymentConfiguration{},
		Cluster: &ServiceDeploymentCluster{
			Id:     types.String{},
			Handle: types.String{},
		},
		Repository: &ServiceDeploymentRepository{
			Id:     types.String{},
			Ref:    types.String{},
			Folder: types.String{},
		},
		Bindings: &ServiceDeploymentBindings{
			Read:  []*ServiceDeploymentPolicyBinding{},
			Write: []*ServiceDeploymentPolicyBinding{},
		},
		SyncConfig: &ServiceDeploymentSyncConfig{
			DiffNormalizer:    &ServiceDeploymentDiffNormalizer{
				Group:       types.String{},
				JsonPatches: types.List{},
				Kind:        types.String{},
				Name:        types.String{},
				Namespace:   types.String{},
			},
			NamespaceMetadata: &ServiceDeploymentNamespaceMetadata{
				Annotations: types.Map{},
				Labels:      types.Map{},
			},
		},
	}
}

func (this *ServiceDeployment) From(response *gqlclient.GetServiceDeployment) {
	this.Id = types.StringValue(response.ServiceDeployment.ID)
	this.Name = types.StringValue(response.ServiceDeployment.Name)
	this.Namespace = types.StringValue(response.ServiceDeployment.Namespace)
	//this.Protect = types.BoolValue(response.ServiceDeployment.Protect)
	this.Version = types.StringValue(response.ServiceDeployment.Version)
	this.Kustomize.From(response.ServiceDeployment.Kustomize)
	this.Configuration = ToServiceDeploymentConfiguration(response.ServiceDeployment.Configuration)
	this.Cluster.From(response.ServiceDeployment.Cluster)
	this.Repository.From(response.ServiceDeployment)
}

func (this *ServiceDeployment) Attributes() gqlclient.ServiceDeploymentAttributes {
	return gqlclient.ServiceDeploymentAttributes{
		Name:          this.Name.ValueString(),
		Namespace:     this.Namespace.ValueString(),
		Version:       this.Version.ValueStringPointer(),
		DocsPath:      this.DocsPath.ValueStringPointer(),
		SyncConfig:    this.SyncConfig.Attributes(),
		Protect:       this.Protect.ValueBoolPointer(),
		RepositoryID:  this.Repository.Id.ValueString(),
		Git:           this.Repository.Attributes(),
		Kustomize:     this.Kustomize.Attributes(),
		Configuration: ToServiceDeploymentConfigAttributes(this.Configuration),
		ReadBindings:  this.Bindings.ReadAttributes(),
		WriteBindings: this.Bindings.WriteAttributes(),
	}
}

func (this *ServiceDeployment) UpdateAttributes() gqlclient.ServiceUpdateAttributes {
	return gqlclient.ServiceUpdateAttributes{
		Version:       this.Version.ValueStringPointer(),
		Protect:       this.Protect.ValueBoolPointer(),
		Git:           lo.ToPtr(this.Repository.Attributes()),
		Configuration: ToServiceDeploymentConfigAttributes(this.Configuration),
		Kustomize:     this.Kustomize.Attributes(),
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
	this.Id = types.StringValue(cluster.ID)
	this.Handle = types.StringPointerValue(cluster.Handle)
}

type ServiceDeploymentRepository struct {
	Id     types.String `tfsdk:"id"`
	Ref    types.String `tfsdk:"ref"`
	Folder types.String `tfsdk:"folder"`
}

func (this *ServiceDeploymentRepository) From(deployment *gqlclient.ServiceDeploymentExtended) {
	this.Id = types.StringValue(deployment.Repository.ID)
	this.Ref = types.StringValue(deployment.Git.Ref)
	this.Folder = types.StringValue(deployment.Git.Folder)
}

func (this *ServiceDeploymentRepository) Attributes() gqlclient.GitRefAttributes {
	return gqlclient.GitRefAttributes{
		Ref:    this.Ref.ValueString(),
		Folder: this.Folder.ValueString(),
	}
}

type ServiceDeploymentKustomize struct {
	Path types.String `tfsdk:"path"`
}

func (this *ServiceDeploymentKustomize) From(kustomize *gqlclient.KustomizeFragment) {
	this.Path = types.StringValue(kustomize.Path)
}

func (this *ServiceDeploymentKustomize) Attributes() *gqlclient.KustomizeAttributes {
	return &gqlclient.KustomizeAttributes{
		Path: this.Path.ValueString(),
	}
}

type ServiceDeploymentBindings struct {
	Read  []*ServiceDeploymentPolicyBinding `tfsdk:"read"`
	Write []*ServiceDeploymentPolicyBinding `tfsdk:"write"`
}

func (this *ServiceDeploymentBindings) ReadAttributes() []*gqlclient.PolicyBindingAttributes {
	return this.attributes(this.Read)
}

func (this *ServiceDeploymentBindings) WriteAttributes() []*gqlclient.PolicyBindingAttributes {
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
	return &gqlclient.SyncConfigAttributes{
		DiffNormalizer:    this.DiffNormalizer.Attributes(),
		NamespaceMetadata: this.NamespaceMetadata.Attributes(),
	}
}

type ServiceDeploymentDiffNormalizer struct {
	Group       types.String `tfsdk:"group"`
	JsonPatches types.List   `tfsdk:"json_patches"`
	Kind        types.String `tfsdk:"kind"`
	Name        types.String `tfsdk:"name"`
	Namespace   types.String `tfsdk:"namespace"`
}

func (this *ServiceDeploymentDiffNormalizer) Attributes() *gqlclient.DiffNormalizerAttributes {
	return &gqlclient.DiffNormalizerAttributes{
		Group:     this.Group.ValueString(),
		Kind:      this.Kind.ValueString(),
		Name:      this.Name.ValueString(),
		Namespace: this.Namespace.ValueString(),
		JSONPatches: algorithms.Map(this.JsonPatches.Elements(), func(v attr.Value) string {
			return v.String()
		}),
	}
}

type ServiceDeploymentNamespaceMetadata struct {
	Annotations types.Map `tfsdk:"annotations"`
	Labels      types.Map `tfsdk:"labels"`
}

func (this *ServiceDeploymentNamespaceMetadata) Attributes() *gqlclient.MetadataAttributes {
	return &gqlclient.MetadataAttributes{
		Annotations: this.toAttributesMap(this.Annotations.Elements()),
		Labels:      this.toAttributesMap(this.Labels.Elements()),
	}
}

func (this *ServiceDeploymentNamespaceMetadata) toAttributesMap(m map[string]attr.Value) map[string]interface{} {
	result := map[string]interface{}{}
	for key, val := range m {
		result[key] = val.String()
	}

	return result
}
