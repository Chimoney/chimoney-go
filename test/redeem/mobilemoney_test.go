package redeem_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/chimoney/chimoney-go/modules/redeem"
)

func TestMobileMoneyRedeem(t *testing.T) {
	tests := []struct {
		name     string
		req      *redeem.MobileMoneyRedeemRequest
		response string
		wantErr  bool
	}{
		{
			name: "successful mobile money redemption",
			req: &redeem.MobileMoneyRedeemRequest{
				ChiRef: "chi_123",
				RedeemOptions: map[string]interface{}{
					"phoneNumber":  "+2348123456789",
					"countryCode": "NG",
					"amount":      50.0,
					"provider":    "mtn",
				},
			},
			response: `{
				"status": "success",
				"data": {
					"id": "red_123",
					"status": "pending",
					"transaction": {
						"phoneNumber": "+2348123456789",
						"amount": 50.0,
						"provider": "mtn",
						"status": "pending"
					}
				}
			}`,
			wantErr: false,
		},
		{
			name: "with subaccount",
			req: &redeem.MobileMoneyRedeemRequest{
				ChiRef: "chi_123",
				RedeemOptions: map[string]interface{}{
					"phoneNumber":  "+2348123456789",
					"countryCode": "NG",
					"amount":      50.0,
					"provider":    "mtn",
				},
				SubAccount: "sub_123",
			},
			response: `{
				"status": "success",
				"data": {
					"id": "red_123",
					"status": "pending",
					"subAccount": "sub_123",
					"transaction": {
						"phoneNumber": "+2348123456789",
						"amount": 50.0,
						"provider": "mtn",
						"status": "pending"
					}
				}
			}`,
			wantErr: false,
		},
		{
			name: "empty chi ref",
			req: &redeem.MobileMoneyRedeemRequest{
				ChiRef: "",
				RedeemOptions: map[string]interface{}{
					"phoneNumber":  "+2348123456789",
					"countryCode": "NG",
					"amount":      50.0,
					"provider":    "mtn",
				},
			},
			response: "",
			wantErr:  true,
		},
		{
			name: "empty redeem options",
			req: &redeem.MobileMoneyRedeemRequest{
				ChiRef:        "chi_123",
				RedeemOptions: map[string]interface{}{},
			},
			response: "",
			wantErr:  true,
		},
		{
			name: "invalid phone number",
			req: &redeem.MobileMoneyRedeemRequest{
				ChiRef: "chi_123",
				RedeemOptions: map[string]interface{}{
					"phoneNumber":  "invalid",
					"countryCode": "NG",
					"amount":      50.0,
					"provider":    "mtn",
				},
			},
			response: `{"status":"error","message":"Invalid phone number"}`,
			wantErr:  true,
		},
		{
			name: "invalid provider",
			req: &redeem.MobileMoneyRedeemRequest{
				ChiRef: "chi_123",
				RedeemOptions: map[string]interface{}{
					"phoneNumber":  "+2348123456789",
					"countryCode": "NG",
					"amount":      50.0,
					"provider":    "invalid",
				},
			},
			response: `{"status":"error","message":"Invalid provider"}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("unexpected method: got %v want POST", r.Method)
				}
				if r.URL.Path != "/redeem/mobile-money" {
					t.Errorf("unexpected path: got %v want /redeem/mobile-money", r.URL.Path)
				}

				var reqBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				if chiRef, ok := reqBody["chiRef"].(string); !ok || chiRef != tt.req.ChiRef {
					t.Errorf("unexpected chiRef: got %v want %v", chiRef, tt.req.ChiRef)
				}

				redeemOptions, ok := reqBody["redeemOptions"].(map[string]interface{})
				if !ok {
					t.Error("redeemOptions field not found in request")
					return
				}

				if len(redeemOptions) != len(tt.req.RedeemOptions) {
					t.Errorf("unexpected number of redeemOptions: got %v want %v", len(redeemOptions), len(tt.req.RedeemOptions))
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

			resp, err := client.Redeem.MobileMoney(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("MobileMoney() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("MobileMoney() got nil response")
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

				transaction, ok := data["transaction"].(map[string]interface{})
				if !ok {
					t.Error("transaction field not found in response")
					return
				}

				if phoneNumber, ok := transaction["phoneNumber"].(string); !ok || phoneNumber == "" {
					t.Error("phoneNumber field not found or empty in transaction")
				}

				if resp.Status != "success" {
					t.Errorf("unexpected status: got %v want success", resp.Status)
				}
			}
		})
	}
}
