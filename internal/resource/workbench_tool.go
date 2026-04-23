package resource

import (
	"context"
	"fmt"

	"terraform-provider-plural/internal/common"
	"terraform-provider-plural/internal/model"
	planmod "terraform-provider-plural/internal/planmodifier"

	"terraform-provider-plural/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
	"github.com/samber/lo"
)

var _ resource.Resource = &WorkbenchToolResource{}
var _ resource.ResourceWithImportState = &WorkbenchToolResource{}

func NewWorkbenchToolResource() resource.Resource {
	return &WorkbenchToolResource{}
}

// WorkbenchToolResource defines the workbench tool resource implementation.
type WorkbenchToolResource struct {
	client *client.Client
}

func (r *WorkbenchToolResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_workbench_tool"
}

func (r *WorkbenchToolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Workbench tool resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Internal identifier of this workbench tool.",
				MarkdownDescription: "Internal identifier of this workbench tool.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description:         "Name of this workbench tool.",
				MarkdownDescription: "Name of this workbench tool.",
				Required:            true,
			},
			"tool": schema.StringAttribute{
				Description:         "Workbench tool type.",
				MarkdownDescription: "Workbench tool type.",
				Required:            true,
				Validators: []validator.String{stringvalidator.OneOf(
					lo.Map(console.AllWorkbenchToolType, func(item console.WorkbenchToolType, index int) string {
						return string(item)
					})...),
				},
			},
			"categories": schema.SetAttribute{
				Description:         "Categories of this workbench tool.",
				MarkdownDescription: "Categories of this workbench tool.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers:       []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
			},
			"project_id": schema.StringAttribute{
				Description:         "ID of the project that this workbench belongs to.",
				MarkdownDescription: "ID of the project that this workbench belongs to.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"mcp_server_id": schema.StringAttribute{
				Description:         "ID of the MCP server referenced by this workbench tool.",
				MarkdownDescription: "ID of the MCP server referenced by this workbench tool.",
				Optional:            true,
			},
			"cloud_connection_id": schema.StringAttribute{
				Description:         "ID of the cloud connection referenced by this workbench tool.",
				MarkdownDescription: "ID of the cloud connection referenced by this workbench tool.",
				Optional:            true,
			},
			"configuration": schema.SingleNestedAttribute{
				Description:         "Configuration of this workbench tool.",
				MarkdownDescription: "Configuration of this workbench tool.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"http": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "HTTP configuration of this workbench tool.",
						MarkdownDescription: "HTTP configuration of this workbench tool.",
						PlanModifiers:       []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
						Attributes: map[string]schema.Attribute{
							"url": schema.StringAttribute{
								Description:         "The request URL.",
								MarkdownDescription: "The request URL.",
								Required:            true,
							},
							"method": schema.StringAttribute{
								Description:         "The HTTP method.",
								MarkdownDescription: "The HTTP method.",
								Required:            true,
								PlanModifiers:       []planmodifier.String{planmod.UppercaseString()},
								Validators: []validator.String{
									stringvalidator.OneOfCaseInsensitive(
										lo.Map(console.AllWorkbenchToolHTTPMethod, func(m console.WorkbenchToolHTTPMethod, _ int) string {
											return string(m)
										})...),
								},
							},
							"headers": schema.MapAttribute{
								Description:         "The request headers.",
								MarkdownDescription: "The request headers.",
								ElementType:         types.StringType,
								Optional:            true,
							},
							"body": schema.StringAttribute{
								Description:         "The request body.",
								MarkdownDescription: "The request body.",
								Optional:            true,
							},
							"input_schema": schema.StringAttribute{
								Description:         "The JSON schema for the tool input.",
								MarkdownDescription: "The JSON schema for the tool input.",
								Required:            true,
							},
						},
					},
					"elastic": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "Elasticsearch connection configuration.",
						MarkdownDescription: "Elasticsearch connection configuration.",
						Attributes: map[string]schema.Attribute{
							"url":      schema.StringAttribute{Required: true},
							"username": schema.StringAttribute{Required: true},
							"password": schema.StringAttribute{Optional: true, Sensitive: true},
							"index":    schema.StringAttribute{Required: true},
						},
					},
					"prometheus": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "Prometheus connection configuration.",
						MarkdownDescription: "Prometheus connection configuration.",
						Attributes: map[string]schema.Attribute{
							"url":       schema.StringAttribute{Required: true},
							"token":     schema.StringAttribute{Optional: true, Sensitive: true},
							"username":  schema.StringAttribute{Optional: true},
							"password":  schema.StringAttribute{Optional: true, Sensitive: true},
							"tenant_id": schema.StringAttribute{Optional: true},
						},
					},
					"loki": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "Loki connection configuration.",
						MarkdownDescription: "Loki connection configuration.",
						Attributes: map[string]schema.Attribute{
							"url":       schema.StringAttribute{Required: true},
							"token":     schema.StringAttribute{Optional: true, Sensitive: true},
							"username":  schema.StringAttribute{Optional: true},
							"password":  schema.StringAttribute{Optional: true, Sensitive: true},
							"tenant_id": schema.StringAttribute{Optional: true},
						},
					},
					"splunk": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "Splunk connection configuration.",
						MarkdownDescription: "Splunk connection configuration.",
						Attributes: map[string]schema.Attribute{
							"url":      schema.StringAttribute{Required: true},
							"token":    schema.StringAttribute{Optional: true, Sensitive: true},
							"username": schema.StringAttribute{Optional: true},
							"password": schema.StringAttribute{Optional: true, Sensitive: true},
						},
					},
					"tempo": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "Tempo connection configuration.",
						MarkdownDescription: "Tempo connection configuration.",
						Attributes: map[string]schema.Attribute{
							"url":       schema.StringAttribute{Required: true},
							"token":     schema.StringAttribute{Optional: true, Sensitive: true},
							"username":  schema.StringAttribute{Optional: true},
							"password":  schema.StringAttribute{Optional: true, Sensitive: true},
							"tenant_id": schema.StringAttribute{Optional: true},
						},
					},
					"jaeger": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "Jaeger connection configuration.",
						MarkdownDescription: "Jaeger connection configuration.",
						Attributes: map[string]schema.Attribute{
							"url":      schema.StringAttribute{Required: true},
							"token":    schema.StringAttribute{Optional: true, Sensitive: true},
							"username": schema.StringAttribute{Optional: true},
							"password": schema.StringAttribute{Optional: true, Sensitive: true},
						},
					},
					"datadog": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "Datadog connection configuration.",
						MarkdownDescription: "Datadog connection configuration.",
						Attributes: map[string]schema.Attribute{
							"site":    schema.StringAttribute{Optional: true},
							"api_key": schema.StringAttribute{Optional: true, Sensitive: true},
							"app_key": schema.StringAttribute{Optional: true, Sensitive: true},
						},
					},
					"dynatrace": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "Dynatrace connection configuration.",
						MarkdownDescription: "Dynatrace connection configuration.",
						Attributes: map[string]schema.Attribute{
							"url":            schema.StringAttribute{Required: true},
							"platform_token": schema.StringAttribute{Required: true, Sensitive: true},
						},
					},
					"cloudwatch": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "CloudWatch connection configuration.",
						MarkdownDescription: "CloudWatch connection configuration.",
						Attributes: map[string]schema.Attribute{
							"region":            schema.StringAttribute{Required: true},
							"log_group_names":   schema.SetAttribute{Optional: true, ElementType: types.StringType},
							"access_key_id":     schema.StringAttribute{Optional: true, Sensitive: true},
							"secret_access_key": schema.StringAttribute{Optional: true, Sensitive: true},
							"role_arn":          schema.StringAttribute{Optional: true},
							"external_id":       schema.StringAttribute{Optional: true, Sensitive: true},
							"role_session_name": schema.StringAttribute{Optional: true},
						},
					},
					"azure": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "Azure Monitor connection configuration.",
						MarkdownDescription: "Azure Monitor connection configuration.",
						Attributes: map[string]schema.Attribute{
							"subscription_id": schema.StringAttribute{Required: true},
							"tenant_id":       schema.StringAttribute{Required: true},
							"client_id":       schema.StringAttribute{Required: true},
							"client_secret":   schema.StringAttribute{Required: true, Sensitive: true},
						},
					},
					"linear": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "Linear connection configuration.",
						MarkdownDescription: "Linear connection configuration.",
						Attributes: map[string]schema.Attribute{
							"access_token": schema.StringAttribute{Optional: true, Sensitive: true},
						},
					},
					"atlassian": schema.SingleNestedAttribute{
						Optional:            true,
						Description:         "Atlassian/Jira connection configuration.",
						MarkdownDescription: "Atlassian/Jira connection configuration.",
						Attributes: map[string]schema.Attribute{
							"service_account": schema.StringAttribute{Optional: true, Sensitive: true},
							"api_token":       schema.StringAttribute{Optional: true, Sensitive: true},
							"email":           schema.StringAttribute{Optional: true},
						},
					},
				},
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *WorkbenchToolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*common.ProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Workbench Tool Resource Configure Type",
			fmt.Sprintf(
				"Expected *common.ProviderData, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)
		return
	}

	r.client = data.Client
}

