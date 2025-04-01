package mobilemoney

import (
	"context"
	"encoding/json"
)

type Client interface {
	Do(ctx context.Context, method, path string, body interface{}, v interface{}, params map[string]string) error
}

type MobileMoney struct {
	client Client
}

func New(client Client) *MobileMoney {
	return &MobileMoney{client: client}
}

type PaymentRequest struct {
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	PhoneNumber string  `json:"phone_number"`
	FullName    string  `json:"fullname"`
	Country     string  `json:"country"`
	Email       string  `json:"email"`
	TxRef       string  `json:"tx_ref"`
	SubAccount  string  `json:"subAccount,omitempty"`
}

type PaymentResponse struct {
	Status  string          `json:"status"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

/**
 * This function initiates a mobile money payment
 * @param {PaymentRequest} req The payment request details
 * @returns The response from the Chimoney API
 */
func (m *MobileMoney) MakePayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	resp := new(PaymentResponse)
	err := m.client.Do(ctx, "POST", "/collections/mobile-money/pay", req, resp, nil)
	return resp, err
}

/**
 * This function verifies a mobile money payment
 * @param {string} id The payment ID to verify
 * @param {string?} subAccount The subAccount of the transaction
 * @returns The response from the Chimoney API
 */
func (m *MobileMoney) VerifyPayment(ctx context.Context, id string, subAccount string) (*PaymentResponse, error) {
	req := map[string]string{
		"id": id,
	}
	if subAccount != "" {
		req["subAccount"] = subAccount
	}
	
	resp := new(PaymentResponse)
	err := m.client.Do(ctx, "POST", "/collections/mobile-money/verify", req, resp, nil)
	return resp, err
}

/**
 * This function gets all mobile money transactions
 * @param {string?} subAccount The subAccount to get transactions for
 * @returns The response from the Chimoney API
 */
func (m *MobileMoney) GetAllTransactions(ctx context.Context, subAccount string) (*PaymentResponse, error) {
	req := make(map[string]string)
	if subAccount != "" {
		req["subAccount"] = subAccount
	}

	resp := new(PaymentResponse)
	err := m.client.Do(ctx, "POST", "/collections/mobile-money/all", req, resp, nil)
	return resp, err
}
