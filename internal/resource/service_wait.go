package resource

import (
	"context"
	"fmt"
	"time"

	"terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/common"
	customvalidator "terraform-provider-plural/internal/validator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
	"k8s.io/apimachinery/pkg/util/wait"
)

type serviceWait struct {
	ServiceID types.String `tfsdk:"service_id"`
	Warmup    types.String `tfsdk:"warmup"`
	Duration  types.String `tfsdk:"duration"`
}

func (in *serviceWait) ParseWarmup() (time.Duration, error) {
	return time.ParseDuration(in.Warmup.ValueString())
}

func (in *serviceWait) ParseDuration() (time.Duration, error) {
	return time.ParseDuration(in.Duration.ValueString())
}

var _ resource.ResourceWithConfigure = &serviceWaitResource{}

func NewServiceWaitResource() resource.Resource {
	return &serviceWaitResource{}
}

type serviceWaitResource struct {
	client *client.Client
}

func (in *serviceWaitResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_service_wait"
}

func (in *serviceWaitResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"service_id": schema.StringAttribute{
				Description:         "ID the service deployment that should be checked.",
				MarkdownDescription: "ID the service deployment that should be checked.",
				Required:            true,
				Validators:          []validator.String{customvalidator.UUID()},
			},
			"warmup": schema.StringAttribute{
				Description:         "Initial delay before checking the service deployment health. Defaults to 5 minutes.",
				MarkdownDescription: "Initial delay before checking the service deployment health. Defaults to 5 minutes.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("5m"),
				Validators:          []validator.String{customvalidator.Duration()},
			},
			"duration": schema.StringAttribute{
				Description:         "Maximum duration to wait for the service deployment to become healthy. Minimum 1 minute. Defaults to 10 minutes.",
				MarkdownDescription: "Maximum duration to wait for the service deployment to become healthy. Minimum 1 minute. Defaults to 10 minutes.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("10m"),
				Validators:          []validator.String{customvalidator.MinDuration(time.Minute)},
			},
		},
	}
}

func (in *serviceWaitResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (in *serviceWaitResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	data := new(serviceWait)
	response.Diagnostics.Append(request.Plan.Get(ctx, data)...)
	if response.Diagnostics.HasError() {
		return
	}

	if err := in.Wait(data); err != nil {
		response.Diagnostics.AddError("Client Error", fmt.Sprintf("Got error while waiting for service: %s", err))
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (in *serviceWaitResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// Ignore.
}

func (in *serviceWaitResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// Ignore.
}

func (in *serviceWaitResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Ignore.
}

func (in *serviceWaitResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("service_id"), req, resp)
}

func (in *serviceWaitResource) Wait(data *serviceWait) error {
	warmup, err := data.ParseWarmup()
	if err != nil {
		return fmt.Errorf("unable to parse warmup duration, got error: %s", err.Error())
	}

	duration, err := data.ParseDuration()
	if err != nil {
		return fmt.Errorf("unable to parse duration, got error: %s", err.Error())
	}

	time.Sleep(warmup)

	var status console.ServiceDeploymentStatus
	if err = wait.PollUntilContextTimeout(context.Background(), 30*time.Second, duration, true,
		func(ctx context.Context) (done bool, err error) {
			service, err := in.client.GetServiceDeployment(ctx, data.ServiceID.ValueString())
			if err != nil {
				return false, nil
			}

			status = service.ServiceDeployment.Status

			return status == console.ServiceDeploymentStatusHealthy, nil
		}); err != nil {
		return fmt.Errorf("service did not became healthy, got error: %s", err.Error())
	}

	return nil
}
