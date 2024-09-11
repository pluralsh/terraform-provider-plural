package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"
	"terraform-provider-plural/internal/model"
)

var _ resource.ResourceWithConfigure = &sharedSecretResource{}

func NewSharedSecretResource() resource.Resource {
	return &sharedSecretResource{}
}

type sharedSecretResource struct {
	client *client.Client
}

func (in *sharedSecretResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_shared_secret"
}

func (in *sharedSecretResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description:         "A one-time-viewable secret shared with a list of eligible users.",
		MarkdownDescription: "A one-time-viewable secret shared with a list of eligible users.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description:         "The name of this shared secret.",
				MarkdownDescription: "The name of this shared secret.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"secret": schema.StringAttribute{
				Description:         "Content of this shared secret.",
				MarkdownDescription: "Content of this shared secret.",
				Required:            true,
				Sensitive:           true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"notification_bindings": schema.SetNestedAttribute{
				Description:         "The users/groups you want this secret to be delivered to.",
				MarkdownDescription: "The users/groups you want this secret to be delivered to.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"group_id": schema.StringAttribute{
							Optional:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"id": schema.StringAttribute{
							Optional:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"user_id": schema.StringAttribute{
							Optional:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
					},
				},
				PlanModifiers: []planmodifier.Set{setplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (in *sharedSecretResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	data, ok := request.ProviderData.(*common.ProviderData)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Project Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
		return
	}

	in.client = data.Client
}

func (in *sharedSecretResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	data := new(model.SharedSecret)
	response.Diagnostics.Append(request.Plan.Get(ctx, data)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := in.client.ShareSecret(ctx, data.Attributes(ctx, response.Diagnostics))
	if err != nil {
		response.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to share a secret, got error: %s", err))
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (in *sharedSecretResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	data := new(model.SharedSecret)
	response.Diagnostics.Append(request.State.Get(ctx, data)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, data)...)
}

func (in *sharedSecretResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// Ignore.
}

func (in *sharedSecretResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Ignore.
}
