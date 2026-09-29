package datasource

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	internalclient "terraform-provider-plural/internal/client"
	"terraform-provider-plural/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gqlclient "github.com/pluralsh/console/go/client"
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

func TestCloudConnectionDataSourceReadBySelector(t *testing.T) {
	for name, testCase := range map[string]struct {
		config    model.CloudConnectionDataSource
		variables map[string]any
	}{
		"by id": {
			config:    model.CloudConnectionDataSource{Id: types.StringValue("cc-1"), Name: types.StringNull(), CloudProvider: types.StringNull()},
			variables: map[string]any{"id": "cc-1", "name": nil},
		},
		"by name": {
			config:    model.CloudConnectionDataSource{Id: types.StringNull(), Name: types.StringValue("aws-prod"), CloudProvider: types.StringNull()},
			variables: map[string]any{"id": nil, "name": "aws-prod"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			var variables map[string]any
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var request struct {
					Variables map[string]any `json:"variables"`
				}
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					t.Errorf("unexpected request body: %v", err)
				}
				variables = request.Variables

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"data":{"cloudConnection":{"id":"cc-1","name":"aws-prod","provider":"AWS","readBindings":[]}}}`))
			}))
			defer server.Close()

			state := readCloudConnection(t, server, testCase.config)

			for key, expected := range testCase.variables {
				if variables[key] != expected {
					t.Fatalf("expected query variable %s=%v, got %v", key, expected, variables[key])
				}
			}
			if state.Id != types.StringValue("cc-1") || state.Name != types.StringValue("aws-prod") {
				t.Fatalf("unexpected id/name in state: %v/%v", state.Id, state.Name)
			}
			if state.CloudProvider != types.StringValue("AWS") {
				t.Fatalf("expected cloud_provider AWS in state, got %v", state.CloudProvider)
			}
		})
	}
}

// readCloudConnection runs the data source Read against server with config and returns the
// resulting state.
func readCloudConnection(t *testing.T, server *httptest.Server, config model.CloudConnectionDataSource) model.CloudConnectionDataSource {
	t.Helper()
	ctx := context.Background()

	cloudConnectionDataSource := &cloudConnectionDataSource{
		client: internalclient.NewClient(gqlclient.NewClient(server.Client(), server.URL, nil)),
	}
	schemaResponse := datasource.SchemaResponse{}
	cloudConnectionDataSource.Schema(ctx, datasource.SchemaRequest{}, &schemaResponse)

	// tfsdk.Config has no setter, so the raw config is built through a state with the same schema.
	configState := tfsdk.State{Schema: schemaResponse.Schema}
	if diags := configState.Set(ctx, &config); diags.HasError() {
		t.Fatalf("unexpected diagnostics building config: %v", diags)
	}

	response := datasource.ReadResponse{State: tfsdk.State{Schema: schemaResponse.Schema}}
	cloudConnectionDataSource.Read(ctx, datasource.ReadRequest{
		Config: tfsdk.Config{Schema: schemaResponse.Schema, Raw: configState.Raw},
	}, &response)
	if response.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics reading cloud connection: %v", response.Diagnostics)
	}

	var state model.CloudConnectionDataSource
	if diags := response.State.Get(ctx, &state); diags.HasError() {
		t.Fatalf("unexpected diagnostics reading state: %v", diags)
	}
	return state
}
