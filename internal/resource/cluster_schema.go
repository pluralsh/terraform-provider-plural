package resource

import (
	"terraform-provider-plural/internal/common"
	resource "terraform-provider-plural/internal/planmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/pluralsh/plural-cli/pkg/console"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *clusterResource) schema() schema.Schema {
	return schema.Schema{
		Description:         "A representation of a cluster you can deploy to.",
		MarkdownDescription: "A representation of a cluster you can deploy to.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Internal identifier of this cluster.",
				MarkdownDescription: "Internal identifier of this cluster.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"inserted_at": schema.StringAttribute{
				Description:         "Creation date of this cluster.",
				MarkdownDescription: "Creation date of this cluster.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description:         "Human-readable name of this cluster, that also translates to cloud resource name.",
				MarkdownDescription: "Human-readable name of this cluster, that also translates to cloud resource name.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"handle": schema.StringAttribute{
				Description:         "A short, unique human-readable name used to identify this cluster. Does not necessarily map to the cloud resource name.",
				MarkdownDescription: "A short, unique human-readable name used to identify this cluster. Does not necessarily map to the cloud resource name.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"project_id": schema.StringAttribute{
				Description:         "ID of the project that this cluster belongs to.",
				MarkdownDescription: "ID of the project that this cluster belongs to.",
				Optional:            true,
			},
			"detach": schema.BoolAttribute{
				Description:         "Determines behavior during resource destruction, if true it will detach resource instead of deleting it.",
				MarkdownDescription: "Determines behavior during resource destruction, if true it will detach resource instead of deleting it.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"metadata": schema.StringAttribute{
				Description:         "Arbitrary JSON metadata to store user-specific state of this cluster (e.g. IAM roles for add-ons). Use 'jsonencode' and 'jsondecode' methods to encode and decode data.",
				MarkdownDescription: "Arbitrary JSON metadata to store user-specific state of this cluster (e.g. IAM roles for add-ons). Use `jsonencode` and `jsondecode` methods to encode and decode data.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("{}"),
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"helm_repo_url": schema.StringAttribute{
				Description:         "Helm repository URL you'd like to use in deployment agent Helm install.",
				MarkdownDescription: "Helm repository URL you'd like to use in deployment agent Helm install.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(console.RepoUrl),
			},
			"helm_values": schema.StringAttribute{
				Description:         "Additional Helm values you'd like to use in deployment agent Helm installs. This is useful for BYOK clusters that need to use custom images or other constructs.",
				MarkdownDescription: "Additional Helm values you'd like to use in deployment agent Helm installs. This is useful for BYOK clusters that need to use custom images or other constructs.",
				Optional:            true,
			},
			"kubeconfig": common.KubeconfigResourceSchema(),
			"protect": schema.BoolAttribute{
				Description:         "If set to \"true\" then this cluster cannot be deleted.",
				MarkdownDescription: "If set to `true` then this cluster cannot be deleted.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"tags": schema.MapAttribute{
				Description:         "Key-value tags used to filter clusters.",
				MarkdownDescription: "Key-value tags used to filter clusters.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"bindings": schema.SingleNestedAttribute{
				Description:         "Read and write policies of this cluster.",
				MarkdownDescription: "Read and write policies of this cluster.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"read": schema.SetNestedAttribute{
						Optional:            true,
						Description:         "Read policies of this cluster.",
						MarkdownDescription: "Read policies of this cluster.",
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
			"created": schema.BoolAttribute{
				Description:         "Whether the cluster was created in the Console API.",
				MarkdownDescription: "Whether the cluster was created in the Console API.",
				Computed:            true,
			},
			"agent_deployed": schema.BoolAttribute{
				Description:         "Whether the agent was deployed to the cluster.",
				MarkdownDescription: "Whether the agent was deployed to the cluster.",
				Computed:            true,
			},
			"reapply_key": schema.Int32Attribute{
				Description:         "Reapply key used to trigger reinstallation of the agent.",
				MarkdownDescription: "Reapply key used to trigger reinstallation of the agent.",
				Computed:            true,
				PlanModifiers:       []planmodifier.Int32{resource.EnsureAgent()},
			},
		},
	}
}
