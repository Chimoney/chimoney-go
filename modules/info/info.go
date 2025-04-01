package info

import (
	"context"
	"encoding/json"
	"errors"
)

var (
	ErrInvalidCountry = errors.New("invalid country code")
	ErrInvalidCurrency = errors.New("invalid currency")
	ErrInvalidAmount = errors.New("invalid amount")
)

type Client interface {
	Do(ctx context.Context, method, path string, body interface{}, v interface{}, params map[string]string) error
}

type Info struct {
	client Client
}

func New(client Client) *Info {
	return &Info{client: client}
}

type InfoResponse struct {
	Status  string          `json:"status"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

/**
 * This function gets a list of supported assets
 * @returns The response from the Chimoney API
 */
func (i *Info) GetSupportedAssets(ctx context.Context) (*InfoResponse, error) {
	resp := new(InfoResponse)
	err := i.client.Do(ctx, "GET", "/info/assets", nil, resp, nil)
	return resp, err
}

/**
 * This function gets a list of supported airtime countries
 * @returns The response from the Chimoney API
 */
func (i *Info) GetAirtimeCountries(ctx context.Context) (*InfoResponse, error) {
	resp := new(InfoResponse)
	err := i.client.Do(ctx, "GET", "/info/airtime-countries", nil, resp, nil)
	return resp, err
}

/**
 * This function gets a list of banks for a country
 * @param {string} countryCode The country code (default: NG)
 * @returns The response from the Chimoney API
 */
func (i *Info) GetBanks(ctx context.Context, countryCode string) (*InfoResponse, error) {
	if countryCode == "" {
		countryCode = "NG" 
	}

	resp := new(InfoResponse)
	params := map[string]string{
		"countryCode": countryCode,
	}
	err := i.client.Do(ctx, "GET", "/info/country-banks", nil, resp, params)
	return resp, err
}

/**
 * This function converts local amount to USD
 * @param {string} currency The origin currency
 * @param {number} amount The amount to convert
 * @returns The response from the Chimoney API
 */
func (i *Info) GetLocalAmountInUSD(ctx context.Context, currency string, amount float64) (*InfoResponse, error) {
	if currency == "" {
		return nil, ErrInvalidCurrency
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	req := map[string]interface{}{
		"originCurrency": currency,
		"amount":        amount,
	}

	resp := new(InfoResponse)
	err := i.client.Do(ctx, "POST", "/info/local-amount-in-usd", req, resp, nil)
	return resp, err
}

/**
 * This function gets mobile money codes
 * @returns The response from the Chimoney API
 */
func (i *Info) GetMobileMoneyCodes(ctx context.Context) (*InfoResponse, error) {
	resp := new(InfoResponse)
	err := i.client.Do(ctx, "GET", "/info/mobile-money-codes", nil, resp, nil)
	return resp, err
}

/**
 * This function converts USD to local amount
 * @param {string} currency The destination currency
 * @param {number} amountInUSD The USD amount to convert
 * @returns The response from the Chimoney API
 */
func (i *Info) GetUSDInLocalAmount(ctx context.Context, currency string, amountInUSD float64) (*InfoResponse, error) {
	if currency == "" {
		return nil, ErrInvalidCurrency
	}
	if amountInUSD <= 0 {
		return nil, ErrInvalidAmount
	}

	req := map[string]interface{}{
		"destinationCurrency": currency,
		"amountInUSD":        amountInUSD,
	}

	resp := new(InfoResponse)
	err := i.client.Do(ctx, "POST", "/info/usd-in-local-amount", req, resp, nil)
	return resp, err
}
