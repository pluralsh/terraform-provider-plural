package model

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	console "github.com/pluralsh/console/go/client"
)

type ServiceContext struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Configuration types.String `tfsdk:"configuration"`
}

func (sc *ServiceContext) From(response *console.ServiceContextFragment, ctx context.Context, d *diag.Diagnostics) {
	configuration, err := json.Marshal(response.Configuration)
	if err != nil {
		d.AddError("Provider Error", fmt.Sprintf("Cannot marshall metadata, got error: %s", err))
		return
	}

	sc.Id = types.StringValue(response.ID)
	sc.Configuration = types.StringValue(string(configuration))
}

type ServiceContextExtended struct {
	ServiceContext
	Secrets types.Map `tfsdk:"secrets"`
}

func (sc *ServiceContextExtended) Attributes(ctx context.Context, d *diag.Diagnostics) console.ServiceContextAttributes {
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
		Configuration: sc.Configuration.ValueStringPointer(),
		Secrets:       configAttributes,
	}
}
