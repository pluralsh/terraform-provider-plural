package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func (r *ProjectResource) schema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Internal identifier of this project.",
				MarkdownDescription: "Internal identifier of this project.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description:         "Human-readable name of this project.",
				MarkdownDescription: "Human-readable name of this project.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"description": schema.StringAttribute{
				Description:         "Description of this project.",
				MarkdownDescription: "Description of this project.",
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"default": schema.StringAttribute{
				Computed: true,
			},
			"bindings": schema.SingleNestedAttribute{
				Description:         "Read and write policies of this project.",
				MarkdownDescription: "Read and write policies of this project.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"read": schema.SetNestedAttribute{
						Description:         "Read policies of this project.",
						MarkdownDescription: "Read policies of this project.",
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
						Description:         "Write policies of this project.",
						MarkdownDescription: "Write policies of this project.",
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
