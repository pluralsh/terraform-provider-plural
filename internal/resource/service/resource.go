package service

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	consoleClient "github.com/pluralsh/console-client-go"
)

var _ resource.Resource = &ServiceDeploymentResource{}
var _ resource.ResourceWithImportState = &ServiceDeploymentResource{}

func NewServiceDeploymentResource() resource.Resource {
	return &ServiceDeploymentResource{}
}

// ServiceDeploymentResource defines the ServiceDeployment resource implementation.
type ServiceDeploymentResource struct {
	client *consoleClient.Client
}

// ServiceDeploymentResourceModel describes the ServiceDeployment resource data model.
type ServiceDeploymentResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Namespace     types.String `tfsdk:"namespace"`
	Configuration types.List   `tfsdk:"configuration"`
	Cluster       types.Object `tfsdk:"cluster"`
	Repository    types.Object `tfsdk:"repository"`
}

type ClusterModel struct {
	Id     string `json:"id"`
	Handle string `json:"handle"`
}

type RepositoryModel struct {
	Id     string `json:"id"`
	Ref    string `json:"ref"`
	Folder string `json:"folder"`
}

func (r *ServiceDeploymentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ServiceDeployment"
}

func (r *ServiceDeploymentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "ServiceDeployment resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal identifier of this ServiceDeployment.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable name of this ServiceDeployment.",
				Required:            true,
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: "Namespace to deploy this ServiceDeployment.",
				Required:            true,
			},
			"configuration": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":  schema.StringAttribute{},
						"value": schema.StringAttribute{},
					},
				},
				MarkdownDescription: "List of [name, value] secrets used to alter this ServiceDeployment configuration.",
				Optional:            true,
			},
			"cluster": schema.ObjectAttribute{
				AttributeTypes: map[string]attr.Type{
					"id":     types.StringType,
					"handle": types.StringType,
				},
				Validators:          []validator.Object{objectvalidator.ExactlyOneOf(path.MatchRelative().AtName("id"), path.MatchRelative().AtName("handle"))},
				MarkdownDescription: "Unique cluster id/handle to deploy this ServiceDeployment",
				Required:            true,
			},
			"repository": schema.ObjectAttribute{
				AttributeTypes: map[string]attr.Type{
					"id":     types.StringType,
					"ref":    types.StringType,
					"folder": types.StringType,
				},
				MarkdownDescription: "Repository information used to pull ServiceDeployment.",
				Required:            true,
			},
		},
	}
}

func (r *ServiceDeploymentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*consoleClient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ServiceDeployment Resource Configure Type",
			fmt.Sprintf("Expected *console.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ServiceDeploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServiceDeploymentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var cluster ClusterModel
	data.Cluster.As(ctx, cluster, basetypes.ObjectAsOptions{})

	var repository RepositoryModel
	data.Repository.As(ctx, repository, basetypes.ObjectAsOptions{})

	var configuration []*consoleClient.ConfigAttributes
	data.Repository.As(ctx, repository, basetypes.ObjectAsOptions{})

	attrs := consoleClient.ServiceDeploymentAttributes{
		Name:         data.Name.String(),
		RepositoryID: repository.Id,
		Git: consoleClient.GitRefAttributes{
			Ref:    repository.Ref,
			Folder: repository.Folder,
		},
		Configuration: configuration,
	}

	_, err := CreateServiceDeployment(ctx, r.client, &cluster.Id, &cluster.Handle, attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create ServiceDeployment, got error: %s", err))
		return
	}

	// TODO: figure out if we need to read response and update state

	tflog.Trace(ctx, "created a ServiceDeployment")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceDeploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ServiceDeploymentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ServiceDeployment, err := r.client.GetServiceDeployment(ctx, data.Id.String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ServiceDeployment, got error: %s", err))
		return
	}

	data.Id = types.StringValue(ServiceDeployment.ServiceDeployment.ID)
	data.Name = types.StringValue(ServiceDeployment.ServiceDeployment.Name)
	// TODO: read rest of the config

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceDeploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ServiceDeploymentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: figure out what can be updated
	//attrs := consoleClient.ServiceDeploymentUpdateAttributes{
	//	Handle: lo.ToPtr(data.Handle.String()),
	//}
	//ServiceDeployment, err := r.client.UpdateServiceDeployment(ctx, data.Id.String(), attrs)
	//if err != nil {
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update ServiceDeployment, got error: %s", err))
	//	return
	//}
	//
	//data.Handle = types.StringValue(*ServiceDeployment.UpdateServiceDeployment.Handle)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceDeploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ServiceDeploymentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteServiceDeployment(ctx, data.Id.String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete ServiceDeployment, got error: %s", err))
		return
	}
}

func (r *ServiceDeploymentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
