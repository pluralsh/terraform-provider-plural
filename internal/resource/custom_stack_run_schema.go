package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *CustomStackRunResource) schema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Internal identifier of this custom run.",
				MarkdownDescription: "Internal identifier of this custom run.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description:         "Human-readable name of this custom run.",
				MarkdownDescription: "Human-readable name of this custom run.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"documentation": schema.StringAttribute{
				Description:         "Extended documentation to explain what this will do.",
				MarkdownDescription: "Extended documentation to explain what this will do.",
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"stack_id": schema.StringAttribute{
				Description:         "The ID of the stack to attach it to.",
				MarkdownDescription: "The ID of the stack to attach it to.",
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"commands": schema.SetNestedAttribute{
				Description:         "The commands for this custom run.",
				MarkdownDescription: "The commands for this custom run.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cmd": schema.StringAttribute{
							Description:         "Command to run.",
							MarkdownDescription: "Command to run.",
							Required:            true,
						},
						"args": schema.SetAttribute{
							Description:         "Arguments to pass to the command when executing it.",
							MarkdownDescription: "Arguments to pass to the command when executing it.",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"dir": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
			"configuration": schema.SetNestedAttribute{
				Description:         "Self-service configuration which will be presented in UI before triggering.",
				MarkdownDescription: "Self-service configuration which will be presented in UI before triggering.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
						},
						"name": schema.StringAttribute{
							Required: true,
						},
						"default": schema.StringAttribute{
							Optional: true,
						},
						"documentation": schema.StringAttribute{
							Optional: true,
						},
						"longform": schema.StringAttribute{
							Optional: true,
						},
						"placeholder": schema.StringAttribute{
							Optional: true,
						},
						"optional": schema.BoolAttribute{
							Optional: true,
						},
						"condition": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"operation": schema.StringAttribute{
									Required: true,
								},
								"field": schema.StringAttribute{
									Required: true,
								},
								"value": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
	}
}
