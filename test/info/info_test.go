package info_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestGetAirtimeCountries(t *testing.T) {
	tests := []struct {
		name     string
		response string
		wantErr  bool
	}{
		{
			name: "successful countries fetch",
			response: `{
				"status": "success",
				"data": {
					"countries": ["Nigeria", "USA", "Ghana"]
				}
			}`,
			wantErr: false,
		},
		{
			name:     "error response",
			response: `{"status":"error","message":"Internal server error"}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("unexpected method: got %v want GET", r.Method)
				}
				if r.URL.Path != "/info/airtime-countries" {
					t.Errorf("unexpected path: got %v want /info/airtime-countries", r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusInternalServerError)
				}
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.Info.GetAirtimeCountries(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAirtimeCountries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("GetAirtimeCountries() got nil response")
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(resp.Data, &data); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				countries, ok := data["countries"].([]interface{})
				if !ok {
					t.Error("response data missing countries array")
					return
				}
				if len(countries) == 0 {
					t.Error("expected non-empty countries array")
				}
			}
		})
	}
}

func TestGetBanks(t *testing.T) {
	tests := []struct {
		name        string
		countryCode string
		response    string
		wantErr     bool
	}{
		{
			name:        "successful banks fetch with default country",
			countryCode: "",
			response: "{\"status\": \"success\", \"data\": {\"banks\": [{\"code\": \"001\", \"name\": \"Access Bank\"}]}}",
			wantErr: false,
		},
		{
			name:        "successful banks fetch with specific country",
			countryCode: "GH",
			response: "{\"status\": \"success\", \"data\": {\"banks\": [{\"code\": \"002\", \"name\": \"Ghana Bank\"}]}}",
			wantErr: false,
		},
		{
			name:        "error response",
			countryCode: "XX",
			response:    "{\"status\":\"error\",\"message\":\"Invalid country code\"}",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("unexpected method: got %v want GET", r.Method)
				}
				if r.URL.Path != "/info/country-banks" {
					t.Errorf("unexpected path: got %v want /info/country-banks", r.URL.Path)
				}

				// Check country code parameter
				expectedCountry := tt.countryCode
				if expectedCountry == "" {
					expectedCountry = "NG"
				}
				if got := r.URL.Query().Get("countryCode"); got != expectedCountry {
					t.Errorf("unexpected countryCode: got %v want %v", got, expectedCountry)
				}

				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
				}
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.Info.GetBanks(context.Background(), tt.countryCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBanks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("GetBanks() got nil response")
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(resp.Data, &data); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				banks, ok := data["banks"].([]interface{})
				if !ok {
					t.Error("response data missing banks array")
					return
				}
				if len(banks) == 0 {
					t.Error("expected non-empty banks array")
				}
			}
		})
	}
}

func TestGetLocalAmountInUSD(t *testing.T) {
	tests := []struct {
		name     string
		currency string
		amount   float64
		response string
		wantErr  bool
	}{
		{
			name:     "successful conversion",
			currency: "NGN",
			amount:   5000,
			response: "{\"status\": \"success\", \"data\": {\"amountInUSD\": 10.50}}",
			wantErr:  false,
		},
		{
			name:     "error - invalid currency",
			currency: "",
			amount:   1000,
			response: "{\"status\":\"error\",\"message\":\"Invalid currency\"}",
			wantErr:  true,
		},
		{
			name:     "error - invalid amount",
			currency: "USD",
			amount:   0,
			response: "{\"status\":\"error\",\"message\":\"Invalid amount\"}",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("unexpected method: got %v want POST", r.Method)
				}
				if r.URL.Path != "/info/local-amount-in-usd" {
					t.Errorf("unexpected path: got %v want /info/local-amount-in-usd", r.URL.Path)
				}

				// Check request body
				var reqBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				if got := reqBody["originCurrency"].(string); got != tt.currency {
					t.Errorf("unexpected currency: got %v want %v", got, tt.currency)
				}
				if got := reqBody["amount"].(float64); got != tt.amount {
					t.Errorf("unexpected amount: got %v want %v", got, tt.amount)
				}

				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
				}
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.Info.GetLocalAmountInUSD(context.Background(), tt.currency, tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLocalAmountInUSD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("GetLocalAmountInUSD() got nil response")
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(resp.Data, &data); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				amountInUSD, ok := data["amountInUSD"].(float64)
				if !ok {
					t.Error("response data missing amountInUSD field or invalid type")
					return
				}
				if amountInUSD <= 0 {
					t.Error("expected positive amountInUSD value")
				}
			}
		})
	}
}

func TestGetMobileMoneyCodes(t *testing.T) {
	tests := []struct {
		name     string
		response string
		wantErr  bool
	}{
		{
			name: "successful mobile money codes fetch",
			response: "{\"status\": \"success\", \"data\": {\"codes\": [{\"code\": \"MTN\", \"name\": \"MTN Mobile Money\"}]}}",
			wantErr: false,
		},
		{
			name:     "error response",
			response: "{\"status\":\"error\",\"message\":\"Internal server error\"}",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("unexpected method: got %v want GET", r.Method)
				}
				if r.URL.Path != "/info/mobile-money-codes" {
					t.Errorf("unexpected path: got %v want /info/mobile-money-codes", r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusInternalServerError)
				}
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.Info.GetMobileMoneyCodes(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMobileMoneyCodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("GetMobileMoneyCodes() got nil response")
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(resp.Data, &data); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				codes, ok := data["codes"].([]interface{})
				if !ok {
					t.Error("response data missing codes array")
					return
				}
				if len(codes) == 0 {
					t.Error("expected non-empty codes array")
				}
			}
		})
	}
}

func TestGetUSDInLocalAmount(t *testing.T) {
	tests := []struct {
		name        string
		currency    string
		amountInUSD float64
		response    string
		wantErr     bool
	}{
		{
			name:        "successful conversion",
			currency:    "NGN",
			amountInUSD: 10.50,
			response:    "{\"status\": \"success\", \"data\": {\"localAmount\": 5000}}",
			wantErr:     false,
		},
		{
			name:        "error - invalid currency",
			currency:    "",
			amountInUSD: 100,
			response:    "{\"status\":\"error\",\"message\":\"Invalid currency\"}",
			wantErr:     true,
		},
		{
			name:        "error - invalid amount",
			currency:    "USD",
			amountInUSD: 0,
			response:    "{\"status\":\"error\",\"message\":\"Invalid amount\"}",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("unexpected method: got %v want POST", r.Method)
				}
				if r.URL.Path != "/info/usd-in-local-amount" {
					t.Errorf("unexpected path: got %v want /info/usd-in-local-amount", r.URL.Path)
				}

				// Check request body
				var reqBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				if got := reqBody["destinationCurrency"].(string); got != tt.currency {
					t.Errorf("unexpected currency: got %v want %v", got, tt.currency)
				}
				if got := reqBody["amountInUSD"].(float64); got != tt.amountInUSD {
					t.Errorf("unexpected amountInUSD: got %v want %v", got, tt.amountInUSD)
				}

				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
				}
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.Info.GetUSDInLocalAmount(context.Background(), tt.currency, tt.amountInUSD)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUSDInLocalAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("GetUSDInLocalAmount() got nil response")
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(resp.Data, &data); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				localAmount, ok := data["localAmount"].(float64)
				if !ok {
					t.Error("response data missing localAmount field or invalid type")
					return
				}
				if localAmount <= 0 {
					t.Error("expected positive localAmount value")
				}
			}
		})
	}
}

func TestGetSupportedAssets(t *testing.T) {
	tests := []struct {
		name     string
		response string
		wantErr  bool
	}{
		{
			name: "successful assets fetch",
			response: `{
				"status": "success",
				"data": {
					"crypto": ["BTC", "ETH", "USDT"]
				}
			}`,
			wantErr: false,
		},
		{
			name:     "error response",
			response: `{"status":"error","message":"Internal server error"}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("unexpected method: got %v want GET", r.Method)
				}
				if r.URL.Path != "/info/assets" {
					t.Errorf("unexpected path: got %v want /info/assets", r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusInternalServerError)
				}
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.Info.GetSupportedAssets(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSupportedAssets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("GetSupportedAssets() got nil response")
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(resp.Data, &data); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				crypto, ok := data["crypto"].([]interface{})
				if !ok {
					t.Error("response data missing crypto array")
					return
				}
				if len(crypto) == 0 {
					t.Error("expected non-empty crypto array")
				}
			}
		})
	}
}
