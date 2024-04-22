package resource

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
				Description:         "A type for the stack, specifies the tool to use to apply it.",
				MarkdownDescription: "A type for the stack, specifies the tool to use to apply it.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:          []validator.String{stringvalidator.OneOf(string(gqlclient.StackTypeAnsible), string(gqlclient.StackTypeTerraform))},
			},
			"approval": schema.BoolAttribute{
				Description:         "Determines whether to require approval.",
				MarkdownDescription: "Determines whether to require approval.",
				Optional:            true,
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
							Default:             booldefault.StaticBool(false),
						},
					},
				},
			},
		},
	}
}
