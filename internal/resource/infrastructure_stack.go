package resource

import (
	"context"
	"fmt"
	"time"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"
	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"k8s.io/apimachinery/pkg/util/wait"
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
	data := new(model.InfrastructureStackExtended)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attr, err := data.Attributes(ctx, &resp.Diagnostics, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get attributes, got error: %s", err))
		return
	}
	sd, err := r.client.CreateStack(ctx, *attr)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create infrastructure stack, got error: %s", err))
		return
	}

	data.From(sd.CreateStack, ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *InfrastructureStackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := new(model.InfrastructureStackExtended)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetInfrastructureStack(ctx, data.Id.ValueStringPointer(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read infrastructure stack, got error: %s", err))
		return
	}

	data.From(response.InfrastructureStack, ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *InfrastructureStackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := new(model.InfrastructureStackExtended)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attr, err := data.Attributes(ctx, &resp.Diagnostics, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get attributes, got error: %s", err))
		return
	}
	_, err = r.client.UpdateStack(ctx, data.Id.ValueString(), *attr)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update infrastructure stack, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *InfrastructureStackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := new(model.InfrastructureStackExtended)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Detach.ValueBool() {
		_, err := r.client.DetachStack(ctx, data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to detach infrastructure stack, got error: %s", err))
			return
		}
	} else {
		_, err := r.client.DeleteStack(ctx, data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete infrastructure stack, got error: %s", err))
			return
		}

		if err = wait.PollUntilContextTimeout(ctx, 10*time.Second, 10*time.Minute, true, func(ctx context.Context) (bool, error) {
			response, err := r.client.GetInfrastructureStack(ctx, data.Id.ValueStringPointer(), nil)
			if client.IsNotFound(err) {
				return true, nil
			}

			if err == nil && (response == nil || response.InfrastructureStack == nil) {
				return true, nil
			}

			return false, err
		}); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during watiting for infrastructure stack to be deleted, got error: %s", err))
			return
		}
	}
}

func (r *InfrastructureStackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
