package payouts_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/chimoney/chimoney-go/modules/payouts"
)

func TestBank(t *testing.T) {
	tests := []struct {
		name       string
		banks      []payouts.BankPayload
		subAccount string
		response   string
		wantErr    bool
	}{
		{
			name: "successful single bank transfer",
			banks: []payouts.BankPayload{
				{
					CountryToSend: "NG",
					AccountBank:   "044",
					AccountNumber: "1234567890",
					ValueInUSD:    100.0,
					Reference:     "ref_123",
				},
			},
			response: `{
				"status": "success",
				"data": {
					"id": "bank_123",
					"status": "pending",
					"transactions": [
						{
							"accountNumber": "1234567890",
							"amount": 100.0,
							"status": "pending"
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name: "successful multiple bank transfers",
			banks: []payouts.BankPayload{
				{
					CountryToSend: "NG",
					AccountBank:   "044",
					AccountNumber: "1234567890",
					ValueInUSD:    100.0,
					Reference:     "ref_123",
				},
				{
					CountryToSend: "GH",
					AccountBank:   "GH001",
					AccountNumber: "0987654321",
					ValueInUSD:    50.0,
					Reference:     "ref_456",
				},
			},
			response: `{
				"status": "success",
				"data": {
					"id": "bank_123",
					"status": "pending",
					"transactions": [
						{
							"accountNumber": "1234567890",
							"amount": 100.0,
							"status": "pending"
						},
						{
							"accountNumber": "0987654321",
							"amount": 50.0,
							"status": "pending"
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name: "with subaccount",
			banks: []payouts.BankPayload{
				{
					CountryToSend: "NG",
					AccountBank:   "044",
					AccountNumber: "1234567890",
					ValueInUSD:    100.0,
					Reference:     "ref_123",
				},
			},
			subAccount: "sub_123",
			response: `{
				"status": "success",
				"data": {
					"id": "bank_123",
					"status": "pending",
					"subAccount": "sub_123",
					"transactions": [
						{
							"accountNumber": "1234567890",
							"amount": 100.0,
							"status": "pending"
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name:     "empty bank list",
			banks:    []payouts.BankPayload{},
			response: `{"status":"error","message":"No bank payouts provided"}`,
			wantErr:  true,
		},
		{
			name: "invalid account number",
			banks: []payouts.BankPayload{
				{
					CountryToSend: "NG",
					AccountBank:   "044",
					AccountNumber: "invalid",
					ValueInUSD:    100.0,
					Reference:     "ref_123",
				},
			},
			response: `{"status":"error","message":"Invalid account number"}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("unexpected method: got %v want POST", r.Method)
				}
				if r.URL.Path != "/payouts/bank" {
					t.Errorf("unexpected path: got %v want /payouts/bank", r.URL.Path)
				}

				var reqBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				banks, ok := reqBody["banks"].([]interface{})
				if !ok {
					t.Error("banks field not found in request")
					return
				}

				if len(banks) != len(tt.banks) {
					t.Errorf("unexpected number of banks: got %v want %v", len(banks), len(tt.banks))
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

			resp, err := client.Payouts.Bank(context.Background(), tt.banks, tt.subAccount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bank() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("Bank() got nil response")
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

				if len(transactions) != len(tt.banks) {
					t.Errorf("unexpected number of transactions: got %v want %v", len(transactions), len(tt.banks))
				}
			}
		})
	}
}
