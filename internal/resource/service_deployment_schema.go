package resource

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *ServiceDeploymentResource) schema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Internal identifier of this ServiceDeployment.",
				MarkdownDescription: "Internal identifier of this ServiceDeployment.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Human-readable name of this ServiceDeployment.",
				MarkdownDescription: "Human-readable name of this ServiceDeployment.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"namespace": schema.StringAttribute{
				Required:            true,
				Description:         "Namespace to deploy this ServiceDeployment.",
				MarkdownDescription: "Namespace to deploy this ServiceDeployment.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Semver version of this service ServiceDeployment.",
				MarkdownDescription: "Semver version of this service ServiceDeployment.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"docs_path": schema.StringAttribute{
				Optional:            true,
				Description:         "Path to the documentation in the target git repository.",
				MarkdownDescription: "Path to the documentation in the target git repository.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protect": schema.BoolAttribute{
				Optional:            true,
				Description:         "If true, deletion of this service is not allowed.",
				MarkdownDescription: "If true, deletion of this service is not allowed.",
			},
			"kustomize":     r.schemaKustomize(),
			"configuration": r.schemaConfiguration(),
			"cluster":       r.schemaCluster(),
			"repository":    r.schemaRepository(),
			"bindings":      r.schemaBindings(),
			"sync_config":   r.schemaSyncConfig(),
			"helm":          r.schemaHelm(),
		},
	}
}

func (r *ServiceDeploymentResource) schemaKustomize() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Description:         "Kustomize related service metadata.",
		MarkdownDescription: "Kustomize related service metadata.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Path to the kustomize file in the target git repository.",
			},
		},
	}
}

func (r *ServiceDeploymentResource) schemaConfiguration() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Optional:            true,
		Description:         "List of [name, value] secrets used to alter this ServiceDeployment configuration.",
		MarkdownDescription: "List of [name, value] secrets used to alter this ServiceDeployment configuration.",
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
	}
}

func (r *ServiceDeploymentResource) schemaCluster() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Required:            true,
		Description:         "Unique cluster id/handle to deploy this ServiceDeployment",
		MarkdownDescription: "Unique cluster id/handle to deploy this ServiceDeployment",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("cluster").AtName("handle")),
					stringvalidator.ExactlyOneOf(path.MatchRoot("cluster").AtName("handle")),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Optional: true,
			},
			"handle": schema.StringAttribute{
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("cluster").AtName("id")),
					stringvalidator.ExactlyOneOf(path.MatchRoot("cluster").AtName("id")),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Optional: true,
			},
		},
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
	}
}

func (r *ServiceDeploymentResource) schemaRepository() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Description:         "Repository information used to pull ServiceDeployment.",
		MarkdownDescription: "Repository information used to pull ServiceDeployment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Optional: true,
			},
			"ref": schema.StringAttribute{
				Optional: true,
			},
			"folder": schema.StringAttribute{
				Optional: true,
			},
		},
		Validators: []validator.Object{
			objectvalidator.AtLeastOneOf(path.MatchRoot("helm")),
		},
	}
}

func (r *ServiceDeploymentResource) schemaBindings() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Description:         "Read and write policies of this ServiceDeployment.",
		MarkdownDescription: "Read and write policies of this ServiceDeployment.",
		Attributes: map[string]schema.Attribute{
			"read": schema.SetNestedAttribute{
				Optional:            true,
				Description:         "Read policies of this ServiceDeployment.",
				MarkdownDescription: "Read policies of this ServiceDeployment.",
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
				Description:         "Write policies of this ServiceDeployment.",
				MarkdownDescription: "Write policies of this ServiceDeployment.",
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
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
	}
}

func (r *ServiceDeploymentResource) schemaSyncConfig() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Description:         "Settings for advanced tuning of the sync process.",
		MarkdownDescription: "Settings for advanced tuning of the sync process.",
		Attributes: map[string]schema.Attribute{
			"namespace_metadata": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"annotations": schema.MapAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"labels": schema.MapAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
				},
			},
		},
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
	}
}

func (r *ServiceDeploymentResource) schemaHelm() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Description:         "Settings defining how Helm charts should be applied.",
		MarkdownDescription: "Settings defining how Helm charts should be applied.",
		Attributes: map[string]schema.Attribute{
			"chart": schema.StringAttribute{
				Optional:            true,
				Description:         "The name of the chart to use.",
				MarkdownDescription: "The name of the chart to use.",
			},
			"repository": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "Resource reference to the flux Helm repository used by this chart.",
				MarkdownDescription: "Resource reference to the flux Helm repository used by this chart.",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Optional:            true,
						Description:         "Name of the flux Helm repository resource used by this chart.",
						MarkdownDescription: "Name of the flux Helm repository resource used by this chart.",
					},
					"namespace": schema.StringAttribute{
						Optional:            true,
						Description:         "Namespace of the flux Helm repository resource used by this chart.",
						MarkdownDescription: "Namespace of the flux Helm repository resource used by this chart.",
					},
				},
			},
			"values": schema.StringAttribute{
				Optional:            true,
				Description:         "Helm values file to use with this service",
				MarkdownDescription: "Helm values file to use with this service",
			},
			"values_files": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Description:         "List of relative paths to values files to use form Helm applies.",
				MarkdownDescription: "List of relative paths to values files to use form Helm applies.",
			},
			"version": schema.StringAttribute{
				Optional:            true,
				Description:         "Chart version to use",
				MarkdownDescription: "Chart version to use",
			},
		},
		Validators: []validator.Object{
			objectvalidator.AtLeastOneOf(path.MatchRoot("repository")),
		},
	}
}
