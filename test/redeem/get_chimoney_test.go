package redeem_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetChimoney(t *testing.T) {
	tests := []struct {
		name       string
		chiRef     string
		subAccount string
		response   string
		wantErr    bool
	}{
		{
			name:   "successful get chimoney",
			chiRef: "chi_123",
			response: `{
				"status": "success",
				"data": {
					"id": "chi_123",
					"status": "completed",
					"amount": 50.0,
					"email": "test@example.com"
				}
			}`,
			wantErr: false,
		},
		{
			name:       "with subaccount",
			chiRef:     "chi_123",
			subAccount: "sub_123",
			response: `{
				"status": "success",
				"data": {
					"id": "chi_123",
					"status": "completed",
					"amount": 50.0,
					"email": "test@example.com",
					"subAccount": "sub_123"
				}
			}`,
			wantErr: false,
		},
		{
			name:     "empty chi ref",
			chiRef:   "",
			response: "",
			wantErr:  true,
		},
		{
			name:     "invalid chi ref",
			chiRef:   "invalid",
			response: `{"status":"error","message":"Invalid chi ref"}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("unexpected method: got %v want POST", r.Method)
				}
				if r.URL.Path != "/redeem/chimoney/get" {
					t.Errorf("unexpected path: got %v want /redeem/chimoney/get", r.URL.Path)
				}

				var reqBody map[string]string
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				if chiRef, ok := reqBody["chiRef"]; !ok || chiRef != tt.chiRef {
					t.Errorf("unexpected chiRef: got %v want %v", chiRef, tt.chiRef)
				}

				if tt.subAccount != "" {
					if subAccount, ok := reqBody["subAccount"]; !ok || subAccount != tt.subAccount {
						t.Errorf("unexpected subAccount: got %v want %v", subAccount, tt.subAccount)
					}
				}

				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
				}
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			})
			defer server.Close()

			resp, err := client.Redeem.GetChimoney(context.Background(), tt.chiRef, tt.subAccount)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChimoney() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("GetChimoney() got nil response")
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(resp.Data, &data); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				if id, ok := data["id"].(string); !ok || id != tt.chiRef {
					t.Errorf("unexpected id in response: got %v want %v", id, tt.chiRef)
				}

				if status, ok := data["status"].(string); !ok || status == "" {
					t.Error("status field not found or empty in response")
				}

				if resp.Status != "success" {
					t.Errorf("unexpected status: got %v want success", resp.Status)
				}
			}
		})
	}
}
