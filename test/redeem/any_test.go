package redeem_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/chimoney/chimoney-go/modules/redeem"
)

func TestAnyRedeem(t *testing.T) {
	tests := []struct {
		name     string
		req      *redeem.AnyRedeemRequest
		response string
		wantErr  bool
	}{
		{
			name: "successful any redemption",
			req: &redeem.AnyRedeemRequest{
				ChiRef: "chi_123",
				RedeemData: []redeem.RedeemDataItem{
					{
						CountryCode:          "US",
						ProductID:            "amazon-us",
						ValueInLocalCurrency: 50.0,
					},
				},
				Meta: map[string]interface{}{
					"note": "Test redemption",
				},
			},
			response: `{
				"status": "success",
				"data": {
					"id": "red_123",
					"status": "pending",
					"redeemData": [
						{
							"countryCode": "US",
							"productId": "amazon-us",
							"valueInLocalCurrency": 50.0
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name: "multiple redeem data items",
			req: &redeem.AnyRedeemRequest{
				ChiRef: "chi_123",
				RedeemData: []redeem.RedeemDataItem{
					{
						CountryCode:          "US",
						ProductID:            "amazon-us",
						ValueInLocalCurrency: 50.0,
					},
					{
						CountryCode:          "GB",
						ProductID:            "amazon-uk",
						ValueInLocalCurrency: 40.0,
					},
				},
			},
			response: `{
				"status": "success",
				"data": {
					"id": "red_123",
					"status": "pending",
					"redeemData": [
						{
							"countryCode": "US",
							"productId": "amazon-us",
							"valueInLocalCurrency": 50.0
						},
						{
							"countryCode": "GB",
							"productId": "amazon-uk",
							"valueInLocalCurrency": 40.0
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name: "with subaccount",
			req: &redeem.AnyRedeemRequest{
				ChiRef: "chi_123",
				RedeemData: []redeem.RedeemDataItem{
					{
						CountryCode:          "US",
						ProductID:            "amazon-us",
						ValueInLocalCurrency: 50.0,
					},
				},
				SubAccount: "sub_123",
			},
			response: `{
				"status": "success",
				"data": {
					"id": "red_123",
					"status": "pending",
					"subAccount": "sub_123",
					"redeemData": [
						{
							"countryCode": "US",
							"productId": "amazon-us",
							"valueInLocalCurrency": 50.0
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name: "empty chi ref",
			req: &redeem.AnyRedeemRequest{
				ChiRef: "",
				RedeemData: []redeem.RedeemDataItem{
					{
						CountryCode:          "US",
						ProductID:            "amazon-us",
						ValueInLocalCurrency: 50.0,
					},
				},
			},
			response: "",
			wantErr:  true,
		},
		{
			name: "empty redeem data",
			req: &redeem.AnyRedeemRequest{
				ChiRef:     "chi_123",
				RedeemData: []redeem.RedeemDataItem{},
			},
			response: "",
			wantErr:  true,
		},
		{
			name: "invalid product id",
			req: &redeem.AnyRedeemRequest{
				ChiRef: "chi_123",
				RedeemData: []redeem.RedeemDataItem{
					{
						CountryCode:          "US",
						ProductID:            "invalid",
						ValueInLocalCurrency: 50.0,
					},
				},
			},
			response: `{
				"status": "error",
				"message": "Invalid product ID"
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
				if r.URL.Path != "/redeem/any" {
					t.Errorf("unexpected path: got %v want /redeem/any", r.URL.Path)
				}

				var reqBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				if chiRef, ok := reqBody["chiRef"].(string); !ok || chiRef != tt.req.ChiRef {
					t.Errorf("unexpected chiRef: got %v want %v", chiRef, tt.req.ChiRef)
				}

				redeemData, ok := reqBody["redeemData"].([]interface{})
				if !ok {
					t.Error("redeemData field not found in request")
					return
				}

				if len(redeemData) != len(tt.req.RedeemData) {
					t.Errorf("unexpected number of redeemData items: got %v want %v", len(redeemData), len(tt.req.RedeemData))
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

			resp, err := client.Redeem.Any(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Any() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("Any() got nil response")
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

				redeemData, ok := data["redeemData"].([]interface{})
				if !ok {
					t.Error("redeemData field not found in response")
					return
				}

				if len(redeemData) != len(tt.req.RedeemData) {
					t.Errorf("unexpected number of redeemData items in response: got %v want %v", len(redeemData), len(tt.req.RedeemData))
				}
			}
		})
	}
}
