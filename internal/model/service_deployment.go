package model

import (
	"context"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
	"github.com/pluralsh/polly/algorithms"
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
	Bindings      *common.Bindings             `tfsdk:"bindings"`
	SyncConfig    *ServiceDeploymentSyncConfig `tfsdk:"sync_config"`
	Helm          *ServiceDeploymentHelm       `tfsdk:"helm"`
}

func (sd *ServiceDeployment) VersionString() *string {
	result := sd.Version.ValueStringPointer()
	if result != nil && len(*result) == 0 {
		result = nil
	}

	return result
}

func (sd *ServiceDeployment) FromCreate(response *gqlclient.ServiceDeploymentExtended, d *diag.Diagnostics) {
	sd.Id = types.StringValue(response.ID)
	sd.Name = types.StringValue(response.Name)
	sd.Namespace = types.StringValue(response.Namespace)
	sd.Protect = types.BoolPointerValue(response.Protect)
	sd.Version = types.StringValue(response.Version)
	sd.Kustomize.From(response.Kustomize)
	sd.Configuration = configFrom(response.Configuration, d)
	sd.Cluster.From(response.Cluster)
	sd.Repository.From(response.Repository, response.Git)
	sd.Templated = types.BoolPointerValue(response.Templated)
}

func (sd *ServiceDeployment) FromGet(response *gqlclient.ServiceDeploymentExtended, d *diag.Diagnostics) {
	sd.Id = types.StringValue(response.ID)
	sd.Name = types.StringValue(response.Name)
	sd.Namespace = types.StringValue(response.Namespace)
	sd.Protect = types.BoolPointerValue(response.Protect)
	sd.Kustomize.From(response.Kustomize)
	sd.Configuration = configFrom(response.Configuration, d)
	sd.Repository.From(response.Repository, response.Git)
	sd.Templated = types.BoolPointerValue(response.Templated)
}

func (sd *ServiceDeployment) Attributes(ctx context.Context, d *diag.Diagnostics) gqlclient.ServiceDeploymentAttributes {
	if sd == nil {
		return gqlclient.ServiceDeploymentAttributes{}
	}

	var repositoryId *string = nil
	if sd.Repository != nil && sd.Repository.Id.ValueStringPointer() != nil {
		repositoryId = sd.Repository.Id.ValueStringPointer()
	}

	return gqlclient.ServiceDeploymentAttributes{
		Name:          sd.Name.ValueString(),
		Namespace:     sd.Namespace.ValueString(),
		Version:       sd.VersionString(),
		DocsPath:      sd.DocsPath.ValueStringPointer(),
		SyncConfig:    sd.SyncConfig.Attributes(d),
		Protect:       sd.Protect.ValueBoolPointer(),
		RepositoryID:  repositoryId,
		Git:           sd.Repository.Attributes(),
		Kustomize:     sd.Kustomize.Attributes(),
		Configuration: sd.ToServiceDeploymentConfigAttributes(ctx, d),
		ReadBindings:  sd.Bindings.ReadAttributes(ctx, d),
		WriteBindings: sd.Bindings.WriteAttributes(ctx, d),
		Helm:          sd.Helm.Attributes(),
		Templated:     sd.Templated.ValueBoolPointer(),
	}
}

func (sd *ServiceDeployment) UpdateAttributes(ctx context.Context, d *diag.Diagnostics) gqlclient.ServiceUpdateAttributes {
	if sd == nil {
		return gqlclient.ServiceUpdateAttributes{}
	}

	return gqlclient.ServiceUpdateAttributes{
		Version:       sd.Version.ValueStringPointer(),
		Protect:       sd.Protect.ValueBoolPointer(),
		Git:           sd.Repository.Attributes(),
		Configuration: sd.ToServiceDeploymentConfigAttributes(ctx, d),
		Kustomize:     sd.Kustomize.Attributes(),
		Helm:          sd.Helm.Attributes(),
		Templated:     sd.Templated.ValueBoolPointer(),
	}
}

type ServiceDeploymentConfiguration struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func configFrom(configuration []*gqlclient.ServiceDeploymentExtended_ServiceDeploymentFragment_Configuration, d *diag.Diagnostics) basetypes.MapValue {
	if len(configuration) == 0 {
		return types.MapNull(types.StringType)
	}

	resultMap := make(map[string]attr.Value, len(configuration))
	for _, c := range configuration {
		resultMap[c.Name] = types.StringValue(c.Value)
	}

	result, diags := types.MapValue(types.StringType, resultMap)
	d.Append(diags...)

	return result
}

func (sd *ServiceDeployment) ToServiceDeploymentConfigAttributes(ctx context.Context, d *diag.Diagnostics) []*gqlclient.ConfigAttributes {
	if sd.Configuration.IsNull() || sd.Configuration.IsUnknown() {
		return nil
	}

	result := make([]*gqlclient.ConfigAttributes, 0)
	elements := make(map[string]types.String, len(sd.Configuration.Elements()))
	d.Append(sd.Configuration.ElementsAs(ctx, &elements, false)...)

	for k, v := range elements {
		result = append(result, &gqlclient.ConfigAttributes{Name: k, Value: v.ValueStringPointer()})
	}

	return result
}

