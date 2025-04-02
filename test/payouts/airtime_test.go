package payouts_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chimoney/chimoney-go/modules/payouts"
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

func TestAirtime(t *testing.T) {
	tests := []struct {
		name       string
		airtimes   []payouts.AirtimePayload
		subAccount string
		response   string
		wantErr    bool
	}{
		{
			name: "successful single airtime",
			airtimes: []payouts.AirtimePayload{
				{
					CountryToSend: "NG",
					PhoneNumber:   "+2348123456789",
					ValueInUSD:    10.0,
				},
			},
			response: `{
				"status": "success",
				"data": {
					"id": "air_123",
					"status": "pending",
					"transactions": [
						{
							"phoneNumber": "+2348123456789",
							"amount": 10.0,
							"status": "pending"
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name: "successful multiple airtimes",
			airtimes: []payouts.AirtimePayload{
				{
					CountryToSend: "NG",
					PhoneNumber:   "+2348123456789",
					ValueInUSD:    10.0,
				},
				{
					CountryToSend: "GH",
					PhoneNumber:   "+233123456789",
					ValueInUSD:    5.0,
				},
			},
			response: `{
				"status": "success",
				"data": {
					"id": "air_123",
					"status": "pending",
					"transactions": [
						{
							"phoneNumber": "+2348123456789",
							"amount": 10.0,
							"status": "pending"
						},
						{
							"phoneNumber": "+233123456789",
							"amount": 5.0,
							"status": "pending"
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name: "with subaccount",
			airtimes: []payouts.AirtimePayload{
				{
					CountryToSend: "NG",
					PhoneNumber:   "+2348123456789",
					ValueInUSD:    10.0,
				},
			},
			subAccount: "sub_123",
			response: `{
				"status": "success",
				"data": {
					"id": "air_123",
					"status": "pending",
					"subAccount": "sub_123",
					"transactions": [
						{
							"phoneNumber": "+2348123456789",
							"amount": 10.0,
							"status": "pending"
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name:     "empty airtime list",
			airtimes: []payouts.AirtimePayload{},
			response: `{
				"status": "error",
				"message": "No airtime payouts provided"
			}`,
			wantErr: true,
		},
		{
			name: "invalid phone number",
			airtimes: []payouts.AirtimePayload{
				{
					CountryToSend: "NG",
					PhoneNumber:   "invalid",
					ValueInUSD:    10.0,
				},
			},
			response: `{
				"status": "error",
				"message": "Invalid phone number format"
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
				if r.URL.Path != "/payouts/airtime" {
					t.Errorf("unexpected path: got %v want /payouts/airtime", r.URL.Path)
				}

				var reqBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				airtimes, ok := reqBody["airtime"].([]interface{})
				if !ok {
					t.Error("airtime field not found in request")
					return
				}

				if len(airtimes) != len(tt.airtimes) {
					t.Errorf("unexpected number of airtimes: got %v want %v", len(airtimes), len(tt.airtimes))
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

			resp, err := client.Payouts.Airtime(context.Background(), tt.airtimes, tt.subAccount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Airtime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("Airtime() got nil response")
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

				if len(transactions) != len(tt.airtimes) {
					t.Errorf("unexpected number of transactions: got %v want %v", len(transactions), len(tt.airtimes))
				}
			}
		})
	}
}
