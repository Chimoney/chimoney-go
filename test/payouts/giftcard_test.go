package payouts_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/chimoney/chimoney-go/modules/payouts"
)

func TestGiftCard(t *testing.T) {
	tests := []struct {
		name       string
		giftCards  []payouts.GiftCardPayload
		subAccount string
		response   string
		wantErr    bool
	}{
		{
			name: "successful single gift card",
			giftCards: []payouts.GiftCardPayload{
				{
					Email:      "test@example.com",
					ValueInUSD: 50.0,
					RedeemData: struct {
						ProductID            string  `json:"productId"`
						CountryCode         string  `json:"countryCode"`
						ValueInLocalCurrency float64 `json:"valueInLocalCurrency"`
					}{
						ProductID:            "amazon-us",
						CountryCode:         "US",
						ValueInLocalCurrency: 50.0,
					},
				},
			},
			response: `{
				"status": "success",
				"data": {
					"id": "gc_123",
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
			name: "successful multiple gift cards",
			giftCards: []payouts.GiftCardPayload{
				{
					Email:      "test1@example.com",
					ValueInUSD: 50.0,
					RedeemData: struct {
						ProductID            string  `json:"productId"`
						CountryCode         string  `json:"countryCode"`
						ValueInLocalCurrency float64 `json:"valueInLocalCurrency"`
					}{
						ProductID:            "amazon-us",
						CountryCode:         "US",
						ValueInLocalCurrency: 50.0,
					},
				},
				{
					Email:      "test2@example.com",
					ValueInUSD: 25.0,
					RedeemData: struct {
						ProductID            string  `json:"productId"`
						CountryCode         string  `json:"countryCode"`
						ValueInLocalCurrency float64 `json:"valueInLocalCurrency"`
					}{
						ProductID:            "amazon-uk",
						CountryCode:         "GB",
						ValueInLocalCurrency: 20.0,
					},
				},
			},
			response: `{
				"status": "success",
				"data": {
					"id": "gc_123",
					"status": "pending",
					"transactions": [
						{
							"email": "test1@example.com",
							"amount": 50.0,
							"status": "pending"
						},
						{
							"email": "test2@example.com",
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
			giftCards: []payouts.GiftCardPayload{
				{
					Email:      "test@example.com",
					ValueInUSD: 50.0,
					RedeemData: struct {
						ProductID            string  `json:"productId"`
						CountryCode         string  `json:"countryCode"`
						ValueInLocalCurrency float64 `json:"valueInLocalCurrency"`
					}{
						ProductID:            "amazon-us",
						CountryCode:         "US",
						ValueInLocalCurrency: 50.0,
					},
				},
			},
			subAccount: "sub_123",
			response: `{
				"status": "success",
				"data": {
					"id": "gc_123",
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
			name:      "empty gift card list",
			giftCards: []payouts.GiftCardPayload{},
			response:  `{"status":"error","message":"No gift card payouts provided"}`,
			wantErr:   true,
		},
		{
			name: "invalid product id",
			giftCards: []payouts.GiftCardPayload{
				{
					Email:      "test@example.com",
					ValueInUSD: 50.0,
					RedeemData: struct {
						ProductID            string  `json:"productId"`
						CountryCode         string  `json:"countryCode"`
						ValueInLocalCurrency float64 `json:"valueInLocalCurrency"`
					}{
						ProductID:            "invalid",
						CountryCode:         "US",
						ValueInLocalCurrency: 50.0,
					},
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
				if r.URL.Path != "/payouts/gift-card" {
					t.Errorf("unexpected path: got %v want /payouts/gift-card", r.URL.Path)
				}

				var reqBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				giftCards, ok := reqBody["giftCards"].([]interface{})
				if !ok {
					t.Error("giftCards field not found in request")
					return
				}

				if len(giftCards) != len(tt.giftCards) {
					t.Errorf("unexpected number of gift cards: got %v want %v", len(giftCards), len(tt.giftCards))
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

			resp, err := client.Payouts.GiftCard(context.Background(), tt.giftCards, tt.subAccount)
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

				transactions, ok := data["transactions"].([]interface{})
				if !ok {
					t.Error("transactions field not found in response")
					return
				}

				if len(transactions) != len(tt.giftCards) {
					t.Errorf("unexpected number of transactions: got %v want %v", len(transactions), len(tt.giftCards))
				}
			}
		})
	}
}