type ServiceDeploymentCluster struct {
	Id     types.String `tfsdk:"id"`
	Handle types.String `tfsdk:"handle"`
}

func (sdc *ServiceDeploymentCluster) From(cluster *gqlclient.BaseClusterFragment) {
	if sdc == nil {
		return
	}

	sdc.Id = types.StringValue(cluster.ID)
	sdc.Handle = types.StringPointerValue(cluster.Handle)
}

type ServiceDeploymentRepository struct {
	Id     types.String `tfsdk:"id"`
	Ref    types.String `tfsdk:"ref"`
	Folder types.String `tfsdk:"folder"`
}

func (sdr *ServiceDeploymentRepository) From(repository *gqlclient.GitRepositoryFragment, git *gqlclient.GitRefFragment) {
	if sdr == nil {
		return
	}

	sdr.Id = types.StringValue(repository.ID)

	if git == nil {
		return
	}

	sdr.Ref = types.StringValue(git.Ref)
	sdr.Folder = types.StringValue(git.Folder)
}

func (sdr *ServiceDeploymentRepository) Attributes() *gqlclient.GitRefAttributes {
	if sdr == nil {
		return nil
	}

	if len(sdr.Ref.ValueString()) == 0 && len(sdr.Folder.ValueString()) == 0 {
		return nil
	}

	return &gqlclient.GitRefAttributes{
		Ref:    sdr.Ref.ValueString(),
		Folder: sdr.Folder.ValueString(),
	}
}

type ServiceDeploymentKustomize struct {
	Path types.String `tfsdk:"path"`
}

func (sdk *ServiceDeploymentKustomize) From(kustomize *gqlclient.KustomizeFragment) {
	if sdk == nil {
		return
	}

	sdk.Path = types.StringValue(kustomize.Path)
}

func (sdk *ServiceDeploymentKustomize) Attributes() *gqlclient.KustomizeAttributes {
	if sdk == nil {
		return nil
	}

	return &gqlclient.KustomizeAttributes{
		Path: sdk.Path.ValueString(),
	}
}

type ServiceDeploymentSyncConfig struct {
	NamespaceMetadata *ServiceDeploymentNamespaceMetadata `tfsdk:"namespace_metadata"`
}

func (sdsc *ServiceDeploymentSyncConfig) Attributes(d *diag.Diagnostics) *gqlclient.SyncConfigAttributes {
	if sdsc == nil {
		return nil
	}

	return &gqlclient.SyncConfigAttributes{
		NamespaceMetadata: sdsc.NamespaceMetadata.Attributes(d),
	}
}

type ServiceDeploymentNamespaceMetadata struct {
	Annotations types.Map `tfsdk:"annotations"`
	Labels      types.Map `tfsdk:"labels"`
}

func (sdnm *ServiceDeploymentNamespaceMetadata) Attributes(d *diag.Diagnostics) *gqlclient.MetadataAttributes {
	if sdnm == nil {
		return nil
	}

	annotations := make(map[string]types.String, len(sdnm.Annotations.Elements()))
	labels := make(map[string]types.String, len(sdnm.Labels.Elements()))

	sdnm.Annotations.ElementsAs(context.Background(), &annotations, false)
	sdnm.Labels.ElementsAs(context.Background(), &labels, false)

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

func (sdh *ServiceDeploymentHelm) Attributes() *gqlclient.HelmConfigAttributes {
	if sdh == nil {
		return nil
	}

	valuesFiles := make([]types.String, len(sdh.ValuesFiles.Elements()))
	sdh.ValuesFiles.ElementsAs(context.Background(), &valuesFiles, false)

	return &gqlclient.HelmConfigAttributes{
		Values: sdh.Values.ValueStringPointer(),
		ValuesFiles: algorithms.Map(valuesFiles, func(v types.String) *string {
			return v.ValueStringPointer()
		}),
		Chart:      sdh.Chart.ValueStringPointer(),
		Version:    sdh.Version.ValueStringPointer(),
		Repository: sdh.Repository.Attributes(),
		URL:        sdh.URL.ValueStringPointer(),
	}
}

type ServiceDeploymentNamespacedName struct {
	Name      types.String `tfsdk:"name"`
	Namespace types.String `tfsdk:"namespace"`
}

func (sdnn *ServiceDeploymentNamespacedName) Attributes() *gqlclient.NamespacedName {
	if sdnn == nil {
		return nil
	}

	return &gqlclient.NamespacedName{
		Name:      sdnn.Name.ValueString(),
		Namespace: sdnn.Namespace.ValueString(),
	}
}
