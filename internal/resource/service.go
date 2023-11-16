package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	"github.com/pluralsh/polly/algorithms"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/model"
)

var _ resource.Resource = &ServiceDeploymentResource{}
var _ resource.ResourceWithImportState = &ServiceDeploymentResource{}

func NewServiceDeploymentResource() resource.Resource {
	return &ServiceDeploymentResource{}
}

// ServiceDeploymentResource defines the ServiceDeployment resource implementation.
type ServiceDeploymentResource struct {
	client *client.Client
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
	Id     types.String `tfsdk:"id"`
	Handle types.String `tfsdk:"handle"`
}

type RepositoryModel struct {
	Id     string `tfsdk:"id"`
	Ref    string `tfsdk:"ref"`
	Folder string `tfsdk:"folder"`
}

type ConfigurationModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func (r *ServiceDeploymentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_deployment"
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
						"name": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Required:  true,
							Sensitive: true,
						},
					},
				},
				MarkdownDescription: "List of [name, value] secrets used to alter this ServiceDeployment configuration.",
				Optional:            true,
			},
			"cluster": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRoot("cluster").AtName("handle")),
							stringvalidator.ExactlyOneOf(path.MatchRoot("cluster").AtName("handle")),
						},
						Optional: true,
					},
					"handle": schema.StringAttribute{
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRoot("cluster").AtName("id")),
							stringvalidator.ExactlyOneOf(path.MatchRoot("cluster").AtName("id")),
						},
						Optional: true,
					},
				},
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

	data, ok := req.ProviderData.(*model.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ServiceDeployment Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = data.Client
}

func (r *ServiceDeploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServiceDeploymentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var cluster ClusterModel
	data.Cluster.As(ctx, &cluster, basetypes.ObjectAsOptions{})

	var repository RepositoryModel
	data.Repository.As(ctx, &repository, basetypes.ObjectAsOptions{})

	var configuration []ConfigurationModel
	data.Configuration.ElementsAs(ctx, &configuration, false)

	attrs := consoleClient.ServiceDeploymentAttributes{
		Name:         data.Name.ValueString(),
		Namespace:    data.Namespace.ValueString(),
		RepositoryID: repository.Id,
		Git: consoleClient.GitRefAttributes{
			Ref:    repository.Ref,
			Folder: repository.Folder,
		},
		Configuration: algorithms.Map(configuration, func(c ConfigurationModel) *consoleClient.ConfigAttributes {
			return &consoleClient.ConfigAttributes{
				Name:  c.Name.ValueString(),
				Value: c.Value.ValueStringPointer(),
			}
		}),
	}

	sd, err := r.client.CreateServiceDeployment(ctx, cluster.Id.ValueStringPointer(), cluster.Handle.ValueStringPointer(), attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create ServiceDeployment, got error: %s", err))
		return
	}

	// TODO: figure out what we need to read from response
	data.Id = types.StringValue(sd.ID)

	tflog.Trace(ctx, "created a ServiceDeployment")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceDeploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ServiceDeploymentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ServiceDeployment, err := r.client.GetServiceDeployment(ctx, data.Id.ValueString())
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
	//	Handle: lo.ToPtr(data.Handle.ValueString()),
	//}
	//ServiceDeployment, err := r.client.UpdateServiceDeployment(ctx, data.Id.ValueString(), attrs)
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

	_, err := r.client.DeleteServiceDeployment(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete ServiceDeployment, got error: %s", err))
		return
	}
}

func (r *ServiceDeploymentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
