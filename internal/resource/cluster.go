package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pluralsh/plural-cli/pkg/console"
	"github.com/samber/lo"
	"k8s.io/apimachinery/pkg/util/wait"
)

var _ resource.Resource = &clusterResource{}
var _ resource.ResourceWithImportState = &clusterResource{}
var _ resource.ResourceWithUpgradeState = &clusterResource{}

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
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("kubeconfig"), &data.Kubeconfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateCluster(ctx, data.Attributes(ctx, &resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create cluster, got error: %s", err))
		return
	}
	data.FromCreate(result, ctx, &resp.Diagnostics)

	if r.kubeClient != nil || data.HasKubeconfig() {
		err = InstallOrUpgradeAgent(ctx, r.client, data.GetKubeconfig(), r.kubeClient, data.HelmRepoUrl.ValueString(),
			data.HelmValues.ValueStringPointer(), r.consoleUrl, lo.FromPtr(result.CreateCluster.DeployToken),
			result.CreateCluster.ID, &resp.Diagnostics)
		if err != nil {
			resp.Diagnostics.AddWarning("Agent Installation Failed", fmt.Sprintf(
				"Unable to install agent, in order to retry run `terraform apply` again. Got error: %s", err))
		} else {
			data.AgentDeployed = types.BoolValue(true)
		}
	}

	resp.Diagnostics.Append(setKubeconfigHost(ctx, resp.Private, data.GetKubeconfig())...)
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
		if err != nil && !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cluster, got error: %s", err))
			return
		}
		if result == nil || result.Cluster == nil || client.IsNotFound(err) {
			// Resource not found, remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		data.From(result.Cluster, ctx, &resp.Diagnostics)
	} else if !data.Handle.IsNull() {
		result, err := r.client.GetClusterByHandle(ctx, data.Handle.ValueStringPointer())
		if err != nil && !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cluster, got error: %s", err))
			return
		}
		if result == nil || result.Cluster == nil || client.IsNotFound(err) {
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
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("kubeconfig"), &data.Kubeconfig)...)
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

	// The plan may leave agent_deployed unknown even when state is true; use
	// state so the value written back after apply is always known.
	if data.AgentDeployed.IsNull() || data.AgentDeployed.IsUnknown() {
		data.AgentDeployed = state.AgentDeployed
	}
	if data.AgentDeployed.IsNull() || data.AgentDeployed.IsUnknown() {
		data.AgentDeployed = types.BoolValue(false)
	}

	kubeconfigChanged, diags := kubeconfigHostChanged(ctx, req.Private, data.GetKubeconfig())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reinstallable := !data.AgentDeployed.ValueBool() || !data.HelmRepoUrl.Equal(state.HelmRepoUrl) || kubeconfigChanged
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

		data.AgentDeployed = types.BoolValue(true)
	}

	resp.Diagnostics.Append(setKubeconfigHost(ctx, resp.Private, data.GetKubeconfig())...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// kubeconfigHostPrivateKey is the private state key used to track the kubeconfig host. Kubeconfig is write-only,
// so it is not available in the state, and it is required to detect kubeconfig changes that require agent reinstall.
const kubeconfigHostPrivateKey = "kubeconfig_host"

type privateStateGetter interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
}

