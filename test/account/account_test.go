package account_test

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

func TestGetTransactionsByIssueID(t *testing.T) {
	tests := []struct {
		name       string
		issueID    string
		subAccount string
		response   string
		wantErr    bool
	}{
		{
			name:       "successful transactions fetch",
			issueID:    "issue123",
			subAccount: "",
			response:   `{"status":"success","message":"Transactions retrieved successfully","data":[{"id":"tx1","issueID":"issue123"}]}`,
			wantErr:    false,
		},
		{
			name:       "with subaccount",
			issueID:    "issue123",
			subAccount: "sub123",
			response:   `{"status":"success","message":"Transactions retrieved successfully","data":[{"id":"tx1","issueID":"issue123"}]}`,
			wantErr:    false,
		},
		{
			name:       "invalid issue ID",
			issueID:    "",
			subAccount: "",
			response:   `{"status":"error","message":"Invalid issue ID"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != "POST" {
					t.Errorf("unexpected method: got %v want POST", r.Method)
				}
				if r.URL.Path != "/accounts/issue-id-transactions" {
					t.Errorf("unexpected path: got %v want /accounts/issue-id-transactions", r.URL.Path)
				}

				// Check query parameters
				if got := r.URL.Query().Get("issueID"); got != tt.issueID {
					t.Errorf("unexpected issueID: got %v want %v", got, tt.issueID)
				}

				// Check request body for subAccount
				if tt.subAccount != "" {
					var reqBody map[string]string
					if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
						t.Errorf("failed to decode request body: %v", err)
						return
					}
					if got := reqBody["subAccount"]; got != tt.subAccount {
						t.Errorf("unexpected subAccount: got %v want %v", got, tt.subAccount)
					}
				}

				// Send response
				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
				}
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.Account.GetTransactionsByIssueID(context.Background(), tt.issueID, tt.subAccount)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionsByIssueID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("GetTransactionsByIssueID() got nil response")
					return
				}

				var transactions []map[string]interface{}
				if err := json.Unmarshal(resp.Data, &transactions); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				if len(transactions) == 0 {
					t.Error("expected non-empty transactions array")
				}
				for _, tx := range transactions {
					if tx["issueID"] != tt.issueID {
						t.Errorf("unexpected issueID in transaction: got %v want %v", tx["issueID"], tt.issueID)
					}
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
			name:       "successful empty transactions",
			subAccount: "",
			response:   `{"status":"success","message":"Transactions retrieved successfully","data":[]}`,
			wantErr:    false,
		},
		{
			name:       "with transactions",
			subAccount: "",
			response:   `{"status":"success","message":"Transactions retrieved successfully","data":[{"id":"1","amount":100}]}`,
			wantErr:    false,
		},
		{
			name:       "with subaccount",
			subAccount: "test-subaccount",
			response:   `{"status":"success","message":"Transactions retrieved successfully","data":[]}`,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != "POST" {
					t.Errorf("unexpected method: got %v want POST", r.Method)
				}
				if r.URL.Path != "/accounts/transactions" {
					t.Errorf("unexpected path: got %v want /accounts/transactions", r.URL.Path)
				}

				// Send response
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.Account.GetAllTransactions(context.Background(), tt.subAccount)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if resp == nil {
				t.Error("GetAllTransactions() got nil response")
			}

			var transactions []interface{}
			if err := json.Unmarshal(resp.Data, &transactions); err != nil {
				t.Errorf("failed to unmarshal response data: %v", err)
			}
		})
	}
}

func TestGetTransactionByID(t *testing.T) {
	tests := []struct {
		name          string
		transactionID string
		subAccount    string
		response      string
		wantErr       bool
	}{
		{
			name:          "successful transaction fetch",
			transactionID: "tx123",
			subAccount:    "",
			response:      `{"status":"success","message":"Transaction retrieved successfully","data":{"id":"tx123","amount":100}}`,
			wantErr:       false,
		},
		{
			name:          "with subaccount",
			transactionID: "tx123",
			subAccount:    "sub123",
			response:      `{"status":"success","message":"Transaction retrieved successfully","data":{"id":"tx123","amount":100}}`,
			wantErr:       false,
		},
		{
			name:          "invalid transaction ID",
			transactionID: "",
			subAccount:    "",
			response:      `{"status":"error","message":"Invalid transaction ID"}`,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != "POST" {
					t.Errorf("unexpected method: got %v want POST", r.Method)
				}
				if r.URL.Path != "/accounts/transaction" {
					t.Errorf("unexpected path: got %v want /accounts/transaction", r.URL.Path)
				}

				// Send response
				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
				}
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.Account.GetTransactionByID(context.Background(), tt.transactionID, tt.subAccount)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && resp == nil {
				t.Error("GetTransactionByID() got nil response")
			}

			if !tt.wantErr {
				var txData map[string]interface{}
				if err := json.Unmarshal(resp.Data, &txData); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
				}
				if txData["id"] != tt.transactionID {
					t.Errorf("unexpected transaction ID: got %v want %v", txData["id"], tt.transactionID)
				}
			}
		})
	}
}

func TestTransfer(t *testing.T) {
	tests := []struct {
		name       string
		chiRef     string
		subAccount string
		response   string
		wantErr    bool
	}{
		{
			name:       "successful transfer",
			chiRef:     "chi_123",
			subAccount: "",
			response:   `{"status":"success","message":"Transfer successful","data":{"id":"tx1","chiRef":"chi_123","status":"completed"}}`,
			wantErr:    false,
		},
		{
			name:       "with subaccount",
			chiRef:     "chi_123",
			subAccount: "sub123",
			response:   `{"status":"success","message":"Transfer successful","data":{"id":"tx1","chiRef":"chi_123","status":"completed"}}`,
			wantErr:    false,
		},
		{
			name:       "invalid chi ref",
			chiRef:     "",
			subAccount: "",
			response:   `{"status":"error","message":"Invalid chi ref"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != "POST" {
					t.Errorf("unexpected method: got %v want POST", r.Method)
				}
				if r.URL.Path != "/accounts/transfer" {
					t.Errorf("unexpected path: got %v want /accounts/transfer", r.URL.Path)
				}

				// Check request body
				var reqBody map[string]string
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				if got := reqBody["chiRef"]; got != tt.chiRef {
					t.Errorf("unexpected chiRef: got %v want %v", got, tt.chiRef)
				}

				if tt.subAccount != "" {
					if got := reqBody["subAccount"]; got != tt.subAccount {
						t.Errorf("unexpected subAccount: got %v want %v", got, tt.subAccount)
					}
				}

				// Send response
				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
				}
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.Account.Transfer(context.Background(), tt.chiRef, tt.subAccount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("Transfer() got nil response")
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(resp.Data, &data); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				if got := data["chiRef"]; got != tt.chiRef {
					t.Errorf("unexpected chiRef in response: got %v want %v", got, tt.chiRef)
				}
			}
		})
	}
}

