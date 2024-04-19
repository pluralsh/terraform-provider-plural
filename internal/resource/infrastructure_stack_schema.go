package resource

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	gqlclient "github.com/pluralsh/console-client-go"
)

func (r *InfrastructureStackResource) schema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Internal identifier of this infrastructure stack.",
				MarkdownDescription: "Internal identifier of this infrastructure stack.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Human-readable name of this infrastructure stack.",
				MarkdownDescription: "Human-readable name of this infrastructure stack.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"type": schema.StringAttribute{
				Required:            true,
				Description:         "A type for the stack, specifies the tool to use to apply it.",
				MarkdownDescription: "A type for the stack, specifies the tool to use to apply it.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:          []validator.String{stringvalidator.OneOf(string(gqlclient.StackTypeAnsible), string(gqlclient.StackTypeTerraform))},
			},
		},
	}
}
