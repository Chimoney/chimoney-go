package payouts_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/chimoney/chimoney-go/modules/payouts"
)

func TestChimoney(t *testing.T) {
	tests := []struct {
		name       string
		chimoneys  []payouts.ChimoneyPayload
		subAccount string
		response   string
		wantErr    bool
	}{
		{
			name: "successful single chimoney",
			chimoneys: []payouts.ChimoneyPayload{
				{
					ValueInUSD: 50.0,
					Email:      "test@example.com",
				},
			},
			response: `{
				"status": "success",
				"data": {
					"id": "chi_123",
					"status": "pending",
					"transactions": [
						{
							"email": "test@example.com",
							"amount": 50.0,
							"status": "pending"
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name: "successful multiple chimoneys",
			chimoneys: []payouts.ChimoneyPayload{
				{
					ValueInUSD: 50.0,
					Email:      "test1@example.com",
				},
				{
					ValueInUSD: 25.0,
					Twitter:    "@testuser",
				},
			},
			response: `{
				"status": "success",
				"data": {
					"id": "chi_123",
					"status": "pending",
					"transactions": [
						{
							"email": "test1@example.com",
							"amount": 50.0,
							"status": "pending"
						},
						{
							"twitter": "@testuser",
							"amount": 25.0,
							"status": "pending"
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name: "with subaccount",
			chimoneys: []payouts.ChimoneyPayload{
				{
					ValueInUSD: 50.0,
					Email:      "test@example.com",
				},
			},
			subAccount: "sub_123",
			response: `{
				"status": "success",
				"data": {
					"id": "chi_123",
					"status": "pending",
					"subAccount": "sub_123",
					"transactions": [
						{
							"email": "test@example.com",
							"amount": 50.0,
							"status": "pending"
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name:      "empty chimoney list",
			chimoneys: []payouts.ChimoneyPayload{},
			response:  `{"status":"error","message":"No chimoney payouts provided"}`,
			wantErr:   true,
		},
		{
			name: "invalid amount",
			chimoneys: []payouts.ChimoneyPayload{
				{
					ValueInUSD: -50.0,
					Email:      "test@example.com",
				},
			},
			response: `{"status":"error","message":"Invalid amount"}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("unexpected method: got %v want POST", r.Method)
				}
				if r.URL.Path != "/payouts/chimoney" {
					t.Errorf("unexpected path: got %v want /payouts/chimoney", r.URL.Path)
				}

				var reqBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				chimoneys, ok := reqBody["chimoneys"].([]interface{})
				if !ok {
					t.Error("chimoneys field not found in request")
					return
				}

				if len(chimoneys) != len(tt.chimoneys) {
					t.Errorf("unexpected number of chimoneys: got %v want %v", len(chimoneys), len(tt.chimoneys))
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

			resp, err := client.Payouts.Chimoney(context.Background(), tt.chimoneys, tt.subAccount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chimoney() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("Chimoney() got nil response")
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(resp.Data, &data); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				transactions, ok := data["transactions"].([]interface{})
				if !ok {
					t.Error("transactions field not found in response")
					return
				}

				if len(transactions) != len(tt.chimoneys) {
					t.Errorf("unexpected number of transactions: got %v want %v", len(transactions), len(tt.chimoneys))
				}
			}
		})
	}
}
