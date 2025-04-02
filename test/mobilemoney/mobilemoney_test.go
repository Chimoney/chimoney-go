package mobilemoney_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chimoney/chimoney-go/modules/mobilemoney"
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

func TestMakePayment(t *testing.T) {
	tests := []struct {
		name     string
		request  *mobilemoney.PaymentRequest
		response string
		wantErr  bool
	}{
		{
			name: "successful payment",
			request: &mobilemoney.PaymentRequest{
				Amount:      100.50,
				Currency:    "GHS",
				PhoneNumber: "+233123456789",
				FullName:    "John Doe",
				Country:     "GH",
				Email:       "john@example.com",
				TxRef:      "tx_123",
			},
			response: `{
				"status": "success",
				"data": {
					"id": "mm_123",
					"amount": 100.50,
					"currency": "GHS",
					"status": "pending"
				}
			}`,
			wantErr: false,
		},
		{
			name: "with subaccount",
			request: &mobilemoney.PaymentRequest{
				Amount:      50.00,
				Currency:    "GHS",
				PhoneNumber: "+233123456789",
				FullName:    "Jane Doe",
				Country:     "GH",
				Email:       "jane@example.com",
				TxRef:      "tx_456",
				SubAccount: "sub_123",
			},
			response: `{
				"status": "success",
				"data": {
					"id": "mm_456",
					"amount": 50.00,
					"currency": "GHS",
					"status": "pending",
					"subAccount": "sub_123"
				}
			}`,
			wantErr: false,
		},
		{
			name: "invalid request",
			request: &mobilemoney.PaymentRequest{
				Amount:      -100,
				Currency:    "INVALID",
				PhoneNumber: "invalid",
			},
			response: `{
				"status": "error",
				"message": "Invalid payment request"
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
				if r.URL.Path != "/collections/mobile-money/pay" {
					t.Errorf("unexpected path: got %v want /collections/mobile-money/pay", r.URL.Path)
				}

				var reqBody mobilemoney.PaymentRequest
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				if reqBody.Amount != tt.request.Amount {
					t.Errorf("unexpected amount: got %v want %v", reqBody.Amount, tt.request.Amount)
				}

				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
				}
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.MobileMoney.MakePayment(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakePayment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("MakePayment() got nil response")
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(resp.Data, &data); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				if data["currency"] != tt.request.Currency {
					t.Errorf("unexpected currency in response: got %v want %v", data["currency"], tt.request.Currency)
				}
			}
		})
	}
}

func TestVerifyPayment(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		subAccount string
		response   string
		wantErr    bool
	}{
		{
			name: "successful verification",
			id:   "mm_123",
			response: `{
				"status": "success",
				"data": {
					"id": "mm_123",
					"status": "completed"
				}
			}`,
			wantErr: false,
		},
		{
			name:       "with subaccount",
			id:         "mm_456",
			subAccount: "sub_123",
			response: `{
				"status": "success",
				"data": {
					"id": "mm_456",
					"status": "completed",
					"subAccount": "sub_123"
				}
			}`,
			wantErr: false,
		},
		{
			name: "payment not found",
			id:   "mm_999",
			response: `{
				"status": "error",
				"message": "Payment not found"
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
				if r.URL.Path != "/collections/mobile-money/verify" {
					t.Errorf("unexpected path: got %v want /collections/mobile-money/verify", r.URL.Path)
				}

				var reqBody map[string]string
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				if reqBody["id"] != tt.id {
					t.Errorf("unexpected id: got %v want %v", reqBody["id"], tt.id)
				}
				if reqBody["subAccount"] != tt.subAccount {
					t.Errorf("unexpected subAccount: got %v want %v", reqBody["subAccount"], tt.subAccount)
				}

				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusNotFound)
				}
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.MobileMoney.VerifyPayment(context.Background(), tt.id, tt.subAccount)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPayment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("VerifyPayment() got nil response")
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(resp.Data, &data); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				if data["id"] != tt.id {
					t.Errorf("unexpected id in response: got %v want %v", data["id"], tt.id)
				}
			}
		})
	}
}

func TestGetAllTransactions(t *testing.T) {
	tests := []struct {
		name       string
		subAccount string
		response   string
		wantErr    bool
	}{
		{
			name: "successful list",
			response: `{
				"status": "success",
				"data": [
					{
						"id": "mm_123",
						"amount": 100.50,
						"status": "completed"
					},
					{
						"id": "mm_456",
						"amount": 50.00,
						"status": "pending"
					}
				]
			}`,
			wantErr: false,
		},
		{
			name:       "with subaccount",
			subAccount: "sub_123",
			response: `{
				"status": "success",
				"data": [
					{
						"id": "mm_789",
						"amount": 75.25,
						"status": "completed",
						"subAccount": "sub_123"
					}
				]
			}`,
			wantErr: false,
		},
		{
			name: "empty list",
			response: `{
				"status": "success",
				"data": []
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
				if r.URL.Path != "/collections/mobile-money/all" {
					t.Errorf("unexpected path: got %v want /collections/mobile-money/all", r.URL.Path)
				}

				var reqBody map[string]string
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				if reqBody["subAccount"] != tt.subAccount {
					t.Errorf("unexpected subAccount: got %v want %v", reqBody["subAccount"], tt.subAccount)
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.MobileMoney.GetAllTransactions(context.Background(), tt.subAccount)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("GetAllTransactions() got nil response")
					return
				}

				var transactions []interface{}
				if err := json.Unmarshal(resp.Data, &transactions); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				if tt.name != "empty list" && len(transactions) == 0 {
					t.Error("expected non-empty transactions array")
				}
			}
		})
	}
}
