package model

import (
	"context"

	"terraform-provider-plural/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
)

type ServiceContext struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Configuration types.Map    `tfsdk:"configuration"`
}

func (sc *ServiceContext) From(response *console.ServiceContextFragment, ctx context.Context, d diag.Diagnostics) {
	sc.Id = types.StringValue(response.ID)
	sc.Configuration = common.MapFrom(response.Configuration, ctx, d)
}

type ServiceContextExtended struct {
	ServiceContext
	Secrets types.Map `tfsdk:"secrets"`
}

func (sc *ServiceContextExtended) Attributes(ctx context.Context, d diag.Diagnostics) console.ServiceContextAttributes {
	configuration := make(map[string]types.String, len(sc.Configuration.Elements()))
	sc.Configuration.ElementsAs(ctx, &configuration, false)

	secrets := make(map[string]types.String, len(sc.Secrets.Elements()))
	sc.Secrets.ElementsAs(ctx, &secrets, false)
	configAttributes := make([]*console.ConfigAttributes, 0)
	for key, val := range secrets {
		configAttributes = append(configAttributes, &console.ConfigAttributes{
			Name:  key,
			Value: val.ValueStringPointer(),
		})
	}

	return console.ServiceContextAttributes{
		Configuration: common.AttributesJson(configuration, d),
		Secrets:       configAttributes,
	}
}
