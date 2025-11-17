package resource

import (
	"context"
	"fmt"
	"strings"
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
	kubeClient *common.KubeClient
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
	r.kubeClient = data.KubeClient
}

func (r *clusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data cluster
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateCluster(ctx, data.Attributes(ctx, &resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create cluster, got error: %s", err))
		return
	}

	data.FromCreate(result, ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.kubeClient != nil || data.HasKubeconfig() {
		if result.CreateCluster.DeployToken == nil {
			resp.Diagnostics.AddError("Client Error", "Unable to fetch cluster deploy token")
			return
		}

		if err = InstallOrUpgradeAgent(ctx, r.client, data.GetKubeconfig(), r.kubeClient, data.HelmRepoUrl.ValueString(),
			data.HelmValues.ValueStringPointer(), r.consoleUrl, lo.FromPtr(result.CreateCluster.DeployToken), result.CreateCluster.ID, &resp.Diagnostics); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to install operator, got error: %s", err))
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data cluster
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Id.IsNull() {
		result, err := r.client.GetCluster(ctx, data.Id.ValueStringPointer())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cluster, got error: %s", err))
			return
		}
		if result == nil || result.Cluster == nil {
			// Resource not found, remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		data.From(result.Cluster, ctx, &resp.Diagnostics)
	} else if !data.Handle.IsNull() {
		result, err := r.client.GetClusterByHandle(ctx, data.Handle.ValueStringPointer())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cluster, got error: %s", err))
			return
		}
		if result == nil || result.Cluster == nil {
			// Resource not found, remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		data.From(result.Cluster, ctx, &resp.Diagnostics)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state cluster
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.ProjectId.Equal(state.ProjectId) && !data.ProjectId.IsNull() {
		resp.Diagnostics.AddError("Invalid Configuration", "Unable to update cluster, project ID must not be modified")
		return
	}

	result, err := r.client.UpdateCluster(ctx, data.Id.ValueString(), data.UpdateAttributes(ctx, &resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update cluster, got error: %s", err))
		return
	}

	kubeconfigChanged := data.HasKubeconfig() && !data.GetKubeconfig().Unchanged(state.GetKubeconfig())
	reinstallable := !data.HelmRepoUrl.Equal(state.HelmRepoUrl) || kubeconfigChanged
	if reinstallable && (r.kubeClient != nil || data.HasKubeconfig()) {
		clusterWithToken, err := r.client.GetClusterWithToken(ctx, data.Id.ValueStringPointer(), nil)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to fetch cluster deploy token, got error: %s", err))
			return
		}

		if err = InstallOrUpgradeAgent(ctx, r.client, data.GetKubeconfig(), r.kubeClient, data.HelmRepoUrl.ValueString(),
			data.HelmValues.ValueStringPointer(), r.consoleUrl, lo.FromPtr(clusterWithToken.Cluster.DeployToken), result.UpdateCluster.ID, &resp.Diagnostics); err != nil {
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
		if _, err := r.client.DetachCluster(ctx, data.Id.ValueString()); err != nil && !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to detach cluster, got error: %s", err))
			return
		}
	} else {
		if _, err := r.client.DeleteCluster(ctx, data.Id.ValueString()); err != nil && !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete cluster, got error: %s", err))
			return
		}

		if err := wait.PollUntilContextTimeout(ctx, 10*time.Second, 10*time.Minute, true, func(ctx context.Context) (bool, error) {
			response, err := r.client.GetCluster(ctx, data.Id.ValueStringPointer())
			if client.IsNotFound(err) {
				return true, nil
			}

			if err == nil && (response == nil || response.Cluster == nil) {
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
	if strings.HasPrefix(req.ID, "@") && len(req.ID) > 1 {
		req.ID = req.ID[1:]
		resource.ImportStatePassthroughID(ctx, path.Root("handle"), req, resp)
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
