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

var _ resource.Resource = &InfrastructureStackResource{}
var _ resource.ResourceWithImportState = &InfrastructureStackResource{}

func NewInfrastructureStackResource() resource.Resource {
	return &InfrastructureStackResource{}
}

// InfrastructureStackResource defines the infrastructure stack resource implementation.
type InfrastructureStackResource struct {
	client *client.Client
}

func (r *InfrastructureStackResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infrastructure_stack"
}

func (r *InfrastructureStackResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema()
}

func (r *InfrastructureStackResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Infrastructure Stack Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = data.Client
}

func (r *InfrastructureStackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := new(infrastructureStack)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sd, err := r.client.CreateStack(ctx, data.Attributes(ctx, resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create infrastructure stack, got error: %s", err))
		return
	}

	data.From(sd.CreateStack, ctx, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *InfrastructureStackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := new(infrastructureStack)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetInfrastructureStack(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read infrastructure stack, got error: %s", err))
		return
	}

	data.From(response.InfrastructureStack, ctx, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *InfrastructureStackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := new(infrastructureStack)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateStack(ctx, data.Id.ValueString(), data.Attributes(ctx, resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update infrastructure stack, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *InfrastructureStackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := new(infrastructureStack)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteStack(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete infrastructure stack, got error: %s", err))
		return
	}

	err = wait.WaitForWithContext(ctx, client.Ticker(5*time.Second), func(ctx context.Context) (bool, error) {
		_, err := r.client.GetInfrastructureStack(ctx, data.Id.ValueString())
		if client.IsNotFound(err) {
			return true, nil
		}

		return false, err
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during watiting for infrastructure stack to be deleted, got error: %s", err))
		return
	}
}

func (r *InfrastructureStackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
