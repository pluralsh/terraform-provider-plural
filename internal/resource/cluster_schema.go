package resource

import (
	"terraform-provider-plural/internal/common"

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
			// "version": schema.StringAttribute{
			// 	Description:         "Kubernetes version to use for this cluster. Leave empty for bring your own cluster. Supported version ranges can be found at https://github.com/pluralsh/console/tree/master/static/k8s-versions.",
			// 	MarkdownDescription: "Kubernetes version to use for this cluster. Leave empty for bring your own cluster. Supported version ranges can be found at https://github.com/pluralsh/console/tree/master/static/k8s-versions.",
			// 	Optional:            true,
			// 	Validators: []validator.String{
			// 		internalvalidator.ConflictsWithIf(internalvalidator.ConflictsIfTargetValueOneOf([]string{common.CloudBYOK.String()}),
			// 			path.MatchRoot("cloud")),
			// 	},
			// },
			// "desired_version": schema.StringAttribute{
			// 	Description:         "Desired Kubernetes version for this cluster.",
			// 	MarkdownDescription: "Desired Kubernetes version for this cluster.",
			// 	Computed:            true,
			// 	Optional:            true,
			// 	PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			// },
			// "provider_id": schema.StringAttribute{
			// 	Description:         "Provider used to create this cluster. Leave empty for bring your own cluster.",
			// 	MarkdownDescription: "Provider used to create this cluster. Leave empty for bring your own cluster.",
			// 	Optional:            true,
			// 	PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			// 	Validators: []validator.String{
			// 		internalvalidator.ConflictsWithIf(internalvalidator.ConflictsIfTargetValueOneOf([]string{common.CloudBYOK.String()}),
			// 			path.MatchRoot("cloud")),
			// 	},
			// },
			"metadata": schema.StringAttribute{
				Description:         "Arbitrary JSON metadata to store user-specific state of this cluster (e.g. IAM roles for add-ons). Use 'jsonencode' and 'jsondecode' methods to encode and decode data.",
				MarkdownDescription: "Arbitrary JSON metadata to store user-specific state of this cluster (e.g. IAM roles for add-ons). Use `jsonencode` and `jsondecode` methods to encode and decode data.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			// "cloud": schema.StringAttribute{
			// 	Description:         "The cloud provider used to create this cluster.",
			// 	MarkdownDescription: "The cloud provider used to create this cluster.",
			// 	Computed:            true,
			// 	Optional:            true,
			// 	Default:             stringdefault.StaticString(common.CloudBYOK.String()),
			// 	PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			// 	Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
			// 		common.CloudBYOK.String(), common.CloudAWS.String(), common.CloudAzure.String(), common.CloudGCP.String()),
			// 		internalvalidator.AlsoRequiresIf(internalvalidator.RequiresIfSourceValueOneOf([]string{
			// 			common.CloudAWS.String(),
			// 			common.CloudAzure.String(),
			// 			common.CloudGCP.String(),
			// 		}), path.MatchRoot("provider_id")),
			// 		internalvalidator.AlsoRequiresIf(internalvalidator.RequiresIfSourceValueOneOf([]string{
			// 			common.CloudAWS.String(),
			// 			common.CloudAzure.String(),
			// 			common.CloudGCP.String(),
			// 		}), path.MatchRoot("cloud_settings")),
			// 	},
			// },
			// "cloud_settings": schema.SingleNestedAttribute{
			// 	Description:         "Cloud-specific settings for this cluster.",
			// 	MarkdownDescription: "Cloud-specific settings for this cluster.",
			// 	Optional:            true,
			// 	Attributes: map[string]schema.Attribute{
			// 		"aws":   r.awsCloudSettingsSchema(),
			// 		"azure": r.azureCloudSettingsSchema(),
			// 		"gcp":   r.gcpCloudSettingsSchema(),
			// 		"byok":  r.byokCloudSettingsSchema(),
			// 	},
			// 	PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
			// },
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
			// "node_pools": schema.MapNestedAttribute{
			// 	Description:         "Experimental, not ready for production use. Map of node pool specs managed by this cluster, where the key is name of the node pool and value contains the spec. Leave empty for bring your own cluster.",
			// 	MarkdownDescription: "**Experimental, not ready for production use.** Map of node pool specs managed by this cluster, where the key is name of the node pool and value contains the spec. Leave empty for bring your own cluster.",
			// 	Optional:            true,
			// 	Computed:            true,
			// 	PlanModifiers:       []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"name": schema.StringAttribute{
			// 				Description:         "Node pool name. Must be unique.",
			// 				MarkdownDescription: "Node pool name. Must be unique.",
			// 				Required:            true,
			// 			},
			// 			"min_size": schema.Int64Attribute{
			// 				Description:         "Minimum number of instances in this node pool.",
			// 				MarkdownDescription: "Minimum number of instances in this node pool.",
			// 				Required:            true,
			// 			},
			// 			"max_size": schema.Int64Attribute{
			// 				Description:         "Maximum number of instances in this node pool.",
			// 				MarkdownDescription: "Maximum number of instances in this node pool.",
			// 				Required:            true,
			// 			},
			// 			"instance_type": schema.StringAttribute{
			// 				Description:         "The type of node to use. Usually cloud-specific.",
			// 				MarkdownDescription: "The type of node to use. Usually cloud-specific.",
			// 				Required:            true,
			// 			},
			// 			"labels": schema.MapAttribute{
			// 				Description:         "Kubernetes labels to apply to the nodes in this pool. Useful for node selectors.",
			// 				MarkdownDescription: "Kubernetes labels to apply to the nodes in this pool. Useful for node selectors.",
			// 				ElementType:         types.StringType,
			// 				Optional:            true,
			// 			},
			// 			"taints": schema.SetNestedAttribute{
			// 				Description:         "Any taints you'd want to apply to a node, i.e. for preventing scheduling on spot instances. See https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/ for more information.",
			// 				MarkdownDescription: "Any taints you'd want to apply to a node, i.e. for preventing scheduling on spot instances. See [Kubernetes docs](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/) for more information.",
			// 				Optional:            true,
			// 				NestedObject: schema.NestedAttributeObject{
			// 					Attributes: map[string]schema.Attribute{
			// 						"key": schema.StringAttribute{
			// 							Description:         "Taint key.",
			// 							MarkdownDescription: "Taint key.",
			// 							Required:            true,
			// 						},
			// 						"value": schema.StringAttribute{
			// 							Description:         "Taint value.",
			// 							MarkdownDescription: "Taint value.",
			// 							Required:            true,
			// 						},
			// 						"effect": schema.StringAttribute{
			// 							Description:         "Taint effect, allowed values include NoExecute, NoSchedule and PreferNoSchedule.",
			// 							MarkdownDescription: "Taint effect, allowed values include `NoExecute`, `NoSchedule` and `PreferNoSchedule`.",
			// 							Required:            true,
			// 						},
			// 					},
			// 				},
			// 			},
			// 			"cloud_settings": schema.SingleNestedAttribute{
			// 				Description:         "Cloud-specific settings for this node pool.",
			// 				MarkdownDescription: "Cloud-specific settings for this node pool.",
			// 				Optional:            true,
			// 				Attributes: map[string]schema.Attribute{
			// 					"aws": schema.SingleNestedAttribute{
			// 						Description:         "AWS node pool customizations.",
			// 						MarkdownDescription: "AWS node pool customizations.",
			// 						Optional:            true,
			// 						Attributes: map[string]schema.Attribute{
			// 							"launch_template_id": schema.StringAttribute{
			// 								Description:         "Custom launch template for your nodes. Useful for Golden AMI setups.",
			// 								MarkdownDescription: "Custom launch template for your nodes. Useful for Golden AMI setups.",
			// 								Optional:            true,
			// 							},
			// 						},
			// 					},
			// 				},
			// 			},
			// 		},
			// 	},
			// },
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
		},
	}
}

