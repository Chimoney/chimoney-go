package account

import (
	"context"
	"encoding/json"
	"errors"
)

var (
	ErrInvalidIssueID = errors.New("invalid issue ID")
	ErrInvalidTransactionID = errors.New("invalid transaction ID")
	ErrInvalidChiRef = errors.New("invalid chi ref")
)

type Client interface {
	Do(ctx context.Context, method, path string, body interface{}, v interface{}, params map[string]string) error
}

type Account struct {
	client Client
}

func New(client Client) *Account {
	return &Account{client: client}
}

type AccountResponse struct {
	ID        string          `json:"id"`
	Status    string          `json:"status"`
	Data      json.RawMessage `json:"data"`
	Message   string          `json:"message"`
	Timestamp int64          `json:"timestamp"`
}


/**
 * This function gets all transactions by issue ID
 * @param {string} issueID The ID of the issue
 * @param {string?} subAccount The subAccount of the transaction
 * @returns The response from the Chimoney API
 */
func (a *Account) GetTransactionsByIssueID(ctx context.Context, issueID string, subAccount string) (*AccountResponse, error) {
	if issueID == "" {
		return nil, ErrInvalidIssueID
	}

	req := map[string]string{}
	if subAccount != "" {
		req["subAccount"] = subAccount
	}

	resp := new(AccountResponse)
	params := map[string]string{
		"issueID": issueID,
	}
	err := a.client.Do(ctx, "POST", "/accounts/issue-id-transactions", req, resp, params)
	return resp, err
}

/**
 * This function gets all transactions
 * @param {string?} subAccount The subAccount of the transaction
 * @returns The response from the Chimoney API
 */
func (a *Account) GetAllTransactions(ctx context.Context, subAccount string) (*AccountResponse, error) {
	req := map[string]string{}
	if subAccount != "" {
		req["subAccount"] = subAccount
	}

	resp := new(AccountResponse)
	err := a.client.Do(ctx, "POST", "/accounts/transactions", req, resp, nil)
	return resp, err
}

/**
 * This function transfers a transaction
 * @param {string} chiRef The ID of the transaction
 * @param {string?} subAccount The subAccount of the transaction
 * @returns The response from the Chimoney API
 */
func (a *Account) Transfer(ctx context.Context, chiRef string, subAccount string) (*AccountResponse, error) {
	if chiRef == "" {
		return nil, ErrInvalidChiRef
	}

	req := map[string]string{
		"chiRef": chiRef,
	}
	if subAccount != "" {
		req["subAccount"] = subAccount
	}

	resp := new(AccountResponse)
	err := a.client.Do(ctx, "POST", "/accounts/transfer", req, resp, nil)
	return resp, err
}

/**
 * This function deletes an unpaid transaction
 * @param {string} chiRef The ID of the transaction
 * @param {string?} subAccount The subAccount of the transaction
 * @returns The response from the Chimoney API
 */
func (a *Account) DeleteUnpaidTransaction(ctx context.Context, chiRef string, subAccount string) (*AccountResponse, error) {
	if chiRef == "" {
		return nil, ErrInvalidChiRef
	}

	req := map[string]string{
		"chiRef": chiRef,
	}
	if subAccount != "" {
		req["subAccount"] = subAccount
	}

	resp := new(AccountResponse)
	err := a.client.Do(ctx, "DELETE", "/accounts/delete-unpaid", req, resp, nil)
	return resp, err
}

func (a *Account) GetTransactionByID(ctx context.Context, transactionID string, subAccount string) (*AccountResponse, error) {
	if transactionID == "" {
		return nil, ErrInvalidTransactionID
	}

	req := map[string]string{
		"id": transactionID,
	}
	if subAccount != "" {
		req["subAccount"] = subAccount
	}

	resp := new(AccountResponse)
	err := a.client.Do(ctx, "POST", "/accounts/transaction", req, resp, nil)
	return resp, err
}
