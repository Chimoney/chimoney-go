package wallet_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chimoney/chimoney-go/test/testclient"
)

func TestGetBalance(t *testing.T) {
	tests := []struct {
		name       string
		subAccount string
		response   string
		wantErr    bool
	}{
		{
			name: "successful get wallet balance",
			response: `{
				"status": "success",
				"data": {
					"balance": 1000.50,
					"currency": "USD",
					"chimoneyBalance": 500.25
				}
			}`,
			wantErr: false,
		},
		{
			name:       "with subaccount",
			subAccount: "sub_123",
			response: `{
				"status": "success",
				"data": {
					"balance": 750.25,
					"currency": "USD",
					"chimoneyBalance": 250.75,
					"subAccount": "sub_123"
				}
			}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("unexpected method: got %v want POST", r.Method)
				}
				if r.URL.Path != "/wallet/balance" {
					t.Errorf("unexpected path: got %v want /wallet/balance", r.URL.Path)
				}

				if tt.subAccount != "" {
					var reqBody map[string]string
					if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
						t.Errorf("failed to decode request body: %v", err)
						return
					}

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

			resp, err := client.Wallet.GetBalance(context.Background(), tt.subAccount)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWalletBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("GetWalletBalance() got nil response")
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(resp.Data, &data); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				if _, ok := data["balance"].(float64); !ok {
					t.Error("balance field not found or not a number in response")
				}

				if currency, ok := data["currency"].(string); !ok || currency == "" {
					t.Error("currency field not found or empty in response")
				}

				if _, ok := data["chimoneyBalance"].(float64); !ok {
					t.Error("chimoneyBalance field not found or not a number in response")
				}

				if resp.Status != "success" {
					t.Errorf("unexpected status: got %v want success", resp.Status)
				}
			}
		})
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name       string
		subAccount string
		response   string
		wantErr    bool
	}{
		{
			name: "successful get wallet transactions",
			response: `{
				"status": "success",
				"data": {
					"transactions": [
						{
							"id": "txn_123",
							"amount": 100.50,
							"type": "credit",
							"status": "completed",
							"date": "2025-04-02T10:00:00Z"
						},
						{
							"id": "txn_456",
							"amount": 50.25,
							"type": "debit",
							"status": "completed",
							"date": "2025-04-02T09:00:00Z"
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name:       "with subaccount",
			subAccount: "sub_123",
			response: `{
				"status": "success",
				"data": {
					"transactions": [
						{
							"id": "txn_789",
							"amount": 75.00,
							"type": "credit",
							"status": "completed",
							"date": "2025-04-02T08:00:00Z",
							"subAccount": "sub_123"
						}
					]
				}
			}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("unexpected method: got %v want POST", r.Method)
				}
				if r.URL.Path != "/wallets/list" {
					t.Errorf("unexpected path: got %v want /wallets/list", r.URL.Path)
				}

				if tt.subAccount != "" {
					var reqBody map[string]string
					if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
						t.Errorf("failed to decode request body: %v", err)
						return
					}

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

			resp, err := client.Wallet.List(context.Background(), tt.subAccount)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWalletTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("GetWalletTransactions() got nil response")
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(resp.Data, &data); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				transactions, ok := data["transactions"].([]interface{})
				if !ok {
					t.Error("transactions field not found or not an array in response")
					return
				}

				for _, txn := range transactions {
					transaction, ok := txn.(map[string]interface{})
					if !ok {
						t.Error("transaction is not an object")
						continue
					}

					if id, ok := transaction["id"].(string); !ok || id == "" {
						t.Error("id field not found or empty in transaction")
					}

					if _, ok := transaction["amount"].(float64); !ok {
						t.Error("amount field not found or not a number in transaction")
					}

					if txnType, ok := transaction["type"].(string); !ok || txnType == "" {
						t.Error("type field not found or empty in transaction")
					}

					if status, ok := transaction["status"].(string); !ok || status == "" {
						t.Error("status field not found or empty in transaction")
					}
				}

				if resp.Status != "success" {
					t.Errorf("unexpected status: got %v want success", resp.Status)
				}
			}
		})
	}
}

func setupTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *testclient.Client) {
	server := httptest.NewServer(handler)
	client := testclient.New(
		testclient.WithTestServer(server.URL),
		testclient.WithAPIKey("test-api-key"),
	)
	return server, client
}
