package redeem_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestChimoneyRedeem(t *testing.T) {
	tests := []struct {
		name       string
		chimoneys  []map[string]interface{}
		subAccount string
		response   string
		wantErr    bool
	}{
		{
			name: "successful single chimoney redemption",
			chimoneys: []map[string]interface{}{
				{
					"id":     "chi_123",
					"amount": 50.0,
					"email":  "test@example.com",
				},
			},
			response: `{
				"status": "success",
				"data": {
					"id": "red_123",
					"status": "pending",
					"transactions": [
						{
							"id": "chi_123",
							"amount": 50.0,
							"email": "test@example.com",
							"status": "pending"
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name: "multiple chimoney redemptions",
			chimoneys: []map[string]interface{}{
				{
					"id":     "chi_123",
					"amount": 50.0,
					"email":  "test1@example.com",
				},
				{
					"id":      "chi_456",
					"amount":  25.0,
					"twitter": "@testuser",
				},
			},
			response: `{
				"status": "success",
				"data": {
					"id": "red_123",
					"status": "pending",
					"transactions": [
						{
							"id": "chi_123",
							"amount": 50.0,
							"email": "test1@example.com",
							"status": "pending"
						},
						{
							"id": "chi_456",
							"amount": 25.0,
							"twitter": "@testuser",
							"status": "pending"
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name: "with subaccount",
			chimoneys: []map[string]interface{}{
				{
					"id":     "chi_123",
					"amount": 50.0,
					"email":  "test@example.com",
				},
			},
			subAccount: "sub_123",
			response: `{
				"status": "success",
				"data": {
					"id": "red_123",
					"status": "pending",
					"subAccount": "sub_123",
					"transactions": [
						{
							"id": "chi_123",
							"amount": 50.0,
							"email": "test@example.com",
							"status": "pending"
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name:      "empty chimoneys list",
			chimoneys: []map[string]interface{}{},
			response:  "",
			wantErr:   true,
		},
		{
			name: "invalid amount",
			chimoneys: []map[string]interface{}{
				{
					"id":     "chi_123",
					"amount": -50.0,
					"email":  "test@example.com",
				},
			},
			response: `{
				"status": "error",
				"message": "Invalid amount"
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("unexpected method: got %v want POST", r.Method)
				}
				if r.URL.Path != "/redeem/chimoney" {
					t.Errorf("unexpected path: got %v want /redeem/chimoney", r.URL.Path)
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

				if tt.subAccount != "" {
					if subAccount, ok := reqBody["subAccount"].(string); !ok || subAccount != tt.subAccount {
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

			resp, err := client.Redeem.Chimoney(context.Background(), tt.chimoneys, tt.subAccount)
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

				if id, ok := data["id"].(string); !ok || id == "" {
					t.Error("id field not found or empty in response")
				}

				if status, ok := data["status"].(string); !ok || status == "" {
					t.Error("status field not found or empty in response")
				}

				transactions, ok := data["transactions"].([]interface{})
				if !ok {
					t.Error("transactions field not found in response")
					return
				}

				if len(transactions) != len(tt.chimoneys) {
					t.Errorf("unexpected number of transactions in response: got %v want %v", len(transactions), len(tt.chimoneys))
				}
			}
		})
	}
}
