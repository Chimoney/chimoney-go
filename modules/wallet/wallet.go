package wallet

import (
	"context"
	"encoding/json"
	"errors"
)

var (
	ErrInvalidID       = errors.New("invalid wallet id")
	ErrInvalidReceiver = errors.New("invalid receiver")
)

type Client interface {
	Do(ctx context.Context, method, path string, body interface{}, v interface{}, params map[string]string) error
}

type Wallet struct {
	client Client
}

func New(client Client) *Wallet {
	return &Wallet{client: client}
}

type WalletResponse struct {
	Status  string          `json:"status"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

/**
 * This function lists all wallets
 * @param {string?} subAccount The subAccount to list wallets for
 * @returns The response from the Chimoney API
 */
func (w *Wallet) List(ctx context.Context, subAccount string) (*WalletResponse, error) {
	req := make(map[string]string)
	if subAccount != "" {
		req["subAccount"] = subAccount
	}

	resp := new(WalletResponse)
	err := w.client.Do(ctx, "POST", "/wallets/list", req, resp, nil)
	return resp, err
}

/**
 * This function gets details of a specific wallet
 * @param {string} id The wallet ID
 * @param {string?} subAccount The subAccount the wallet belongs to
 * @returns The response from the Chimoney API
 */
func (w *Wallet) Details(ctx context.Context, id, subAccount string) (*WalletResponse, error) {
	if id == "" {
		return nil, ErrInvalidID
	}

	req := map[string]string{
		"id": id,
	}
	if subAccount != "" {
		req["subAccount"] = subAccount
	}

	resp := new(WalletResponse)
	err := w.client.Do(ctx, "POST", "/wallets/lookup", req, resp, nil)
	return resp, err
}

/**
 * This function transfers funds from one wallet to another
 * @param {string} receiver The receiver's wallet address
 * @param {string} walletType The type of wallet to transfer to
 * @returns The response from the Chimoney API
 */
func (w *Wallet) Transfer(ctx context.Context, receiver, walletType string) (*WalletResponse, error) {
	if receiver == "" {
		return nil, ErrInvalidReceiver
	}

	req := map[string]string{
		"receiver": receiver,
		"wallet":   walletType,
	}

	resp := new(WalletResponse)
	err := w.client.Do(ctx, "POST", "/wallets/transfer", req, resp, nil)
	return resp, err
}

/**
 * This function gets the balance of a wallet
 * @param {string?} subAccount The subAccount to get balance for
 * @returns The response from the Chimoney API
 */
func (w *Wallet) GetBalance(ctx context.Context, subAccount string) (*WalletResponse, error) {
	req := make(map[string]string)
	if subAccount != "" {
		req["subAccount"] = subAccount
	}

	resp := new(WalletResponse)
	err := w.client.Do(ctx, "POST", "/wallet/balance", req, resp, nil)
	return resp, err
}
