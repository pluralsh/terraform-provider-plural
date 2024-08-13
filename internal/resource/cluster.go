package resource

import (
	"context"
	"fmt"
	"time"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/samber/lo"
	"k8s.io/apimachinery/pkg/util/wait"
)

var _ resource.Resource = &clusterResource{}
var _ resource.ResourceWithImportState = &clusterResource{}

func NewClusterResource() resource.Resource {
	return &clusterResource{}
}

// ClusterResource defines the cluster resource implementation.
type clusterResource struct {
	client     *client.Client
	consoleUrl string
}

func (r *clusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (r *clusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema()
}

func (r *clusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Cluster Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = data.Client
	r.consoleUrl = data.ConsoleUrl
}

func (r *clusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data cluster
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateCluster(ctx, data.Attributes(ctx, resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create cluster, got error: %s", err))
		return
	}

	if data.HasKubeconfig() {
		if result.CreateCluster.DeployToken == nil {
			resp.Diagnostics.AddError("Client Error", "Unable to fetch cluster deploy token")
			return
		}

		handler, err := NewOperatorHandler(ctx, r.client, data.GetKubeconfig(), data.HelmRepoUrl.ValueString(), data.HelmValues.ValueStringPointer(), r.consoleUrl)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to init operator handler, got error: %s", err))
			return
		}

		err = handler.InstallOrUpgrade(*result.CreateCluster.DeployToken)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to install operator, got error: %s", err))
			return
		}
	}

	data.FromCreate(result, ctx, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data cluster
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetCluster(ctx, data.Id.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cluster, got error: %s", err))
		return
	}
	if result == nil || result.Cluster == nil {
		resp.Diagnostics.AddError("Not Found", "Unable to find cluster, it looks like it was deleted manually")
		return
	}

	data.From(result.Cluster, ctx, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state cluster
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.ProjectId.Equal(state.ProjectId) {
		resp.Diagnostics.AddError("Invalid Configuration", "Unable to update cluster, project ID must not be modified")
		return
	}

	_, err := r.client.UpdateCluster(ctx, data.Id.ValueString(), data.UpdateAttributes(ctx, resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update cluster, got error: %s", err))
		return
	}

	if !data.HelmRepoUrl.Equal(state.HelmRepoUrl) && data.HasKubeconfig() {
		clusterWithToken, err := r.client.GetClusterWithToken(ctx, data.Id.ValueStringPointer(), nil)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to fetch cluster deploy token, got error: %s", err))
			return
		}

		handler, err := NewOperatorHandler(ctx, r.client, data.GetKubeconfig(), data.HelmRepoUrl.ValueString(), data.HelmValues.ValueStringPointer(), r.consoleUrl)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to init operator handler, got error: %s", err))
			return
		}

		err = handler.InstallOrUpgrade(lo.FromPtr(clusterWithToken.Cluster.DeployToken))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to install operator, got error: %s", err))
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data cluster
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Detach.ValueBool() {
		_, err := r.client.DetachCluster(ctx, data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to detach cluster, got error: %s", err))
			return
		}
	} else {
		_, err := r.client.DeleteCluster(ctx, data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete cluster, got error: %s", err))
			return
		}

		if err = wait.PollUntilContextTimeout(ctx, 10*time.Second, 10*time.Minute, true, func(ctx context.Context) (bool, error) {
			response, err := r.client.GetCluster(ctx, data.Id.ValueStringPointer())
			if client.IsNotFound(err) || response == nil || response.Cluster == nil {
				return true, nil
			}

			return false, err
		}); err != nil {
			resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Error while watiting for cluster to be deleted, got error: %s", err))

			_, err = r.client.DetachCluster(ctx, data.Id.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to detach cluster, got error: %s", err))
				return
			}
		}
	}
}

func (r *clusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
