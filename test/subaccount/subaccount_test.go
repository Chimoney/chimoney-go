package subaccount_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chimoney/chimoney-go/modules/subaccount"
	"github.com/chimoney/chimoney-go/test/testclient"
)

func setupTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *testclient.Client) {
	server := httptest.NewServer(handler)
	client := testclient.New(
		testclient.WithTestServer(server.URL),
		testclient.WithAPIKey("test-api-key"),
	)
	return server, client
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name     string
		request  map[string]interface{}
		response string
		wantErr  bool
	}{
		{
			name: "successful creation",
			request: map[string]interface{}{
				"name":        "Test Account",
				"email":       "test@example.com",
				"description": "Test sub-account",
			},
			response: `{
				"status": "success",
				"data": {
					"id": "sub_123",
					"name": "Test Account",
					"email": "test@example.com",
					"description": "Test sub-account"
				}
			}`,
			wantErr: false,
		},
		{
			name: "without description",
			request: map[string]interface{}{
				"name":  "Test Account",
				"email": "test@example.com",
			},
			response: `{
				"status": "success",
				"data": {
					"id": "sub_123",
					"name": "Test Account",
					"email": "test@example.com"
				}
			}`,
			wantErr: false,
		},
		{
			name: "missing required fields",
			request: map[string]interface{}{
				"name": "Test Account",
			},
			response: `{
				"status": "error",
				"message": "Email is required"
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != "POST" {
					t.Errorf("unexpected method: got %v want POST", r.Method)
				}
				if r.URL.Path != "/sub-account" {
					t.Errorf("unexpected path: got %v want /sub-account", r.URL.Path)
				}

				// Check request body
				var reqBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				for k, v := range tt.request {
					if got := reqBody[k]; got != v {
						t.Errorf("unexpected %s: got %v want %v", k, got, v)
					}
				}

				// Send response
				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
				}
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			req := &subaccount.CreateRequest{}
			if name, ok := tt.request["name"].(string); ok {
				req.Name = name
			}
			if email, ok := tt.request["email"].(string); ok {
				req.Email = email
			}
			if desc, ok := tt.request["description"].(string); ok {
				req.Description = desc
			}

			resp, err := client.SubAccount.Create(context.Background(), req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("Create() got nil response")
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(resp.Data, &data); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				if got := data["name"].(string); got != tt.request["name"] {
					t.Errorf("unexpected name in response: got %v want %v", got, tt.request["name"])
				}
				if got := data["email"].(string); got != tt.request["email"] {
					t.Errorf("unexpected email in response: got %v want %v", got, tt.request["email"])
				}
			}
		})
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name     string
		response string
		wantErr  bool
	}{
		{
			name: "successful list",
			response: `{
				"status": "success",
				"data": [
					{
						"id": "sub_123",
						"name": "Test Account 1",
						"email": "test1@example.com"
					},
					{
						"id": "sub_456",
						"name": "Test Account 2",
						"email": "test2@example.com"
					}
				]
			}`,
			wantErr: false,
		},
		{
			name: "empty list",
			response: `{
				"status": "success",
				"data": []
			}`,
			wantErr: false,
		},
		{
			name:     "error response",
			response: `{"status":"error","message":"Internal server error"}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("unexpected method: got %v want GET", r.Method)
				}
				if r.URL.Path != "/sub-account/list" {
					t.Errorf("unexpected path: got %v want /sub-account/list", r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusInternalServerError)
				}
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.SubAccount.List(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("List() got nil response")
					return
				}

				var accounts []interface{}
				if err := json.Unmarshal(resp.Data, &accounts); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				if !tt.wantErr && len(accounts) == 0 && tt.name != "empty list" {
					t.Error("expected non-empty accounts array")
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		response string
		wantErr  bool
	}{
		{
			name:     "successful deletion",
			id:       "sub_123",
			response: `{"status":"success","message":"Sub-account deleted successfully"}`,
			wantErr:  false,
		},
		{
			name:     "invalid id",
			id:       "",
			response: `{"status":"error","message":"Invalid sub-account ID"}`,
			wantErr:  true,
		},
		{
			name:     "non-existent id",
			id:       "sub_999",
			response: `{"status":"error","message":"Sub-account not found"}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("unexpected method: got %v want DELETE", r.Method)
				}
				if r.URL.Path != "/sub-account" {
					t.Errorf("unexpected path: got %v want /sub-account", r.URL.Path)
				}

				// Check request body
				var reqBody map[string]string
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				if got := reqBody["id"]; got != tt.id {
					t.Errorf("unexpected id: got %v want %v", got, tt.id)
				}

				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
				}
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.SubAccount.Delete(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("Delete() got nil response")
					return
				}
				if resp.Status != "success" {
					t.Errorf("unexpected status: got %v want success", resp.Status)
				}
			}
		})
	}
}
