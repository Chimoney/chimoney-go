package redeem_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/chimoney/chimoney-go/modules/redeem"
)

func TestGiftCardRedeem(t *testing.T) {
	tests := []struct {
		name     string
		req      *redeem.GiftCardRedeemRequest
		response string
		wantErr  bool
	}{
		{
			name: "successful gift card redemption",
			req: &redeem.GiftCardRedeemRequest{
				ChiRef: "chi_123",
				RedeemOptions: map[string]interface{}{
					"email":      "test@example.com",
					"productId":  "amazon-us",
					"amount":     50.0,
					"countryCode": "US",
				},
			},
			response: `{
				"status": "success",
				"data": {
					"id": "red_123",
					"status": "pending",
					"giftCard": {
						"code": "XXXX-XXXX-XXXX",
						"pin": "1234",
						"amount": 50.0,
						"email": "test@example.com"
					}
				}
			}`,
			wantErr: false,
		},
		{
			name: "with subaccount",
			req: &redeem.GiftCardRedeemRequest{
				ChiRef: "chi_123",
				RedeemOptions: map[string]interface{}{
					"email":      "test@example.com",
					"productId":  "amazon-us",
					"amount":     50.0,
					"countryCode": "US",
				},
				SubAccount: "sub_123",
			},
			response: `{
				"status": "success",
				"data": {
					"id": "red_123",
					"status": "pending",
					"subAccount": "sub_123",
					"giftCard": {
						"code": "XXXX-XXXX-XXXX",
						"pin": "1234",
						"amount": 50.0,
						"email": "test@example.com"
					}
				}
			}`,
			wantErr: false,
		},
		{
			name: "empty chi ref",
			req: &redeem.GiftCardRedeemRequest{
				ChiRef: "",
				RedeemOptions: map[string]interface{}{
					"email":      "test@example.com",
					"productId":  "amazon-us",
					"amount":     50.0,
					"countryCode": "US",
				},
			},
			response: "",
			wantErr:  true,
		},
		{
			name: "empty redeem options",
			req: &redeem.GiftCardRedeemRequest{
				ChiRef:        "chi_123",
				RedeemOptions: map[string]interface{}{},
			},
			response: "",
			wantErr:  true,
		},
		{
			name: "invalid product id",
			req: &redeem.GiftCardRedeemRequest{
				ChiRef: "chi_123",
				RedeemOptions: map[string]interface{}{
					"email":      "test@example.com",
					"productId":  "invalid",
					"amount":     50.0,
					"countryCode": "US",
				},
			},
			response: `{"status":"error","message":"Invalid product ID"}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("unexpected method: got %v want POST", r.Method)
				}
				if r.URL.Path != "/redeem/gift-card" {
					t.Errorf("unexpected path: got %v want /redeem/gift-card", r.URL.Path)
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

			resp, err := client.Redeem.GiftCard(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GiftCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("GiftCard() got nil response")
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

				giftCard, ok := data["giftCard"].(map[string]interface{})
				if !ok {
					t.Error("giftCard field not found in response")
					return
				}

				if code, ok := giftCard["code"].(string); !ok || code == "" {
					t.Error("code field not found or empty in giftCard")
				}

				if resp.Status != "success" {
					t.Errorf("unexpected status: got %v want success", resp.Status)
				}
			}
		})
	}
}
