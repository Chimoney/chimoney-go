package payouts

import (
	"context"
	"encoding/json"
)

type Client interface {
	Do(ctx context.Context, method, path string, body interface{}, v interface{}, params map[string]string) error
}

type Payouts struct {
	client Client
}

func New(client Client) *Payouts {
	return &Payouts{client: client}
}

type AirtimePayload struct {
	CountryToSend string  `json:"countryToSend"`
	PhoneNumber   string  `json:"phoneNumber"`
	ValueInUSD    float64 `json:"valueInUSD"`
}

type BankPayload struct {
	CountryToSend  string  `json:"countryToSend"`
	AccountBank    string  `json:"account_bank"`
	AccountNumber  string  `json:"account_number"`
	ValueInUSD     float64 `json:"valueInUSD"`
	Reference      string  `json:"reference"`
}

type ChimoneyPayload struct {
	ValueInUSD float64 `json:"valueInUSD"`
	Email      string  `json:"email,omitempty"`
	Twitter    string  `json:"twitter,omitempty"`
}

type GiftCardPayload struct {
	Email      string `json:"email"`
	ValueInUSD float64 `json:"valueInUSD"`
	RedeemData struct {
		ProductID            string  `json:"productId"`
		CountryCode         string  `json:"countryCode"`
		ValueInLocalCurrency float64 `json:"valueInLocalCurrency"`
	} `json:"redeemData"`
}

type PayoutResponse struct {
	Status  string          `json:"status"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

/**
 * This function sends airtime payouts
 * @param {AirtimePayload[]} airtimes Array of airtime payouts
 * @param {string?} subAccount The subAccount for the transaction
 * @returns The response from the Chimoney API
 */
func (p *Payouts) Airtime(ctx context.Context, airtimes []AirtimePayload, subAccount string) (*PayoutResponse, error) {
	req := map[string]interface{}{
		"airtime": airtimes,
	}
	if subAccount != "" {
		req["subAccount"] = subAccount
	}

	resp := new(PayoutResponse)
	err := p.client.Do(ctx, "POST", "/payouts/airtime", req, resp, nil)
	return resp, err
}

/**
 * This function sends bank payouts
 * @param {BankPayload[]} banks Array of bank payouts
 * @param {string?} subAccount The subAccount for the transaction
 * @returns The response from the Chimoney API
 */
func (p *Payouts) Bank(ctx context.Context, banks []BankPayload, subAccount string) (*PayoutResponse, error) {
	req := map[string]interface{}{
		"banks": banks,
	}
	if subAccount != "" {
		req["subAccount"] = subAccount
	}

	resp := new(PayoutResponse)
	err := p.client.Do(ctx, "POST", "/payouts/bank", req, resp, nil)
	return resp, err
}

/**
 * This function sends Chimoney payouts
 * @param {ChimoneyPayload[]} chimoneys Array of Chimoney payouts
 * @param {string?} subAccount The subAccount for the transaction
 * @returns The response from the Chimoney API
 */
func (p *Payouts) Chimoney(ctx context.Context, chimoneys []ChimoneyPayload, subAccount string) (*PayoutResponse, error) {
	req := map[string]interface{}{
		"chimoneys": chimoneys,
	}
	if subAccount != "" {
		req["subAccount"] = subAccount
	}

	resp := new(PayoutResponse)
	err := p.client.Do(ctx, "POST", "/payouts/chimoney", req, resp, nil)
	return resp, err
}

/**
 * This function sends gift card payouts
 * @param {GiftCardPayload[]} giftCards Array of gift card payouts
 * @param {string?} subAccount The subAccount for the transaction
 * @returns The response from the Chimoney API
 */
func (p *Payouts) GiftCard(ctx context.Context, giftCards []GiftCardPayload, subAccount string) (*PayoutResponse, error) {
	req := map[string]interface{}{
		"giftCards": giftCards,
	}
	if subAccount != "" {
		req["subAccount"] = subAccount
	}

	resp := new(PayoutResponse)
	err := p.client.Do(ctx, "POST", "/payouts/gift-card", req, resp, nil)
	return resp, err
}

/**
 * This function gets the status of a payout
 * @param {string} chiRef The ID of the transaction
 * @param {string?} subAccount The subAccount for the transaction
 * @returns The response from the Chimoney API
 */
func (p *Payouts) Status(ctx context.Context, chiRef, subAccount string) (*PayoutResponse, error) {
	req := map[string]string{
		"chiRef": chiRef,
	}
	if subAccount != "" {
		req["subAccount"] = subAccount
	}

	resp := new(PayoutResponse)
	err := p.client.Do(ctx, "POST", "/payouts/status", req, resp, nil)
	return resp, err
}

type CryptoPayment struct {
	XRPL struct {
		Address  string `json:"address"`
		Issuer   string `json:"issuer"`
		Currency string `json:"currency"`
	} `json:"xrpl,omitempty"`
}

/**
 * This function initiates a Chimoney payout with optional crypto payments
 * @param {ChimoneyPayload[]} chimoneys Array of Chimoney payouts
 * @param {boolean} turnOffNotification Whether to turn off notifications
 * @param {CryptoPayment[]} cryptoPayments Optional array of crypto payment details
 * @param {string?} subAccount The subAccount for the transaction
 * @returns The response from the Chimoney API
 */
func (p *Payouts) InitiateChimoney(ctx context.Context, chimoneys []ChimoneyPayload, turnOffNotification bool, cryptoPayments []CryptoPayment, subAccount string) (*PayoutResponse, error) {
	req := map[string]interface{}{
		"chimoneys":          chimoneys,
		"turnOffNotification": turnOffNotification,
	}
	if len(cryptoPayments) > 0 {
		req["crypto_payments"] = cryptoPayments
	}
	if subAccount != "" {
		req["subAccount"] = subAccount
	}

	resp := new(PayoutResponse)
	err := p.client.Do(ctx, "POST", "/payouts/initiate", req, resp, nil)
	return resp, err
}
