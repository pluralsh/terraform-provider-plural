package resource

import (
	"testing"

	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDashboardResourceSchemaMatchesModel(t *testing.T) {
	datasource := &model.DashboardDatasource{
		Type:  types.StringValue("METRICS"),
		Tool:  types.StringValue("prometheus"),
		Input: types.StringValue(`{"query":"up"}`),
	}
	dashboard := &model.Dashboard{
		Id:          types.StringValue("dashboard-1"),
		WorkbenchID: types.StringValue("workbench-1"),
		Name:        types.StringValue("overview"),
		Description: types.StringNull(),
		Graphs: []*model.DashboardGraph{
			{
				Identifier:  types.StringValue("section"),
				Title:       types.StringValue("Errors"),
				Description: types.StringNull(),
				Type:        types.StringValue("SECTION"),
				SectionID:   types.StringNull(),
				Markdown:    types.StringNull(),
				Options:     types.StringValue(`{"collapsed":false}`),
				Layout:      model.DashboardGraphLayout{X: types.Int64Value(0), Y: types.Int64Value(0), W: types.Int64Value(12), H: types.Int64Value(1)},
			},
			{
				Identifier:  types.StringValue("errors"),
				Title:       types.StringNull(),
				Description: types.StringNull(),
				Type:        types.StringValue("TIMESERIES"),
				SectionID:   types.StringValue("section"),
				Markdown:    types.StringNull(),
				Options:     types.StringNull(),
				Layout:      model.DashboardGraphLayout{X: types.Int64Value(0), Y: types.Int64Value(1), W: types.Int64Value(6), H: types.Int64Value(4)},
				Datasource:  datasource,
			},
		},
		Inputs: []*model.DashboardInput{{
			Name:        types.StringValue("namespace"),
			Label:       types.StringValue("Namespace"),
			Description: types.StringNull(),
			Type:        types.StringValue("SELECT"),
			Default:     types.StringValue("prod"),
			Options:     types.ListValueMust(types.StringType, []attr.Value{types.StringValue("prod"), types.StringValue("dev")}),
			Required:    types.BoolValue(false),
			Datasource:  datasource,
		}},
	}

	assertSchemaMatchesModel(t, NewDashboardResource(), dashboard, new(model.Dashboard))
	assertSchemaMatchesModel(t, NewDashboardResource(), &model.Dashboard{
		Id:          types.StringValue("dashboard-2"),
		WorkbenchID: types.StringValue("workbench-1"),
		Name:        types.StringValue("empty"),
		Description: types.StringNull(),
	}, new(model.Dashboard))
}
