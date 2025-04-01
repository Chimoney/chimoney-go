package redeem

import (
	"context"
	"encoding/json"
	"errors"
)

var (
	ErrInvalidChiRef = errors.New("invalid chi ref")
	ErrInvalidRedeemData = errors.New("invalid redeem data")
)

type Client interface {
	Do(ctx context.Context, method, path string, body interface{}, v interface{}, params map[string]string) error
}

type Redeem struct {
	client Client
}

func New(client Client) *Redeem {
	return &Redeem{client: client}
}

type RedeemResponse struct {
	Status  string          `json:"status"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

type AirtimeRedeemRequest struct {
	ChiRef       string                 `json:"chiRef"`
	PhoneNumber  string                 `json:"phoneNumber"`
	CountryToSend string                `json:"countryToSend"`
	Meta         map[string]interface{} `json:"meta,omitempty"`
	SubAccount   string                 `json:"subAccount,omitempty"`
}

type RedeemDataItem struct {
	CountryCode          string  `json:"countryCode"`
	ProductID            string  `json:"productId"`
	ValueInLocalCurrency float64 `json:"valueInLocalCurrency"`
}

type AnyRedeemRequest struct {
	ChiRef     string          `json:"chiRef"`
	RedeemData []RedeemDataItem `json:"redeemData"`
	Meta       map[string]interface{} `json:"meta,omitempty"`
	SubAccount string                 `json:"subAccount,omitempty"`
}

type ChimoneyRedeemRequest struct {
	Chimoneys  []map[string]interface{} `json:"chimoneys"`
	SubAccount string                   `json:"subAccount,omitempty"`
}

type GiftCardRedeemRequest struct {
	ChiRef        string                 `json:"chiRef"`
	RedeemOptions map[string]interface{} `json:"redeemOptions"`
	SubAccount    string                 `json:"subAccount,omitempty"`
}

type MobileMoneyRedeemRequest struct {
	ChiRef        string                 `json:"chiRef"`
	RedeemOptions map[string]interface{} `json:"redeemOptions"`
	SubAccount    string                 `json:"subAccount,omitempty"`
}

/**
 * This function redeems airtime
 * @param {string} chiRef The ID of the transaction
 * @param {string} phoneNumber The phone number to send airtime to
 * @param {string} countryToSend The country code to send airtime to
 * @param {object?} meta Additional metadata
 * @param {string?} subAccount The subAccount of the transaction
 * @returns The response from the Chimoney API
 */
func (r *Redeem) Airtime(ctx context.Context, req *AirtimeRedeemRequest) (*RedeemResponse, error) {
	if req.ChiRef == "" {
		return nil, ErrInvalidChiRef
	}

	resp := new(RedeemResponse)
	err := r.client.Do(ctx, "POST", "/redeem/airtime", req, resp, nil)
	return resp, err
}

/**
 * This function redeems any transaction
 * @param {string} chiRef The ID of the transaction
 * @param {RedeemDataItem[]} redeemData The redeem data
 * @param {object?} meta Additional metadata
 * @param {string?} subAccount The subAccount of the transaction
 * @returns The response from the Chimoney API
 */
func (r *Redeem) Any(ctx context.Context, req *AnyRedeemRequest) (*RedeemResponse, error) {
	if req.ChiRef == "" {
		return nil, ErrInvalidChiRef
	}
	if len(req.RedeemData) == 0 {
		return nil, ErrInvalidRedeemData
	}

	resp := new(RedeemResponse)
	err := r.client.Do(ctx, "POST", "/redeem/any", req, resp, nil)
	return resp, err
}

/**
 * This function redeems Chimoney
 * @param {object[]} chimoneys Array of Chimoney transactions
 * @param {string?} subAccount The subAccount of the transaction
 * @returns The response from the Chimoney API
 */
func (r *Redeem) Chimoney(ctx context.Context, chimoneys []map[string]interface{}, subAccount string) (*RedeemResponse, error) {
	if len(chimoneys) == 0 {
		return nil, ErrInvalidRedeemData
	}

	req := ChimoneyRedeemRequest{
		Chimoneys:  chimoneys,
		SubAccount: subAccount,
	}

	resp := new(RedeemResponse)
	err := r.client.Do(ctx, "POST", "/redeem/chimoney", req, resp, nil)
	return resp, err
}

/**
 * This function gets Chimoney transaction details
 * @param {string} chiRef The ID of the transaction
 * @param {string?} subAccount The subAccount of the transaction
 * @returns The response from the Chimoney API
 */
func (r *Redeem) GetChimoney(ctx context.Context, chiRef string, subAccount string) (*RedeemResponse, error) {
	if chiRef == "" {
		return nil, ErrInvalidChiRef
	}

	req := map[string]string{
		"chiRef": chiRef,
	}
	if subAccount != "" {
		req["subAccount"] = subAccount
	}

	resp := new(RedeemResponse)
	err := r.client.Do(ctx, "POST", "/redeem/chimoney/get", req, resp, nil)
	return resp, err
}

/**
 * This function redeems a gift card
 * @param {string} chiRef The ID of the transaction
 * @param {object} redeemOptions The gift card redeem options
 * @param {string?} subAccount The subAccount of the transaction
 * @returns The response from the Chimoney API
 */
func (r *Redeem) GiftCard(ctx context.Context, req *GiftCardRedeemRequest) (*RedeemResponse, error) {
	if req.ChiRef == "" {
		return nil, ErrInvalidChiRef
	}
	if len(req.RedeemOptions) == 0 {
		return nil, ErrInvalidRedeemData
	}

	resp := new(RedeemResponse)
	err := r.client.Do(ctx, "POST", "/redeem/gift-card", req, resp, nil)
	return resp, err
}

/**
 * This function redeems mobile money
 * @param {string} chiRef The ID of the transaction
 * @param {object} redeemOptions The mobile money redeem options
 * @param {string?} subAccount The subAccount of the transaction
 * @returns The response from the Chimoney API
 */
func (r *Redeem) MobileMoney(ctx context.Context, req *MobileMoneyRedeemRequest) (*RedeemResponse, error) {
	if req.ChiRef == "" {
		return nil, ErrInvalidChiRef
	}
	if len(req.RedeemOptions) == 0 {
		return nil, ErrInvalidRedeemData
	}

	resp := new(RedeemResponse)
	err := r.client.Do(ctx, "POST", "/redeem/mobile-money", req, resp, nil)
	return resp, err
}
