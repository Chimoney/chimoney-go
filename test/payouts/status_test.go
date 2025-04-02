package payouts_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestStatus(t *testing.T) {
	tests := []struct {
		name       string
		issueID    string
		subAccount string
		response   string
		wantErr    bool
	}{
		{
			name:    "successful status check",
			issueID: "issue_123",
			response: `{
				"status": "success",
				"data": {
					"id": "issue_123",
					"status": "completed",
					"transactions": [
						{
							"id": "txn_123",
							"status": "completed",
							"amount": 50.0
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name:       "with subaccount",
			issueID:    "issue_123",
			subAccount: "sub_123",
			response: `{
				"status": "success",
				"data": {
					"id": "issue_123",
					"status": "completed",
					"subAccount": "sub_123",
					"transactions": [
						{
							"id": "txn_123",
							"status": "completed",
							"amount": 50.0
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name:     "invalid issue id",
			issueID:  "invalid",
			response: `{"status":"error","message":"Invalid issue ID"}`,
			wantErr:  true,
		},
		{
			name:     "empty issue id",
			issueID:  "",
			response: `{"status":"error","message":"Issue ID is required"}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("unexpected method: got %v want POST", r.Method)
				}
				if r.URL.Path != "/payouts/status" {
					t.Errorf("unexpected path: got %v want /payouts/status", r.URL.Path)
				}

				var reqBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				if issueID, ok := reqBody["issueID"].(string); ok {
					if issueID != tt.issueID {
						t.Errorf("unexpected issueID: got %v want %v", issueID, tt.issueID)
					}
				}

				if subAccount, ok := reqBody["subAccount"].(string); ok {
					if subAccount != tt.subAccount {
						t.Errorf("unexpected subAccount: got %v want %v", subAccount, tt.subAccount)
					}
				}

				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
				}
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.Payouts.Status(context.Background(), tt.issueID, tt.subAccount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Status() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("Status() got nil response")
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(resp.Data, &data); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				if id, ok := data["id"].(string); !ok || id != tt.issueID {
					t.Errorf("unexpected issue ID in response: got %v want %v", id, tt.issueID)
				}

				if status, ok := data["status"].(string); !ok || status == "" {
					t.Error("status field not found or empty in response")
				}

				if transactions, ok := data["transactions"].([]interface{}); !ok || len(transactions) == 0 {
					t.Error("transactions field not found or empty in response")
				}
			}
		})
	}
}
