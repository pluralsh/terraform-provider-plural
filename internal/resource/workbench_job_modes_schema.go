package resource

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
)

// workbenchJobModesSchema returns the schema of mode-specific options for workbench jobs started by a resource.
func workbenchJobModesSchema(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description:         description,
		MarkdownDescription: description,
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"plan": schema.BoolAttribute{
				Description:         "Whether planning mode is enabled for the job.",
				MarkdownDescription: "Whether planning mode is enabled for the job.",
				Optional:            true,
			},
			"verification": schema.BoolAttribute{
				Description:         "Whether verification mode is enabled for the job.",
				MarkdownDescription: "Whether verification mode is enabled for the job.",
				Optional:            true,
			},
			"model": schema.SingleNestedAttribute{
				Description:         "AI model override for the job.",
				MarkdownDescription: "AI model override for the job.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"provider": schema.StringAttribute{
						Description:         withAllowedValues("AI provider for the job.", enumValues(gqlclient.AllAiProvider), false),
						MarkdownDescription: withAllowedValues("AI provider for the job.", enumValues(gqlclient.AllAiProvider), true),
						Required:            true,
						Validators:          []validator.String{stringvalidator.OneOf(enumValues(gqlclient.AllAiProvider)...)},
					},
					"model": schema.StringAttribute{
						Description:         "Model name for the job.",
						MarkdownDescription: "Model name for the job.",
						Required:            true,
						Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
					},
				},
			},
			"coding": schema.SingleNestedAttribute{
				Description:         "Coding mode options for the job.",
				MarkdownDescription: "Coding mode options for the job.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"babysit": schema.BoolAttribute{
						Description:         "Whether babysit mode is enabled for coding agent runs.",
						MarkdownDescription: "Whether babysit mode is enabled for coding agent runs.",
						Optional:            true,
					},
					"approval": schema.BoolAttribute{
						Description:         "Whether coding agent runs require approval before continuing.",
						MarkdownDescription: "Whether coding agent runs require approval before continuing.",
						Optional:            true,
					},
					"review": schema.BoolAttribute{
						Description:         "Whether pull request review mode is enabled for coding agent runs.",
						MarkdownDescription: "Whether pull request review mode is enabled for coding agent runs.",
						Optional:            true,
					},
				},
			},
			"budget": schema.SingleNestedAttribute{
				Description:         "Budget limits for the job.",
				MarkdownDescription: "Budget limits for the job.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"cost": schema.Float64Attribute{
						Description:         "Maximum cost budget for the job.",
						MarkdownDescription: "Maximum cost budget for the job.",
						Optional:            true,
					},
					"tokens": schema.Int64Attribute{
						Description:         "Maximum token budget for the job.",
						MarkdownDescription: "Maximum token budget for the job.",
						Optional:            true,
					},
				},
			},
			"kubernetes": schema.SingleNestedAttribute{
				Description:         "Kubernetes action options for the job.",
				MarkdownDescription: "Kubernetes action options for the job.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"update": schema.BoolAttribute{
						Description:         "Whether Kubernetes update actions are enabled. Defaults to false.",
						MarkdownDescription: "Whether Kubernetes update actions are enabled. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"delete": schema.BoolAttribute{
						Description:         "Whether Kubernetes delete actions are enabled. Defaults to false.",
						MarkdownDescription: "Whether Kubernetes delete actions are enabled. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"exec": schema.BoolAttribute{
						Description:         "Whether Kubernetes exec actions are enabled. Defaults to false.",
						MarkdownDescription: "Whether Kubernetes exec actions are enabled. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"drain": schema.BoolAttribute{
						Description:         "Whether Kubernetes node drain actions are enabled. Defaults to false.",
						MarkdownDescription: "Whether Kubernetes node drain actions are enabled. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"exclude_namespaces": schema.SetAttribute{
						Description:         "Namespaces the agent can never act in.",
						MarkdownDescription: "Namespaces the agent can never act in.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"require_namespaces": schema.SetAttribute{
						Description:         "If set, the agent can only act in these namespaces.",
						MarkdownDescription: "If set, the agent can only act in these namespaces.",
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},
		},
	}
}
