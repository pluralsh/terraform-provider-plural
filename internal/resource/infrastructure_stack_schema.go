package resource

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
		},
	}
}
