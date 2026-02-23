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
	gqlclient "github.com/pluralsh/console/go/client"
)

func (r *InfrastructureStackResource) schema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Infrastructure stack provides a scalable framework to manage infrastructure as code with a K8s-friendly, API-driven approach. It declaratively defines a stack with a type, Git repository location, and target cluster for execution. On each commit to the tracked repository, a run is created which the Plural deployment operator detects and executes on the targeted cluster, enabling fine-grained permissions and network location control for IaC runs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Internal identifier of this stack.",
				MarkdownDescription: "Internal identifier of this stack.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description:         "Name of this stack.",
				MarkdownDescription: "Name of this stack.",
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
			"actor": schema.StringAttribute{
				Description:         "The User email to use for default Plural authentication in this stack.",
				MarkdownDescription: "The User email to use for default Plural authentication in this stack.",
				Optional:            true,
			},
			"project_id": schema.StringAttribute{
				Description:         "ID of the project that this stack belongs to.",
				MarkdownDescription: "ID of the project that this stack belongs to.",
				Computed:            true,
				Optional:            true,
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
					"tag": schema.StringAttribute{
						Description:         "The docker image tag you wish to use if you're customizing the version.",
						MarkdownDescription: "The docker image tag you wish to use if you're customizing the version.",
						Optional:            true,
					},
					"hooks": schema.SetNestedAttribute{
						Description:         "The hooks to customize execution for this stack.",
						MarkdownDescription: "The hooks to customize execution for this stack.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"cmd": schema.StringAttribute{
									Required: true,
								},
								"args": schema.ListAttribute{
									Description:         "Arguments to pass to the command when executing it.",
									MarkdownDescription: "Arguments to pass to the command when executing it.",
									Optional:            true,
									ElementType:         types.StringType,
								},
								"after_stage": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
					"terraform": schema.SingleNestedAttribute{
						Description:         "The terraform configuration for this stack.",
						MarkdownDescription: "The terraform configuration for this stack.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"parallelism": schema.Int64Attribute{
								Description:         "Equivalent to the -parallelism flag in terraform.",
								MarkdownDescription: "Equivalent to the -parallelism flag in terraform.",
								Optional:            true,
							},
							"refresh": schema.BoolAttribute{
								Description:         "Equivalent to the -refresh flag in terraform.",
								MarkdownDescription: "Equivalent to the -refresh flag in terraform.",
								Optional:            true,
							},
							"approve_empty": schema.BoolAttribute{
								Description:         "Whether to auto-approve a plan if there are no changes, preventing a stack from being blocked.",
								MarkdownDescription: "Whether to auto-approve a plan if there are no changes, preventing a stack from being blocked.",
								Optional:            true,
							},
						},
					},
					"ansible": schema.SingleNestedAttribute{
						Description:         "The ansible configuration for this stack.",
						MarkdownDescription: "The ansible configuration for this stack.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"playbook": schema.StringAttribute{
								Description:         "The playbook to run.",
								MarkdownDescription: "The playbook to run.",
								Optional:            true,
							},
							"inventory": schema.StringAttribute{
								Description:         "The ansible inventory file to use. We recommend checking this into git alongside your playbook files.",
								MarkdownDescription: "The ansible inventory file to use. We recommend checking this into git alongside your playbook files.",
								Optional:            true,
							},
							"additional_args": schema.ListAttribute{
								Description:         "Additional args for the playbook.",
								MarkdownDescription: "Additional args for the playbook.",
								Optional:            true,
								ElementType:         types.StringType,
							},
						},
					},
					"ai_approval": schema.SingleNestedAttribute{
						Description:         "The ai approval configuration for this stack.",
						MarkdownDescription: "The ai approval configuration for this stack.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description:         "Whether ai approval is enabled for this stack.",
								MarkdownDescription: "Whether ai approval is enabled for this stack.",
								Required:            true,
							},
							"ignore_cancel": schema.BoolAttribute{
								Description:         "Whether to ignore the cancellation of a stack run by ai, this allows human approval to override.",
								MarkdownDescription: "Whether to ignore the cancellation of a stack run by ai, this allows human approval to override.",
								Optional:            true,
							},
							"git": schema.SingleNestedAttribute{
								Description:         "The git reference to use for the ai approval rules.",
								MarkdownDescription: "The git reference to use for the ai approval rules.",
								Required:            true,
								Attributes: map[string]schema.Attribute{
									"ref": schema.StringAttribute{
										Description:         "A general git ref, either a branch name or commit sha.",
										MarkdownDescription: "A general git ref, either a branch name or commit sha.",
										Required:            true,
									},
									"folder": schema.StringAttribute{
										Description:         "The subdirectory in the git repository to use.",
										MarkdownDescription: "The subdirectory in the git repository to use.",
										Required:            true,
									},
								},
							},
							"file": schema.StringAttribute{
								Description:         "The rules file to use alongside the git reference.",
								MarkdownDescription: "The rules file to use alongside the git reference.",
								Required:            true,
							},
						},
					},
				},
			},
			"files": schema.MapAttribute{
				Description:         "File path-content map.",
				MarkdownDescription: "File path-content map.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"environment": schema.SetNestedAttribute{
				Description:         "Defines environment variables for the stack.",
				MarkdownDescription: "Defines environment variables for the stack.",
				Optional:            true,
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
							stringvalidator.LengthAtLeast(1),
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("containers")),
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("containers")),
						},
					},
					"labels": schema.MapAttribute{
						Description:         "Kubernetes labels applied to the job.",
						MarkdownDescription: "Kubernetes labels applied to the job.",
						ElementType:         types.StringType,
						Optional:            true,
						Validators:          []validator.Map{mapvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("raw"))},
					},
					"annotations": schema.MapAttribute{
						Description:         "Kubernetes annotations applied to the job.",
						MarkdownDescription: "Kubernetes annotations applied to the job.",
						ElementType:         types.StringType,
						Optional:            true,
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
							setvalidator.SizeAtLeast(1),
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