func TestDeleteUnpaidTransaction(t *testing.T) {
	tests := []struct {
		name       string
		chiRef     string
		subAccount string
		response   string
		wantErr    bool
	}{
		{
			name:       "successful deletion",
			chiRef:     "chi_123",
			subAccount: "",
			response:   `{"status":"success","message":"Transaction deleted successfully","data":{"chiRef":"chi_123"}}`,
			wantErr:    false,
		},
		{
			name:       "with subaccount",
			chiRef:     "chi_123",
			subAccount: "sub123",
			response:   `{"status":"success","message":"Transaction deleted successfully","data":{"chiRef":"chi_123"}}`,
			wantErr:    false,
		},
		{
			name:       "invalid chi ref",
			chiRef:     "",
			subAccount: "",
			response:   `{"status":"error","message":"Invalid chi ref"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != "DELETE" {
					t.Errorf("unexpected method: got %v want DELETE", r.Method)
				}
				if r.URL.Path != "/accounts/delete-unpaid" {
					t.Errorf("unexpected path: got %v want /accounts/delete-unpaid", r.URL.Path)
				}

				// Check request body
				var reqBody map[string]string
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
					return
				}

				if got := reqBody["chiRef"]; got != tt.chiRef {
					t.Errorf("unexpected chiRef: got %v want %v", got, tt.chiRef)
				}

				if tt.subAccount != "" {
					if got := reqBody["subAccount"]; got != tt.subAccount {
						t.Errorf("unexpected subAccount: got %v want %v", got, tt.subAccount)
					}
				}

				// Send response
				w.Header().Set("Content-Type", "application/json")
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
				}
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			resp, err := client.Account.DeleteUnpaidTransaction(context.Background(), tt.chiRef, tt.subAccount)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUnpaidTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("DeleteUnpaidTransaction() got nil response")
					return
				}

				var data map[string]interface{}
				if err := json.Unmarshal(resp.Data, &data); err != nil {
					t.Errorf("failed to unmarshal response data: %v", err)
					return
				}

				if got := data["chiRef"]; got != tt.chiRef {
					t.Errorf("unexpected chiRef in response: got %v want %v", got, tt.chiRef)
				}
			}
		})
	}
}