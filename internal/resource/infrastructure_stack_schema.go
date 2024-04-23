package resource

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console-client-go"
)

func (r *InfrastructureStackResource) schema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Internal identifier of this stack.",
				MarkdownDescription: "Internal identifier of this stack.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description:         "Human-readable name of this stack.",
				MarkdownDescription: "Human-readable name of this stack.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"type": schema.StringAttribute{
				Description:         fmt.Sprintf("A type for the stack, specifies the tool to use to apply it. Allowed values include \"%s\" and \"%s\".", gqlclient.StackTypeAnsible, gqlclient.StackTypeTerraform),
				MarkdownDescription: fmt.Sprintf("A type for the stack, specifies the tool to use to apply it. Allowed values include `%s` and `%s`.", gqlclient.StackTypeAnsible, gqlclient.StackTypeTerraform),
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:          []validator.String{stringvalidator.OneOf(string(gqlclient.StackTypeAnsible), string(gqlclient.StackTypeTerraform))},
			},
			"approval": schema.BoolAttribute{
				Description:         "Determines whether to require approval.",
				MarkdownDescription: "Determines whether to require approval.",
				Optional:            true,
			},
			"detach": schema.BoolAttribute{
				Description:         "Determines behavior during resource destruction, if true it will detach resource instead of deleting it.",
				MarkdownDescription: "Determines behavior during resource destruction, if true it will detach resource instead of deleting it.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"cluster_id": schema.StringAttribute{
				Description:         "The cluster on which the stack will be applied.",
				MarkdownDescription: "The cluster on which the stack will be applied.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"repository": schema.SingleNestedAttribute{
				Description:         "Repository information used to pull stack.",
				MarkdownDescription: "Repository information used to pull stack.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         "ID of the repository to pull from.",
						MarkdownDescription: "ID of the repository to pull from.",
						Required:            true,
					},
					"ref": schema.StringAttribute{
						Description:         "A general git ref, either a branch name or commit sha understandable by `git checkout <ref>`.",
						MarkdownDescription: "A general git ref, either a branch name or commit sha understandable by \"git checkout <ref>\".",
						Required:            true,
					},
					"folder": schema.StringAttribute{
						Description:         "The folder where manifests live.",
						MarkdownDescription: "The folder where manifests live.",
						Required:            true,
					},
				},
			},
			"configuration": schema.SingleNestedAttribute{
				Description:         "Stack configuration.",
				MarkdownDescription: "Stack configuration.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"image": schema.StringAttribute{
						Description:         "Optional custom image you might want to use.",
						MarkdownDescription: "Optional custom image you might want to use.",
						Optional:            true,
					},
					"version": schema.StringAttribute{
						Description:         "The semver of the tool you wish to use.",
						MarkdownDescription: "The semver of the tool you wish to use.",
						Required:            true,
					},
				},
			},
			"files": schema.MapAttribute{
				MarkdownDescription: "File path-content map.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"environment": schema.SetNestedAttribute{
				Description:         "Defines environment variables for the stack.",
				MarkdownDescription: "Defines environment variables for the stack.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "Environment variable name.",
							MarkdownDescription: "Environment variable name.",
							Required:            true,
						},
						"value": schema.StringAttribute{
							Description:         "Environment variable value.",
							MarkdownDescription: "Environment variable value.",
							Required:            true,
						},
						"secret": schema.BoolAttribute{
							Description:         "Indicates if environment variable is secret.",
							MarkdownDescription: "Indicates if environment variable is secret.",
							Optional:            true,
						},
					},
				},
			},
			"job_spec": schema.SingleNestedAttribute{
				Description:         "Repository information used to pull stack.",
				MarkdownDescription: "Repository information used to pull stack.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"namespace": schema.StringAttribute{
						Description:         "Namespace where job will be deployed.",
						MarkdownDescription: "Namespace where job will be deployed.",
						Required:            true,
					},
					"raw": schema.StringAttribute{
						Description:         "If you'd rather define the job spec via straight Kubernetes YAML.",
						MarkdownDescription: "If you'd rather define the job spec via straight Kubernetes YAML.",
						Optional:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Validators: []validator.String{
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("containers")),
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("containers")),
						},
					},
					"labels": schema.MapAttribute{
						Description:         "Kubernetes labels applied to the job.",
						MarkdownDescription: "Kubernetes labels applied to the job.",
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
						Validators:          []validator.Map{mapvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("raw"))},
					},
					"annotations": schema.MapAttribute{
						Description:         "Kubernetes annotations applied to the job.",
						MarkdownDescription: "Kubernetes annotations applied to the job.",
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
						Validators:          []validator.Map{mapvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("raw"))},
					},
					"service_account": schema.StringAttribute{
						Description:         "Kubernetes service account for this job.",
						MarkdownDescription: "Kubernetes service account for this job.",
						Optional:            true,
						Validators:          []validator.String{stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("raw"))},
					},
					"containers": schema.SetNestedAttribute{
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"image": schema.StringAttribute{
									Required: true,
								},
								"args": schema.ListAttribute{
									Description:         "Arguments to pass to the command when executing it.",
									MarkdownDescription: "Arguments to pass to the command when executing it.",
									Optional:            true,
									ElementType:         types.StringType,
								},
								"env": schema.MapAttribute{
									Description:         "Defines environment variables to expose to the process.",
									MarkdownDescription: "Defines environment variables to expose to the process.",
									Optional:            true,
									ElementType:         types.StringType,
								},
								"env_from": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"secret": schema.StringAttribute{
												Required: true,
											},
											"config_map": schema.StringAttribute{
												Required: true,
											},
										},
									},
								},
							},
						},
						Validators: []validator.Set{
							setvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("raw")),
							setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("raw")),
						},
					},
				},
			},
			"bindings": schema.SingleNestedAttribute{
				Description:         "Read and write policies of this stack.",
				MarkdownDescription: "Read and write policies of this stack.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"read": schema.SetNestedAttribute{
						Description:         "Read policies of this stack.",
						MarkdownDescription: "Read policies of this stack.",
						Optional:            true,
						Computed:            true,
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
						Description:         "Write policies of this stack.",
						MarkdownDescription: "Write policies of this stack.",
						Optional:            true,
						Computed:            true,
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
