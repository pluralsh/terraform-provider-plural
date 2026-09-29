package datasource

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestCloudConnectionDataSourceSchemaLookupAndComputedFields(t *testing.T) {
	cloudConnectionDataSource := &cloudConnectionDataSource{}
	response := datasource.SchemaResponse{}
	cloudConnectionDataSource.Schema(context.Background(), datasource.SchemaRequest{}, &response)

	assertDataSourceStringFlags(t, response.Schema.Attributes, "id", false, true, true)
	assertDataSourceStringFlags(t, response.Schema.Attributes, "name", false, true, true)
	assertDataSourceStringFlags(t, response.Schema.Attributes, "cloud_provider", false, false, true)

	for _, name := range []string{"configuration", "read_bindings"} {
		if _, ok := response.Schema.Attributes[name]; ok {
			t.Fatalf("expected no %q attribute on the data source", name)
		}
	}
}