type privateStateSetter interface {
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

// kubeconfigHostChanged checks if the kubeconfig host has changed compared to the one tracked in the private state.
// If the host was not tracked yet, i.e. state was upgraded from older provider version, it is treated as changed.
// State upgrades cannot write private state, so the previous host is lost, and it is not possible to tell if it has
// changed. Reinstalling the agent is idempotent, while skipping it could leave a new cluster without the agent.
func kubeconfigHostChanged(ctx context.Context, private privateStateGetter, kubeconfig *common.Kubeconfig) (bool, diag.Diagnostics) {
	if kubeconfig == nil {
		return false, nil
	}

	value, diags := private.GetKey(ctx, kubeconfigHostPrivateKey)
	if diags.HasError() {
		return false, diags
	}

	if value == nil {
		return true, diags
	}

	var host *string
	if err := json.Unmarshal(value, &host); err != nil {
		diags.AddError("Provider Error", fmt.Sprintf("Cannot unmarshal kubeconfig host from private state, got error: %s", err))
		return false, diags
	}

	return host == nil || *host != kubeconfig.Host.ValueString(), diags
}

// setKubeconfigHost stores the kubeconfig host in the private state. If kubeconfig is not set, JSON null is stored.
func setKubeconfigHost(ctx context.Context, private privateStateSetter, kubeconfig *common.Kubeconfig) diag.Diagnostics {
	var host *string
	if kubeconfig != nil {
		host = lo.ToPtr(kubeconfig.Host.ValueString())
	}

	value, err := json.Marshal(host)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Provider Error", fmt.Sprintf("Cannot marshal kubeconfig host to private state, got error: %s", err))
		return diags
	}

	return private.SetKey(ctx, kubeconfigHostPrivateKey, value)
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

func (r *clusterResource) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	// Version 1 schema is the same as the current one, except kubeconfig was not write-only.
	priorSchemaV1 := r.schema()
	priorSchemaV1.Version = 1
	priorSchemaV1.Attributes["kubeconfig"] = common.KubeconfigResourceSchema(false)

	return map[int64]resource.StateUpgrader{
		// State upgrade from 1 to 2. Kubeconfig became write-only, so it has to be removed from the state.
		// It cannot be done in the same schema version, as the framework passes such state through as is.
		1: {
			PriorSchema: &priorSchemaV1,
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var data cluster
				resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
				if resp.Diagnostics.HasError() {
					return
				}

				data.Kubeconfig = nil
				resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
			},
		},
		// State upgrade from 0 to 2
		0: {
			PriorSchema: &schema.Schema{
				Version: 0,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:      true,
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
					},
					"inserted_at": schema.StringAttribute{
						Computed:      true,
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
					},
					"name": schema.StringAttribute{
						Required:      true,
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
					},
					"handle": schema.StringAttribute{
						Optional:      true,
						Computed:      true,
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
					},
					"project_id": schema.StringAttribute{
						Optional: true,
					},
					"detach": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"metadata": schema.StringAttribute{
						Optional:      true,
						Computed:      true,
						Default:       stringdefault.StaticString("{}"),
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
					},
					"helm_repo_url": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(console.RepoUrl),
					},
					"helm_values": schema.StringAttribute{
						Optional: true,
					},
					"kubeconfig": common.KubeconfigResourceSchema(false),
					"protect": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"tags": schema.MapAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
					"bindings": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"read": schema.SetNestedAttribute{
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"group_id": schema.StringAttribute{Optional: true},
										"id":       schema.StringAttribute{Optional: true},
										"user_id":  schema.StringAttribute{Optional: true},
									},
								},
							},
							"write": schema.SetNestedAttribute{
								Optional:            true,
								Description:         "Write policies of this cluster.",
								MarkdownDescription: "Write policies of this cluster.",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"group_id": schema.StringAttribute{
											Optional: true,
										},
										"id": schema.StringAttribute{
											Optional: true,
										},
										"user_id": schema.StringAttribute{
											Optional: true,
										},
									},
								},
							},
						},
						PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var priorStateData struct {
					Id          types.String       `tfsdk:"id"`
					InsertedAt  types.String       `tfsdk:"inserted_at"`
					Name        types.String       `tfsdk:"name"`
					Handle      types.String       `tfsdk:"handle"`
					ProjectId   types.String       `tfsdk:"project_id"`
					Detach      types.Bool         `tfsdk:"detach"`
					Protect     types.Bool         `tfsdk:"protect"`
					Tags        types.Map          `tfsdk:"tags"`
					Metadata    types.String       `tfsdk:"metadata"`
					Bindings    *common.Bindings   `tfsdk:"bindings"`
					HelmRepoUrl types.String       `tfsdk:"helm_repo_url"`
					HelmValues  types.String       `tfsdk:"helm_values"`
					Kubeconfig  *common.Kubeconfig `tfsdk:"kubeconfig"`
				}

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)
				if resp.Diagnostics.HasError() {
					return
				}

				// Kubeconfig is write-only now, so it is not copied to the upgraded state.
				upgradedStateData := cluster{
					Id:            priorStateData.Id,
					InsertedAt:    priorStateData.InsertedAt,
					Name:          priorStateData.Name,
					Handle:        priorStateData.Handle,
					ProjectId:     priorStateData.ProjectId,
					Detach:        priorStateData.Detach,
					Protect:       priorStateData.Protect,
					Tags:          priorStateData.Tags,
					Metadata:      priorStateData.Metadata,
					Bindings:      priorStateData.Bindings,
					HelmRepoUrl:   priorStateData.HelmRepoUrl,
					HelmValues:    priorStateData.HelmValues,
					AgentDeployed: types.BoolValue(true),
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}
}