func (r *WorkbenchToolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := new(model.WorkbenchTool)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs, err := data.Attributes(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get attributes, got error: %s", err))
		return
	}

	response, err := r.client.CreateWorkbenchTool(ctx, *attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workbench tool, got error: %s", err))
		return
	}

	data.From(response.CreateWorkbenchTool, ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *WorkbenchToolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := new(model.WorkbenchTool)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetWorkbenchTool(ctx, data.Id.ValueStringPointer(), nil)
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get workbench tool, got error: %s", err))
		return
	}
	if response == nil || response.WorkbenchTool == nil || client.IsNotFound(err) {
		// Resource not found, remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	data.From(response.WorkbenchTool, ctx, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *WorkbenchToolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := new(model.WorkbenchTool)
	resp.Diagnostics.Append(req.Plan.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs, err := data.Attributes(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get attributes, got error: %s", err))
		return
	}

	_, err = r.client.UpdateWorkbenchTool(ctx, data.Id.ValueString(), *attrs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update workbench tool, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *WorkbenchToolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := new(model.WorkbenchTool)
	resp.Diagnostics.Append(req.State.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.DeleteWorkbenchTool(ctx, data.Id.ValueString()); err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete workbench tool, got error: %s", err))
		return
	}
}

func (r *WorkbenchToolResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
