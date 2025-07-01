package resource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"
	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &CloudConnectionResource{}
var _ resource.ResourceWithImportState = &CloudConnectionResource{}

func NewCloudConnectionResource() resource.Resource {
	return &CloudConnectionResource{}
}

// CloudConnectionResource defines the cloud connection resource implementation.
type CloudConnectionResource struct {
	client *client.Client
}

func (r *CloudConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_connection"
}

func (r *CloudConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A cloud connection resource for connecting cloud provider accounts.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Internal identifier of this cloud connection.",
				MarkdownDescription: "Internal identifier of this cloud connection.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description:         "Name of this cloud connection.",
				MarkdownDescription: "Name of this cloud connection.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"cloud_provider": schema.StringAttribute{
				Description:         "Cloud provider type (AWS, GCP, etc).",
				MarkdownDescription: "Cloud provider type (AWS, GCP, etc).",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"configuration": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Cloud provider configuration block.",
				Attributes: map[string]schema.Attribute{
					"aws": schema.ObjectAttribute{
						Optional:            true,
						MarkdownDescription: "AWS-specific configuration.",
						AttributeTypes:      r.awsCloudConnectionAttrTypes(),
					},
					"gcp": schema.ObjectAttribute{
						Optional:            true,
						MarkdownDescription: "GCP-specific configuration.",
						AttributeTypes:      r.gcpCloudConnectionAttrTypes(),
					},
					"azure": schema.ObjectAttribute{
						Optional:            true,
						MarkdownDescription: "Azure-specific configuration.",
						AttributeTypes:      r.azureCloudConnectionAttrTypes(),
					},
				},
			},
			"read_bindings": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"group_id": schema.StringAttribute{Optional: true},
						"user_id":  schema.StringAttribute{Optional: true},
						"id":       schema.StringAttribute{Optional: true},
					},
				},
				PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *CloudConnectionResource) awsCloudConnectionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"access_key_id":     types.StringType,
		"secret_access_key": types.StringType,
		"region":            types.StringType,
	}
}

func (r *CloudConnectionResource) gcpCloudConnectionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"service_account_key": types.StringType,
		"project_id":          types.StringType,
	}
}

func (r *CloudConnectionResource) azureCloudConnectionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{

		"subscription_id": types.StringType,
		"tenant_id":       types.StringType,
		"client_id":       types.StringType,
		"client_secret":   types.StringType,
	}
}

func (r *CloudConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Cloud Connection Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = data.Client
}

func (r *CloudConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data model.CloudConnection
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.UpsertCloudConnection(ctx, data.Attributes(ctx, &resp.Diagnostics))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create cloud connection, got error: %s", err))
		return
	}

	data.From(result.UpsertCloudConnection, ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *CloudConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data model.CloudConnection
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetCloudConnection(ctx, data.Id.ValueStringPointer(), data.Name.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cloud connection, got error: %s", err))
		return
	}

	data.From(response.CloudConnection, ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *CloudConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := new(model.CloudConnection)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attr := data.Attributes(ctx, &resp.Diagnostics)

	_, err := r.client.UpsertCloudConnection(ctx, attr)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update cloud connection, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *CloudConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := new(model.CloudConnection)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteCloudConnection(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete cloud connection, got error: %s", err))
		return
	}
}

func (r *CloudConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
