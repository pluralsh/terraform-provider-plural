package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gqltransport "github.com/Yamashou/gqlgenc/clientv2"
	gqlclient "github.com/pluralsh/console/go/client"
)

func TestPolicyOperationsUseSchemaFaithfulGraphQLDocuments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}

		var request struct {
			OperationName string         `json:"operationName"`
			Query         string         `json:"query"`
			Variables     map[string]any `json:"variables"`
		}
		if err := json.Unmarshal(body, &request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if request.OperationName != "TerraformGetPolicy" {
			t.Fatalf("expected TerraformGetPolicy operation, got %q", request.OperationName)
		}
		if !strings.Contains(request.Query, "policy(id: $id, name: $name)") {
			t.Fatalf("expected schema policy query, got %q", request.Query)
		}
		if request.Variables["id"] != "policy-1" {
			t.Fatalf("expected ID variable, got %#v", request.Variables["id"])
		}

		_, _ = w.Write([]byte(`{"data":{"policy":{"id":"policy-1","name":"allow","type":"WORKBENCH","policy":"package policy","project":{"id":"project-1"}}}}`))
	}))
	defer server.Close()

	raw := &gqlclient.Client{Client: gqltransport.NewClient(server.Client(), server.URL, nil)}
	client := NewClient(raw)
	response, err := client.GetPolicy(context.Background(), stringPointer("policy-1"), nil)
	if err != nil {
		t.Fatalf("GetPolicy returned error: %v", err)
	}
	if response.Policy == nil || response.Policy.ID != "policy-1" {
		t.Fatalf("expected policy response, got %#v", response.Policy)
	}
}

func TestBindingPolicyMutationDocumentsIncludeRequiredFields(t *testing.T) {
	for _, document := range []string{createBindingPolicyDocument, updateBindingPolicyDocument, deleteBindingPolicyDocument} {
		if !strings.Contains(document, "BindingPolicy") {
			t.Fatalf("expected binding policy operation in document %q", document)
		}
		if !strings.Contains(document, "policy {") || !strings.Contains(document, "bindPolicy {") {
			t.Fatalf("expected policy ID fields in document %q", document)
		}
	}
}

func stringPointer(value string) *string {
	return &value
}
