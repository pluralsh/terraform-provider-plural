package resource

import (
	"context"
	"fmt"
	"time"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"k8s.io/apimachinery/pkg/util/wait"

	"terraform-provider-plural/internal/client"
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

func (r *ServiceDeploymentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_deployment"
}

func (r *ServiceDeploymentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema()
}

func (r *ServiceDeploymentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ServiceDeployment Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = data.Client
}

func (r *ServiceDeploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := new(ServiceDeployment)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs := data.Attributes(ctx, resp.Diagnostics)
	sd, err := r.client.CreateServiceDeployment(ctx, data.Cluster.Id.ValueStringPointer(), data.Cluster.Handle.ValueStringPointer(), attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create ServiceDeployment, got error: %s", err))
		return
	}

	data.FromCreate(sd, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *ServiceDeploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := new(ServiceDeployment)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetServiceDeployment(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ServiceDeployment, got error: %s", err))
		return
	}

	data.FromGet(response.ServiceDeployment, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *ServiceDeploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := new(ServiceDeployment)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateServiceDeployment(ctx, data.Id.ValueString(), data.UpdateAttributes(ctx, resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update ServiceDeployment, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *ServiceDeploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := new(ServiceDeployment)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteServiceDeployment(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete ServiceDeployment, got error: %s", err))
		return
	}

	err = wait.PollUntilContextCancel(ctx, 5*time.Second, true, func(ctx context.Context) (bool, error) {
		_, err := r.client.GetServiceDeployment(ctx, data.Id.ValueString())
		if client.IsNotFound(err) {
			return true, nil
		}

		return false, err
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during watiting for ServiceDeployment to be deleted, got error: %s", err))
		return
	}
}

func (r *ServiceDeploymentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
