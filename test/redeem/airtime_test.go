package redeem_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chimoney/chimoney-go/modules/redeem"
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

func TestAirtimeRedeem(t *testing.T) {
	tests := []struct {
		name    string
		req     *redeem.AirtimeRedeemRequest
		response string
		wantErr bool
	}{
		{
			name: "successful airtime redemption",
			req: &redeem.AirtimeRedeemRequest{
				ChiRef:       "chi_123",
				PhoneNumber:  "+2348123456789",
				CountryToSend: "NG",
				Meta: map[string]interface{}{
					"note": "Test redemption",
				},
			},
			response: `{
				"status": "success",
				"data": {
					"id": "red_123",
					"status": "pending",
					"phoneNumber": "+2348123456789",
					"countryToSend": "NG"
				}
			}`,
			wantErr: false,
		},
		{
			name: "with subaccount",
			req: &redeem.AirtimeRedeemRequest{
				ChiRef:       "chi_123",
				PhoneNumber:  "+2348123456789",
				CountryToSend: "NG",
				SubAccount:   "sub_123",
			},
			response: `{
				"status": "success",
				"data": {
					"id": "red_123",
					"status": "pending",
					"phoneNumber": "+2348123456789",
					"countryToSend": "NG",
					"subAccount": "sub_123"
				}
			}`,
			wantErr: false,
		},
		{
			name: "empty chi ref",
			req: &redeem.AirtimeRedeemRequest{
				ChiRef:       "",
				PhoneNumber:  "+2348123456789",
				CountryToSend: "NG",
			},
			response: "",
			wantErr: true,
		},
		{
			name: "invalid phone number",
			req: &redeem.AirtimeRedeemRequest{
				ChiRef:       "chi_123",
				PhoneNumber:  "invalid",
				CountryToSend: "NG",
			},
			response: `{
				"status": "error",
				"message": "Invalid phone number"
			}`,
			wantErr: true,
		},
		{
			name: "invalid country code",
			req: &redeem.AirtimeRedeemRequest{
				ChiRef:       "chi_123",
				PhoneNumber:  "+2348123456789",
				CountryToSend: "XX",
			},
			response: `{
				"status": "error",
				"message": "Invalid country code"
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
				if r.URL.Path != "/redeem/airtime" {
					t.Errorf("unexpected path: got %v want /redeem/airtime", r.URL.Path)
				}

				var reqBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				if chiRef, ok := reqBody["chiRef"].(string); !ok || chiRef != tt.req.ChiRef {
					t.Errorf("unexpected chiRef: got %v want %v", chiRef, tt.req.ChiRef)
				}

				if phoneNumber, ok := reqBody["phoneNumber"].(string); !ok || phoneNumber != tt.req.PhoneNumber {
					t.Errorf("unexpected phoneNumber: got %v want %v", phoneNumber, tt.req.PhoneNumber)
				}

				if countryToSend, ok := reqBody["countryToSend"].(string); !ok || countryToSend != tt.req.CountryToSend {
					t.Errorf("unexpected countryToSend: got %v want %v", countryToSend, tt.req.CountryToSend)
				}

				if tt.req.SubAccount != "" {
					if subAccount, ok := reqBody["subAccount"].(string); !ok || subAccount != tt.req.SubAccount {
						t.Errorf("unexpected subAccount: got %v want %v", subAccount, tt.req.SubAccount)
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

			resp, err := client.Redeem.Airtime(context.Background(), tt.req)
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

				if id, ok := data["id"].(string); !ok || id == "" {
					t.Error("id field not found or empty in response")
				}

				if status, ok := data["status"].(string); !ok || status == "" {
					t.Error("status field not found or empty in response")
				}

				if phoneNumber, ok := data["phoneNumber"].(string); !ok || phoneNumber != tt.req.PhoneNumber {
					t.Errorf("unexpected phoneNumber in response: got %v want %v", phoneNumber, tt.req.PhoneNumber)
				}
			}
		})
	}
}