//func (r *clusterResource) awsCloudSettingsSchema() schema.SingleNestedAttribute {
//	return schema.SingleNestedAttribute{
//		Optional: true,
//		Attributes: map[string]schema.Attribute{
//			"region": schema.StringAttribute{
//				Description:         "AWS region to deploy the cluster to.",
//				MarkdownDescription: "AWS region to deploy the cluster to.",
//				Required:            true,
//			},
//		},
//	}
//}

//func (r *clusterResource) azureCloudSettingsSchema() schema.SingleNestedAttribute {
//	return schema.SingleNestedAttribute{
//		Optional: true,
//		Attributes: map[string]schema.Attribute{
//			"resource_group": schema.StringAttribute{
//				Description:         "Name of the Azure resource group for this cluster.",
//				MarkdownDescription: "Name of the Azure resource group for this cluster.",
//				Required:            true,
//			},
//			"network": schema.StringAttribute{
//				Description:         "Name of the Azure virtual network for this cluster.",
//				MarkdownDescription: "Name of the Azure virtual network for this cluster.",
//				Required:            true,
//			},
//			"subscription_id": schema.StringAttribute{
//				Description:         "GUID of the Azure subscription to hold this cluster.",
//				MarkdownDescription: "GUID of the Azure subscription to hold this cluster.",
//				Required:            true,
//			},
//			"location": schema.StringAttribute{
//				Description:         "String matching one of the canonical Azure region names, i.e. eastus.",
//				MarkdownDescription: "String matching one of the canonical Azure region names, i.e. eastus.",
//				Required:            true,
//			},
//		},
//	}
//}

//func (r *clusterResource) gcpCloudSettingsSchema() schema.SingleNestedAttribute {
//	return schema.SingleNestedAttribute{
//		Optional: true,
//		Attributes: map[string]schema.Attribute{
//			"project": schema.StringAttribute{
//				Required:            true,
//				Description:         "GCP project id to deploy cluster to.",
//				MarkdownDescription: "GCP project id to deploy cluster to.",
//			},
//			"network": schema.StringAttribute{
//				Required:            true,
//				Description:         "GCP network id to use when creating the cluster.",
//				MarkdownDescription: "GCP network id to use when creating the cluster.",
//			},
//			"region": schema.StringAttribute{
//				Required:            true,
//				Description:         "GCP region to deploy cluster to.",
//				MarkdownDescription: "GCP region to deploy cluster to.",
//			},
//		},
//	}
//}

//func (r *clusterResource) byokCloudSettingsSchema() schema.SingleNestedAttribute {
//	return schema.SingleNestedAttribute{
//		Optional: true,
//		Attributes: map[string]schema.Attribute{
//			"kubeconfig": r.kubeconfigSchema(true),
//		},
//	}
//}
